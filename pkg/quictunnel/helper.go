package quictunnel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/quic-go/quic-go"
	"infinitoon.dev/infinitoon/pkg/logger"
	"infinitoon.dev/infinitoon/shared/packets"
)

func sendMessage(
	ctx context.Context,
	encoder *json.Encoder,
	decoder *json.Decoder,
	msg *packets.Message,
) (*packets.Message, error) {
	var response packets.Message
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// if encoder is nil, sleep for a while before retrying
			if encoder == nil {
				fmt.Println("retrying send message: encoder is nil")
				time.Sleep(1 * time.Second)
				continue
			}

			if err := encoder.Encode(msg); err != nil {
				fmt.Printf("retrying send message: %v\n", err)
				// sleep for a while before retrying
				time.Sleep(1 * time.Second)
				continue
			}

			if err := decoder.Decode(&response); err != nil {
				return nil, errors.Join(err, errors.New("error decoding response"))
			}

			return &response, nil
		}
	}
}

func handleConnError(sess quic.Connection, err error) {

	switch err := err.(type) {
	case *quic.ApplicationError:
		if err.ErrorCode == 0x0 {
			log.Printf("[%v] Client closed session", sess.RemoteAddr())
			return
		}

	case *quic.StreamError:
		if err.ErrorCode == 0x0 {
			log.Printf("[%v] Client closed stream", sess.RemoteAddr())
			return
		}

	// handle idle connection
	case *quic.IdleTimeoutError:
		log.Printf("[%v] Client idle timeout", sess.RemoteAddr())
		return

	default:
		if err == io.EOF {
			log.Printf("[%v] Client closed stream", sess.RemoteAddr())
			return
		}

		log.Printf("[%v] Client Error: %v", sess.RemoteAddr(), err)
	}
}

func handleStreamError(log *logger.Logger, connAddr string, err error) {

	switch err := err.(type) {
	case *quic.ApplicationError:
		if err.ErrorCode == 0x0 {
			log.Warn().Err(err).Msgf("[%v] Client closed session", connAddr)
			return
		}

	case *quic.StreamError:
		if err.ErrorCode == 0x0 {
			log.Warn().Err(err).Msgf("[%v] Client closed stream", connAddr)
			return
		}

	default:
		if err == io.EOF {
			log.Warn().Err(err).Msgf("[%v] Client closed stream", connAddr)
			return
		}

		log.Warn().Err(err).Msgf("[%v] Client Error", connAddr)
	}
}
