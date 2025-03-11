package quictunnel

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/quic-go/quic-go"
	"infinitoon.dev/infinitoon/shared/packets"
)

func sendMessage(
	ctx context.Context,
	encoder *json.Encoder,
	decoder *json.Decoder,
	msg *packets.Message,
) (*packets.Message, error) {
	if err := encoder.Encode(msg); err != nil {
		return nil, err
	}

	var response packets.Message
	if err := decoder.Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
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

	default:
		if err == io.EOF {
			log.Printf("[%v] Client closed stream", sess.RemoteAddr())
			return
		}

		log.Printf("[%v] Client Error: %v", sess.RemoteAddr(), err)
	}
}

func handleStreamError(sess quic.Stream, err error) {

	switch err := err.(type) {
	case *quic.ApplicationError:
		if err.ErrorCode == 0x0 {
			log.Printf("[%v] Client closed session", sess.StreamID())
			return
		}

	case *quic.StreamError:
		if err.ErrorCode == 0x0 {
			log.Printf("[%v] Client closed stream", sess.StreamID())
			return
		}

	default:
		if err == io.EOF {
			log.Printf("[%v] Client closed stream", sess.StreamID())
			return
		}

		log.Printf("[%v] Client Error: %v", sess.StreamID(), err)
	}
}
