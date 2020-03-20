package main

import (
	"flag"
	"fmt"
	"github.com/tech-xiwi/win-client/util"
	"github.com/tech-xiwi/win-client/winapi"
	"golang.org/x/sys/windows"
	"os"
	"time"
)

var (
	quit        = flag.Bool("quit", false, "quit")
	serviceName = "ebe9e163-0eaf-47f1-9239-17f10113539b"
)

func main() {
	flag.Parse()
	event, err := winapi.OpenEvent(serviceName)
	if *quit {
		if err != nil {
			fmt.Println("service quit open event failed:", err.Error())
			return
		}
		_ = event.Set()
		_ = event.Close()
		return
	}

	//single process
	if err == nil {
		fmt.Println("current process event exist;have another same process!!!")
		_ = event.Close()
		return
	}

	//此后能正常启动
	liveMgr := util.NewServerLiveMgr(serviceName)
	if liveMgr == nil {
		panic("new server live mgr failed")
	}

	if err := winapi.Setup(); err != nil {
		fmt.Println("winapi setup failed")
	}

	//todo start your business
	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Println("Hi, Current Time is:", time.Now())
			fmt.Println(windows.Getenv("GOPROXY"))
		}
	}()

	//quit
	liveMgr.RegExitCallback(func() {
		fmt.Println("exitCallback")
		for {
			fmt.Println("-------------Hi, Current Time is:", time.Now())
			time.Sleep(10 * time.Millisecond)
			os.Exit(0)
		}
	})

	_ = liveMgr.WaitToQuit()
}
