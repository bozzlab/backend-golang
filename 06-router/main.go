package main

import (
	"net/http"
)

func main() {
	router := http.NewServeMux()
	// router := &Router{}
	router.Handle("/", http.HandlerFunc(index))
	router.Handle("/about", http.HandlerFunc(about))
	http.ListenAndServe(":3333", router)
}

// Router struct
type Router struct {
	data map[string]http.Handler
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := router.data[r.URL.Path]
	// if h == nil {
	// 	http.router.data["/"]
	// }
	if h == nil {
		http.NotFound(w, r)
		return
	}
	h.ServeHTTP(w, r)
}

// Handle is function
func (router *Router) Handle(path string, h http.Handler) {
	if router.data == nil {
		router.data = make(map[string]http.Handler)
	}
	router.data[path] = h
}

// func router(w http.ResponseWriter, r *http.Request) {
// 	// implement router
// }

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("index"))
}

func about(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("about"))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound) // 404
	w.Write([]byte("404 page not found"))
}
