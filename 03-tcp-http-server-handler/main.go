package main

import (
	"net/http"
)

func main() {
	// create tcp listener at :3333
	// lis, err := net.Listen("tcp", ":3333")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // don't forget to close tcp listener when done
	// defer lis.Close()
	http.ListenAndServe(":3333", http.HandlerFunc(handler))
	// for {
	// 	// accept connection from listener
	// 	conn, _ := lis.Accept()

	// 	go func() {
	// 		// don't forget to close connection when done
	// 		defer conn.Close()

	// 		// create bufio.NewReader from connection
	// 		r, _ := http.ReadRequest(bufio.NewReader(conn))
	// 		w := &responseWriter{conn: conn}
	// 		handler(w, r)
	// 		// response HTTP to connection
	// 	}()
	// 	// don't forget to close connection when done

	// 	// create bufio.NewReader from connection

	// 	// use http.ReadRequest to parse HTTP request

	// 	// create new responseWriter

	// 	// call handler with responseWriter and request
	// }
}

// type responseWriter struct {
// 	conn   net.Conn
// 	header http.Header
// }

// func (w *responseWriter) Header() http.Header {
// 	if w.header == nil {
// 		w.header = make(http.Header)
// 	}
// 	return w.header
// }

// func (w *responseWriter) WriteHeader(statusCode int) {
// 	statusText := http.StatusText(statusCode)
// 	fmt.Fprintf(w.conn, "HTTP1.1 %d %s\n", statusCode, statusText)
// 	w.Header().Write(w.conn)
// 	w.conn.Write([]byte("\n"))
// }

// func (w *responseWriter) Write(b []byte) (int, error) {
// 	return w.conn.Write(b)
// }

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	w.Write([]byte("<h1>123ddd456</h1>"))
}

// func (w *responseWriter) Write(p []byte) (int, error) {
// 	if true { // check is header not written
// 		// write header with status code 200 (default)
// 	}
// 	// write p to connection
// 	return 0, nil // return wrote bytes and/or error
// }

// func (w *responseWriter) WriteHeader(code int) {
// 	// DO NOT write header > 1 time

// 	// write HTTP version and status code
// 	// write HTTP headers
// 	// write empty line to start HTTP body
// }

// func handler(w http.ResponseWriter, r *http.Request) {
// 	// write data to responseWriter
// }
