package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	// create tcp listener at :3333
	lis, _ := net.Listen("tcp", ":3333")
	// don't forget to close tcp listener when done
	defer lis.Close()
	for {
		// accept connection from listener
		conn, _ := lis.Accept()
		// don't forget to close connection when done
		go func() {
			defer conn.Close()
		}()

		// read data from connection
		res, _ := io.Copy(os.Stdout, conn)
		// print data out to console
		log.Println(res)
	}
}
