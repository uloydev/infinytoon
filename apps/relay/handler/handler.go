package handler

import (
	"context"
	"encoding/json"
	"time"

	// "log"

	"github.com/quic-go/quic-go"
	"github.com/rs/zerolog/log"
	"infinitoon.dev/infinitoon/apps/relay/config"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/database"
	"infinitoon.dev/infinitoon/pkg/logger"
	"infinitoon.dev/infinitoon/pkg/quictunnel"
	"infinitoon.dev/infinitoon/shared/packets"
	"infinitoon.dev/infinitoon/shared/schema"
)

func RootHandler(appCtx *appctx.AppContext, conn quic.Connection, stream quic.Stream, encoder *json.Encoder, msg packets.Message) {
	cfg := appCtx.Get(appctx.ConfigKey).(*config.Config)
	log := appCtx.Get(appctx.LoggerKey).(*logger.Logger)
	tun := appCtx.Get(appctx.QuicTunnelKey).(quictunnel.QuicTunnel)
	srv := tun.GetServer(quictunnel.QuicServerKey(cfg.Server.Name))
	kv := appCtx.Get(appctx.KVClientKey).(*database.KVClient)

	log.Debug().Any("payload", msg).Msg("message received")
	switch msg.Type {
	case packets.EchoRequest:
		payload := packets.NewEchoRs(msg.ClientID)
		resp, err := payload.Encode()
		if err != nil {
			log.Error().Err(err).Any("client", msg.ClientID).Msg("error encoding echo response")
			sendInvalidPayloadResponse(encoder, &msg)
			return
		}
		sendSuccessResponse(encoder, resp)
	case packets.ServeRequest:
		rqPayload := packets.ServeRq{}
		if err := rqPayload.Decode(&msg); err != nil {
			log.Error().Err(err).Msg("error decoding serve request")
			sendInvalidPayloadResponse(encoder, &msg)
			return
		}

		serveData := &schema.KVServeData{
			ClientID:    msg.ClientID,
			BaseURL:     "http://localhost:1337/" + msg.ClientID,
			ConnectedAt: time.Now().Unix(),
			Host:        rqPayload.Host,
			Port:        rqPayload.Port,
			Protocol:    packets.ProtocolHTTP,
			Status:      schema.KVServeStatusConnected,
			Stats:       schema.KVServeStats{},
		}

		log.Info().Any("serve_data", serveData).Msg("setting serve data")

		srv.AddClientSession(msg.ClientID, conn)

		kvRes := kv.Client().Set(context.Background(), serveData.BaseURL, serveData, 0)
		if kvRes.Err() != nil {
			log.Error().Err(kvRes.Err()).Any("client", msg.ClientID).Msg("error setting kv data")
			sendInvalidPayloadResponse(encoder, &msg)
			return
		}

		payload := packets.NewServeRs(true, msg.ClientID, serveData.BaseURL)
		resp, err := payload.Encode()
		if err != nil {
			log.Error().Err(err).Any("client", msg.ClientID).Msg("error encoding serve response")
			sendInvalidPayloadResponse(encoder, &msg)
			return
		}
		sendSuccessResponse(encoder, resp)

	case packets.HttpRequest:
		rqPayload := packets.HttpRq{}
		if err := rqPayload.Decode(&msg); err != nil {
			log.Error().Err(err).Msg("error decoding http request")
			sendInvalidPayloadResponse(encoder, &msg)
			return
		}

		srvData := &schema.KVServeData{}

		status := kv.Client().Get(context.Background(), rqPayload.BaseURL)
		if status.Err() != nil {
			log.Error().Err(status.Err()).Msg("error getting serve data")
			sendInvalidPayloadResponse(encoder, &msg)
			return
		}

		if err := status.Scan(srvData); err != nil {
			log.Error().Err(err).Msg("error scanning serve data")
			sendInvalidPayloadResponse(encoder, &msg)
			return
		}

		rqPayload.Host = srvData.Host + ":" + srvData.Port
		rqPayload.Header["Host"] = rqPayload.Host

		msg, _ := rqPayload.Encode()

		log.Info().Any("payload", rqPayload).Msg("sending http request payload to cli client")
		resp, err := srv.SendMessage(context.Background(), srvData.ClientID, msg)
		if err != nil {
			log.Error().Err(err).Msg("error sending message to cli-relay-client")
			sendInvalidPayloadResponse(encoder, msg)
			return
		}

		if resp.Type != packets.HttpResponse {
			log.Error().Msg("invalid response type")
			sendInvalidPayloadResponse(encoder, msg)
			return
		}

		sendSuccessResponse(encoder, resp)
	}
}

func sendInvalidPayloadResponse(encoder *json.Encoder, msg *packets.Message) {
	msg.Type = packets.ErrInvalidPayload
	if err := encoder.Encode(msg); err != nil {
		log.Error().Err(err).Msg("error sending invalid payload response")
	}
}

func sendSuccessResponse(encoder *json.Encoder, msg *packets.Message) {
	if err := encoder.Encode(msg); err != nil {
		log.Error().Err(err).Msg("error sending success response")
	}
}
