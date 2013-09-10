package main

import (
	"log"
	"net/http"

	"github.com/kuba--/yag/pkg/config"
	"github.com/kuba--/yag/webserver/render"
)

func main() {
	http.HandleFunc("/render", render.Handler)

	log.Println("ListenAndServe", config.Cfg.Webserver.Addr)
	log.Fatal(http.ListenAndServe(config.Cfg.Webserver.Addr, nil))
}
