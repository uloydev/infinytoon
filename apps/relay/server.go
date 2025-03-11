package main

import (
	"infinitoon.dev/infinitoon/apps/relay/handler"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/quictunnel"
)

func InitServers(appCtx *appctx.AppContext) []quictunnel.QuicServer {
	cfg := appCtx.Get(appctx.ConfigKey).(*Config)
	servers := []quictunnel.QuicServer{
		quictunnel.NewQuicServer(appCtx, quictunnel.QuicServerConfig{
			Name:       cfg.Server.Name,
			IP:         cfg.Server.IP,
			Port:       cfg.Server.Port,
			TLSConfing: generateTLSConfig(),
		}, handler.RootHandler),
	}
	return servers
}
