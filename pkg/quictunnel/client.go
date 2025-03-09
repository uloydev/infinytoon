package quictunnel

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	"net"
	"sync"

	"github.com/quic-go/quic-go"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	packets "infinitoon.dev/infinitoon/shared/packet"
)

type StreamHandler func(quic.Stream, *json.Encoder, packets.Message)
type ConnHandler func(quic.Connection)

type QuicClientKey string
type QuicServerKey string

type QuicTunnel interface {
	AddClient(ctx context.Context, key QuicClientKey, client QuicClient)
	AddServer(ctx context.Context, key QuicServerKey, server QuicServer)
	GetClient(key QuicClientKey) QuicClient
	GetServer(key QuicServerKey) QuicServer
}

type quicTunnel struct {
	clients map[QuicClientKey]QuicClient
	servers map[QuicServerKey]QuicServer
}

type QuicClient interface {
	Setup(ctx context.Context) error
	SendMessage(ctx context.Context, msg *packets.Message) (*packets.Message, error)
	Stream(context.Context, StreamHandler)
	ShutdownClient(ctx context.Context) error
}

type QuicServer interface {
	StartServer(ctx context.Context, handler StreamHandler) error
	SendMessage(ctx context.Context, connKey string, msg *packets.Message) (*packets.Message, error)
	ShutdownServer(ctx context.Context) error
}

func NewQuicTunnel() QuicTunnel {
	return &quicTunnel{
		clients: make(map[QuicClientKey]QuicClient),
		servers: make(map[QuicServerKey]QuicServer),
	}
}

func (qt *quicTunnel) AddClient(ctx context.Context, key QuicClientKey, client QuicClient) {
	qt.clients[key] = client
}

func (qt *quicTunnel) AddServer(ctx context.Context, key QuicServerKey, server QuicServer) {
	qt.servers[key] = server
}

func (qt *quicTunnel) GetClient(key QuicClientKey) QuicClient {
	return qt.clients[key]
}

func (qt *quicTunnel) GetServer(key QuicServerKey) QuicServer {
	return qt.servers[key]
}

type QuicClientConfig struct {
	IP         string
	Port       int
	TLSConfing *tls.Config
}

type QuicServerConfig struct {
	IP         string
	Port       int
	TLSConfing *tls.Config
}

type quicClient struct {
	cfg      QuicClientConfig
	udpConn  *net.UDPConn
	quicConn quic.Connection
	stream   quic.Stream
	encoder  *json.Encoder
	decoder  *json.Decoder
}

func NewQuicClient(appCtx *appctx.AppContext, cfg QuicClientConfig) QuicClient {
	return &quicClient{
		cfg: cfg,
	}
}

func (qc *quicClient) Setup(ctx context.Context) error {
	addr := &net.UDPAddr{
		IP:   net.ParseIP(qc.cfg.IP),
		Port: qc.cfg.Port,
	}
	udpConn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}

	tr := &quic.Transport{
		Conn: udpConn,
	}

	quicConn, err := tr.Dial(ctx, addr, qc.cfg.TLSConfing, &quic.Config{})
	if err != nil {
		return err
	}

	quicConn.CloseWithError(0, "close normal")

	qc.udpConn = udpConn
	qc.quicConn = quicConn

	stream, err := quicConn.OpenStream()
	if err != nil {
		return err
	}

	qc.stream = stream

	qc.encoder = json.NewEncoder(stream)
	qc.decoder = json.NewDecoder(stream)

	return nil

}

func (qc *quicClient) SendMessage(ctx context.Context, msg *packets.Message) (*packets.Message, error) {
	return sendMessage(ctx, qc.encoder, qc.decoder, msg)
}

func (qc *quicClient) Stream(ctx context.Context, handler StreamHandler) {
	for {
		stream, err := qc.quicConn.AcceptStream(ctx)
		if err != nil {
			handleStreamError(stream, err)
		}

		encoder := json.NewEncoder(stream)
		decoder := json.NewDecoder(stream)

		var (
			req, res packets.Message
		)
		if err := decoder.Decode(&req); err != nil {
			log.Println(err)
			res.Type = packets.ErrInvalidPayload
			res.ClientID = stream.StreamID().InitiatedBy().String()
			encoder.Encode(res)
		}
		go handler(stream, encoder, req)
	}
}

func (qc *quicClient) ShutdownClient(ctx context.Context) error {
	if err := qc.stream.Close(); err != nil {
		return err
	}
	return qc.udpConn.Close()
}

type quicServer struct {
	cfg          QuicServerConfig
	udpConn      *net.UDPConn
	quicListener *quic.Listener
	Clients      *sync.Map
	handler      StreamHandler
}

func NewQuicServer(appCtx *appctx.AppContext, cfg QuicServerConfig) QuicServer {
	return &quicServer{
		cfg:     cfg,
		Clients: &sync.Map{},
	}
}

func (qs *quicServer) SendMessage(ctx context.Context, connKey string, msg *packets.Message) (*packets.Message, error) {

	val, ok := qs.Clients.Load(connKey)
	if !ok {
		return nil, errors.New("client connection not found")
	}

	conn, ok := val.(quic.Connection)
	if !ok {
		return nil, errors.New("invalid client connection")
	}

	stream, err := conn.OpenStream()
	if err != nil {
		return nil, err
	}

	return sendMessage(ctx, json.NewEncoder(stream), json.NewDecoder(stream), msg)
}

func (qs *quicServer) StartServer(ctx context.Context, handler StreamHandler) error {
	addr := &net.UDPAddr{
		IP:   net.ParseIP(qs.cfg.IP),
		Port: qs.cfg.Port,
	}

	udpConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	qs.udpConn = udpConn

	quicLs, err := quic.Listen(udpConn, qs.cfg.TLSConfing, &quic.Config{})
	if err != nil {
		return err
	}
	qs.quicListener = quicLs

	log.Printf("Listening on %v", quicLs.Addr())

	qs.handler = handler

	for {
		clientConn, err := qs.quicListener.Accept(ctx)
		if err != nil {
			handleConnError(clientConn, err)
			continue
		}

		qs.Clients.Store(clientConn.RemoteAddr().String(), clientConn)

		go qs.connHandler(clientConn)
	}
}

func (qs *quicServer) connHandler(conn quic.Connection) {
	log.Printf("New session: %v", conn.RemoteAddr())

	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			panic(err)
		}

		decoder := json.NewDecoder(stream)
		encoder := json.NewEncoder(stream)

		var (
			req, res packets.Message
		)
		if err := decoder.Decode(&req); err != nil {
			log.Println(err)
			res.Type = packets.ErrInvalidPayload
			res.ClientID = stream.StreamID().InitiatedBy().String()
			encoder.Encode(res)
		}

		go qs.handler(stream, encoder, req)
	}
}

func (qs *quicServer) ShutdownServer(ctx context.Context) error {
	return qs.udpConn.Close()
}
