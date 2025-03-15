package handler

import (
	"context"
	"encoding/json"

	// "log"

	"github.com/quic-go/quic-go"
	"github.com/rs/zerolog/log"
	"infinitoon.dev/infinitoon/apps/relay/config"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/logger"
	"infinitoon.dev/infinitoon/pkg/quictunnel"
	"infinitoon.dev/infinitoon/shared/packets"
)

func RootHandler(appCtx *appctx.AppContext, conn quic.Connection, stream quic.Stream, encoder *json.Encoder, msg packets.Message) {
	cfg := appCtx.Get(appctx.ConfigKey).(*config.Config)
	log := appCtx.Get(appctx.LoggerKey).(*logger.Logger)
	tun := appCtx.Get(appctx.QuicTunnelKey).(quictunnel.QuicTunnel)
	srv := tun.GetServer(quictunnel.QuicServerKey(cfg.Server.Name))

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

		srv.AddClientSession(msg.ClientID, conn)

		payload := packets.NewServeRs(true, msg.ClientID, "http://localhost:1337/")
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

		rqPayload.Header["Host"] = "localhost:3000"
		rqPayload.Host = "localhost:3000"

		msg, _ := rqPayload.Encode()

		log.Info().Any("payload", rqPayload).Msg("sending http request payload to cli client")
		resp, err := srv.SendMessage(context.Background(), "cli-relay-client", msg)
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
