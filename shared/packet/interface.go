package packets

type BasePayload struct {
	ClientID string `json:"-"`
}

type Payload interface {
	EncodeRq() (*Message, error)
	EncodeRs() (*Message, error)
	DecodeRq(*Message) error
	DecodeRs(*Message) error
}
