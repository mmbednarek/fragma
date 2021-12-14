package linux

import (
	"fmt"
	"syscall"
)

type Error struct {
	Errno syscall.Errno
}

func (e Error) Error() string {
	return fmt.Sprintf("linux syscall error: %s", e.Errno.Error())
}
