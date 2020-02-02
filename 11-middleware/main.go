package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

func main() {
	r := router{}
	r.Get("/", http.HandlerFunc(index))
	r.Get("/about", http.HandlerFunc(about))

	var h http.Handler

	// wrap router with middleware
	h = requestLogger(&r)
	http.ListenAndServe(":3333", h)
}

type logRecord struct {
	Time         time.Time `json:"time"`
	RemoteIP     string    `json:"remote_ip"`
	Host         string    `json:"host"`
	Method       string    `json:"method"`
	URI          string    `json:"uri"`
	Status       int       `json:"status"`
	Latency      int64     `json:"latency"`
	LatencyHuman string    `json:"latency_human"`
	BytesIn      int64     `json:"bytein"`
	BytesOut     int64     `json:"byteout"`
}

type logResponseWritter struct {
	w           http.ResponseWriter
	code        int
	lenByte     int64
	wroteHeader bool
}

func (w *logResponseWritter) Header() http.Header {
	return w.w.Header()
}

func (w *logResponseWritter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.code = code
	w.w.WriteHeader(code)
	w.wroteHeader = true
}

func (w *logResponseWritter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(200)
	}
	lenByte, err := w.w.Write(b)
	w.lenByte += int64(lenByte)
	return lenByte, err

}

func requestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &logRecord{}
		rec.Time = start
		rec.Method = r.Method
		rec.Host = r.Host
		rec.URI = r.RequestURI
		rec.RemoteIP = r.RemoteAddr
		rec.BytesIn = r.ContentLength

		newWrite := &logResponseWritter{w: w}
		h.ServeHTTP(newWrite, r)
		diff := time.Since(start)

		rec.Latency = int64(diff)
		rec.Status = newWrite.code
		rec.BytesOut = newWrite.lenByte
		json.NewEncoder(os.Stdout).Encode(*rec)
	})
}

type router struct {
	// path => method => handler
	path map[string]map[string]http.Handler
}

type path struct {
	Method  string
	Path    string
	Handler http.Handler
}

func (router *router) Add(m, p string, h http.Handler) {
	if router.path == nil {
		router.path = make(map[string]map[string]http.Handler)
	}

	if router.path[p] == nil {
		router.path[p] = make(map[string]http.Handler)
	}

	router.path[p][m] = h
}

func (router *router) Get(p string, h http.Handler) {
	router.Add(http.MethodGet, p, h)
}

func (router *router) Post(p string, h http.Handler) {
	router.Add(http.MethodPost, p, h)
}

func (router *router) Put(p string, h http.Handler) {
	router.Add(http.MethodPut, p, h)
}

func (router *router) Patch(p string, h http.Handler) {
	router.Add(http.MethodPatch, p, h)
}

func (router *router) Delete(p string, h http.Handler) {
	router.Add(http.MethodDelete, p, h)
}

func (router *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if router.path == nil {
		http.NotFound(w, r)
		return
	}

	p := router.path[r.URL.Path]
	if p == nil {
		http.NotFound(w, r)
		return
	}

	h := p[r.Method]
	if h == nil {
		http.NotFound(w, r)
		return
	}

	h.ServeHTTP(w, r)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("index"))
}

func about(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("about"))
}
