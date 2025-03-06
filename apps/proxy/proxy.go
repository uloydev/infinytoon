package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
	packets "infinitoon.dev/infinitoon/shared/packet"
)

type ContextKey string

const (
	ClientIDKey ContextKey = "clientId"
)

func startClient(ctx context.Context) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}

	conn, err := quic.DialAddr(context.Background(), "localhost:12345", tlsConf, nil)
	if err != nil {
		panic(err)
	}
	defer conn.CloseWithError(0, "close normal")
	clientId := ctx.Value(ClientIDKey).(string)
	log.Printf("[%v] Starting client", clientId)
	stream, err := conn.OpenStream()
	if err != nil {
		log.Fatalf("[%v] Failed to accept stream: %v", clientId, err)
	}

	log.Printf("[%v] Accepted stream", clientId)
	defer stream.Close()

	msg, err := packets.NewEchoPayload(clientId).EncodeRq()
	if err != nil {
		panic(err)
	}

	for i := 0; i < 3; i++ {
		sendMessage(ctx, stream, *msg)
		time.Sleep(1 * time.Second)
	}

	log.Printf("[%v] Closing stream", clientId)
	ctx.Done()
}

func main() {

	wg := &sync.WaitGroup{}
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := strconv.Itoa(i)
			startClient(context.WithValue(context.Background(), ClientIDKey, "client-"+id))
		}()
		time.Sleep(200 * time.Millisecond)
	}

	wg.Wait()
}

func sendMessage(ctx context.Context, stream quic.Stream, message packets.Message) *packets.Message {
	clientID := ctx.Value(ClientIDKey).(string)
	encoder := json.NewEncoder(stream)
	if err := encoder.Encode(message); err != nil {
		panic(err)
	}

	log.Printf("[%v] Sent payload: %v", clientID, string(message.Payload))
	log.Printf("[%v] Sent message: %+v", clientID, message)

	decoder := json.NewDecoder(stream)
	var res packets.Message
	if err := decoder.Decode(&res); err != nil {
		panic(err)
	}

	log.Printf("[%v] Received message: %+v", clientID, res)
	log.Printf("[%v] Received payload: %v", clientID, string(res.Payload))

	return &res
}
