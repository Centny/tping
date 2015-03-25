package main

import (
	"bufio"
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"time"
)

var tsig = make(chan os.Signal, 1)
var ping_l bool = true

type Val struct {
	V string `json:"v"`
}

type Ping_C struct {
	Rcm *impl.RC_Runner_m_j
}

func (p *Ping_C) Ping(val string) (string, error) {
	if p.Rcm == nil {
		return "", util.Err("RCM not inital,call RunC first")
	}
	var vs Val
	_, err := p.Rcm.VExec("ping", map[string]interface{}{
		"tv": val,
	}, &vs)
	return vs.V, err
}

func Ping_S(rc *impl.RCM_Cmd) (interface{}, error) {
	var tv string
	err := rc.ValidF(`
		tv,R|S,L:0,
		`, &tv)
	if err == nil {
		if ping_l {
			log.D("receive ping(%v) from %v", tv, rc.RemoteAddr().String())
		}
		return Val{
			V: tv,
		}, nil
	} else {
		return nil, err
	}
}

func RunC(addr string, t int, delay time.Duration) error {
	netw.ShowLog = true
	bp := pool.NewBytePool(8, 1024) //memory pool.
	pc := Ping_C{}
	pc.Rcm = impl.NewRC_Runner_m_j(addr, bp)
	pc.Rcm.Start()
	defer pc.Rcm.Stop()
	for t != 0 {
		if t > 0 {
			t--
		}
		tv := fmt.Sprintf("%d", rand.Intn(10000))
		bv, err := pc.Ping(tv)
		if err != nil {
			log.E("ping(%v) err:%v", tv, err.Error())
		} else if tv == bv {
			if ping_l {
				log.D("ping(%v),ret(%v)", tv, bv)
			}
		} else {
			log.E("ping(%v),ret(%v)", tv, bv)
		}
		time.Sleep(delay * time.Millisecond)
	}
	log.D("all done...")
	return nil
}
func RunS(port string) error {
	netw.ShowLog = true
	p := pool.NewBytePool(8, 1024) //memory pool.
	l, cc, cms := impl.NewChanExecListener_m_j(p, port, netw.NewCWH(true))
	cms.AddHFunc("ping", Ping_S)
	cc.Run(runtime.NumCPU() - 1) //start the chan distribution, if not start, sub handler will not receive message
	err := l.Run()               //run the listen server
	if err != nil {
		return err
	}
	log.D("start listener on %v", port)
	go func() {
		signal.Notify(tsig, os.Interrupt, os.Kill)
		<-tsig
		l.Close()
	}()
	l.Wait()
	log.D("Srv done...")
	return nil
}
func Usage() {
	fmt.Println(`Usage:
	tping -m S [-l log file] [-p listen port, default :9100] [-a show all long, default Y]
	tping [-m C] [-l log file] [-d ping delay,default 1s] [-t ping times, default -1] [-a show all long, default Y]
		`)
}
func Run(args []string) {
	var p string = ":9100"
	var t int = -1
	var d time.Duration = 1000
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
			d = time.Duration(delay)
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
	if m == "S" {
		err = RunS(p)
	} else {
		if len(h) < 1 {
			Usage()
			return
		}
		err = RunC(h, t, d)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	Run(os.Args)
}
