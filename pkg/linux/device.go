package linux

import (
	"os"
	"syscall"
)

func SetupDeviceSymlinks() error {
	if err := CreateSymlink("/dev/pts/ptmx", "/dev/ptmx"); err != nil {
		return err
	}
	if err := CreateSymlink("/proc/self/fd/0", "/dev/stdin"); err != nil {
		return err
	}
	if err := CreateSymlink("/proc/self/fd/1", "/dev/stdout"); err != nil {
		return err
	}
	if err := CreateSymlink("/proc/self/fd/2", "/dev/stderr"); err != nil {
		return err
	}

	return nil
}

func CreateSymlink(source string, target string) error {
	if _, err := os.Stat(target); err == nil {
		// TODO: maybe check if the device is ok and leave it
		if err := os.Remove(target); err != nil {
			return err
		}
	}
	if err := syscall.Symlink(source, target); err != nil {
		return err
	}
	return nil
}
