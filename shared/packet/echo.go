package packets

import "encoding/json"

type EchoPayload struct {
	BasePayload
	Message string `json:"message"`
}

func NewEchoPayload(clientId string) Payload {
	return &EchoPayload{
		Message: "echo",
		BasePayload: BasePayload{
			ClientID: clientId,
		},
	}
}

func (p *EchoPayload) EncodeRq() (*Message, error) {
	payload, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	return &Message{
		ClientID: p.ClientID,
		Type:     EchoRequest,
		Payload:  payload,
	}, nil
}

func (p *EchoPayload) EncodeRs() (*Message, error) {
	payload, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	return &Message{
		ClientID: p.ClientID,
		Type:     EchoResponse,
		Payload:  payload,
	}, nil
}

func (p *EchoPayload) DecodeRq(m *Message) error {
	p.ClientID = m.ClientID
	if m.Type != EchoRequest {
		return ErrInvalidMessageType
	}

	return json.Unmarshal(m.Payload, p)
}

func (p *EchoPayload) DecodeRs(m *Message) error {
	p.ClientID = m.ClientID
	if m.Type != EchoResponse {
		return ErrInvalidMessageType
	}

	return json.Unmarshal(m.Payload, p)
}
