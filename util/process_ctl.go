package util

import (
	"fmt"
	"github.com/tech-xiwi/win-client/winapi"
)

type ServerExitCallback func()

type ServerLiveMgr struct {
	event     *winapi.Event
	callbacks []ServerExitCallback
}

func NewServerLiveMgr(name string) *ServerLiveMgr {
	event, err := winapi.NewEvent(name)
	if err != nil {
		fmt.Println("new event failed err:", err.Error())
		return nil
	}
	return &ServerLiveMgr{
		event: event,
	}
}

func (mgr *ServerLiveMgr) RegExitCallback(callback ServerExitCallback) {
	mgr.callbacks = append(mgr.callbacks, callback)
}

func (mgr *ServerLiveMgr) WaitToQuit() error {
	err := mgr.event.Wait()
	if err != nil {
		return err
	}

	for _, callback := range mgr.callbacks {
		callback()
	}

	return mgr.event.Close()
}
