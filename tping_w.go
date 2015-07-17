package main

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"os"
	"time"
)

func RunW(h string, d time.Duration, t int) (int64, error) {
	beg := util.Now()
	spath := fmt.Sprintf("%v/%v.tping", os.TempDir(), beg)
	defer os.Remove(spath)
	var err error
	var ended bool = false
	end_c := make(chan int, 2)
	go func() {
		for i := 0; i < t; i++ {
			err = util.DLoad(spath, h)
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
