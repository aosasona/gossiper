package main

import (
	"fmt"
	"unsafe"
)

type PayloadType string

type Message string

const (
	MSG  PayloadType = "MSG"
	ACK  PayloadType = "ACK"
	PING PayloadType = "PING"
)

// TODO: implement functions to make different message types and SEPARATE the encapsulation and to-byte conversion process

func encapsulate(message string, msgType PayloadType) []byte {
	rawPayload := Message(message)

	switch msgType {
	case MSG:
		rawPayload.toMessage()
	case PING:
		rawPayload.toPing()
	default:
		panic(
			"invalid message type received",
		) // returning here just breaks the switch, we don't want it to continue at all
	}

	payloadSize := len(rawPayload) * int(unsafe.Sizeof(byte(0)))
	payloadWithTrailer := rawPayload + Message(fmt.Sprintf("|%d", payloadSize))

	return toByte(payloadWithTrailer)
}

func toByte(data Message) []byte {
	return []byte(data)
}

func (m *Message) toMessage() {
	msgID := generateID(GeneratorArgs{NumOnly: true, Max: 999999})
	newMsg := fmt.Sprintf("%s|%s|%s|%s", MSG, clientID, msgID, string(*m))
	*m = Message(newMsg)
}

func (m *Message) toPing() {
	newMsg := fmt.Sprintf("%s|%s", MSG, clientID)
	*m = Message(newMsg)
}
