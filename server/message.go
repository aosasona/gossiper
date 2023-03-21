package main

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

type PayloadType string

const (
	MSG  PayloadType = "MSG"
	ACK  PayloadType = "ACK"
	PING PayloadType = "PING"
)

type Payload struct {
	ID         string
	Type       PayloadType
	Message    string
	RawMessage []byte
	ClientID   string
	AckCount   int
	TotalBytes uint
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

func getMessageType(parts []string) (PayloadType, error) {
	if len(parts) < 2 {
		return "", fmt.Errorf("failed to get type for invalid payload")
	}

	switch PayloadType(parts[0]) {
	case MSG:
		return MSG, nil
	case ACK:
		return ACK, nil
	case PING:
		return PING, nil
	default:
		return "", fmt.Errorf(
			"unable to determine payload type for payload with header `%s`",
			parts[0],
		)
	}
}
