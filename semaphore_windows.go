package semaphore

import (
	"fmt"
	"syscall"
	"unsafe"
)

type Semaphore struct {
	hdl  uintptr
	Name string
}

var (
	kernel32             = syscall.NewLazyDLL("kernel32.dll")
	procCreateSemaphore  = kernel32.NewProc("CreateSemaphoreW")
	procOpenSemaphore    = kernel32.NewProc("OpenSemaphoreW")
	procCloseHandle      = kernel32.NewProc("CloseHandle")
	procReleaseSemaphore = kernel32.NewProc("ReleaseSemaphore")
)

const (
	SEMAPHORE_ALL_ACCESS = uint32(0x1F0003)
)

func NewSemaphore(name string, cnt int) (*Semaphore, error) {
	var err error
	var p *Semaphore
	p = &Semaphore{}
	p.Name = name
	p.hdl, _, err = procCreateSemaphore.Call(0, uintptr(cnt), uintptr(cnt), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(name))))
	switch int(err.(syscall.Errno)) {
	case 0:
		return p, nil
	}

	p.hdl, _, err = procOpenSemaphore.Call(uintptr(SEMAPHORE_ALL_ACCESS), 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(name))))
	switch int(err.(syscall.Errno)) {
	case 0:
		return p, nil
	}
	return nil, err
}

func (psema *Semaphore) Wait(mills int) error {
	var err error
	var dmills uint32
	var evt uint32

	if psema.hdl == uintptr(0) {
		return fmt.Errorf("not valid hdl")
	}

	if mills < 0 {
		dmills = syscall.INFINITE
	} else {
		dmills = uint32(mills)
	}

	evt, err = syscall.WaitForSingleObject(syscall.Handle(psema.hdl), dmills)
	if evt != syscall.WAIT_OBJECT_0 {
		err = fmt.Errorf("wait error evt[%v] [%v]", evt, err)
		return err
	}
	return nil
}

func (psema *Semaphore) Release() error {
	var err error
	var ret uintptr
	if psema.hdl == uintptr(0) {
		return fmt.Errorf("not valid hdl")
	}
	ret, _, err = procReleaseSemaphore.Call(uintptr(psema.hdl), 1, uintptr(0))
	if ret == 0 {
		err = fmt.Errorf("release [%s] error[%v]", psema.Name, int(err.(syscall.Errno)))
		return err
	}
	return nil
}

func (psema *Semaphore) Close() {
	if psema.hdl != uintptr(0) {
		procCloseHandle.Call(uintptr(psema.hdl))
	}
	psema.hdl = uintptr(0)
}
