run-server:
	go run ./server/*.go

run-client-sender:
	go run ./client/*.go

run-client:
	go run ./client/*.go -no-input

PHONY: run-server run-client-sender run-client