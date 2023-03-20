package main

import (
	"fmt"
	"unsafe"
)

const (
	MSG  PayloadType = "MSG"
	ACK  PayloadType = "ACK"
	PING PayloadType = "PING"
)

type PayloadType string

type Message string

type RawPayload struct {
	Message     Message
	PayloadType PayloadType
	MessageID   string
}

// TODO: implement functions to make different message types and SEPARATE the encapsulation and to-byte conversion process

func encapsulate(args RawPayload) []byte {

	switch args.PayloadType {
	case MSG:
		args.toMessage()
	case PING:
		args.toPing()
	case ACK:
		args.toAck(args.MessageID)
	default:
		panic(
			"invalid message type received",
		) // returning here just breaks the switch, we don't want it to continue at all
	}

	payloadSize := len(args.Message) * int(unsafe.Sizeof(byte(0)))
	payloadWithTrailer := args.Message + Message(fmt.Sprintf("|%d", payloadSize))

	return toByte(payloadWithTrailer)
}

func toByte(data Message) []byte {
	return []byte(data)
}

func (e *RawPayload) toMessage() {
	msgID := generateID(GeneratorArgs{NumOnly: true, Max: 999999})
	newMsg := fmt.Sprintf("%s|%s|%s|%s", MSG, clientID, msgID, string(e.Message))
	e.Message = Message(newMsg)
}

func (e *RawPayload) toPing() {
	newMsg := fmt.Sprintf("%s|%s", PING, clientID)
	e.Message = Message(newMsg)
}

func (e *RawPayload) toAck(msgID string) {
	newMsg := fmt.Sprintf("%s|%s|%s", ACK, clientID, msgID)
	e.Message = Message(newMsg)
}
