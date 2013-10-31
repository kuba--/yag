package metrics

import (
	"encoding/json"
	"log"
	"math"
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

type Metrics struct {
	Key        string
	Target     string
	Datapoints [][2]float64
}

func newMetrics(key, target string, datapoints [][2]float64) (m *Metrics) {
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
			switch config.Cfg.Metrics.ConsolidationFunc {
			case "avg":
				m.Datapoints = consolidateByAvg(datapoints, from, to, config.Cfg.Metrics.ConsolidationStep)
			case "sum":
				m.Datapoints = consolidateBySum(datapoints, from, to, config.Cfg.Metrics.ConsolidationStep)
			case "max":
				m.Datapoints = consolidateByMax(datapoints, from, to, config.Cfg.Metrics.ConsolidationStep)
			case "min":
				m.Datapoints = consolidateByMin(datapoints, from, to, config.Cfg.Metrics.ConsolidationStep)

			default:
				for _, dp := range datapoints {
					var pt [2]float64
					if err := json.Unmarshal([]byte(dp.(string)), &pt); err != nil {
						log.Println(err)
						continue
					}
					m.Datapoints = append(m.Datapoints, pt)
				}
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

func consolidateBySum(data []interface{}, from, to int64, step int) (datapoints [][2]float64) {
	for i := 0; from <= to; from += int64(step) {
		var sum float64 = 0.0
		for ; i < len(data); i++ {
			var pt [2]float64
			if err := json.Unmarshal([]byte(data[i].(string)), &pt); err != nil {
				log.Println(err)
				continue
			}
			if int64(pt[1]) >= from && int64(pt[1]) < from+int64(step) {
				sum += pt[0]
			} else {
				break
			}
		}
		datapoints = append(datapoints, [2]float64{sum, float64(from)})
	}
	return
}

func consolidateByAvg(data []interface{}, from, to int64, step int) (datapoints [][2]float64) {
	for i := 0; from <= to; from += int64(step) {
		var sum float64 = 0.0
		var n int = 0

		for ; i < len(data); i++ {
			var pt [2]float64
			if err := json.Unmarshal([]byte(data[i].(string)), &pt); err != nil {
				log.Println(err)
				continue
			}
			if int64(pt[1]) >= from && int64(pt[1]) < from+int64(step) {
				sum += pt[0]
				n++
			} else {
				break
			}
		}

		if n > 0 {
			sum /= float64(n)
		}
		datapoints = append(datapoints, [2]float64{sum, float64(from)})
	}
	return
}

func consolidateByMax(data []interface{}, from, to int64, step int) (datapoints [][2]float64) {
	for i := 0; from <= to; from += int64(step) {
		var max float64 = 0.0
		for ; i < len(data); i++ {
			var pt [2]float64
			if err := json.Unmarshal([]byte(data[i].(string)), &pt); err != nil {
				log.Println(err)
				continue
			}
			if int64(pt[1]) >= from && int64(pt[1]) < from+int64(step) {
				if pt[0] > max {
					max = pt[0]
				}
			} else {
				break
			}
		}
		datapoints = append(datapoints, [2]float64{max, float64(from)})
	}
	return
}

func consolidateByMin(data []interface{}, from, to int64, step int) (datapoints [][2]float64) {
	for i := 0; from <= to; from += int64(step) {
		var min float64 = math.MaxFloat64
		var isset bool = false
		for ; i < len(data); i++ {
			var pt [2]float64
			if err := json.Unmarshal([]byte(data[i].(string)), &pt); err != nil {
				log.Println(err)
				continue
			}
			if int64(pt[1]) >= from && int64(pt[1]) < from+int64(step) {
				if pt[0] < min {
					min = pt[0]
				}
				isset = true
			} else {
				break
			}
		}
		if !isset {
			min = 0
		}

		datapoints = append(datapoints, [2]float64{min, float64(from)})
	}
	return
}
