package linux

/*
#include <linux/loop.h>

typedef struct loop_info64 loop_info64_t;
typedef struct loop_config loop_config_t;
*/
import "C"
import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	LOOP_CONFIGURE    = 0x4C0A
	LOOP_CTL_GET_FREE = 0x4C82
	LOOP_SET_FD       = 0x4C00
	LOOP_SET_STATUS64 = 0x4C04
	LOOP_CLR_FD       = 0x4C01
)

type LoopInfo struct {
	Device         uint64 /* ioctl r/o */
	Inode          uint64 /* ioctl r/o */
	RDevice        uint64 /* ioctl r/o */
	Offset         uint64
	SizeLimit      uint64 /* bytes, 0 == max available */
	Number         uint32 /* ioctl r/o */
	EncryptType    uint32
	EncryptKeySize uint32 /* ioctl w/o */
	Flags          uint32
	FileName       [64]byte
	CryptName      [64]byte
	EncryptKey     [32]byte /* ioctl w/o */
	Init           [2]uint64
}

func (l LoopInfo) ToCStruct() C.loop_info64_t {
	result := C.loop_info64_t{}
	result.lo_device = C.__u64(l.Device)
	result.lo_inode = C.__u64(l.Inode)
	result.lo_rdevice = C.__u64(l.RDevice)
	result.lo_offset = C.__u64(l.Offset)
	result.lo_sizelimit = C.__u64(l.SizeLimit)
	result.lo_number = C.__u32(l.Number)
	result.lo_encrypt_type = C.__u32(l.EncryptType)
	result.lo_encrypt_key_size = C.__u32(l.EncryptKeySize)
	result.lo_flags = C.__u32(l.Flags)
	for i := 0; i < 64; i++ {
		result.lo_file_name[i] = C.__u8(l.FileName[i])
		result.lo_crypt_name[i] = C.__u8(l.CryptName[i])
	}
	for i := 0; i < 32; i++ {
		result.lo_encrypt_key[i] = C.__u8(l.EncryptKey[i])
	}
	result.lo_init[0] = C.__u64(l.Init[0])
	result.lo_init[1] = C.__u64(l.Init[1])
	return result
}

type LoopConfig struct {
	Fd        uint32
	BlockSize uint32
	Info      LoopInfo
}

func (lc LoopConfig) ToCStruct() C.loop_config_t {
	result := C.loop_config_t{}
	result.fd = C.__u32(lc.Fd)
	result.block_size = C.__u32(lc.BlockSize)
	result.info = lc.Info.ToCStruct()
	for i := 0; i < 8; i++ {
		result.__reserved[i] = 0
	}
	return result
}

func LoopGetFreeDevice() (string, error) {
	ctlFd, err := syscall.Open("/dev/loop-control", syscall.O_RDWR, 0644)
	if err != nil {
		return "", err
	}
	defer syscall.Close(ctlFd)

	r1, _, errno := syscall.RawSyscall(syscall.SYS_IOCTL, uintptr(ctlFd), LOOP_CTL_GET_FREE, 0)
	if errno != 0 {
		return "", Error{Errno: errno}
	}

	return fmt.Sprintf("/dev/loop%d", int(r1)), nil
}

func LoopSetFd(deviceFd int, fileFd int) error {
	_, _, errno := syscall.RawSyscall(syscall.SYS_IOCTL, uintptr(deviceFd), LOOP_SET_FD, uintptr(fileFd))
	if errno != 0 {
		return Error{Errno: errno}
	}
	return nil
}

func LoopConfigure(deviceFd int, cfg LoopConfig) error {
	ccfg := cfg.ToCStruct()
	_, _, errno := syscall.RawSyscall(syscall.SYS_IOCTL, uintptr(deviceFd), LOOP_CONFIGURE, uintptr(unsafe.Pointer(&ccfg)))
	if errno != 0 {
		return Error{Errno: errno}
	}
	return nil
}

func LoopSetupDevice(file string) (string, error) {
	fileFd, err := syscall.Open(file, syscall.O_RDWR, 0644)
	if err != nil {
		return "", err
	}
	defer syscall.Close(fileFd)

	devicePath, err := LoopGetFreeDevice()
	if err != nil {
		return "", err
	}

	deviceFd, err := syscall.Open(devicePath, syscall.O_RDWR, 0644)
	if err != nil {
		return "", err
	}
	defer syscall.Close(deviceFd)

	if err := LoopSetFd(deviceFd, fileFd); err != nil {
		return "", err
	}

	return devicePath, nil
}

func LoopClear(device string) error {
	deviceFd, err := syscall.Open(device, syscall.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer syscall.Close(deviceFd)

	if _, _, errno := syscall.RawSyscall(syscall.SYS_IOCTL, uintptr(deviceFd), LOOP_CLR_FD, 0); errno != 0 {
		return Error{Errno: errno}
	}

	return nil
}
