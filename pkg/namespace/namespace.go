package namespace

import (
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/mmbednarek/fragma/pkg/linux"
	"github.com/mmbednarek/fragma/pkg/process"
)

type Resource int

const (
	ResourceNetwork Resource = iota
	ResourceIPC
	ResourceUTS
	ResourceMount
	ResourceProcess
	ResourceUser
	ResourceCGroup
)

func (r Resource) NsType() int {
	switch r {
	case ResourceNetwork:
		return syscall.CLONE_NEWNET
	case ResourceIPC:
		return syscall.CLONE_NEWIPC
	case ResourceUTS:
		return syscall.CLONE_NEWUTS
	case ResourceMount:
		return syscall.CLONE_NEWNS
	case ResourceProcess:
		return syscall.CLONE_NEWPID
	case ResourceUser:
		return syscall.CLONE_NEWUSER
	case ResourceCGroup:
		return linux.CLONE_NEWCGROUP
	}
	return 0
}

func (r Resource) String() string {
	switch r {
	case ResourceNetwork:
		return "net"
	case ResourceIPC:
		return "ipc"
	case ResourceUTS:
		return "uts"
	case ResourceProcess:
		return "pid"
	case ResourceMount:
		return "mnt"
	case ResourceUser:
		return "user"
	case ResourceCGroup:
		return "cgroup"
	}
	return ""
}

type Namespace struct {
	Path string
	Id   uint64
	Kind Resource
}

func FindNamespaceByPid(resource Resource, pid process.Pid) (Namespace, error) {
	path := fmt.Sprintf("/proc/%d/ns/%s", pid, resource.String())
	fileInfo, err := os.Stat(path)
	if err != nil {
		return Namespace{}, err
	}

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return Namespace{}, errors.New("not a linux process")
	}

	return Namespace{Path: path, Id: stat.Ino, Kind: resource}, nil
}

func SetNamespace(ns *Namespace) error {
	fd, err := syscall.Open(ns.Path, syscall.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	if err := linux.Setns(fd, 0); err != nil {
		return err
	}

	return nil
}
