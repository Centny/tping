package main

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/tutil"
	"github.com/Centny/gwf/util"
	"time"
)

type Task struct {
	Name  string `m2s:"name",json:"name"`
	Type  string `m2s:"type",json:"type"`
	Host  string `m2s:"host",json:"host"`
	Cmds  string `m2s:"cmds",json:"cmds"`
	Delay int64  `m2s:"delay",json:"delay"`
	Times int    `m2s:"times",json:"times"`
}

func RunJ(in, e string) error {
	var tasks []Task
	err := util.J2Ss_f(in, &tasks)
	if err != nil {
		return err
	}
	for _, task := range tasks {
		if task.Times < 1 {
			task.Times = 1
		}
		switch task.Type {
		case TP_W:
			if len(task.Host) < 1 {
				log.W("run task(%v) by type(%v),delay(%v),times(%v) err:host is empty ", task.Name, task.Type, task.Delay, task.Times)
				break
			}
			log.D("running task(%v) by type(%v),host(%v),delay(%v),times(%v)", task.Name, task.Type, task.Host, task.Delay, task.Times)
			delay, err := RunW(task.Host, time.Duration(task.Delay)*time.Millisecond, task.Times)
			var class, block, line string = "1/1", "", ""
			if err == nil {
				if task.Delay < 1 {
					line = fmt.Sprintf("%v/%v", delay, delay)
				} else {
					if task.Delay < delay {
						task.Delay = delay
					}
					line = fmt.Sprintf("%v/%v", task.Delay-delay, task.Delay)
				}
				block = fmt.Sprintf("%v/%v", delay, delay)
			} else {
				log.E("run task(%v) by type(%v) err:%v", task.Name, task.Type, err.Error())
				class, block, line = "0/1", "0/1", "0/1"
			}
			err = tutil.Emma(e, task.Name, class, class, block, line)
			if err != nil {
				log.E("run task(%v) by type(%v),delay(%v),times(%v) err:append emma report err(%v)", task.Name, task.Type, task.Delay, task.Times, err.Error())
				return err
			}
			log.D("task(%v) done by type(%v),host(%v), delay:%v", task.Name, task.Type, task.Host, delay)
		case TP_R:
			if len(task.Cmds) < 1 {
				log.W("run task(%v) by type(%v),delay(%v),times(%v) err:cmds is empty ", task.Name, task.Type, task.Delay, task.Times)
				break
			}
			log.D("running task(%v) by type(%v),cmds(%v),delay(%v),times(%v)", task.Name, task.Type, task.Cmds, task.Delay, task.Times)
			delay, err := RunR(task.Cmds, time.Duration(task.Delay)*time.Millisecond, task.Times)
			var class, block, line string = "1/1", "", ""
			if err == nil {
				if task.Delay < 1 {
					line = fmt.Sprintf("%v/%v", delay, delay)
				} else {
					if task.Delay < delay {
						task.Delay = delay
					}
					line = fmt.Sprintf("%v/%v", task.Delay-delay, task.Delay)
				}
				block = fmt.Sprintf("%v/%v", delay, delay)
			} else {
				log.E("run task(%v) by type(%v),delay(%v),times(%v) err:%v", task.Name, task.Type, task.Delay, task.Times, err.Error())
				class, block, line = "0/1", "0/1", "0/1"
			}
			err = tutil.Emma(e, task.Name, class, class, block, line)
			if err != nil {
				log.E("run task(%v) by type(%v) err:append emma report err(%v)", task.Name, task.Type, err.Error())
				return err
			}
			log.D("task(%v) done by type(%v),cmds(%v), delay:%v", task.Name, task.Type, task.Cmds, delay)
		default:
			log.W("run task(%v) by type(%v) err:unknow type ", task.Name, task.Type)
		}
	}
	return nil
}
