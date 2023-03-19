package main

import (
	"math/rand"
	"strconv"
	"time"
)

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
