package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io"
	"log"
	"math/big"
	"net"

	"github.com/quic-go/quic-go"
	packets "infinitoon.dev/infinitoon/shared/packet"
)

func handleSession(sess quic.Connection) {
	log.Printf("New session: %v", sess.RemoteAddr())

	stream, err := sess.AcceptStream(context.Background())
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(stream)
	encoder := json.NewEncoder(stream)

	for {
		// check if client has closed the session
		select {
		case <-sess.Context().Done():
			log.Printf("Session closed: %v", sess.RemoteAddr())
			return
		default:
		}

		var (
			req packets.Message
			res *packets.Message
		)
		if err := decoder.Decode(&req); err != nil {
			handleApplicationError(sess, err)
			return
		}

		log.Printf("[%v] Received message: %+v", sess.RemoteAddr(), req)

		switch req.Type {
		case packets.EchoRequest:
			var echoReq packets.EchoPayload
			if err := echoReq.DecodeRq(&req); err != nil {
				handleApplicationError(sess, err)
				return
			}
			log.Printf("[%v] Received echo payload: %v", sess.RemoteAddr(), string(req.Payload))

			res, err = echoReq.EncodeRs()
			if err != nil {
				handleApplicationError(sess, err)
				return
			}
			if err := encoder.Encode(res); err != nil {
				handleApplicationError(sess, err)
				return
			}
			log.Printf("[%v] Sent echo response: %v", sess.RemoteAddr(), string(req.Payload))
		default:
		}

		log.Printf("[%v] Sent message: %+v", sess.RemoteAddr(), req)
	}
}

func handleApplicationError(sess quic.Connection, err error) {

	switch err := err.(type) {
	case *quic.ApplicationError:
		if err.ErrorCode == 0x0 {
			log.Printf("[%v] Client closed session", sess.RemoteAddr())
			return
		}

	case *quic.StreamError:
		if err.ErrorCode == 0x0 {
			log.Printf("[%v] Client closed stream", sess.RemoteAddr())
			return
		}

	default:
		if err == io.EOF {
			log.Printf("[%v] Client closed stream", sess.RemoteAddr())
			return
		}

		log.Printf("[%v] Client Error: %v", sess.RemoteAddr(), err)
	}
}

func main() {
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{
		Port: 12345,
		IP:   net.IPv4(127, 0, 0, 1),
	})
	if err != nil {
		panic(err)
	}

	tr := quic.Transport{
		Conn: udpConn,
	}

	quicConn, err := tr.Listen(generateTLSConfig(), &quic.Config{
		MaxIncomingStreams: 1000,
	})
	if err != nil {
		panic(err)
	}

	log.Printf("Listening on %v", quicConn.Addr())

	for {
		sess, err := quicConn.Accept(context.Background())
		if err != nil {
			handleApplicationError(sess, err)
			panic(err)
		}

		go handleSession(sess)
	}
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
