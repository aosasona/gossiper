package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
)

func ping(conn *net.UDPConn, killChannel chan os.Signal) {
	PING_RETRY_COUNT := 0

	if PING_INTERVAL == 0 {
		PING_INTERVAL = 500
	}

	for range time.Tick(time.Millisecond * time.Duration(PING_INTERVAL)) {
		_, err := conn.Write(encapsulate(RawPayload{
			PayloadType: PING,
		}))

		if err != nil {
			if PING_RETRY_COUNT < 3 {
				fmt.Printf("\n[ERROR] Failed to ping, retrying (%d)", PING_RETRY_COUNT+1)
				PING_RETRY_COUNT++
				continue
			}

			fmt.Println("\n[FATAL] Unable to reach server, shutting down...")
			killChannel <- syscall.SIGTERM
			return
		}

		if !INITIAL_PING_COMPLETE {
			INITIAL_PING_COMPLETE = true
		}

		if PING_RETRY_COUNT > 0 {
			PING_RETRY_COUNT = 0
		}
	}
}
