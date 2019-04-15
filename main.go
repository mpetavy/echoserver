package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
)

var (
	port = flag.Int("p", 9000, "echo server port ")
)

func main() {
	flag.Parse()
	server, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if server == nil {
		panic("couldn't start listening: " + err.Error())
	}
	conns := clientConns(server)
	for {
		go handleConn(<-conns)
	}
}

func clientConns(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	i := 0
	go func() {
		for {
			client, err := listener.Accept()
			if client == nil {
				fmt.Printf("couldn't accept: " + err.Error())
				continue
			}
			i++
			fmt.Printf("%d: %v <-> %v\n", i, client.LocalAddr(), client.RemoteAddr())
			ch <- client
		}
	}()
	return ch
}

func handleConn(client net.Conn) {
	b := make([]byte, 8192)
	for {
		n, err := client.Read(b)
		if err != nil { // EOF, or worse
			break
		}
		client.Write(b[:n])
	}
}
