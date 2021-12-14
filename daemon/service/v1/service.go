package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	core "github.com/mmbednarek/fragma/api/fragma/core/v1"
	"github.com/mmbednarek/fragma/pkg/linux"
	"github.com/mmbednarek/fragma/pkg/log"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) RunApplication(ctx context.Context, volume *core.Volume, application *core.Application, options *core.RunOptions) error {
	loopPath, err := linux.LoopSetupDevice(volume.Path)
	if err != nil {
		return err
	}
	defer linux.LoopClear(loopPath)

	mount, err := linux.MountDevice(loopPath)
	if err != nil {
		return err
	}
	defer mount.Unmount()

	cmd := exec.Command(application.Path)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	cmd.Env = []string{
		"PS1=[fragma] # ",
		"TERM=xterm",
		"HOME=/root",
	}
	for key, value := range options.Environment {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	gid := syscall.Getgid()
	uid := syscall.Getuid()

	log.With(ctx, "gid", gid, "uid", uid).Info("running with")

	if err := syscall.Mount("/dev", mount.Path+"/dev", "", syscall.MS_BIND, ""); err != nil {
		os.Exit(1)
	}
	defer syscall.Unmount(mount.Path+"/dev", 0)

	if err := syscall.Mount("/proc", mount.Path+"/proc", "", syscall.MS_BIND, ""); err != nil {
		os.Exit(1)
	}
	defer syscall.Unmount(mount.Path+"/proc", 0)

	cmd.Dir = "/root"
	cmd.Args = options.Arguments
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Chroot: mount.Path,
		Credential: &syscall.Credential{
			Uid: 0,
			Gid: 0,
		},
		Cloneflags:   syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID,
		Unshareflags: linux.CLONE_FS,
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cmd.Run: %w", err)
	}

	//ns, err := namespace.FindNamespaceByPid(namespace.ResourceNetwork, process.Pid(cmd.Process.Pid))
	//if err != nil {
	//	return err
	//}
	//
	//log.Printf("process namespace: %v\n", ns.Id)
	//
	//if err := namespace.SetNamespace(&ns); err != nil {
	//	return err
	//}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("cmd.Wait: %w", err)
	}
	return nil
}
