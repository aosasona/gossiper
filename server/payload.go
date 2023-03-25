package main

import (
	"errors"
	"fmt"
	"net"
)

func (p *Payload) handleMsg(broadcastChannel *chan []byte, server *Server) error {
	fmt.Printf("[%s] %s\n", p.ClientID, p.Message)
	*broadcastChannel <- p.RawMessage[:p.Length]
	return nil
}

func (p *Payload) handleAck(server *Server) error {
	return nil
}

func (p *Payload) handlePing(server *Server) error {
	client, exists := server.FindClient(p.ClientID)
	if !exists {
		return errors.New("received ping from unregistered client")
	}

	_, err := server.Conn().WriteToUDP(
		[]byte(fmt.Sprintf("PONG|%d", server.GetPingInterval())),
		&net.UDPAddr{IP: net.ParseIP(client.Addr), Port: client.Port},
	)
	if err != nil {
		return err
	}
	return nil
}
