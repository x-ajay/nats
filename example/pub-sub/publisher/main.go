package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("[error] while connecting nats ", err)
	}
	defer nc.Close()

	for {
		msg := "Hello World!"
		err = nc.Publish("intros", []byte(msg))
		if err != nil {
			fmt.Println("[error] message failed to published", err)
		}
		log.Println("[info] message published", msg)
		time.Sleep(1 * time.Second)
	}

}
