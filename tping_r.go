package main

import (
	"github.com/Centny/gwf/util"
	"os/exec"
	"runtime"
	"time"
)

func RunR(args string, d time.Duration, t int) (int64, error) {
	beg := util.Now()
	var err error
	var ended bool = false
	end_c := make(chan int)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", args)
	} else {
		cmd = exec.Command("/bin/bash", "-c", args)
	}
	go func() {
		for i := 0; i < t; i++ {
			_, err = cmd.Output()
			if err != nil {
				break
			}
		}
		ended = true
		end_c <- 0
	}()
	if d > 0 {
		go func() {
			time.Sleep(d)
			end_c <- 0
		}()
	}
	<-end_c
	if err != nil {
		return 0, err
	}
	if ended {
		return util.Now() - beg, nil
	} else {
		return 0, util.Err("timeout")
	}
}
