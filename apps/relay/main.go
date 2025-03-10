package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"log"
	"math/big"
	"os"

	"github.com/quic-go/quic-go"
	"infinitoon.dev/infinitoon/pkg/cmd"
	"infinitoon.dev/infinitoon/pkg/container"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/quictunnel"
	packets "infinitoon.dev/infinitoon/shared/packet"
)

var serverSessionHandler quictunnel.StreamHandler = func(q quic.Stream, e *json.Encoder, m packets.Message) {
	log.Printf("Received message: %v", m)
	if m.Type == packets.EchoRequest {
		resp := packets.Message{
			Type:    packets.EchoResponse,
			Payload: m.Payload,
		}
		log.Printf("Sending response: %v", resp)
		if err := e.Encode(resp); err != nil {
			log.Printf("Error sending response: %v", err)
		}
	}

}

func main() {
	appCtx := appctx.NewAppContext()

	ctr := container.NewContainer(appCtx)

	defer ctr.Shutdown()

	ctr.RegisterCommand(cmd.NewQuicCommand(appCtx, cmd.QuicCommandConfig{
		Servers: []quictunnel.QuicServer{
			quictunnel.NewQuicServer(appCtx, quictunnel.QuicServerConfig{
				Name:       "relay server",
				IP:         "127.0.0.1",
				Port:       54321,
				TLSConfing: generateTLSConfig(),
			}, serverSessionHandler),
		},
	}))

	if err := ctr.Run(); err != nil {
		log.Fatalf("Error running command: %v", err)
	}

	signalChan := make(chan os.Signal, 1)
	<-signalChan
}

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}
