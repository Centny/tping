package main

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"os"
	"time"
)

func RunW(h string, d time.Duration, t int) error {
	beg := util.Now()
	spath := fmt.Sprintf("%v/%v.tping", os.TempDir(), beg)
	defer os.Remove(spath)
	var err error
	var ended bool = false
	end_c := make(chan int)
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
		time.Sleep(d)
	} else {
		<-end_c
	}
	if err != nil {
		return err
	}
	if ended {
		fmt.Println(util.Now() - beg)
		return nil
	} else {
		return util.Err("timeout")
	}
}
