package main

import (
	"compress/gzip"
	"io/ioutil"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/large", large)
	http.ListenAndServe(":3333", gzipMiddleware(mux))
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := []byte(`
		<!doctype html>
		<title>Compression Test</title>
		<h1>Index</h1>
		<p>Hello, World!</p>
	`)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
}

func large(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := []byte(`
		<!doctype html>
		<title>Compression Test</title>
		<h1>Large</h1>
		<p>The quick, brown fox jumps over a lazy dog. DJs flock by when MTV ax quiz prog. Junk MTV quiz graced by fox whelps. Bawds jog, flick quartz, vex nymphs. Waltz, bad nymph, for quick jigs vex! Fox nymphs grab quick-jived waltz. Brick quiz whangs jumpy veldt fox. Bright vixens jump; dozy fowl quack. Quick wafting zephyrs vex bold Jim. Quick zephyrs blow, vexing daft Jim. Sex-charged fop blew my junk TV quiz. How quickly daft jumping zebras vex. Two driven jocks help fax my big quiz. Quick, Baz, get my woven flax jodhpurs! "Now fax quiz Jack!" my brave ghost pled. Five quacking zephyrs jolt my wax bed. Flummoxed by job, kvetching W. zaps Iraq. Cozy sphinx waves quart jug of bad milk. A very bad quack might jinx zippy fowls. Few quips galvanized the mock jury box. Quick brown dogs jump over the lazy fox. The jay, pig, fox, zebra, and my wolves quack! Blowzy red vixens fight for a quick jump. Joaquin Phoenix was gazed by MTV for luck. A wizard’s job is to vex chumps quickly in fog. Watch "Jeopardy!", Alex Trebek's fun TV quiz game. Woven silk pyjamas exchanged for blue quartz. Brawny gods just</p>
		<p>Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque eu, pretium quis, sem. Nulla consequat massa quis enim. Donec pede justo, fringilla vel, aliquet nec, vulputate eget, arcu. In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam dictum felis eu pede mollis pretium. Integer tincidunt. Cras dapibus. Vivamus elementum semper nisi. Aenean vulputate eleifend tellus. Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim. Aliquam lorem ante, dapibus in, viverra quis, feugiat a, tellus. Phasellus viverra nulla ut metus varius laoreet. Quisque rutrum. Aenean imperdiet. Etiam ultricies nisi vel augue. Curabitur ullamcorper ultricies nisi. Nam eget dui. Etiam rhoncus. Maecenas tempus, tellus eget condimentum rhoncus, sem quam semper libero, sit amet adipiscing sem neque sed ipsum. Nam quam nunc, blandit vel, luctus pulvinar, hendrerit id, lorem. Maecenas nec odio et ante tincidunt tempus. Donec vitae sapien ut libero venenatis faucibus. Nullam quis ante. Etiam sit amet orci eget eros faucibus tincidunt. Duis leo. Sed fringilla mauris sit amet nibh. Donec sodales sagittis magna. Sed consequat, leo eget bibendum sodales, augue velit cursus nunc,</p>
	`)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func gzipMiddleware(h http.Handler) http.Handler {
	pool := &sync.Pool{
		New: func() interface{} {
			return gzip.NewWriter(ioutil.Discard) //dev/null/
		},
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check is browser support gzip
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}
		// check is request is web socket
		if r.Header.Get("Sec-WebSocket-Key") != "" {
			h.ServeHTTP(w, r)
			return
		}
		// check is response already encoded
		if w.Header().Get("Content-Encoding") != "" {
			h.ServeHTTP(w, r)
			return
		}
		// add header `Vary: Accept-Encoding` to response
		// if w.Header().Get("Vary") != "Accept-Encoding" {
		w.Header().Set("Vary", "Accept-Encoding")
		// }
		// wrap responseWriter with gzipResponseWriter
		nw := &gzipResponseWriter{
			ResponseWriter: w,
			pool:           pool,
		}
		// don't forget to close gzip writer when done
		defer nw.Close()
		// call next handler
		h.ServeHTTP(nw, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	gw          *gzip.Writer
	wroteHeader bool
	pool        *sync.Pool
}

var allowCompressType = map[string]bool{
	"text/plain":               true,
	"text/html":                true,
	"text/css":                 true,
	"text/xml":                 true,
	"text/javascript":          true,
	"application/x-javascript": true,
	"application/xml":          true,
}

func (w *gzipResponseWriter) init() {
	// check is response already encoded
	if w.Header().Get("Content-Encoding") != "" {
		return
	}
	// if content length not init, retrive from response's header `Content-Length`
	// skip if content length less than 860 bytes
	if length, _ := strconv.Atoi(w.Header().Get("Content-Length")); length < 860 {
		return
	}

	// skip if no match type
	ct := w.Header().Get("Content-Type")
	mt, _, _ := mime.ParseMediaType(ct)
	if !allowCompressType[mt] {
		return
	}
	// create new gzip writer
	// w.gw = gzip.NewWriter(w.ResponseWriter)
	w.gw = w.pool.Get().(*gzip.Writer)
	w.gw.Reset(w.ResponseWriter)
	// remove response header `Content-Length`
	w.Header().Del("Content-Length")
	// set response header `Content-Encoding: gzip`
	w.Header().Set("Content-Encoding", "gzip")
}

func (w *gzipResponseWriter) Write(p []byte) (int, error) {
	// write header if not written
	if !w.wroteHeader {
		w.WriteHeader(200)
	}
	// if gzip writer inited write to gzip writer
	if w.gw != nil {
		return w.gw.Write(p)
	}
	// if not bypass gzip writer to original response writer
	return w.ResponseWriter.Write(p)
}

func (w *gzipResponseWriter) Close() {
	if w.gw != nil {
		w.gw.Close()
		w.pool.Put(w.gw)
	}
	// close gzip writer if inited
}

func (w *gzipResponseWriter) WriteHeader(code int) {
	// do not write header twice
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	// init gzip writer
	w.init()
	// write header
	w.ResponseWriter.WriteHeader(code)
}
