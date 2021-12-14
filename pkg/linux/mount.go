package linux

import (
	"os"
	"syscall"

	"github.com/mmbednarek/fragma/pkg/util"
)

type Mount struct {
	Name string
	Path string
}

func MountDevice(path string) (Mount, error) {
	name := util.String(6)
	mountPath := "/opt/frama/mount/" + name + "/"

	if err := os.Mkdir(mountPath, 0755); err != nil {
		if !os.IsExist(err) {
			return Mount{}, err
		}
	}

	err := syscall.Mount(path, mountPath, "ext4", 0, "")
	if err != nil {
		return Mount{}, err
	}

	return Mount{Name: name, Path: mountPath}, nil
}

func (m *Mount) Unmount() error {
	if err := syscall.Unmount(m.Path, 0); err != nil {
		return err
	}
	return nil
}
