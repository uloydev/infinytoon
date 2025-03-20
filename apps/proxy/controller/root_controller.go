package controller

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"infinitoon.dev/infinitoon/apps/proxy/config"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/database"
	"infinitoon.dev/infinitoon/pkg/logger"
	"infinitoon.dev/infinitoon/pkg/quictunnel"
	"infinitoon.dev/infinitoon/pkg/rest"
	"infinitoon.dev/infinitoon/shared/packets"
	"infinitoon.dev/infinitoon/shared/schema"
)

type RootController struct {
	appCtx *appctx.AppContext
	kv     *database.KVClient
	log    *logger.Logger
	cfg    *config.Config
	tun    quictunnel.QuicTunnel
}

func NewRootController(appCtx *appctx.AppContext, cfg *config.Config, kv *database.KVClient) IController {
	return &RootController{
		appCtx: appCtx,
		kv:     kv,
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

var ErrInvalidPath = errors.New("invalid path")

func getBaseURL(ctx *fiber.Ctx) (string, string, error) {
	path := ctx.Path()
	fmt.Println("path", path)
	if path == "" {
		return "", "", ErrInvalidPath
	}

	pathArr := strings.Split(path, "/")
	if len(pathArr) < 2 {
		return "", "", ErrInvalidPath
	}

	path = "/" + pathArr[1]

	return ctx.BaseURL() + path, "/" + strings.TrimLeft(ctx.Path(), path), nil
}

func (c *RootController) serveDataExists(baseUrl string) bool {
	cmd := c.kv.Client().Exists(context.Background(), baseUrl)
	if cmd.Err() != nil {
		c.log.Error().Err(cmd.Err()).Msg("error checking if serve data exists")
		return false
	}

	return cmd.Val() == 1
}

func (c *RootController) rootHandler(ctx *fiber.Ctx) error {
	c.log.Info().Str("url", ctx.OriginalURL()).Str("method", ctx.Method()).Msg("http request received")

	bodyRaw := ctx.Body()
	headers := map[string]string{}
	ctx.Request().Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = string(value)
	})
	queryParam := ctx.Queries()

	baseURL, path, err := getBaseURL(ctx)
	c.log.Info().Str("base_url", baseURL).Msg("base url")
	if err != nil {
		c.log.Error().Err(err).Msg("error getting base url")
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if !c.serveDataExists(baseURL) {
		c.log.Error().Str("base_url", baseURL).Msg("serve data not found")
		return ctx.Status(fiber.StatusNotFound).SendString("serve data not found")
	}

	srvData := &schema.KVServeData{}
	cmd := c.kv.Client().Get(context.Background(), baseURL)
	if cmd.Err() != nil {
		c.log.Error().Err(cmd.Err()).Msg("error getting serve data")
		return ctx.Status(fiber.StatusInternalServerError).SendString("error getting serve data")
	}

	if err := cmd.Scan(srvData); err != nil {
		c.log.Error().Err(err).Msg("error scanning serve data")
		return ctx.Status(fiber.StatusInternalServerError).SendString("error scanning serve data")
	}

	srvData.Stats.RequestCount++
	srvData.Stats.BytesReceived += int64(len(bodyRaw))

	status := c.kv.Client().Set(context.Background(), baseURL, srvData, 0)
	if status.Err() != nil {
		c.log.Error().Err(status.Err()).Msg("error setting serve data")
		return ctx.Status(fiber.StatusInternalServerError).SendString("error setting serve data")
	}

	rq := packets.HttpRq{
		BaseURL: baseURL,
		BaseHttpPayload: packets.BaseHttpPayload{
			Protocol: packets.HttpProtocol(ctx.Protocol()),
			Method:   ctx.Method(),
			Host:     ctx.Hostname(),
			Path:     path,
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

	cmd = c.kv.Client().Get(context.Background(), baseURL)
	if cmd.Err() != nil {
		c.log.Error().Err(cmd.Err()).Msg("error getting serve data")
		return ctx.Status(fiber.StatusInternalServerError).SendString("error getting serve data")
	}

	if err := cmd.Scan(srvData); err != nil {
		c.log.Error().Err(err).Msg("error scanning serve data")
		return ctx.Status(fiber.StatusInternalServerError).SendString("error scanning serve data")
	}

	srvData.Stats.BytesSent += int64(len(rs.Body))
	status = c.kv.Client().Set(context.Background(), baseURL, srvData, 0)
	if status.Err() != nil {
		c.log.Error().Err(status.Err()).Msg("error setting serve data")
		return ctx.Status(fiber.StatusInternalServerError).SendString("error setting serve data")
	}

	c.log.Info().Int("status_code", rs.StatusCode).Msg("sending http response")
	return ctx.Status(rs.StatusCode).Send(rs.Body)
}
