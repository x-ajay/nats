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

	_, err = js.AddConsumer("ORDERS", &nats.ConsumerConfig{
		Durable:      "durable-consumer",
		Description:  "durable-consumer description",
		ReplayPolicy: nats.ReplayInstantPolicy,
	})

	if err != nil {
		log.Fatal("[exit] failed to add consumer", err)
	}

	sub, err := js.PullSubscribe("orders.us", "durable-consumer")
	if err != nil {
		log.Fatal("[exit] failed to subscribe to orders durable-consumer", err)
	}

	go processMessage(sub)

	time.Sleep(time.Minute)
	sub.Unsubscribe()

	fmt.Println("[exit] shutting down")
}

func processMessage(sub *nats.Subscription) {
	for sub.IsValid() {
		msgs, err := sub.Fetch(1)
		msg := msgs[0]
		if err == nil {
			fmt.Println("[info] received message: ", string(msg.Data))
		}
	}
}
