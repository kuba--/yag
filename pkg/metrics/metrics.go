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

type Pt [2]*float64

type Metrics struct {
	Key        string
	Target     string
	Datapoints []Pt
}

func newMetrics(key, target string, datapoints []Pt) (m *Metrics) {
	m = new(Metrics)
	m.Key, m.Target, m.Datapoints = key, target, datapoints
	return
}

/*
 * Get queries for metrics which matches to the key pattern (e.g.: status.*)
 *
 * [
 *  {"target": "status.200", "datapoints": [[1720.0, 1370846820], ...], },
 *  {"target": "status.204", "datapoints": [[1.0, 1370846820], ..., ]}
 * ]
 */
func Get(key string, from int64, to int64) (ms []*Metrics) {
	var js []byte
	var data []map[string]interface{}

	if client, err := db.Client(); err != nil {
		log.Println(err)
	} else {
		defer db.Release(client)

		if js, err = client.Cmd("EVALSHA", getSha, 1, key, from, to).Bytes(); err != nil {
			log.Println(err)
			if js, err = client.Cmd("EVAL", config.Cfg.Metrics.GetScript, 1, key, from, to).Bytes(); err != nil {
				log.Println(err)
			}
		}

		if err = json.Unmarshal(js, &data); err != nil {
			log.Println(err)
		}
	}

	for _, d := range data {
		m := new(Metrics)
		m.Key = key
		if target, ok := d["target"].(string); ok {
			m.Target = target
		}

		if datapoints, ok := d["datapoints"].([]interface{}); ok {
			if config.Cfg.Metrics.ConsolidationStep < 1 || len(config.Cfg.Metrics.ConsolidationFunc) < 1 {
				for _, dp := range datapoints {
					var pt Pt
					if err := json.Unmarshal([]byte(dp.(string)), &pt); err != nil {
						log.Println(err)
						continue
					}
					m.Datapoints = append(m.Datapoints, pt)
				}
			} else {
				m.Datapoints = consolidateBy(datapoints, from, to, config.Cfg.Metrics.ConsolidationStep, config.Cfg.Metrics.ConsolidationFunc)
			}
		}
		ms = append(ms, m)
	}
	return
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

/*
 * Valid consolidation function names are 'sum', 'avg', 'min', and 'max'
 */
func consolidateBy(data []interface{}, from, to int64, step int, fn string) (datapoints []Pt) {
	for i := 0; from <= to; from += int64(step) {
		var isset bool = false
		var n int = 0

		var sum, max, min, ts *float64 = new(float64), nil, nil, new(float64)
		*sum, *ts = 0.0, float64(from)

		for ; i < len(data); i++ {
			var pt Pt
			if err := json.Unmarshal([]byte(data[i].(string)), &pt); err != nil {
				log.Println(err)
				continue
			}
			if pt[1] != nil && int64(*pt[1]) >= from && int64(*pt[1]) < from+int64(step) {
				*sum = *sum + *pt[0]

				if max == nil {
					max = new(float64)
					*max = *pt[0]
				} else {
					if *max < *pt[0] {
						*max = *pt[0]
					}
				}

				if min == nil {
					min = new(float64)
					*min = *pt[0]
				} else {
					if *min > *pt[0] {
						*min = *pt[0]
					}
				}

				n++
				isset = true
			} else {
				break
			}
		}

		var value *float64 = nil
		if isset {
			value = new(float64)
			switch fn {
			case "sum":
				*value = *sum
			case "avg":
				*value = *sum / float64(n)
			case "max":
				*value = *max
			case "min":
				*value = *min
			}
		}
		datapoints = append(datapoints, Pt{value, ts})
	}
	return
}
