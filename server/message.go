package main

import (
	"errors"
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

const (
	MSG_LENGTH  = 5
	ACK_LENGTH  = 4
	PING_LENGTH = 3
)

type Payload struct {
	ID         string
	Type       PayloadType
	Message    string
	RawMessage []byte
	ClientID   string
	AckCount   int
	TotalBytes uint
	Length     int
}

func decodeMessage(message []byte, length int) (*Payload, error) {
	payload := new(Payload)
	message = message[:length]
	messageString := string(message)

	messageParts := strings.Split(messageString, "|")

	messageType, err := getMessageType(messageParts)
	if err != nil {
		return payload, err
	}

	if err = validateMessageLength(messageType, messageParts); err != nil {
		return payload, fmt.Errorf("invalid payload received")
	}

	messageSize, err := strconv.Atoi(messageParts[len(messageParts)-1])
	if err != nil {
		return payload, fmt.Errorf("unable to parse total bytes: %s", err.Error())
	}

	originalMsg := messageParts[0]
	for i := 1; i < len(messageParts)-1; i++ {
		originalMsg += fmt.Sprintf("|%s", messageParts[i])
	}
	reconstructedMessageSize := len(originalMsg) * int(unsafe.Sizeof(byte(0)))

	if reconstructedMessageSize != messageSize {
		return payload, fmt.Errorf("message has been corrupted in transit")
	}

	id, msg := extractMessageMeta(messageType, messageParts)

	payload.ClientID = messageParts[1]
	payload.ID = id
	payload.Message = string(msg)
	payload.RawMessage = message
	payload.Type = messageType
	payload.TotalBytes = uint(messageSize)
	payload.Length = length

	return payload, nil
}

func extractMessageMeta(msgType PayloadType, parts []string) (id, message string) {
	if msgType == MSG {
		id = parts[2]
		message = parts[3]
	}

	return
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

func validateMessageLength(msgType PayloadType, parts []string) error {
	var length int

	switch msgType {
	case MSG:
		length = MSG_LENGTH
	case ACK:
		length = ACK_LENGTH
	case PING:
		length = PING_LENGTH
	default:
		return errors.New("invalid message type received")
	}

	if len(parts) != length {
		return errors.New("invalid payload received")
	}

	return nil
}
