package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
)

var Cfg struct {
	DB struct {
		Addr       string
		Timeout    int
		MaxClients int
	}

	Metrics struct {
		TTL int

		GetScript string
		AddScript string
		TtlScript string
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
	flag.StringVar(&f, "f", ".", "Specify a file hierarchy for config and script files")
	flag.Parse()

	if cfg, err := ioutil.ReadFile(f + "/config.json"); err != nil {
		log.Panic(err)
	} else {
		json.Unmarshal(cfg, &Cfg)
		log.Printf("%v", string(cfg))
	}

	if script, err := ioutil.ReadFile(f + "/add.lua"); err != nil {
		log.Panic(err)
	} else {
		Cfg.Metrics.AddScript = string(script)
	}

	if script, err := ioutil.ReadFile(f + "/get.lua"); err != nil {
		log.Panic(err)
	} else {
		Cfg.Metrics.GetScript = string(script)
	}

	if script, err := ioutil.ReadFile(f + "/ttl.lua"); err != nil {
		log.Panic(err)
	} else {
		Cfg.Metrics.TtlScript = string(script)
	}
}
