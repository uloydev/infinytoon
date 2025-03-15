package controller

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"infinitoon.dev/infinitoon/apps/proxy/config"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/logger"
	"infinitoon.dev/infinitoon/pkg/quictunnel"
	"infinitoon.dev/infinitoon/pkg/rest"
	"infinitoon.dev/infinitoon/shared/packets"
)

type RootController struct {
	appCtx *appctx.AppContext
	log    *logger.Logger
	cfg    *config.Config
	tun    quictunnel.QuicTunnel
}

func NewRootController(appCtx *appctx.AppContext, cfg *config.Config) IController {
	return &RootController{
		appCtx: appCtx,
		log:    appCtx.Get(appctx.LoggerKey).(*logger.Logger),
		tun:    appCtx.Get(appctx.QuicTunnelKey).(quictunnel.QuicTunnel),
		cfg:    cfg,
	}
}

func (c *RootController) Route() *rest.RestRoute {
	route := rest.NewRestRoute()
	route.SetRoot().Handler(func(router fiber.Router) {
		router.All("*", c.rootHandler)
	})
	return route
}

func (c *RootController) rootHandler(ctx *fiber.Ctx) error {
	c.log.Info().Str("url", ctx.OriginalURL()).Str("method", ctx.Method()).Msg("http request received")

	bodyRaw := ctx.Body()
	headers := map[string]string{}
	ctx.Request().Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = string(value)
	})
	queryParam := ctx.Queries()

	rq := packets.HttpRq{
		BaseHttpPayload: packets.BaseHttpPayload{
			Protocol: packets.HttpProtocol(ctx.Protocol()),
			Method:   ctx.Method(),
			Host:     ctx.Hostname(),
			Path:     ctx.Path(),
			Queries:  queryParam,
			Body:     bodyRaw,
			Header:   headers,
		},
	}

	msg, err := rq.Encode()
	if err != nil {
		c.log.Error().Err(err).Msg("error encoding http request")
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	relay := c.tun.GetClient("relay-client")
	if relay == nil {
		c.log.Error().Msg("error getting relay client")
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	c.log.Info().Msg("sending http request to relay")
	msg, err = relay.SendMessage(context.Background(), msg)
	if err != nil {
		c.log.Error().Stack().Err(err).Msg("error sending http request")
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	c.log.Info().Any("message", msg).Msg("received http response from relay")
	if msg.Type != packets.HttpResponse {
		c.log.Error().Str("type", string(msg.Type)).Msg("invalid response type")
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	rs := &packets.HttpRs{}

	if err := rs.Decode(msg); err != nil {
		c.log.Error().Err(err).Msg("error decoding http response")
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	for key, value := range rs.Header {
		ctx.Set(key, value)
	}

	c.log.Info().Int("status_code", rs.StatusCode).Msg("sending http response")
	return ctx.Status(rs.StatusCode).Send(rs.Body)
}
