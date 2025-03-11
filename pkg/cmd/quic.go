package cmd

import (
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/quictunnel"
)

type QuicCommand struct {
	appCtx *appctx.AppContext
	app    quictunnel.QuicTunnel
	cfg    QuicCommandConfig
}

type QuicCommandConfig struct {
	Clients []quictunnel.QuicClient
	Servers []quictunnel.QuicServer
}

func NewQuicCommand(appCtx *appctx.AppContext, cfg QuicCommandConfig) Command {
	c := &QuicCommand{
		appCtx: appCtx,
		app:    quictunnel.NewQuicTunnel(),
		cfg:    cfg,
	}

	appCtx.Set(appctx.QuicTunnelKey, c.app)
	return c
}

func (q *QuicCommand) Name() string {
	return "quic tunnel"
}

func (q *QuicCommand) Run() error {
	for _, client := range q.cfg.Clients {
		q.app.AddClient(quictunnel.QuicClientKey(client.Name()), client)
	}

	for _, server := range q.cfg.Servers {
		q.app.AddServer(quictunnel.QuicServerKey(server.Name()), server)
	}

	q.appCtx.Set(appctx.QuicTunnelKey, q.app)

	q.app.Start()
	return nil
}

func (q *QuicCommand) Shutdown() error {
	q.app.Shutdown()
	return nil
}
