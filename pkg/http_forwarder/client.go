package http_forwarder

import (
	"net/http"

	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/logger"
	"infinitoon.dev/infinitoon/shared/packets"
	"resty.dev/v3"
)

type HttpForwarder interface {
	Do(*packets.HttpRq) (*packets.HttpRs, error)
}

type httpForwarder struct {
	appCtx *appctx.AppContext
	log    *logger.Logger
	client *resty.Client
}

func InitHttpForwarder(appCtx *appctx.AppContext) HttpForwarder {
	c := &httpForwarder{
		appCtx: appCtx,
		log:    appCtx.Get(appctx.LoggerKey).(*logger.Logger),
		client: resty.New(),
	}

	appCtx.Set(appctx.HttpForwarderKey, c)

	return c
}

func (c *httpForwarder) Do(payload *packets.HttpRq) (res *packets.HttpRs, err error) {
	url := packets.GetURLHttpPayload(&payload.BaseHttpPayload)
	req := c.client.R().
		SetMethod(payload.Method).
		SetHeaders(payload.Header).
		SetURL(url).
		SetBody(payload)

	c.log.Info().Any("payload", payload).Msg("sending http request to client app")
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	res = &packets.HttpRs{}
	res.ClientID = payload.ClientID
	res.Protocol = payload.Protocol
	res.Method = payload.Method
	res.Host = payload.Host
	res.Path = payload.Path
	res.Queries = payload.Queries
	res.Body = resp.Bytes()
	res.Header = c.headerToMap(resp.Header())
	res.StatusCode = resp.StatusCode()
	return res, nil
}

func (c *httpForwarder) headerToMap(header http.Header) map[string]string {
	headerMap := make(map[string]string)
	for key, value := range header {
		headerMap[key] = value[0]
	}
	return headerMap
}
