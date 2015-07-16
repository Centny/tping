package main

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"time"
)

var tsig = make(chan os.Signal, 1)

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
		time.Sleep(delay)
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
