package linux

import (
	"fmt"
	"syscall"
)

const (
	SYS_SETNS = 308

	CLONE_NEWCGROUP = 0x02000000
	CLONE_FS        = 0x00000200
	CLONE_FILES     = 0x00000400
	CLONE_NEWNS     = 0x00020000
	CLONE_VM        = 0x00000100
)

func Setns(fd int, nstype int) error {
	_, _, err := syscall.RawSyscall(SYS_SETNS, uintptr(fd), uintptr(nstype), 0)
	if err != 0 {
		return fmt.Errorf("could not call setns: %v", err)
	}
	return nil
}
