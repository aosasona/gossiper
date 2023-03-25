package main

import (
	"fmt"
	"net"
)

func handleIncomingPayload(conn *net.UDPConn, server *Server, broadcastChannel *chan []byte) {
	for {
		msg := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(msg)

		payload, err := decodeMessage(msg, n)
		if err != nil {
			fmt.Printf("\n[ERROR] %s\n", err.Error())
		}

		if _, exists := server.FindClient(payload.ClientID); !exists {
			server.AddClient(Client{
				ID:    payload.ClientID,
				Addr:  addr.IP.String(),
				Port:  addr.Port,
				Alive: true,
			})
		}

		switch payload.Type {
		case MSG:
			err = payload.handleMsg(broadcastChannel, server)
		case ACK:
			err = payload.handleAck(server)
		case PING:
			err = payload.handlePing(server)
		default:
			fmt.Printf("\n[ERROR] unable to handle message: invalid message type\n")
			continue
		}

		if err != nil {
			fmt.Printf("\n[ERROR] unable to handle message: %s\n", err.Error())
			continue
		}
	}
}

func handleBroadcast(conn *net.UDPConn, server *Server, broadcastChannel *chan []byte) {
	select {
	case msg := <-*broadcastChannel:
		for _, client := range server.GetClients() {
			_, err := conn.WriteToUDP(msg, &net.UDPAddr{
				IP:   net.ParseIP(client.Addr),
				Port: client.Port,
			})
			if err != nil {
				fmt.Printf("[ERROR] unable to broadcast message: %s", err.Error())
			}
		}
	}
}
