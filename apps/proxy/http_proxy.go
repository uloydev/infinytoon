package main

import (
	"infinitoon.dev/infinitoon/apps/proxy/controller"
	"infinitoon.dev/infinitoon/pkg/cmd"
	appctx "infinitoon.dev/infinitoon/pkg/context"
)

func InitHttpProxy(appCtx *appctx.AppContext) cmd.Command {
	cfg := appCtx.Get(appctx.ConfigKey).(*Config)

	// register middleware here
	// cfg.HttpProxy.Middlewares = append(cfg.HttpProxy.Middlewares, cmd.NewLoggerMiddleware(appCtx, &cfg.Logger))

	// register controller here
	ctrl := controller.InitController(appCtx)
	cfg.HttpProxy.Routes = ctrl.Routes()

	return cmd.NewRestCommand(appCtx, &cfg.HttpProxy)
}
