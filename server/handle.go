package main

import (
	"fmt"
	"net"
)

func handleIncomingMsg(conn *net.UDPConn, server *Server, broadcastChannel chan []byte) {
	// TODO: handle ACKs and PINGs

	for {
		msg := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(msg)

		payload, err := decodeMessage(msg[:n])
		if err != nil {
			fmt.Printf("[ERROR] %s\n", err.Error())
		}

		server.AddClient(Client{
			ID:    payload.ClientID,
			Addr:  addr.IP.String(),
			Port:  addr.Port,
			Alive: true,
		})

		fmt.Printf("[%s] %s\n", payload.ClientID, payload.Message)
		if payload.Type == MSG {
			broadcastChannel <- msg[:n]
		}
	}
}

func handleBroadcast(conn *net.UDPConn, clients *[]Client, broadcastChannel chan []byte) {
	select {
	case msg := <-broadcastChannel:

		for _, client := range *clients {
			addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", client.Addr, client.Port))
			if err != nil {
				fmt.Printf("[ERROR] unable to resolve address for client (%s): %s", client.ID, err.Error())
			}

			_, err = conn.WriteToUDP(msg, addr)
			if err != nil {
				fmt.Printf("[ERROR] unable to broadcast message: %s", err.Error())
			}
		}
	}
}
