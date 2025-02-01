package cmd

import (
	"github.com/gofiber/fiber/v2"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/rest"
)

type RestCommand struct {
	appCtx *appctx.AppContext
	app    *fiber.App
	cfg    *RestCommandConfig
}

type RestCommandConfig struct {
	Name        string
	Host        string
	Port        string
	BasePath    string
	Middlewares []fiber.Handler
	Routes      []*rest.RestRoute
}

func NewRestCommand(appCtx *appctx.AppContext, cfg *RestCommandConfig) Command {
	return &RestCommand{
		appCtx: appCtx,
		app:    fiber.New(),
		cfg:    cfg,
	}
}

func (w *RestCommand) Name() string {
	return w.cfg.Name
}

func (w *RestCommand) Run() error {

	for _, middleware := range w.cfg.Middlewares {
		w.app.Use(middleware)
	}

	for _, route := range w.cfg.Routes {
		route.Register(w.cfg.BasePath, w.app)
	}

	return w.app.Listen(w.cfg.Host + ":" + w.cfg.Port)
}

func (w *RestCommand) Shutdown() error {
	return w.app.Shutdown()
}
