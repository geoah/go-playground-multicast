package main

import (
	"fmt"
	"time"

	multicast "github.com/geoah/go-playground-multicast"
)

func main() {
	g, err := multicast.New()
	if err != nil {
		panic(err)
	}

	l, err := g.Listen()
	if err != nil {
		panic(err)
	}

	go func() {
		for p := range l {
			fmt.Printf("Got packet from %s: %s\n", p.Source, string(p.Body))
		}
	}()

	for {
		_, err := g.Write([]byte(time.Now().Format(time.RFC3339)))
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 5)
	}
}
