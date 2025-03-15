package packets

import "errors"

var (
	ErrInvalidMessageType        = errors.New("invalid message type")
	ErrInvalidMessageSource      = errors.New("invalid message source")
	ErrInvalidMessageDestination = errors.New("invalid message destination")
)
