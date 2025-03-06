package packets

type MessageType string

const (
	EchoRequest  MessageType = "echo_rq"
	EchoResponse MessageType = "echo_rs"

	AuthRequest  MessageType = "auth_rq"
	AuthResponse MessageType = "auth_rs"

	HttpRequest  MessageType = "http_rq"
	HttpResponse MessageType = "http_rs"

	TCPRequest  MessageType = "tcp_rq"
	TCPResponse MessageType = "tcp_rs"

	UDPRequest  MessageType = "udp_rq"
	UDPResponse MessageType = "udp_rs"
)

type Message struct {
	Type     MessageType `json:"type"`
	ClientID string      `json:"client_id"`
	Payload  []byte      `json:"payload"`
}
