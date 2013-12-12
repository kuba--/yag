package main

import (
	"bufio"
	"net"
	"os"
	"os/signal"
	"runtime/pprof"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/kuba--/yag/pkg/config"
	"github.com/kuba--/yag/pkg/metrics"
)

func main() {
	if config.Pprof.Cpu != "" {
		// Start Cpu Profiler
		f, err := os.Create(config.Pprof.Cpu)
		if err != nil {
			glog.Fatal(err)
		}
		defer f.Close()

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	ln, err := net.Listen("tcp", config.Cfg.Listener.Addr)
	if err != nil {
		glog.Fatalln(err)
	}
	glog.Infoln("ListenAndServe", config.Cfg.Listener.Addr)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				glog.Errorln(err)
				continue
			}
			go handle(conn)
		}
	}()

	// Handle SIGINT and SIGTERM.
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	glog.Errorf("Syscall: %v - Exit\n", <-sig)

	if config.Pprof.Mem != "" {
		// Start Mem Profiler
		f, err := os.Create(config.Pprof.Mem)
		if err != nil {
			glog.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
	}
}

// Handle TCP Connections
func handle(conn net.Conn) {
	glog.Infof("Connection: %s -> %s\n", conn.RemoteAddr(), conn.LocalAddr())
	defer func() {
		glog.Infof("Closing connection: %s\n", conn.RemoteAddr())
		conn.Close()
		glog.Flush()
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if m := strings.Split(scanner.Text(), " "); len(m) > 2 {
			if ts, err := strconv.ParseInt(m[2], 10, 0); err != nil {
				glog.Warningln(err)
				continue
			} else {
				metrics.Add(m[0], m[1], ts)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		glog.Errorln(err)
	}
}
