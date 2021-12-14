package main

import (
	"os"
	"os/exec"
	"syscall"
)

func main() {
	if len(os.Args) < 3 {
		os.Exit(1)
	}

	binPath := os.Args[1]
	rootPath := os.Args[2]

	if err := syscall.Chdir(rootPath); err != nil {
		os.Exit(1)
	}

	if err := syscall.PivotRoot(".", "/opt/root-fs"); err != nil {
		os.Exit(1)
	}

	cmd := exec.Command(binPath)
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}
