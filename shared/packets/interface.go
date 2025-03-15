package packets

type BasePayload struct {
	ClientID string `json:"-"`
}

type Payload interface {
	Encode() (*Message, error)
	Decode(*Message) error
}
