package controller

import (
	"github.com/gofiber/fiber/v2"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/logger"
	"infinitoon.dev/infinitoon/pkg/rest"
)

type RootController struct {
	appCtx *appctx.AppContext
	log    *logger.Logger
}

func NewRootController(appCtx *appctx.AppContext) IController {
	return &RootController{
		appCtx: appCtx,
		log:    appCtx.Get(appctx.LoggerKey).(*logger.Logger),
	}
}

func (c *RootController) Route() *rest.RestRoute {
	route := rest.NewRestRoute()
	route.SetRoot().Handler(func(router fiber.Router) {
		router.All("/", c.rootHandler)
	})
	return route
}

func (c *RootController) rootHandler(ctx *fiber.Ctx) error {
	return ctx.SendString("Hello, World!")
}
