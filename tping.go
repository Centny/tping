package main

import (
	"bufio"
	"fmt"
	"github.com/Centny/gwf/log"
	"io"
	"os"
	"runtime"
	"strconv"
	"time"
)

var ping_l bool = true

func Usage() {
	fmt.Println(`Usage:
	tping -m S [-l log file] [-p listen port, default :9100] [-a show all long, default Y]
	tping [-m C] [-h the target host] [-l log file] [-d ping delay,default 1000ms] [-t ping times, default -1] [-a show all log, default Y]
	tping [-m W] [-h the target host] [-l log file] [-d ping delay,default 0ms] [-t ping times, default 1] [-a show all log, default Y]
		`)
}
func Run(args []string) {
	var p string = ":9100"
	var t int = -1
	var d time.Duration = 0
	var l string = ""
	var h string = ""
	var m string = "C"
	alen := len(args) - 1
	for i := 1; i < alen; i++ {
		switch args[i] {
		case "-l":
			l = args[i+1]
		case "-d":
			delay, err := strconv.ParseInt(args[i+1], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			d = time.Duration(delay) * time.Millisecond
		case "-t":
			tc, err := strconv.ParseInt(args[i+1], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			t = int(tc)
		case "-p":
			p = args[i+1]
		case "-h":
			h = args[i+1]
		case "-m":
			m = args[i+1]
		case "-a":
			ping_l = "Y" == args[i+1]
		}
	}
	if len(l) > 0 {
		f, err := os.OpenFile(l, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		w := bufio.NewWriter(f)
		defer func() {
			w.Flush()
			f.Close()
		}()
		log.SetWriter(io.MultiWriter(f, os.Stdout))
	}
	var err error
	switch m {
	case "S":
		err = RunS(p)
	case "C":
		if len(h) < 1 {
			Usage()
			return
		}
		if d < 1 {
			d = time.Second
		}
		err = RunC(h, t, d)
	case "W":
		if t < 1 {
			t = 1
		}
		err = RunW(h, d, t)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	Run(os.Args)
}
