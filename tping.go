package main

import (
	"bufio"
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	TP_S = "S"
	TP_C = "C"
	TP_W = "W"
	TP_R = "R"
	TP_J = "J"
)

var exit func(code int) = os.Exit
var ping_l bool = true

func Usage() {
	fmt.Println(`Usage:
 args:
  -l log file
  -a show all log, default Y

 run mode:
  tping -m S [-p listen port, default :9100]
  tping -m C [-h the target host] [-d ping delay,default 1000ms] [-t ping times, default -1]
  tping -m W [-h the target host] [-d ping delay,default 0ms] [-t ping times, default 1]
  tping -m R [-r target run command] [-d ping delay,default 0ms] [-t ping times, default 1]
  tping -m J [-j json configure file path] [-e emma report file path]

 example:
  'tping [-m S]' will listen tcp port on :9100
  'tping [-m W] http://www.bing.com' testing connect to www.bing.com by http and delay 1s.

 json format for J mode:
  [{
  	"name":"task name",
  	"type":"W",   //task type on W|R,W is http task,R is command task.
  	"delay":1000, //millisecond
  	"times":1,    //exec times.
  	"host":"http://www.bing.com"
  },{
  	"name":"task name2",
  	"type":"R",   //task type on W|R,W is http task,R is command task.
  	"delay":1000, //millisecond
  	"times":1,    //exec times.
  	"cmds":"echo abc" //command and arguments.
  }]
		`)
}
func Run(args []string) {
	var p string = ":9100"
	var t int = -1
	var d time.Duration = 0
	var l string = ""
	var h string = ""
	var m string = ""
	var e string = ""
	var r string = ""
	var j string = ""
	alen := len(args) - 1
	for i := 1; i < alen; i++ {
		switch args[i] {
		case "-l":
			l = args[i+1]
			i++
		case "-d":
			delay, err := strconv.ParseInt(args[i+1], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			d = time.Duration(delay) * time.Millisecond
			i++
		case "-t":
			tc, err := strconv.ParseInt(args[i+1], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			t = int(tc)
			i++
		case "-p":
			p = args[i+1]
			i++
		case "-h":
			h = args[i+1]
			i++
		case "-m":
			m = args[i+1]
			i++
		case "-a":
			ping_l = "Y" == args[i+1]
			i++
		case "-r":
			r = args[i+1]
			i++
		case "-e":
			e = args[i+1]
			i++
		case "-j":
			j = args[i+1]
			i++
		default:
			if len(h) < 1 {
				h = args[i]
			}
		}
	}
	if alen == 1 && len(h) < 1 {
		h = args[alen]
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
	if len(m) < 1 {
		if len(h) > 0 {
			if strings.HasPrefix(h, "http://") {
				m = TP_W
			} else {
				m = TP_C
			}
		} else {
			m = TP_S
		}
	}
	var err error
	var delay int64
	switch m {
	case TP_S:
		err = RunS(p)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			exit(1)
		}
	case TP_C:
		if len(h) < 1 {
			Usage()
			return
		}
		if d < 1 {
			d = time.Second
		}
		RunC(h, t, d)
	case TP_W:
		if len(h) < 1 {
			Usage()
			return
		}
		if t < 1 {
			t = 1
		}
		delay, err = RunW(h, d, t)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			exit(1)
		} else {
			fmt.Println(delay)
		}
	case TP_R:
		if len(r) < 1 {
			Usage()
			return
		}
		if t < 1 {
			t = 1
		}
		delay, err = RunR(r, d, t)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			exit(1)
		} else {
			fmt.Println(delay)
		}
	case TP_J:
		if len(j) < 1 || len(e) < 1 {
			Usage()
			return
		}
		err = RunJ(j, e)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			exit(1)
		} else {
			fmt.Println(delay)
		}
	}
}
func main() {
	runtime.GOMAXPROCS(util.CPU())
	Run(os.Args)
}
