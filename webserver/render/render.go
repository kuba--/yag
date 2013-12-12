package render

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/kuba--/yag/pkg/api"
	"github.com/kuba--/yag/pkg/config"
	"github.com/kuba--/yag/pkg/metrics"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	glog.Infoln(r.RequestURI)

	// ResponseWriter wrapper
	w.Header().Set("Server", "YAG")
	w.Header().Set("Content-Type", "application/json")
	rw := &RenderResponseWriter{w: w}

	// Handler composition
	http.TimeoutHandler(&RenderHandler{}, time.Duration(config.Cfg.Webserver.Timeout)*time.Second,
		http.StatusText(http.StatusRequestTimeout)).ServeHTTP(rw, r)

	glog.Infof("[%v] in %v\n", rw.Code, time.Now().Sub(t))
}

// RenderResponseWriter retrieves StatusCode from ResponseWriter
type RenderResponseWriter struct {
	Code int // the HTTP response code from WriteHeader
	w    http.ResponseWriter
}

// Header returns the response headers.
func (w *RenderResponseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *RenderResponseWriter) Write(buf []byte) (int, error) {
	return w.w.Write(buf)
}

// WriteHeader sets Code.
func (w *RenderResponseWriter) WriteHeader(code int) {
	w.Code = code
	w.w.WriteHeader(code)
}

// Render Handler
type RenderHandler struct {
	jsonp         string
	target        string
	from          int64
	to            int64
	maxDataPoints int
}

// GET: /render?target=my.key&from=-1h[&to=...&jsonp=...]
func (h *RenderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		glog.Infof("%v\n", r.URL)
		err := h.parseQuery(r)
		if err != nil {
			glog.Warningln(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.jsonResponse(w, api.Eval(h.target, h.from, h.to, &metrics.Api{h.maxDataPoints}))
		glog.Flush()

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

/*
 * JSON Response:
 * [
 *  {"target": "status.200", "datapoints": [[1720.0, 1370846820], ...], },
 *  {"target": "status.204", "datapoints": [[null, 1370846820], ..., ]}
 * ]s
 */
func (h *RenderHandler) jsonResponse(w http.ResponseWriter, data interface{}) {
	if m, ok := data.([]*metrics.Metrics); ok {
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, "%s([", h.jsonp)
		for i, mi := range m {
			if i > 0 {
				fmt.Fprintf(w, ",")
			}

			fmt.Fprintf(w, `{"target":"%s","datapoints":[`, mi.Target)
			for ii := 0; ii < len(mi.Datapoints); ii++ {
				if ii > 0 {
					fmt.Fprintf(w, ",")
				}

				val := "null"
				if mi.Datapoints[ii][0] != nil {
					val = fmt.Sprintf("%.2f", *mi.Datapoints[ii][0])
				}

				fmt.Fprintf(w, "[%s, %.0f]", val, *mi.Datapoints[ii][1])
			}
			fmt.Fprintf(w, "]}")
		}
		fmt.Fprintf(w, "])")
		glog.Infof("%v\n", data)
	} else {
		glog.Errorf("%v\n", data)
		http.Error(w, fmt.Sprintf("%v", data), http.StatusBadRequest)
	}
}

// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func (h *RenderHandler) parseQuery(r *http.Request) error {
	parseDuration := func(duration string) (time.Duration, error) {
		return time.ParseDuration(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(duration, "seconds", "s", -1), "sec", "s", -1), "minutes", "m", -1), "min", "m", -1), "hours", "h", -1))
	}

	f, err := parseDuration(r.FormValue("from"))
	if err != nil {
		return err
	}
	h.from = time.Now().Add(f).Unix()

	t := r.FormValue("to")
	if len(t) < 1 {
		h.to = time.Now().Unix()
	} else {
		t, err := parseDuration(t)
		if err != nil {
			return err
		}
		h.to = time.Now().Add(t).Unix()
	}
	h.maxDataPoints, _ = strconv.Atoi(r.FormValue("maxDataPoints"))
	h.jsonp = r.FormValue("jsonp")
	h.target = fmt.Sprintf("_(%s)", strings.Join(r.URL.Query()["target"], ","))

	return nil
}
