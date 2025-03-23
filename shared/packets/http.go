package packets

import (
	"encoding/json"

	"infinitoon.dev/infinitoon/pkg/encryption"
)

type HttpProtocol string

const (
	ProtocolHTTP  HttpProtocol = "http"
	ProtocolHTTPS HttpProtocol = "https"
)

type BaseHttpPayload struct {
	BasePayload
	Protocol HttpProtocol      `json:"protocol"`
	Method   string            `json:"method"`
	Host     string            `json:"host"`
	Path     string            `json:"path"`
	Queries  map[string]string `json:"queries"`
	Body     []byte            `json:"body"`
	Header   map[string]string `json:"header"`
}

type HttpRq struct {
	BaseURL string `json:"base_url"`
	BaseHttpPayload
}

type HttpRs struct {
	BaseHttpPayload
	StatusCode int `json:"status_code"`
}

func GetURLHttpPayload(payload *BaseHttpPayload) string {
	queries := ""
	for key, value := range payload.Queries {
		queries += key + "=" + value + "&"
	}
	return string(payload.Protocol) + "://" + payload.Host + payload.Path + "?" + queries
}

func (rq *HttpRq) Encode() (*Message, error) {
	payload, err := json.Marshal(rq)
	if err != nil {
		return nil, err
	}
	return &Message{
		Type:     HttpRequest,
		ClientID: rq.ClientID,
		Payload:  payload,
	}, nil
}

func (rq *HttpRq) Decode(msg *Message) error {
	if msg.Type != HttpRequest {
		return ErrInvalidMessageType
	}
	err := json.Unmarshal(msg.Payload, rq)
	if err != nil {
		return err
	}
	return nil
}

func (rq *HttpRq) EncodeAndEncrypt(enc encryption.Encryption) (*Message, error) {
	payload, err := json.Marshal(rq)
	if err != nil {
		return nil, err
	}

	encryptedPayload, err := enc.Encrypt(payload)
	if err != nil {
		return nil, err
	}

	return &Message{
		Type:     HttpRequest,
		ClientID: rq.ClientID,
		Payload:  encryptedPayload,
	}, nil
}

func (rq *HttpRq) DecryptAndDecode(enc encryption.Encryption, msg *Message) error {
	if msg.Type != HttpRequest {
		return ErrInvalidMessageType
	}

	decryptedPayload, err := enc.Decrypt(msg.Payload)
	if err != nil {
		return err
	}

	err = json.Unmarshal(decryptedPayload, rq)
	if err != nil {
		return err
	}

	rq.ClientID = msg.ClientID

	return nil
}

func (rs *HttpRs) Encode() (*Message, error) {
	payload, err := json.Marshal(rs)
	if err != nil {
		return nil, err
	}
	return &Message{
		Type:     HttpResponse,
		ClientID: rs.ClientID,
		Payload:  payload,
	}, nil
}

func (rs *HttpRs) Decode(msg *Message) error {
	if msg.Type != HttpResponse {
		return ErrInvalidMessageType
	}

	err := json.Unmarshal(msg.Payload, rs)
	if err != nil {
		return err
	}

	rs.ClientID = msg.ClientID

	return nil
}

func (rs *HttpRs) EncodeAndEncrypt(enc encryption.Encryption) (*Message, error) {
	payload, err := json.Marshal(rs)
	if err != nil {
		return nil, err
	}

	encryptedPayload, err := enc.Encrypt(payload)
	if err != nil {
		return nil, err
	}

	return &Message{
		Type:     HttpResponse,
		ClientID: rs.ClientID,
		Payload:  encryptedPayload,
	}, nil
}

func (rs *HttpRs) DecryptAndDecode(enc encryption.Encryption, msg *Message) error {
	if msg.Type != HttpResponse {
		return ErrInvalidMessageType
	}

	decryptedPayload, err := enc.Decrypt(msg.Payload)
	if err != nil {
		return err
	}

	err = json.Unmarshal(decryptedPayload, rs)
	if err != nil {
		return err
	}

	rs.ClientID = msg.ClientID

	return nil
}
