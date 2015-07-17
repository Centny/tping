package main

import (
	"errors"
	"fmt"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"os"
	"runtime"
	"testing"
	"time"
)

var ping_s_c int = 0

func ping_s(rc *impl.RCM_Cmd) (interface{}, error) {
	ping_s_c++
	switch ping_s_c % 2 {
	case 0:
		return nil, errors.New("--->")
	default:
		return Val{
			V: "xxxx",
		}, nil
	}
}
func run_ts(port string) error {
	netw.ShowLog = true
	p := pool.NewBytePool(8, 1024) //memory pool.
	l, cc, cms := impl.NewChanExecListener_m_j(p, port, netw.NewCWH(true))
	cms.AddHFunc("ping", ping_s)
	cc.Run(runtime.NumCPU() - 1) //start the chan distribution, if not start, sub handler will not receive message
	err := l.Run()               //run the listen server
	if err != nil {
		return err
	}
	defer l.Close()
	l.Wait()
	return nil
}
func TestPing(t *testing.T) {
	exit = func(code int) {
		fmt.Println("calling exit by code:", code)
	}
	os.Remove("e.xml")
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	go Run([]string{"tping", "-m", "S", "-p", ":9910", "-l", "t.log", "-a", "Y"})
	time.Sleep(100 * time.Millisecond)
	go Run([]string{"tping", "-m", "S", "-p", ":9910"})
	go run_ts(":9920")
	go Run([]string{"tping", "-m", "C", "-h", "127.0.0.1:9910", "-t", "-1", "-d", "200"})
	go Run([]string{"tping", "-m", "C", "-h", "127.0.0.1:9910", "-t", "1", "-d", "10"})
	go Run([]string{"tping", "-m", "C", "-h", "127.0.0.1:9910", "-t", "3", "-d", "200"})
	go Run([]string{"tping", "-m", "C", "-h", "127.0.0.1:9920", "-t", "10", "-d", "200"})
	go Run([]string{"tping", "-m", "C", "-h", "127.0.0.1:9920", "-t", "1", "-d", "100"})
	go Run([]string{"tping", "-m", "C", "-h", "127.0.0.1:9920", "-t", "10", "-d", "0"})
	go Run([]string{"tping", "-m", "W", "-h", "http://www.bing.com", "-t", "1", "-d", "2000"})
	go Run([]string{"tping", "-m", "W", "-h", "http://www.bing.com"})
	go Run([]string{"tping", "-m", "W", "-h", "http://www.bing.com", "-d", "1"})
	go Run([]string{"tping", "-m", "W", "-h", "http://127.0.0.1:234"})
	go Run([]string{"tping", "-m", "R", "-r", "echo abc"})
	go Run([]string{"tping", "-m", "R", "-r", "xddd abc"})
	go Run([]string{"tping", "127.0.0.1:9910x"})
	go func() {
		lc, rc, _ := impl.ExecDail_m_j(pool.NewBytePool(8, 1024), "127.0.0.1:9910", netw.NewCWH(true))
		rc.Start()
		var mv map[string]interface{}
		rc.Exec("ping", map[string]interface{}{}, &mv)
		lc.Close()
	}()
	go func() {
		Run([]string{"tping", "-m", "C", "-h", "127.0.0.1:9910", "-t", "xx", "-d", "200"})
		Run([]string{"tping", "-m", "C", "-h", "127.0.0.1:9910", "-l", "t.log", "-t", "-1", "-d", "xxx"})
		Run([]string{"tping", "-m", "C", "-h", "127.0.0.1:9910", "-l", "/ss/t.log", "-t", "-1", "-d", "200"})
		Run([]string{"tping"})
	}()
	time.Sleep(3 * time.Second)
	tsig <- os.Kill
	time.Sleep(time.Second)
	v := Ping_C{}
	v.Ping("ss")
	os.Args = []string{"tping", "-m", "SS"}
	main()
	//
	Run([]string{"tping", "-m", "J", "-j", "t.json", "-e", "e.xml"})
	Run([]string{"tping", "-m", "J", "-j", "t.json", "-help", "-e", "e.xml"})
	Run([]string{"tping", "-m", "J", "-j", "t1.json", "-e", "/e.xml"})
	Run([]string{"tping", "-m", "J", "-j", "t2.json", "-e", "/e.xml"})
	Run([]string{"tping", "-m", "J", "-j", "xx.json", "-e", "/e.xml"})
	Run([]string{"tping", "-a", "Y", "http://www.bing.com"})
	Run([]string{"tping", "http://www.bing.com", "-a", "Y"})
	Run([]string{"tping", "-m", "C"})
	Run([]string{"tping", "-m", "R"})
	Run([]string{"tping", "-m", "W"})
	Run([]string{"tping", "-m", "J"})
	Run([]string{"tping", "-help"})
	//
}
