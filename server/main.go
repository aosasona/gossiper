package main

import (
	"flag"
	"net"
	"sync"
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

	server := NewServer(conn)

	server.GeneratePingInterval()
	broadcastChan := make(chan []byte)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		handleIncomingPayload(conn, server, broadcastChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		handleBroadcast(conn, server, broadcastChan)
	}()

	wg.Wait()
	close(broadcastChan)
}
