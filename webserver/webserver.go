package main

import (
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"

	"github.com/golang/glog"
	"github.com/kuba--/yag/pkg/config"
	"github.com/kuba--/yag/webserver/render"
)

func main() {
	if config.Pprof.Cpu != "" {
		// Start Profiler
		f, err := os.Create(config.Pprof.Cpu)
		if err != nil {
			glog.Fatal(err)
		}
		defer f.Close()

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	http.HandleFunc("/render", render.Handler)

	glog.Infoln("ListenAndServe", config.Cfg.Webserver.Addr)
	go func() {
		glog.Fatal(http.ListenAndServe(config.Cfg.Webserver.Addr, nil))
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
