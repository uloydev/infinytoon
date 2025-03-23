package packets

import "encoding/json"

type EchoPayload struct {
	*BasePayload
	Message string `json:"message"`
}

type EchoRq struct {
	EchoPayload
}

type EchoRs struct {
	EchoPayload
}

func NewEchoRq(clientID string) Payload {
	return &EchoRq{
		EchoPayload: EchoPayload{
			BasePayload: &BasePayload{
				ClientID: clientID,
			},
			Message: "echo request",
		},
	}
}

func NewEchoRs(clientID string) Payload {
	return &EchoRs{
		EchoPayload: EchoPayload{
			BasePayload: &BasePayload{
				ClientID: clientID,
			},
			Message: "echo response",
		},
	}
}

func (rq *EchoRq) Encode() (*Message, error) {
	payload, err := json.Marshal(rq)
	if err != nil {
		return nil, err
	}
	return &Message{
		Type:     EchoRequest,
		ClientID: rq.ClientID,
		Payload:  payload,
	}, nil
}

func (rq *EchoRq) Decode(msg *Message) error {
	if msg.Type != EchoRequest {
		return ErrInvalidMessageType
	}
	err := json.Unmarshal(msg.Payload, rq)
	if err != nil {
		return err
	}
	return nil
}

func (rs *EchoRs) Encode() (*Message, error) {
	payload, err := json.Marshal(rs)
	if err != nil {
		return nil, err
	}
	return &Message{
		Type:     EchoResponse,
		ClientID: rs.ClientID,
		Payload:  payload,
	}, nil
}

func (rs *EchoRs) Decode(msg *Message) error {
	if msg.Type != EchoResponse {
		return ErrInvalidMessageType
	}
	err := json.Unmarshal(msg.Payload, rs)
	if err != nil {
		return err
	}
	return nil
}
