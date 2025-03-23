package packets

import "infinitoon.dev/infinitoon/pkg/encryption"

type BasePayload struct {
	ClientID string `json:"-"`
}

func (p *BasePayload) EncodeAndEncrypt(encryption.Encryption) (*Message, error) {
	return nil, nil
}

func (p *BasePayload) DecryptAndDecode(encryption.Encryption, *Message) error {
	return nil
}

type Payload interface {
	Encode() (*Message, error)
	Decode(*Message) error
	EncodeAndEncrypt(encryption.Encryption) (*Message, error)
	DecryptAndDecode(encryption.Encryption, *Message) error
}
