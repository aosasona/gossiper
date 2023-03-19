package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
)

const (
	MAX_CLIENT_ID = 99999
	MIN_CLIENT_ID = 999
)

var clientID string

type GeneratorArgs struct {
	Min       int
	Max       int
	NumOnly   bool
	Delimiter string
}

func init() {
	clientID = generateID(GeneratorArgs{Min: MIN_CLIENT_ID, Max: MAX_CLIENT_ID})
	fmt.Printf("-> Client ID: %s\n\n", clientID)
}

func main() {
	noInput := flag.Bool("no-input", false, "if you want the client to take input from std-in")
	serverPort := flag.String("server-port", "8080", "port for the server")
	flag.Parse()

	remoteAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%s", *serverPort))
	if err != nil {
		log.Printf("failed to resolve address: %s", err.Error())
		return
	}

	conn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		log.Printf("unable to dial server: %s", err.Error())
		return
	}

	defer conn.Close()

	killChan := make(chan os.Signal)
	var wg sync.WaitGroup

	if !*noInput {
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleOutgoingMsg(conn, remoteAddr)
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		handleIncomingMsg(conn, remoteAddr)
	}()

	signal.Notify(killChan, os.Interrupt)
	wg.Wait()
}
