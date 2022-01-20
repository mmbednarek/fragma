package linux

import (
	"os"
	"syscall"
)

type DeviceType int

const (
	// Basic utility character device
	NullDevice    DeviceType = (1 << 8) | (3 & 0xff) | ((3 & 0xfff00) << 12)
	ZeroDevice    DeviceType = (1 << 8) | (5 & 0xff) | ((5 & 0xfff00) << 12)
	FullDevice    DeviceType = (1 << 8) | (7 & 0xff) | ((7 & 0xfff00) << 12)
	RandomDevice  DeviceType = (1 << 8) | (8 & 0xff) | ((8 & 0xfff00) << 12)
	URandomDevice DeviceType = (1 << 8) | (9 & 0xff) | ((9 & 0xfff00) << 12)

	// Console character devices
	TTYDevice     DeviceType = (5 << 8) | (0 & 0xff) | ((0 & 0xfff00) << 12)
	ConsoleDevice DeviceType = (5 << 8) | (1 & 0xff) | ((1 & 0xfff00) << 12)
	PTMXDevice    DeviceType = (5 << 8) | (2 & 0xff) | ((2 & 0xfff00) << 12)
)

func (d DeviceType) Path() string {
	switch d {
	case NullDevice:
		return "/dev/null"
	case ZeroDevice:
		return "/dev/zero"
	case FullDevice:
		return "/dev/full"
	case RandomDevice:
		return "/dev/random"
	case URandomDevice:
		return "/dev/urandom"
	case TTYDevice:
		return "/dev/tty"
	case ConsoleDevice:
		return "/dev/console"
	case PTMXDevice:
		return "/dev/ptmx"
	}
	return ""
}

func MakeCharacterDevice(d DeviceType) error {
	return syscall.Mknod(d.Path(), 0666|syscall.S_IFCHR, int(d))
}

func SetupCharacterDevices() error {
	characterDevices := []DeviceType{
		NullDevice,
		ZeroDevice,
		FullDevice,
		RandomDevice,
		URandomDevice,
		TTYDevice,
		ConsoleDevice,
	}

	for _, deviceType := range characterDevices {
		if _, err := os.Stat(deviceType.Path()); err == nil {
			// TODO: maybe check if the device is ok and leave it
			if err := os.Remove(deviceType.Path()); err != nil {
				return err
			}
		}

		if err := MakeCharacterDevice(deviceType); err != nil {
			return err
		}
	}
	return nil
}
