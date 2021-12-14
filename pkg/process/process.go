package process

import (
	"os/exec"
)

type RunEnvironment struct {
	Command *exec.Cmd
}

type Pid int

func NewRunEnvironment(path string, args []string) RunEnvironment {
	//syscall.Open()
	return RunEnvironment{}
}
