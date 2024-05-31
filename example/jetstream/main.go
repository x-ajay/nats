package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL, nats.UserInfo("trilio", "trilio"))
	if err != nil {
		log.Fatal("[exit] failed to connect NETS app", err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal("[exit] failed to connect NATS JetStream", err)
	}

	sub, err := js.Subscribe("orders.us", processMessage, nats.BindStream("ORDERS"))
	if err != nil {
		log.Fatal("[exit] failed to subscribe to ORDERS", err)
	}
	defer sub.Unsubscribe()

	time.Sleep(time.Minute)
	fmt.Println("[exit] shutting down")
}

func processMessage(msg *nats.Msg) {
	fmt.Println("[info] received message: ", string(msg.Data))
}
