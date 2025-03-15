package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"

	"infinitoon.dev/infinitoon/apps/cli/cmd"
	"infinitoon.dev/infinitoon/apps/cli/config"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/http_forwarder"
	"infinitoon.dev/infinitoon/pkg/logger"
)

func main() {

	appCtx := appctx.NewAppContext()
	cfg := config.InitConfig(appCtx)
	cfg.TunnelClient.TLSConfing = generateTLSConfig()
	logger.NewLogger(appCtx, cfg.Logger)
	http_forwarder.InitHttpForwarder(appCtx)

	cmd.Execute(appCtx)
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
		Certificates:       []tls.Certificate{tlsCert},
		NextProtos:         []string{"quic-echo-example"},
		InsecureSkipVerify: true,
	}
}
