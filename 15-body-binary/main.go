package main

import (
	"log"
	"mime"
	"net/http"
)

func main() {
	http.ListenAndServe(":3333", http.HandlerFunc(handler))
}

// func handler(w http.ResponseWriter, r *http.Request) {
// 	// print request body to console
// 	// user := r.PostFormValue("username")
// 	// password := r.PostFormValue("password")
// 	// log.Println(user, password)
// 	log.Print(r.Header.Get("Content-Type"))
// 	io.Copy(os.Stdout, r.Body)
// 	w.Write([]byte("ok\n"))
// }

func handler(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("content-type")
	mt, _, _ := mime.ParseMediaType(ct)
	if mt != "application/json" {
		log.Print("eueu")
	} else {
		log.Print("hehe")
	}
}
