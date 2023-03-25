package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	MAX_CLIENT_ID = 99999
	MIN_CLIENT_ID = 999
)

var (
	clientID string

	INITIAL_PING_COMPLETE bool
	PING_INTERVAL         int64
)

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

	killChan := make(chan os.Signal, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		ping(conn, killChan)
	}()

	if !*noInput {
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleOutgoingPayload(conn, remoteAddr, killChan)
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		handleIncomingPayload(conn, remoteAddr)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case msg := <-killChan:
			code := 0

			if msg == os.Interrupt {
				fmt.Println("\nreceived kill signal, exiting...")
				code = 1
			}

			close(killChan)
			os.Exit(code)
		}
	}()

	signal.Notify(killChan, os.Interrupt, syscall.SIGTERM)
	wg.Wait()
}
