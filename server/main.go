package main

import (
	"flag"
	"net"
)

func main() {
	port := flag.Int("port", 8080, "port to listen on")
	addr := net.UDPAddr{
		Port: *port,
		IP:   net.ParseIP("127.0.0.1"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
}

func handleMessage(message []byte) error {
	return nil
}

func decodeMessage() error {
	return nil
}
