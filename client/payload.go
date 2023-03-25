package main

import (
	"errors"
	"strings"
)

type Payload struct {
	Type      PayloadType
	ClientID  string
	Message   string
	MessageID string
}

func decodePayload(bytePayload []byte) (Payload, error) {
	var (
		payload Payload
		err     error
	)

	rawPayload := string(bytePayload)

	payloadParts := strings.Split(rawPayload, "|")
	payloadType, err := getPayloadType(payloadParts)
	if err != nil {
		return Payload{}, err
	}

	switch payloadType {
	case MSG:
		payload, err = decodeMsg(payloadParts)
	case PONG:
		payload, err = decodePong(payloadParts)
	}

	if err != nil {
		return Payload{}, err
	}

	return payload, nil
}

func getPayloadType(parts []string) (PayloadType, error) {
	if len(parts) == 0 {
		return "", errors.New("invalid payload provided")
	}

	switch PayloadType(parts[0]) {
	case MSG:
		return MSG, nil
	case PONG:
		return PONG, nil
	default:
		return "", errors.New("unable to read payload type")
	}
}

func decodeMsg(parts []string) (Payload, error) {
	var payload Payload

	if len(parts) != 5 {
		return payload, errors.New("invalid message payload detected")
	}

	payload.Type = MSG
	payload.ClientID = parts[1]
	payload.Message = string(parts[3])
	payload.MessageID = parts[2]

	return payload, nil
}

func decodePong(parts []string) (Payload, error) {
	var payload Payload

	if len(parts) != 2 {
		return payload, errors.New("invalid message payload detected")
	}

	payload.Type = PONG
	payload.Message = parts[1]

	return payload, nil
}
