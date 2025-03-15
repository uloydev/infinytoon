package packets

import "encoding/json"

type ServeRq struct {
	BasePayload
	Host     string `json:"host"`
	Port     string `json:"port"`
	ClientID string `json:"client_id"`
}

func NewServeRq(host, port, clientID string) *ServeRq {
	return &ServeRq{
		Host: host,
		Port: port,
		BasePayload: BasePayload{
			ClientID: clientID,
		},
		ClientID: clientID,
	}
}

func (p *ServeRq) Encode() (*Message, error) {
	payload, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return &Message{
		Type:        ServeRequest,
		ClientID:    p.ClientID,
		Source:      CLIENT_INSTANCE,
		Destination: RELAY_INSTANCE,
		Payload:     payload,
	}, nil
}

func (p *ServeRq) Decode(msg *Message) error {
	if msg.Type != ServeRequest {
		return ErrInvalidMessageType
	}

	if msg.Source != CLIENT_INSTANCE {
		return ErrInvalidMessageSource
	}

	err := json.Unmarshal(msg.Payload, p)
	if err != nil {
		return err
	}
	return nil
}

type ServeRs struct {
	BasePayload
	ClientID string `json:"client_id"`
	Status   bool   `json:"status"`
	URL      string `json:"url"`
}

func NewServeRs(status bool, clientID, url string) *ServeRs {
	return &ServeRs{
		ClientID: clientID,
		Status:   status,
		URL:      url,
		BasePayload: BasePayload{
			ClientID: clientID,
		},
	}
}

func (p *ServeRs) Encode() (*Message, error) {
	payload, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return &Message{
		Type:        ServeResponse,
		ClientID:    p.ClientID,
		Source:      RELAY_INSTANCE,
		Destination: CLIENT_INSTANCE,
		Payload:     payload,
	}, nil
}

func (p *ServeRs) Decode(msg *Message) error {
	if msg.Type != ServeResponse {
		return ErrInvalidMessageType
	}
	if msg.Source != RELAY_INSTANCE {
		return ErrInvalidMessageSource
	}
	err := json.Unmarshal(msg.Payload, p)
	if err != nil {
		return err
	}
	return nil
}

// func NewServePayload(host, port, clientID string) Payload {
// 	return &ServePayload{
// 		Host: host,
// 		Port: port,
// 		BasePayload: BasePayload{
// 			ClientID: clientID,
// 		},
// 	}
// }

// func (p *ServePayload) EncodeRq() (*Message, error) {
// 	payload, err := json.Marshal(p)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Message{
// 		Type:     ServeRequest,
// 		ClientID: p.ClientID,
// 		Payload:  payload,
// 	}, nil
// }

// func (p *ServePayload) EncodeRs() (*Message, error) {
// 	payload, err := json.Marshal(p)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Message{
// 		Type:     ServeResponse,
// 		ClientID: p.ClientID,
// 		Payload:  payload,
// 	}, nil
// }

// func (p *ServePayload) DecodeRq(msg *Message) error {
// 	if msg.Type != ServeRequest {
// 		return ErrInvalidMessageType
// 	}
// 	err := json.Unmarshal(msg.Payload, p)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (p *ServePayload) DecodeRs(msg *Message) error {
// 	if msg.Type != ServeResponse {
// 		return ErrInvalidMessageType
// 	}
// 	err := json.Unmarshal(msg.Payload, p)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
