package main

import (
	"infinitoon.dev/infinitoon/apps/proxy/config"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/quictunnel"
)

func InitClients(appCtx *appctx.AppContext) []quictunnel.QuicClient {
	cfg := appCtx.Get(appctx.ConfigKey).(*config.Config)
	clients := []quictunnel.QuicClient{}

	for _, clientCfg := range cfg.Clients {
		client := quictunnel.NewQuicClient(appCtx, quictunnel.QuicClientConfig{
			Name:       clientCfg.Name,
			IP:         clientCfg.IP,
			Port:       clientCfg.Port,
			TLSConfing: generateTLSConfig(),
		})

		clients = append(clients, client)
	}
	return clients
}
