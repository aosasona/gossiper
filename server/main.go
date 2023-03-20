package main

import (
	"flag"
	"net"
	"sync"
	"time"
)

type Client struct {
	ID       string
	Addr     string
	Port     int
	Alive    bool
	LastPing time.Time
}

/**
|==> MESSAGE FORMATS <==|
*
* Every message has 3 common parts; type|the client ID and the tail; total_bytes (used by the server to verify the data is in good shape)
*
* message: MSG|client_id|message_id|message|total_bytes
* ack: ACK|client_id|message_id|total_bytes
* ping: PING|client_id|total_bytes
*
* Clients reach out to the server at an interval provided by the SERVER and it is periodically checked to ensure that the client is still connected
*/

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

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	handleBroadcast(conn, broadcastChan)
	// }()

	wg.Wait()
	close(broadcastChan)

}
