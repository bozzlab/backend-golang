package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	// create tcp listener at :3333
	lis, _ := net.Listen("tcp", ":3333")

	// don't forget to close tcp listener when done
	defer lis.Close()

	for {
		// accept connection from listener
		conn, _ := lis.Accept()
		r := bufio.NewReader(conn)
		go func() {
			// don't forget to close connection when done
			defer conn.Close()

			// create bufio.NewReader from connection
			for {
				// read line from reader
				lineBytes, _, _ := r.ReadLine()
				line := string(lineBytes)
				// print data out to console
				fmt.Print(line)
				// check is data is empty string
				if line == "" {
					// response HTTP to connection
					conn.Write([]byte("HTTP/1.1 200 OK\n"))
					conn.Write([]byte("Content-Type: text/html\n"))
					conn.Write([]byte("Content-Length: 13\n"))
					conn.Write([]byte("\n"))
					conn.Write([]byte("<h1>eiei</h1>"))
					break
				}
			}
		}()
	}
}
