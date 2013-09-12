package db

import (
	"log"
	"time"

	"github.com/fzzy/radix/redis"
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
				log.Println("PING error: ", r.Err)
				if err := c.Close(); err != nil {
					log.Println("Close error: ", err)
				}
				continue
			}
			return c, nil

		case <-time.After(time.Duration(config.Cfg.DB.Timeout) * time.Second):
			log.Println("DB Client timed out")
		}
	}
	return newClient()
}

// Release pushes back client to pool (if number of available clients in the pool is  < MaxClients),
// otherwise closes client
func Release(client *redis.Client) {
	if client != nil {
		if len(clients) < config.Cfg.DB.MaxClients {
			log.Println("Releasing Db Client")
			clients <- client
			log.Println("Number of idle Db Clients: ", len(clients))
		} else {
			log.Println("Closing Db Client")
			if err := client.Close(); err != nil {
				log.Println("Close error: ", err)
			}
		}
	}
}
