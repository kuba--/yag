package main

import (
	"bufio"
	"net"
	"strconv"
	"strings"

	"github.com/kuba--/glog"
	"github.com/kuba--/yag/pkg/config"
	"github.com/kuba--/yag/pkg/metrics"
)

func main() {
	ln, err := net.Listen("tcp", config.Cfg.Listener.Addr)
	if err != nil {
		glog.Fatalln(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			glog.Warningln(err)
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	glog.Infof("Connection: %s -> %s\n", conn.RemoteAddr(), conn.LocalAddr())
	defer func() {
		glog.Infof("Closing connection: %s\n", conn.RemoteAddr())
		conn.Close()
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
	glog.Flush()
}
