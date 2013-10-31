package metrics

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Api struct{}

func (api *Api) Value(name string, from int64, to int64) interface{} {
	log.Printf("Api.Value(%s, %d, %d)", name, from, to)

	// check constant value first
	if f, err := strconv.ParseFloat(name, 10); err == nil {
		return f
	}
	return Get(name, from, to)
}

func (api *Api) Call(name string, argv interface{}) interface{} {
	if arr, ok := argv.([]interface{}); ok {
		if n := len(arr); n > 0 {
			if m0, ok := arr[0].([]*Metrics); ok {
				for i := 1; i < n; i++ {
					if m, ok := arr[i].([]*Metrics); ok {
						for _, mi := range m {
							m0 = append(m0, mi)
						}
					}
				}
				switch strings.ToLower(name) {
				case "sum", "sumseries":
					return api.sumSeries(m0)
				case "div", "divseries", "divideseries":
					return api.divSeries(m0)
				case "diff", "diffseries":
					return api.diffSeries(m0)
				case "_":
					return m0
				default:
					log.Println("[ ! ]\tFunction not supported: ", name)
				}
			}
		}
	}
	return nil
}

func (api *Api) sumSeries(m []*Metrics) []*Metrics {

	sm := make([]*Metrics, 0)
	{
		if n := len(m); n > 0 {
			key := m[0].Key
			dp0 := m[0].Datapoints

			for i := 1; i < n; i++ {
				datapoints := make([][2]float64, 0)
				dp1 := m[i].Datapoints
				i0, n0 := 0, len(dp0)
				i1, n1 := 0, len(dp1)

				for i0 < n0 && i1 < n1 {
					v0, ts0 := dp0[i0][0], dp0[i0][1]
					v1, ts1 := dp1[i1][0], dp1[i1][1]
					its0, its1 := int64(ts0), int64(ts1)

					switch {
					case its0 == its1:
						datapoints = append(datapoints, [2]float64{v0 + v1, ts0})
						i0++
						i1++
						continue

					case its0 < its1:
						datapoints = append(datapoints, [2]float64{v0, ts0})
						i0++
						continue

					case its0 > its1:
						datapoints = append(datapoints, [2]float64{v1, ts1})
						i1++
					}
				}

				for ; i0 < n0; i0++ {
					datapoints = append(datapoints, [2]float64{dp0[i0][0], dp0[i0][1]})
				}

				for ; i1 < n1; i1++ {
					datapoints = append(datapoints, [2]float64{dp1[i1][0], dp1[i1][1]})
				}

				dp0 = datapoints
			}
			sm = append(sm, newMetrics(key, fmt.Sprintf("sum(%s)", key), dp0))
		}
	}
	return sm
}

func (api *Api) divSeries(m []*Metrics) []*Metrics {
	sm := make([]*Metrics, 0)
	{
		if n := len(m); n > 0 {
			key := m[0].Key
			dp0 := m[0].Datapoints

			for i := 1; i < n; i++ {
				datapoints := make([][2]float64, 0)
				dp1 := m[i].Datapoints
				i0, n0 := 0, len(dp0)
				i1, n1 := 0, len(dp1)

				for i0 < n0 && i1 < n1 {
					v0, ts0 := dp0[i0][0], dp0[i0][1]
					v1, ts1 := dp1[i1][0], dp1[i1][1]
					its0, its1 := int64(ts0), int64(ts1)

					switch {
					case its0 == its1:
						datapoints = append(datapoints, [2]float64{v0 / v1, ts0})
						i0++
						i1++
						continue

					case its0 < its1:
						i0++
						continue

					case its0 > its1:
						i1++
					}
				}

				dp0 = datapoints
			}
			sm = append(sm, newMetrics(key, fmt.Sprintf("div(%s)", key), dp0))
		}
	}
	return sm
}

func (api *Api) diffSeries(m []*Metrics) []*Metrics {
	sm := make([]*Metrics, 0)
	{
		if n := len(m); n > 0 {
			key := m[0].Key
			dp0 := m[0].Datapoints

			for i := 1; i < n; i++ {
				datapoints := make([][2]float64, 0)
				dp1 := m[i].Datapoints
				i0, n0 := 0, len(dp0)
				i1, n1 := 0, len(dp1)

				for i0 < n0 && i1 < n1 {
					v0, ts0 := dp0[i0][0], dp0[i0][1]
					v1, ts1 := dp1[i1][0], dp1[i1][1]
					its0, its1 := int64(ts0), int64(ts1)

					switch {
					case its0 == its1:
						datapoints = append(datapoints, [2]float64{v0 - v1, ts0})
						i0++
						i1++
						continue

					case its0 < its1:
						datapoints = append(datapoints, [2]float64{v0, ts0})
						i0++
						continue

					case its0 > its1:
						datapoints = append(datapoints, [2]float64{v1, ts1})
						i1++
					}
				}

				for ; i0 < n0; i0++ {
					datapoints = append(datapoints, [2]float64{dp0[i0][0], dp0[i0][1]})
				}

				for ; i1 < n1; i1++ {
					datapoints = append(datapoints, [2]float64{dp1[i1][0], dp1[i1][1]})
				}

				dp0 = datapoints
			}
			sm = append(sm, newMetrics(key, fmt.Sprintf("diff(%s)", key), dp0))
		}
	}
	return sm
}
