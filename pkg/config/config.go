package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"path/filepath"

	"github.com/golang/glog"
)

var Cfg struct {
	DB struct {
		Addr       string
		Timeout    int
		MaxClients int
	}

	Metrics struct {
		GetScript         string
		AddScript         string
		TtlScript         string
		TTL               int
		ConsolidationStep int
		ConsolidationFunc string
	}

	Listener struct {
		Addr string
	}

	Webserver struct {
		Addr    string
		Timeout int
	}

	TTL struct {
		Tick int
	}
}

var Pprof struct {
	Cpu string
	Mem string
}

func init() {
	// Max log file size - rotate log file if exceed 10MB
	glog.MaxSize = 10 * 1024 * 1024

	var f = flag.String("f", "config.json", "read configuration from file")

	flag.StringVar(&Pprof.Cpu, "cpuprofile", "", "write cpu profile to file")
	flag.StringVar(&Pprof.Mem, "memprofile", "", "write memory profile to this file")

	flag.Parse()

	if cfg, err := ioutil.ReadFile(*f); err != nil {
		glog.Errorln(err)
		return
	} else {
		if err := json.Unmarshal(cfg, &Cfg); err != nil {
			glog.Fatal(err)
		}
		glog.Infof("%v", string(cfg))
	}

	dir := filepath.Dir(*f) + "/"
	if script, err := ioutil.ReadFile(dir + Cfg.Metrics.AddScript); err != nil {
		glog.Errorln(err)
		Cfg.Metrics.AddScript = ""
	} else {
		Cfg.Metrics.AddScript = string(script)
	}

	if script, err := ioutil.ReadFile(dir + Cfg.Metrics.GetScript); err != nil {
		glog.Errorln(err)
		Cfg.Metrics.GetScript = ""
	} else {
		Cfg.Metrics.GetScript = string(script)
	}

	if script, err := ioutil.ReadFile(dir + Cfg.Metrics.TtlScript); err != nil {
		glog.Errorln(err)
		Cfg.Metrics.TtlScript = ""
	} else {
		Cfg.Metrics.TtlScript = string(script)
	}
}
