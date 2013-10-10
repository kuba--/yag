package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
)

var Cfg struct {
	DB struct {
		Addr       string
		Timeout    int
		MaxClients int
	}

	Metrics struct {
		GetScript string
		AddScript string
		TtlScript string
		TTL       int
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
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var f string
	flag.StringVar(&f, "f", "config.json", "Specify a path to the config file")
	flag.Parse()

	if cfg, err := ioutil.ReadFile(f); err != nil {
		log.Println(err)
	} else {
		if err := json.Unmarshal(cfg, &Cfg); err != nil {
			log.Fatal(err)
		}
		log.Printf("%v", string(cfg))
	}

	dir := filepath.Dir(f) + "/"
	if script, err := ioutil.ReadFile(dir + Cfg.Metrics.AddScript); err != nil {
		log.Println(err)
		Cfg.Metrics.AddScript = ""
	} else {
		Cfg.Metrics.AddScript = string(script)
	}

	if script, err := ioutil.ReadFile(dir + Cfg.Metrics.GetScript); err != nil {
		log.Println(err)
		Cfg.Metrics.GetScript = ""
	} else {
		Cfg.Metrics.GetScript = string(script)
	}

	if script, err := ioutil.ReadFile(dir + Cfg.Metrics.TtlScript); err != nil {
		log.Println(err)
		Cfg.Metrics.TtlScript = ""
	} else {
		Cfg.Metrics.TtlScript = string(script)
	}
}
