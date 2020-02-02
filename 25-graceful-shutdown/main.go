package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	http.HandleFunc("/", index)

	// create http.Server
	server := &http.Server{
		Addr: ":3333",
	}

	// start server on another goroutine
	go server.ListenAndServe()

	// create buffered (size=1) channel for os.Signal
	stop := make(chan os.Signal, 1)

	// call signal.Notify to notify channel when received syscall.SIGTERM
	signal.Notify(stop, syscall.SIGTERM, os.Interrupt)

	// block until receive signal
	<-stop
	log.Println("Pending")
	<-stop
	log.Println("Receive signal")

	// shutdown server
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 5*time.Minute)

	server.Shutdown(ctx)
	log.Println("Server already shutdown")
}

func index(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	w.Write([]byte("ok"))
}
