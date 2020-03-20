package winapi

import (
	"errors"
	"fmt"
	"sync"
	"syscall"
	"unsafe"
)

type DriveType int

const (
	DRIVE_UNKNOWN     DriveType = 0
	DRIVE_NO_ROOT_DIR DriveType = 1
	DRIVE_REMOVABLE   DriveType = 2
	DRIVE_FIXED       DriveType = 3
	DRIVE_REMOTE      DriveType = 4
	DRIVE_CDROM       DriveType = 5
	DRIVE_RAMDISK     DriveType = 6
)

type winApiMgr struct {
	kernel32      *syscall.DLL
	diskSpaceProc *syscall.Proc
	driveTypeProc *syscall.Proc
}

var _winApiMgrIns *winApiMgr = nil

var _once sync.Once

func Setup() error {
	ins := WinApiMgrIns()
	if ins.kernel32 = syscall.MustLoadDLL("kernel32.dll"); ins.kernel32 == nil {
		return errors.New("win api load kernel32 dll failed")
	}

	if ins.diskSpaceProc = ins.kernel32.MustFindProc("GetDiskFreeSpaceExW"); ins.diskSpaceProc == nil {
		return errors.New("win api load GetDiskFreeSpaceExW proc failed")
	}

	if ins.driveTypeProc = ins.kernel32.MustFindProc("GetDriveTypeW"); ins.driveTypeProc == nil {
		return errors.New("win api load GetDriveTypeW proc failed")
	}

	return nil
}

func WinApiMgrIns() *winApiMgr {
	_once.Do(func() {
		_winApiMgrIns = &winApiMgr{}
	})

	return _winApiMgrIns
}

func (mgr *winApiMgr) GetDiskFreeSpace(path string) (freeBytes int64, err error) {
	if mgr.diskSpaceProc == nil {
		return 0, errors.New("diskSpaceProc invalid")
	}

	freeBytes, totNumOfBytes, totNumOfFreeBytes := int64(0), int64(0), int64(0)
	ret, _, lasterr := mgr.diskSpaceProc.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))),
		uintptr(unsafe.Pointer(&freeBytes)), uintptr(unsafe.Pointer(&totNumOfBytes)), uintptr(unsafe.Pointer(&totNumOfFreeBytes)))
	if lasterr != nil {
		fmt.Println("getDiskFreeSpace proc call lasterr:", lasterr.Error())
	}

	if ret == 0 {
		return 0, errors.New("getDiskFreeSpace failed")
	}

	return
}

func (mgr *winApiMgr) GetDriveType(path string) (driType DriveType, err error) {
	if mgr.driveTypeProc == nil {
		return DRIVE_UNKNOWN, errors.New("driveTypeProc invalid")
	}

	ret, _, lasterr := mgr.driveTypeProc.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))))
	if lasterr != nil {
		fmt.Println("getDriveType proc call lasterr:", lasterr.Error())
	}

	if ret == 0 {
		return DRIVE_UNKNOWN, errors.New("getDriveType failed")
	}

	switch ret {
	case 0:
		driType = DRIVE_UNKNOWN
	case 1:
		driType = DRIVE_NO_ROOT_DIR
	case 2:
		driType = DRIVE_REMOVABLE
	case 3:
		driType = DRIVE_FIXED
	case 4:
		driType = DRIVE_REMOTE
	case 5:
		driType = DRIVE_CDROM
	case 6:
		driType = DRIVE_RAMDISK
	}

	return
}
