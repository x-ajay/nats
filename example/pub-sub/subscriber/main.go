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
		log.Fatal("error while connecting nats ", err)
	}
	defer nc.Close()

	nc.Subscribe("intros", func(msg *nats.Msg) {
		fmt.Println("[info] received message: ", string(msg.Data))
	})

	time.Sleep(1 * time.Minute)
}
