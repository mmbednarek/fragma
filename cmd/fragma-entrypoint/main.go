package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/mmbednarek/fragma/pkg/linux"
)

func die(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
func main() {
	if len(os.Args) < 2 {
		die("invalid number arguments: %s", os.Args)
		os.Exit(1)
	}

	binPath := os.Args[1]

	if err := syscall.Chdir("/"); err != nil {
		die("could not change dir: %s", err)
	}

	if err := linux.SetupCharacterDevices(); err != nil {
		die("could not setup character devices: %s", err)
	}

	if err := syscall.Mount("devpts", "/dev/pts", "devpts", 0, ""); err != nil {
		die("could not mount devpts: %s", err)
	}

	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		die("could not mount proc: %s", err)
	}

	if err := syscall.Mount("sysfs", "/sys", "sysfs", 0, ""); err != nil {
		die("could not mount sysfs: %s", err)
	}

	syscall.Sync()

	if err := linux.SetupDeviceSymlinks(); err != nil {
		fmt.Printf("cannot create symlink: %s\n", err)
		//die("could not setup device links: %s", err)
	}

	if err := syscall.Sethostname([]byte("fragma-pod")); err != nil {
		die("could not set hostname: %s", err)
	}

	terminal, err := linux.NewTerminal()
	if err != nil {
		die("could not create terminal: %s", err)
	}
	defer terminal.Close()
	fmt.Printf("%s\n", terminal.SlavePath)

	slaveFile, err := terminal.OpenSlave()
	if err != nil {
		die("could not open terminal file: %s", err)
	}
	defer slaveFile.Close()

	hostTTYAttrib, err := linux.Attr(os.Stdin)
	if err != nil {
		die("could not get host terminal attrib: %s", err)
	}
	hostTTYAttrib.Raw()

	if err := hostTTYAttrib.Winsz(os.Stdin); err != nil {
		die("could not get host terminal window size: %s", err)
	}

	fmt.Printf("terminal info (w = %d, h = %d))\n", hostTTYAttrib.Wz.WsRow, hostTTYAttrib.Wz.WsCol)

	if err := hostTTYAttrib.Set(slaveFile); err != nil {
		die("could not set slave terminal attrib: %s", err)
	}

	if err := hostTTYAttrib.Setwinsz(os.Stdin); err != nil {
		die("could not set slave terminal window size: %s", err)
	}

	go func() {
		for {
			n, err := io.Copy(os.Stdout, terminal.MasterFile)
			if err != nil {
				fmt.Printf("copy error: %s\n", err)
			}
			if n > 0 {
				fmt.Printf("wrote %d bytes to stdout\n", n)
			}
		}
	}()
	go func() {
		for {
			n, err := io.Copy(terminal.MasterFile, os.Stdin)
			if err != nil {
				fmt.Printf("copy error: %s\n", err)
			}
			if n > 0 {
				fmt.Printf("wrote %d bytes to term\n", n)
			}
		}
	}()

	cmd := exec.Command(binPath)
	//cmd.Stderr = os.Stderr
	//cmd.Stdin = os.Stdin
	//cmd.Stdout = os.Stdout
	cmd.Stderr = slaveFile
	cmd.Stdin = slaveFile
	cmd.Stdout = slaveFile
	if err := cmd.Run(); err != nil {
		die("could not execute the command: %s", err)
	}
}
