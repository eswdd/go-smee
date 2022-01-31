package main

import (
	"fmt"
	"log"
	"os"

	"github.com/eswdd/go-smee/smee"
)

func main() {
	var source *string
	var err error
	if len(os.Args) > 1 {
		source = &os.Args[1]
	} else {
		source, err = smee.CreateChannel()
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Subscribing to smee source: " + *source)

	logger := log.Logger{}

	target := make(chan smee.SSEvent)
	client := smee.NewClient(source, target, &logger)

	fmt.Println("Client initialised")

	sub, err := client.Start()
	if err != nil {
		panic(err)
	}

	fmt.Println("Client running")

	for ev := range target {
		// do what you want with the event
		fmt.Printf("Received event: id=%v, name=%v, payload=%v\n", ev.ID, ev.Name, string(ev.Data))
	}

	sub.Stop()
}
