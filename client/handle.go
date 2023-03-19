package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func handleIncomingMsg(conn *net.UDPConn, serverAddr *net.UDPAddr) {
	buf := make([]byte, 1024)
	for {
		_, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("error reading incoming message: %v\n", err)
			continue
		}

		fmt.Printf("[RECEIVED] %s", buf)
	}
}

func handleOutgoingMsg(conn *net.UDPConn, serverAddr *net.UDPAddr) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s>", clientID)
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("unable to read input: %s", err.Error())
			continue
		}

		msg = msg[:len(msg)-1]
		_, err = conn.Write(encapsulate(msg))
		if err != nil {
			log.Printf("error sending message: %s", err.Error())
			continue
		}
	}
}
