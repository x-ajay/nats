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
		log.Fatal("[error] while connecting nats ", err)
	}
	defer nc.Close()

	count := 0
	timeout := time.Millisecond * 500

	payload := &Payload{
		Data: "hello world",
	}
	for {
		payload.Count = count
		msg, _ := json.Marshal(payload)
		reply, err := nc.Request("intros", msg, timeout)
		time.Sleep(1 * time.Second)
		if err != nil {
			fmt.Println("[error] failed to send message=", payload.Data, "count=", payload.Count, "err=", err)
			continue
		}
		count++
		log.Println("[info] message=", payload.Data, "count=", payload.Count, "reply=", string(reply.Data))
	}

}
