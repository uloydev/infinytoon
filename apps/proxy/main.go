package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"os"

	"infinitoon.dev/infinitoon/pkg/cmd"
	"infinitoon.dev/infinitoon/pkg/container"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/logger"
)

func main() {
	appCtx := appctx.NewAppContext()
	cfg := InitConfig(appCtx)
	log := logger.NewLogger(appCtx, cfg.Logger)

	ctr := container.NewContainer(appCtx)

	defer ctr.Shutdown()

	ctr.RegisterCommand(cmd.NewQuicCommand(appCtx, cmd.QuicCommandConfig{
		Clients: InitClients(appCtx),
	}))

	if err := ctr.Run(); err != nil {
		log.Fatal().Err(err).Msg("error running command")
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
