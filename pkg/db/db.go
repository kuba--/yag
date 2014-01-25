package db

import (
	"time"

	"github.com/fzzy/radix/redis"
	"github.com/kuba--/glog"
	"github.com/kuba--/yag/pkg/config"
)

// Pool of clients
var clients chan *redis.Client

func init() {
	clients = make(chan *redis.Client, config.Cfg.DB.MaxClients)
}

func newClient() (*redis.Client, error) {
	return redis.DialTimeout("tcp", config.Cfg.DB.Addr, time.Duration(config.Cfg.DB.Timeout)*time.Second)
}

// Client tries to get first available client from pool,
// otherwise creates new instance of client
func Client() (*redis.Client, error) {
	for i := 0; i < len(clients); i++ {
		select {
		case c := <-clients:
			r := c.Cmd("PING")
			if r.Err != nil {
				glog.Warningln("PING error: ", r.Err)
				if err := c.Close(); err != nil {
					glog.Warningln("Close error: ", err)
				}
				continue
			}
			return c, nil

		case <-time.After(time.Duration(config.Cfg.DB.Timeout) * time.Second):
			glog.Warningln("DB Client timed out")
		}
	}
	return newClient()
}

// Release pushes back client to pool (if number of available clients in the pool is  < MaxClients),
// otherwise closes client
func Release(client *redis.Client) {
	if client != nil {
		if len(clients) < config.Cfg.DB.MaxClients {
			glog.Infoln("Releasing Db Client")
			clients <- client
			glog.Infoln("Number of idle Db Clients: ", len(clients))
		} else {
			glog.Infoln("Closing Db Client")
			if err := client.Close(); err != nil {
				glog.Warningln("Close error: ", err)
			}
		}
	}
}
