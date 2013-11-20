package main

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/kuba--/yag/pkg/config"
	"github.com/kuba--/yag/webserver/render"
)

func main() {
	http.HandleFunc("/render", render.Handler)

	glog.Infoln("ListenAndServe", config.Cfg.Webserver.Addr)
	glog.Fatal(http.ListenAndServe(config.Cfg.Webserver.Addr, nil))
}
