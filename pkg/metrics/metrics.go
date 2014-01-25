package metrics

import (
	"encoding/json"
	"math"
	"time"

	"github.com/kuba--/glog"
	"github.com/kuba--/yag/pkg/config"
	"github.com/kuba--/yag/pkg/db"
)

var (
	addSha string
	getSha string
	ttlSha string
)

func init() {
	if client, err := db.Client(); err != nil {
		glog.Errorln(err)
	} else {
		defer db.Release(client)
		{
			if len(config.Cfg.Metrics.AddScript) > 0 {
				if addSha, err = client.Cmd("SCRIPT", "LOAD", config.Cfg.Metrics.AddScript).Str(); err != nil {
					glog.Errorln(err)
				} else {
					glog.Infoln("ADD SHA", addSha)
				}
			}

			if len(config.Cfg.Metrics.GetScript) > 0 {
				if getSha, err = client.Cmd("SCRIPT", "LOAD", config.Cfg.Metrics.GetScript).Str(); err != nil {
					glog.Errorln(err)
				} else {
					glog.Infoln("GET SHA", getSha)
				}
			}

			if len(config.Cfg.Metrics.TtlScript) > 0 {
				if ttlSha, err = client.Cmd("SCRIPT", "LOAD", config.Cfg.Metrics.TtlScript).Str(); err != nil {
					glog.Errorln(err)
				} else {
					glog.Infoln("TTL SHA", ttlSha)
				}
			}
		}
	}
}

type Pt [2]*float64

func makePt(dp interface{}) (pt Pt) {
	if err := json.Unmarshal([]byte(dp.(string)), &pt); err != nil {
		glog.Warningln(err)
	}
	return
}

type Metrics struct {
	Key        string
	Target     string
	Datapoints []Pt
}

/*
 * Get queries for metrics which matches to the key pattern (e.g.: status.*)
 *
 * [
 *  {"target": "status.200", "datapoints": [[1720.0, 1370846820], ...], },
 *  {"target": "status.204", "datapoints": [[1.0, 1370846820], ..., ]}
 * ]
 */
func Get(key string, from int64, to int64, maxDataPoints int) (ms []*Metrics) {
	var js []byte
	var data []map[string]interface{}

	if client, err := db.Client(); err != nil {
		glog.Errorln(err)
	} else {
		defer db.Release(client)

		if js, err = client.Cmd("EVALSHA", getSha, 1, key, from, to).Bytes(); err != nil {
			glog.Warningln(err)
			if js, err = client.Cmd("EVAL", config.Cfg.Metrics.GetScript, 1, key, from, to).Bytes(); err != nil {
				glog.Errorln(err)
			}
		}

		if err = json.Unmarshal(js, &data); err != nil {
			glog.Errorln(err)
		}
	}

	for _, d := range data {
		m := new(Metrics)
		m.Key = key
		if target, ok := d["target"].(string); ok {
			m.Target = target
		}

		datapoints, ok := d["datapoints"].([]interface{})
		if !ok {
			datapoints = make([]interface{}, 0)
		}

		if config.Cfg.Metrics.ConsolidationStep < 1 || len(config.Cfg.Metrics.ConsolidationFunc) < 1 {
			for _, dp := range datapoints {
				if pt := makePt(dp); pt[1] != nil {
					m.Datapoints = append(m.Datapoints, pt)
				}
			}
		} else {
			step := consolidationStep(from, to, config.Cfg.Metrics.ConsolidationStep, maxDataPoints)
			m.Datapoints = consolidateBy(datapoints, from, to, step, config.Cfg.Metrics.ConsolidationFunc)
		}
		ms = append(ms, m)
	}
	return
}

func Add(key string, value string, timestamp int64) {
	if client, err := db.Client(); err != nil {
		glog.Errorln(err)
	} else {
		defer db.Release(client)

		if r := client.Cmd("EVALSHA", addSha, 1, key, value, timestamp); r.Err != nil {
			glog.Warningln(r.Err)

			if r = client.Cmd("EVAL", config.Cfg.Metrics.AddScript, 1, key, value, timestamp); r.Err != nil {
				glog.Errorln(r.Err)
			}
		} else {
			glog.Infof("[OK: %v]\t(%s %s)|%d", r, key, value, timestamp)
		}
	}
}

func Ttl(from int64, to int64) {
	if client, err := db.Client(); err != nil {
		glog.Errorln(err)
	} else {
		t0 := time.Now()
		defer db.Release(client)

		if r := client.Cmd("EVALSHA", ttlSha, 1, "*", from, to); r.Err != nil {
			glog.Warningln(r.Err)

			if r = client.Cmd("EVAL", config.Cfg.Metrics.TtlScript, 1, "*", from, to); r.Err != nil {
				glog.Errorln(r.Err)
			}
		} else {
			t1 := time.Now()
			glog.Infof("ZREMRANGEBYSCORE(%d, %d): %v in %v", from, to, r, t1.Sub(t0))
		}
	}
}

func consolidationStep(from, to int64, step, maxDataPoints int) int {
	numberOfDataPoints := int(math.Ceil(float64(to-from) / float64(step)))

	if maxDataPoints > 0 && numberOfDataPoints > maxDataPoints {
		step *= int(numberOfDataPoints / maxDataPoints)
	}
	return step
}

/*
 * Valid consolidation function names are 'sum', 'avg', 'min', and 'max'
 */
func consolidateBy(data []interface{}, from, to int64, step int, fn string) (datapoints []Pt) {
	for i := 0; from <= to; from += int64(step) {
		var (
			isset             bool     = false
			n                 int      = 0
			sum, max, min, ts *float64 = new(float64), nil, nil, new(float64)
		)
		*sum, *ts = 0.0, float64(from)

		for ; i < len(data); i++ {
			pt := makePt(data[i])
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
