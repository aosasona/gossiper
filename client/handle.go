package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"syscall"
)

func handleIncomingPayload(conn *net.UDPConn, serverAddr *net.UDPAddr) {
	buf := make([]byte, 1024)
	for {
		if INITIAL_PING_COMPLETE {
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				fmt.Printf("\n[ERROR] error reading incoming payload: %v\n", err)
				continue
			}

			raw := buf[:n]
			payload, err := decodePayload(raw)
			if err != nil {
				fmt.Printf("\n[ERROR] error decode incoming payload: %v\n", err)
				continue
			}

			switch payload.Type {
			case MSG:
				handleMsg(conn, payload)
			case PONG:
				handlePong(payload)
			}

			break
		}
	}
}

func handleOutgoingPayload(conn *net.UDPConn, serverAddr *net.UDPAddr, killChannel chan os.Signal) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s> ", clientID)
		msg, err := reader.ReadString('\n')

		if err != nil {
			fmt.Printf("\n[ERROR] unable to read input: %s\n", err.Error())
			continue
		}

		msg = msg[:len(msg)-1]
		if msg == ":exit" || msg == ":q" {
			fmt.Println("Received exit command, exiting now...")
			killChannel <- syscall.SIGTERM
			break
		}
		_, err = conn.Write(encapsulate(RawPayload{
			Message:     Message(msg),
			PayloadType: MSG,
		}))
		if err != nil {
			fmt.Printf("\n[ERROR] error sending message: %s\n", err.Error())
			break
		}

		continue
	}
}

func handleMsg(conn *net.UDPConn, payload Payload) error {
	fmt.Printf("\n[RECEIVED] %s\n", payload.Message)

	ack := encapsulate(RawPayload{
		MessageID:   payload.MessageID,
		Message:     Message(payload.Message),
		PayloadType: payload.Type,
	})

	retryCount := 0
	for {
		// we need to keep retrying to send the ack until it goes through or terminate after 5 tries
		_, err := conn.Write(ack)
		if err != nil {
			if retryCount >= 5 {
				break
			}
			fmt.Printf("\n[ERROR] error sending ack: %s\n", err.Error())
			retryCount += 1
			continue
		}
		break
	}

	return nil
}

func handlePong(payload Payload) error {
	interval, err := strconv.Atoi(payload.Message)
	if err != nil {
		return nil
	}

	PING_INTERVAL = int64(interval)

	return nil
}
