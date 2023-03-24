package main

import (
	"fmt"
	"time"
)

type Client struct {
	ID       string
	Addr     string
	Port     int
	Alive    bool
	LastPing time.Time
}

type Server struct {
	clients []Client
	queue   []Payload // obviously you'd use a more efficient in-memory database here
}

func (s *Server) FindClient(id string) (Client, bool) {
	for _, client := range s.clients {
		if client.ID == id {
			return client, true
		}
	}

	return Client{}, false
}

func (s *Server) AddClient(client Client) {
	if _, exists := s.FindClient(client.ID); !exists {
		s.clients = append(s.clients, client)
		fmt.Printf(
			"[JOIN] Client %s (%s) joined\n",
			client.ID,
			fmt.Sprintf("%s:%d", client.Addr, client.Port),
		)
	}
}
