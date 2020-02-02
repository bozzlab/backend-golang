package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":3333", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	i := 0

	// read i from cookie named `data`
	cookie, _ := r.Cookie("data")
	if cookie != nil {
		i, _ = strconv.Atoi(cookie.Value)
	}
	i++
	// w.Header().Add("Set-Cookie", "data")
	http.SetCookie(w, &http.Cookie{
		Name:     "data",
		Value:    strconv.Itoa(i),
		Path:     "/",
		MaxAge:   10,
		HttpOnly: true,
	})

	// set i to new cookie named `data`

	fmt.Fprintf(w, "new cookie: %d", i)
}
