package config

import (
	"encoding/json"
	"flag"
	"github.com/golang/glog"
	"io/ioutil"
	"path/filepath"
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

func init() {
	var f string
	flag.StringVar(&f, "f", "config.json", "Specify a path to the config file")
	flag.Parse()

	if cfg, err := ioutil.ReadFile(f); err != nil {
		glog.Fatalln(err)
	} else {
		if err := json.Unmarshal(cfg, &Cfg); err != nil {
			glog.Fatal(err)
		}
		glog.Infof("%v", string(cfg))
	}

	dir := filepath.Dir(f) + "/"
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
