package main

import (
	"fmt"
	"net"
)

func handleIncomingMsg(conn *net.UDPConn, server *Server, broadcastChannel chan []byte) {
	for {
		msg := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(msg)

		payload, err := decodeMessage(msg[:n])
		if err != nil {
			fmt.Printf("[ERROR] %s\n", err.Error())
		}

		server.AddClient(Client{
			ID:   payload.ClientID,
			Addr: addr.IP.String(),
			Port: addr.Port,
		})

		fmt.Printf("[%s] %s\n", payload.ClientID, payload.Message)
		// broadcastChannel <- msg[:n]
	}
}

func handleBroadcast(conn *net.UDPConn, broadcastChannel chan []byte) {
	select {
	case msg := <-broadcastChannel:
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Printf("[ERROR] unable to broadcast message: %s", err.Error())
		}
	}
}
