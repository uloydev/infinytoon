package schema

import (
	"encoding/json"

	"infinitoon.dev/infinitoon/shared/packets"
)

type KVServeStatus string

const (
	KVServeStatusConnected    KVServeStatus = "connected"
	KVServeStatusDisconnected KVServeStatus = "disconnected"
)

type KVServeData struct {
	ClientID       string               `json:"client_id"`
	ConnectedAt    int64                `json:"connected_at"`
	DisconnectedAt *int64               `json:"disconnected_at"`
	BaseURL        string               `json:"base_url"`
	Host           string               `json:"host"`
	Port           string               `json:"port"`
	Protocol       packets.HttpProtocol `json:"protocol"`
	Status         KVServeStatus        `json:"status"`
	Stats          KVServeStats         `json:"stats"`
}

type KVServeStats struct {
	BytesReceived int64 `json:"bytes_received"`
	BytesSent     int64 `json:"bytes_sent"`
	RequestCount  int64 `json:"request_count"`
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (s *KVServeData) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (s *KVServeData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}
