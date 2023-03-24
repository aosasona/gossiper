package main

import (
	"flag"
	"net"
	"sync"
)

func main() {
	server := new(Server)

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

	broadcastChan := make(chan []byte)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		handleIncomingMsg(conn, server, broadcastChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		handleBroadcast(conn, &server.clients, broadcastChan)
	}()

	wg.Wait()
	close(broadcastChan)

}
