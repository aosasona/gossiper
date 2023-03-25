package main

import (
	"fmt"
	"math/rand"
	"net"
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
	conn          *net.UDPConn
	clients       []Client
	queue         []Payload // obviously you'd use a more efficient in-memory database here
	ping_interval int64
}

func NewServer(conn *net.UDPConn) *Server {
	return &Server{
		conn:          conn,
		clients:       []Client{},
		queue:         []Payload{},
		ping_interval: 0,
	}
}

func (s *Server) Conn() *net.UDPConn {
	return s.conn
}

func (s *Server) HasPingInterval() bool {
	return s.ping_interval != 0
}

func (s *Server) GeneratePingInterval() {
	if !s.HasPingInterval() {
		rand.Seed(int64(time.Millisecond))
		s.ping_interval = int64(rand.Intn((1000-500)+1) + 500)
	}
}

func (s *Server) GetPingInterval() int64 {
	return s.ping_interval
}

func (s *Server) GetClients() []Client {
	return s.clients
}

func (s *Server) FindClient(id string) (Client, bool) {
	for _, client := range s.clients {
		if client.ID == id {
			return client, true
		}
	}

	return Client{}, false
}

func (s *Server) FindClientIdx(id string) (int, bool) {
	for idx, client := range s.clients {
		if client.ID == id {
			return idx, true
		}
	}

	return -1, false
}

func (s *Server) AddClient(client Client) {
	s.clients = append(s.clients, client)
	fmt.Printf(
		"[CONNECTION] Client %s (%s) joined\n",
		client.ID,
		fmt.Sprintf("%s:%d", client.Addr, client.Port),
	)
}

func (s *Server) RemoveClient(id string) {
	if idx, exists := s.FindClientIdx(id); exists {
		s.clients[idx] = s.clients[len(s.clients)-1]
		s.clients = s.clients[:len(s.clients)-1]
	}
}

func (s *Server) FindPayloadIdx(id string) (int, bool) {
	for idx, payload := range s.queue {
		if payload.ID == id {
			return idx, true
		}
	}

	return -1, false
}

func (s *Server) AddAck(id string) bool {
	for _, payload := range s.queue {
		if payload.ID == id {
			payload.AckCount += 1
			return true
		}
	}
	return false
}

func (s *Server) FindInQueue(id string) (Payload, bool) {
	for _, payload := range s.queue {
		if payload.ID == id {
			return payload, true
		}
	}

	return Payload{}, false
}

func (s *Server) AddToQueue(payload Payload) {
	if _, exists := s.FindInQueue(payload.ID); !exists {
		s.queue = append(s.queue, payload)
	}
}

func (s *Server) RemoveFromQueue(id string) {
	if idx, exists := s.FindPayloadIdx(id); exists {
		s.queue[idx] = s.queue[len(s.queue)-1]
		s.queue = s.queue[:len(s.queue)-1]
	}
}
