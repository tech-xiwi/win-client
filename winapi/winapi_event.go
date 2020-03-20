package winapi

import (
	"errors"
	"syscall"

	"golang.org/x/sys/windows"
)

var EVENT_ALL_ACCESS uint32 = 0x1F0003

type Event struct {
	h windows.Handle
}

func NewEvent(name string) (*Event, error) {
	h, err := windows.CreateEvent(nil, 0, 0, syscall.StringToUTF16Ptr(name))
	if err != nil {
		return nil, err
	}
	return &Event{h: h}, nil
}

func OpenEvent(name string) (*Event, error) {
	h, err := windows.OpenEvent(EVENT_ALL_ACCESS, false, syscall.StringToUTF16Ptr(name))
	if err != nil {
		return nil, err
	}
	return &Event{h: h}, nil
}

func (e *Event) Close() error {
	return windows.CloseHandle(e.h)
}

func (e *Event) Set() error {
	return windows.SetEvent(e.h)
}

func (e *Event) Wait() error {
	s, err := windows.WaitForSingleObject(e.h, windows.INFINITE)
	switch s {
	case windows.WAIT_OBJECT_0:
		break
	case windows.WAIT_FAILED:
		return err
	default:
		return errors.New("unexpected result from WaitForSingleObject")
	}
	return nil
}
