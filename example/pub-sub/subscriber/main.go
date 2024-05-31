package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Payload struct {
	Data  string `json:"data,omitempty"`
	Count int    `json:"count,omitempty"`
}

func main() {

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("error while connecting nats ", err)
	}
	defer nc.Close()

	subscription, err := nc.QueueSubscribe("intros", "zip", processMsg)
	if err != nil {
		log.Fatal("error while subscribing zip queue", err)
	}
	defer subscription.Unsubscribe()

	time.Sleep(1 * time.Minute)
}

func processMsg(msg *nats.Msg) {
	payload := &Payload{}
	_ = json.Unmarshal(msg.Data, payload)
	replyMsg := fmt.Sprintf("received payload %d", payload.Count)
	msg.Respond([]byte(replyMsg))
	fmt.Println("[info] received message: ", string(msg.Data))
}
