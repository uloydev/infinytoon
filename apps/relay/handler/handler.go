package handler

import (
	"encoding/json"
	// "log"

	"github.com/quic-go/quic-go"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/logger"
	"infinitoon.dev/infinitoon/shared/packets"
)

func RootHandler(appCtx *appctx.AppContext, stream quic.Stream, encoder *json.Encoder, msg packets.Message) {
	log := appCtx.Get(appctx.LoggerKey).(*logger.Logger)

	log.Info().Any("payload", msg).Msg("message received")
	if msg.Type == packets.EchoRequest {
		resp := packets.Message{
			Type:    packets.EchoResponse,
			Payload: msg.Payload,
		}
		log.Info().Any("payload", resp).Msg("sending response")
		if err := encoder.Encode(resp); err != nil {
			log.Error().Err(err).Msg("error sending response")
		}
	}
}
