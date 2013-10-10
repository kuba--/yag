package metrics

import (
	"encoding/json"
	"log"
	"time"

	"github.com/kuba--/yag/pkg/config"
	"github.com/kuba--/yag/pkg/db"
)

var addSha, getSha, ttlSha string

func init() {
	if client, err := db.Client(); err != nil {
		log.Println(err)
	} else {
		defer db.Release(client)
		{
			if len(config.Cfg.Metrics.AddScript) > 0 {
				if addSha, err = client.Cmd("SCRIPT", "LOAD", config.Cfg.Metrics.AddScript).Str(); err != nil {
					log.Println(err)
				} else {
					log.Println("ADD SHA", addSha)
				}
			}

			if len(config.Cfg.Metrics.GetScript) > 0 {
				if getSha, err = client.Cmd("SCRIPT", "LOAD", config.Cfg.Metrics.GetScript).Str(); err != nil {
					log.Println(err)
				} else {
					log.Println("GET SHA", getSha)
				}
			}

			if len(config.Cfg.Metrics.TtlScript) > 0 {
				if ttlSha, err = client.Cmd("SCRIPT", "LOAD", config.Cfg.Metrics.TtlScript).Str(); err != nil {
					log.Println(err)
				} else {
					log.Println("TTL SHA", ttlSha)
				}
			}
		}
	}
}

/*
 * [
 *  {"target": "status.200", "datapoints": [[1720.0, 1370846820], ...], },
 *  {"target": "status.204", "datapoints": [[1.0, 1370846820], ..., ]}
 * ]
 */
type Metrics struct {
	Key        string
	Target     string
	Datapoints [][2]float64
}

func (m1 *Metrics) isEqual(m2 *Metrics) bool {
	if m2 == nil {
		return false
	}

	d1, d2 := m1.Datapoints, m2.Datapoints
	ld1, ld2 := len(d1), len(d2)
	if ld1 != ld2 {
		return false
	}

	for i := 0; i < ld1; i++ {
		v1, t1 := d1[i][0], d1[i][1]
		v2, t2 := d2[i][0], d2[i][1]
		if t1 != t2 || v1 != v2 {
			return false
		}
	}

	return true
}

func newMetrics(key string, m []map[string]interface{}) []*Metrics {
	ms := make([]*Metrics, 0)

	for _, mi := range m {
		mm := new(Metrics)
		mm.Key = key
		if target, ok := mi["target"].(string); ok {
			mm.Target = target
		}
		mm.Datapoints = make([][2]float64, 0)
		if datapoints, ok := mi["datapoints"].([]interface{}); ok {
			for _, dp := range datapoints {
				dpi := dp.(string)

				var pt [2]float64
				err := json.Unmarshal([]byte(dpi), &pt)
				if err != nil {
					log.Println(err)
				}
				mm.Datapoints = append(mm.Datapoints, pt)
			}
		}
		ms = append(ms, mm)
	}
	return ms
}

func Add(key string, value string, timestamp int64) {
	if client, err := db.Client(); err != nil {
		log.Println(err)
	} else {
		defer db.Release(client)

		if r := client.Cmd("EVALSHA", addSha, 1, key, value, timestamp); r.Err != nil {
			log.Println(r.Err)

			if r = client.Cmd("EVAL", config.Cfg.Metrics.AddScript, 1, key, value, timestamp); r.Err != nil {
				log.Println(r.Err)
			}
		} else {
			log.Printf("[OK: %v]\t(%s %s)|%d", r, key, value, timestamp)
		}
	}
}

// Get queries for metrics which matches to the key pattern
func Get(key string, from int64, to int64, limit int) []*Metrics {
	var m []map[string]interface{}

	if client, err := db.Client(); err != nil {
		log.Println(err)
	} else {
		defer db.Release(client)

		if data, err := client.Cmd("EVALSHA", getSha, 1, key, from, to, limit).Str(); err != nil {
			log.Println(err)

			if data, err = client.Cmd("EVAL", config.Cfg.Metrics.GetScript, 1, key, from, to, limit).Str(); err != nil {
				log.Println(err)
			}
		} else {
			json.Unmarshal([]byte(data), &m)
		}
	}

	return newMetrics(key, m)
}

func Ttl(from int64, to int64) {
	if client, err := db.Client(); err != nil {
		log.Println(err)
	} else {
		t0 := time.Now()
		defer db.Release(client)

		if r := client.Cmd("EVALSHA", ttlSha, 1, "*", from, to); r.Err != nil {
			log.Println(r.Err)

			if r = client.Cmd("EVAL", config.Cfg.Metrics.TtlScript, 1, "*", from, to); r.Err != nil {
				log.Println(r.Err)
			}
		} else {
			t1 := time.Now()
			log.Printf("ZREMRANGEBYSCORE(%d, %d): %v in %v", from, to, r, t1.Sub(t0))
		}
	}
}
