package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

const (
	MAX_CLIENT_ID = 99999
	MIN_CLIENT_ID = 999
)

var clientID string

type GeneratorArgs struct {
	Min       int
	Max       int
	NumOnly   bool
	Delimiter string
}

func init() {
	clientID = generateID(GeneratorArgs{Min: MIN_CLIENT_ID, Max: MAX_CLIENT_ID})
	fmt.Printf("-> Client ID: %s\n\n", clientID)
}

func main() {
	noInput := flag.Bool("no-input", false, "if you want the client to take input from std-in")
	flag.Parse()

	remoteAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
	if err != nil {
		log.Printf("failed to resolve address: %s", err.Error())
		return
	}

	conn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		log.Printf("unable to dial server: %s", err.Error())
		return
	}

	defer conn.Close()

	var wg sync.WaitGroup

	if !*noInput {
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleOutgoingMsg(conn, remoteAddr)
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		handleIncomingMsg(conn, remoteAddr)
	}()

	wg.Wait()
}

func generateID(args GeneratorArgs) string {
	var id string

	if !args.NumOnly && args.Delimiter == "" {
		args.Delimiter = "_"
	}

	if args.Max == 0 {
		args.Max = 9999
	}

	if args.Min == 0 {
		args.Min = 99
	}

	prefixes := []rune{
		'a',
		'b',
		'c',
		'd',
		'e',
		'f',
		'h',
		'k',
		'm',
		'n',
		'p',
		'q',
		'r',
		's',
		't',
		'u',
		'w',
		'x',
	}
	rand.Seed(time.Now().UnixNano())
	suffix := rand.Intn((args.Max-args.Min)+1) + args.Min

	if !args.NumOnly {
		for i := 0; i < 6; i++ {
			prefix := prefixes[rand.Intn(len(prefixes)-1)]
			id += string(prefix)
		}
	}
	id += args.Delimiter + strconv.Itoa(suffix)
	return id
}

func encapsulate(message string) []byte {
	var payload []byte

	msgID := generateID(GeneratorArgs{NumOnly: true, Max: 999999})

	rawPayload := fmt.Sprintf("%s|%s|%s", clientID, msgID, message)
	payloadSize := len(rawPayload) * int(unsafe.Sizeof(byte(0)))
	payloadWithTail := rawPayload + fmt.Sprintf("|%d", payloadSize)

	payload = []byte(payloadWithTail)

	return payload
}
