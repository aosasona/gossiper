package main

import (
	"fmt"
	"unsafe"
)

type PayloadType string

const (
	MSG  PayloadType = "MSG"
	ACK  PayloadType = "ACK"
	PING PayloadType = "PING"
)

// TODO: implement functions to make different message types and SEPARATE the encapsulation and to-byte conversion process

func encapsulate(message string, msgType PayloadType) []byte {
	var payload []byte

	msgID := generateID(GeneratorArgs{NumOnly: true, Max: 999999})

	rawPayload := fmt.Sprintf("%s|%s|%s|%s", msgType, clientID, msgID, message)
	payloadSize := len(rawPayload) * int(unsafe.Sizeof(byte(0)))
	payloadWithTail := rawPayload + fmt.Sprintf("|%d", payloadSize)

	payload = []byte(payloadWithTail)

	return payload
}

func makeMsg() string {
	return ""
}
