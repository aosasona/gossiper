package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

type Client struct {
	ID   string
	Addr string
	Port int
}

type Payload struct {
	ID         string
	Message    string
	RawMessage []byte
	ClientID   string
	AckCount   int
	TotalBytes uint
}

// message format: clientID|messageID|message|totalBytes
// ack format: ACK|clientID|mesageID

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

func decodeMessage(message []byte) (*Payload, error) {
	payload := new(Payload)
	messageString := string(message)

	messageParts := strings.Split(messageString, "|")

	if len(messageParts) != 4 {
		return payload, fmt.Errorf("invalid payload received")
	}

	messageSize, err := strconv.Atoi(messageParts[3])
	if err != nil {
		return payload, fmt.Errorf("unable to parse total bytes: %s", err.Error())
	}

	originalMsg := fmt.Sprintf("%s|%s|%s", messageParts[0], messageParts[1], messageParts[2])
	reconstructedMessageSize := len(originalMsg) * int(unsafe.Sizeof(byte(0)))

	if reconstructedMessageSize != messageSize {
		return payload, fmt.Errorf("message has been corrupted in transit")
	}

	payload.ClientID = messageParts[0]
	payload.ID = messageParts[1]
	payload.Message = string(messageParts[2])
	payload.RawMessage = message

	return payload, nil
}
