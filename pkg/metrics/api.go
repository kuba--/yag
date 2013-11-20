package metrics

import (
	"fmt"

	"github.com/golang/glog"
	"strconv"
	"strings"
)

type Api struct {
	maxDataPoints int
}

func NewApi(maxDataPoints int) *Api {
	api := new(Api)
	api.maxDataPoints = maxDataPoints
	return api
}

func (api *Api) Value(name string, from int64, to int64) interface{} {
	glog.Infof("Api.Value(%s, %d, %d)", name, from, to)

	// check constant value first
	if f, err := strconv.ParseFloat(name, 10); err == nil {
		return f
	}
	return Get(name, from, to, api.maxDataPoints)
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
					glog.Warningln("[ ! ]\tFunction not supported: ", name)
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
				datapoints := make([]Pt, 0)
				dp1 := m[i].Datapoints
				i0, n0 := 0, len(dp0)
				i1, n1 := 0, len(dp1)

				for i0 < n0 && i1 < n1 {
					pt0 := dp0[i0]
					pt1 := dp1[i1]
					ts0, ts1 := int64(*pt0[1]), int64(*pt1[1])

					pt := Pt{nil, new(float64)}
					switch {
					case ts0 == ts1:
						*pt[1] = *pt0[1]

						if pt0[0] != nil || pt1[0] != nil {
							pt[0] = new(float64)
							if pt0[0] != nil {
								*pt[0] = *pt[0] + *pt0[0]
							}
							if pt1[0] != nil {
								*pt[0] = *pt[0] + *pt1[0]
							}
						}

						datapoints = append(datapoints, pt)
						i0++
						i1++
						continue

					case ts0 < ts1:
						*pt[1] = *pt0[1]

						if pt0[0] != nil {
							pt[0] = new(float64)
							*pt[0] = *pt0[0]
						}

						datapoints = append(datapoints, pt)
						i0++
						continue

					case ts0 > ts1:
						*pt[1] = *pt1[1]

						if pt1[0] != nil {
							pt[0] = new(float64)
							*pt[0] = *pt1[0]
						}

						datapoints = append(datapoints, pt)
						i1++
					}
				}

				for ; i0 < n0; i0++ {
					datapoints = append(datapoints, dp0[i0])
				}

				for ; i1 < n1; i1++ {
					datapoints = append(datapoints, dp1[i1])
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
				datapoints := make([]Pt, 0)
				dp1 := m[i].Datapoints
				i0, n0 := 0, len(dp0)
				i1, n1 := 0, len(dp1)

				for i0 < n0 && i1 < n1 {
					pt0 := dp0[i0]
					pt1 := dp1[i1]
					ts0, ts1 := int64(*pt0[1]), int64(*pt1[1])

					switch {
					case ts0 == ts1:
						pt := Pt{nil, new(float64)}
						*pt[1] = *pt0[1]

						if pt0[0] != nil && pt1[0] != nil && *pt1[0] != 0 {
							pt[0] = new(float64)
							*pt[0] = *pt0[0] / *pt1[0]
						}

						datapoints = append(datapoints, pt)
						i0++
						i1++
						continue

					case ts0 < ts1:
						i0++
						continue

					case ts0 > ts1:
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
				datapoints := make([]Pt, 0)
				dp1 := m[i].Datapoints
				i0, n0 := 0, len(dp0)
				i1, n1 := 0, len(dp1)

				for i0 < n0 && i1 < n1 {
					pt0 := dp0[i0]
					pt1 := dp1[i1]
					ts0, ts1 := int64(*pt0[1]), int64(*pt1[1])

					pt := Pt{nil, new(float64)}
					switch {
					case ts0 == ts1:
						*pt[1] = *pt0[1]

						if pt0[0] != nil || pt1[0] != nil {
							pt[0] = new(float64)
							if pt0[0] != nil {
								*pt[0] = *pt0[0]
							}
							if pt1[0] != nil {
								*pt[0] = *pt[0] - *pt1[0]
							}
						}

						datapoints = append(datapoints, pt)
						i0++
						i1++
						continue

					case ts0 < ts1:
						*pt[1] = *pt0[1]
						if pt0[0] != nil {
							pt[0] = new(float64)
							*pt[0] = *pt0[0]
						}

						datapoints = append(datapoints, pt)
						i0++
						continue

					case ts0 > ts1:
						*pt[1] = *pt1[1]
						if pt1[0] != nil {
							pt[0] = new(float64)
							*pt[0] = *pt1[0]
						}

						datapoints = append(datapoints, pt)
						i1++
					}
				}

				for ; i0 < n0; i0++ {
					datapoints = append(datapoints, dp0[i0])
				}

				for ; i1 < n1; i1++ {
					datapoints = append(datapoints, dp1[i1])
				}

				dp0 = datapoints
			}
			sm = append(sm, newMetrics(key, fmt.Sprintf("diff(%s)", key), dp0))
		}
	}
	return sm
}
