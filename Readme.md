# NATS: Microservice Connectivity 🚀

## Why NATS? 🤔

### The Problem 🛑
- **Microservices need to communicate** effectively.
- **Applications need to be resilient against failure**.
- **Applications should scale seamlessly**.
- **New services should be added/removed without disruption**.

### Desired Solution 🛠️
- Connect many decoupled application services.
- Make it easy and secure.
- Provide flexibility in messaging patterns.
- Ensure scalability.

## NATS for Microservices 🧩

### Architecture 🏗️
```
s1   s2   s3   s3 (all about messages)
|     |    |    | 
----------------- [bus, common interface]
```

### Features 🌟
- **Location independent addressing** 📍
- **Pub/Sub messaging** 📬
- **Request/Reply pattern** 🔄
- **Streaming** 📽️
- **Persistence** 💾
- **Secure by default** 🔒
- **Global scale** 🌍
- **Multi-tenancy** 🏢
- **Small binary** 🪶
- **Many client libraries** 📚
- **High speed** ⚡
- **High fan-out** 🔧
- **Double subscriber** ➡️➡️
- **Encryption at rest** 🔐
- **Subject limits** 🚧
- **Scalable clustering** 🌐
- **Edge & IoT** 📡
- **Easy configuration** 🛠️
- **Open-source software (OSS)** 🔓

## NATS Server Installation and Setup ⚙️
- **Build in Go** 🐹
- **Types**: Binary, Docker, Kubernetes 🐋
- **Start Server**: `nats-server`

### Installation 🔧
- Visit [NATS Installation](https://nats.io) for details.

#### Run Docker 🐳
```bash
docker run --rm -p 4222:4222 nats:latest
```

#### Binary Install 📦
Refer to [NATS Installation](https://nats.io) for different options.

## How NATS Works ⚙️

- **Service-to-Service Communication**: Clients communicate by topics, not by direct service names.
  ```
  s1 => message[to:s3, from:s1, message:"message"]

  client1 (s1) <----> [NATS] <-----> client2(s3)
  ```

### Message Exchange 📧
- **Subscription**: Clients subscribe to subjects to receive messages.
  ```
  subject      subscribed client 
  foo.bar      s3, s4, s7
  payment.post s8 
  ```

- **No Subscribers**: Messages are dropped if no subscribers are present (can be persisted using configurations like JetStream).

## Messaging Patterns 📨

### Publish/Subscribe 📢
- **Description**: Decouples publishers and subscribers, allowing multiple subscribers to listen to the same topic.

#### Publisher Example 📤
```go
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
            fmt.Println("[error] message failed to publish", err)
        }
        log.Println("[info] message published", msg)
        time.Sleep(1 * time.Second)
    }
}
```
- **Explanation**:
    - Connects to NATS server.
    - Publishes "Hello World!" message to the "intros" subject every second.

#### Subscriber Example 📥
```go
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
```
- **Explanation**:
    - Connects to NATS server.
    - Subscribes to the "intros" subject and prints received messages.

### Request/Reply 🔄
- **Description**: Ensures the message is received and a response is returned.

#### Publisher Example 📤
```go
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
```
- **Explanation**:
    - Connects to NATS server.
    - Publishes a payload with a count to the "intros" subject and waits for a reply.

#### Subscriber Example 📥
```go
func main() {
    nc, err := nats.Connect(nats.DefaultURL)
    if err != nil {
        log.Fatal("error while connecting nats ", err)
    }
    defer nc.Close()

    nc.Subscribe("intros", func(msg *nats.Msg) {
        payload := &Payload{}
        _ = json.Unmarshal(msg.Data, payload)
        replyMsg := fmt.Sprintf("received payload %d", payload.Count)
        msg.Respond([]byte(replyMsg))
        fmt.Println("[info] received message: ", string(msg.Data))
    })

    time.Sleep(1 * time.Minute)
}
```
- **Explanation**:
    - Connects to NATS server.
    - Subscribes to the "intros" subject, processes the message, and sends a reply.

### Queue 📋
- **Description**: Ensures at least one subscriber processes the message, but only one.

#### Subscriber Example 📥
```go
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
```
- **Explanation**:
    - Connects to NATS server.
    - Subscribes to the "intros" subject in the "zip" queue group, ensuring one subscriber processes each message.

## NATS Tools 🔧

### Fan-out Pattern 🌐
- One message is sent and received by multiple receivers.

#### Commands
```bash
nats pub "intros" "hello world {{Count}}" --count=100 --sleep 1s
nats sub "intros"
nats sub "intros" --queue zip
nats request "intros" ""
nats reply "intros" "got your request # {{Count}}"
```

## NATS Subjects 📜

### Subject Naming 📛
- Subjects are strings that form names which publishers and subscribers use to find each other.

#### Recommendations
- **Characters**: (a-z, A-Z, 0-9)
- **Special Characters**: (., >, *)

### Special Characters ✨
- **"." (Dot)**: Creates hierarchy
    - Examples: `weather.us`, `weather.us.east`, `weather.gy.west`
- **"*" (Wildcard)**: Matches a single token
    - Example: `weather.*.east`
- **">" (Wildcard)**: Matches multiple tokens, must appear at the end of the subject
    - Example: `weather.us.east.>`

### Reverse Subjects 🔙
- Reverse subjects start with `$`.

## Security 🔒

## JetStream 🛠️
- **Description**: Built-in distributed persistence system enabling new functionalities and higher qualities of service on top of core NATS functionalities.

### JetStream Capabilities 🚀
- **Streaming**: Temporary decoupling between publishers and subscribers.
    - **Key/Value Store**: Similar to Redis.
    - **Document Store**: Similar to S3.
- **Stream**: Captures and stores messages published to one or more subjects.

#### Stream Illustration 📊
```
s1 ---[pub]---> [(order.new)]--------------------------+-----[sub]----> s4
                                                       |
                                +---------------------------+
                                |   stream             |    |
                                | Name: ORDER          |    |
                                | subject: order.*     |    |
s2---[sub]<---------------------|------[.....]  <------+    |
 |                              |                           |
 +---[pub]---> [(order.filled)]-|----->[.....]  -------+    |
                                +----------------------|----+ 
                                                       |
 s3 <------[sub]---------------------------------------+

```

### Ephemeral/Durable Consumer 📈
1. **Setup Config NATS**: Setup user, password in `~/.config/nats/context/<context>.json`.
2. **Setup Auth with NATS-Server**: Refer to `example/jetstream/js.conf` file.
   ```bash
   nats-server -c js.conf
   ```
3. **Start Publishing**:
   ```bash
   nats pub orders.us --count=1000 --sleep 2s "US order # {{Count}}"
   ```