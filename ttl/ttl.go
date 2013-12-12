package main

import (
	"time"
	"github.com/kuba--/yag/pkg/config"
	"github.com/kuba--/yag/pkg/metrics"
)

func main() {
	for _ = range time.Tick(time.Second * time.Duration(config.Cfg.TTL.Tick)) {
		metrics.Ttl(0, time.Now().Add(-1*time.Second*time.Duration(config.Cfg.Metrics.TTL)).Unix())
	}
}
