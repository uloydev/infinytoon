package quictunnel

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
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
	AddClient(key QuicClientKey, client QuicClient)
	AddServer(key QuicServerKey, server QuicServer)
	GetClient(key QuicClientKey) QuicClient
	GetServer(key QuicServerKey) QuicServer
	Start()
	Shutdown()
}

type quicTunnel struct {
	clients map[QuicClientKey]QuicClient
	servers map[QuicServerKey]QuicServer
}

type QuicClient interface {
	Name() string
	Setup(ctx context.Context) error
	SendMessage(ctx context.Context, msg *packets.Message) (*packets.Message, error)
	Stream(context.Context, StreamHandler)
	ShutdownClient(ctx context.Context) error
}

type QuicServer interface {
	Name() string
	StartServer(ctx context.Context) error
	SendMessage(ctx context.Context, connKey string, msg *packets.Message) (*packets.Message, error)
	ShutdownServer(ctx context.Context) error
}

func NewQuicTunnel() QuicTunnel {
	return &quicTunnel{
		clients: make(map[QuicClientKey]QuicClient),
		servers: make(map[QuicServerKey]QuicServer),
	}
}

func (qt *quicTunnel) AddClient(key QuicClientKey, client QuicClient) {
	qt.clients[key] = client
}

func (qt *quicTunnel) AddServer(key QuicServerKey, server QuicServer) {
	qt.servers[key] = server
}

func (qt *quicTunnel) GetClient(key QuicClientKey) QuicClient {
	return qt.clients[key]
}

func (qt *quicTunnel) GetServer(key QuicServerKey) QuicServer {
	return qt.servers[key]
}

func (qt *quicTunnel) Start() {
	for k, client := range qt.clients {
		go func() {
			log.Printf("starting quic client %s\n", k)
			if err := client.Setup(context.Background()); err != nil {
				log.Println(errors.Join(fmt.Errorf("failed to start quic client %s", k), err))
			}
			log.Printf("quic client %s started\n", k)
		}()
	}

	for k, server := range qt.servers {
		go func() {
			log.Printf("starting quic server %s\n", k)
			if err := server.StartServer(context.Background()); err != nil {
				log.Println(errors.Join(fmt.Errorf("failed to start quic server %s", k), err))
			}
		}()
	}
}

func (qt *quicTunnel) Shutdown() {
	for k, client := range qt.clients {
		log.Printf("shutting down quic client %s\n", k)
		if err := client.ShutdownClient(context.Background()); err != nil {
			log.Println(errors.Join(fmt.Errorf("failed to shutdown quic client %s", k), err))
			continue
		}
		log.Printf("quic client %s shutdown\n", k)
	}

	for k, server := range qt.servers {
		log.Printf("shutting down quic server %s\n", k)
		if err := server.ShutdownServer(context.Background()); err != nil {
			log.Println(errors.Join(fmt.Errorf("failed to shutdown quic server %s", k), err))
			continue
		}
		log.Printf("quic server %s shutdown\n", k)
	}
}

type QuicClientConfig struct {
	Name       string
	IP         string
	Port       int
	TLSConfing *tls.Config
}

type QuicServerConfig struct {
	Name       string
	IP         string
	Port       int
	TLSConfing *tls.Config
}

type quicClient struct {
	appCtx   *appctx.AppContext
	cfg      QuicClientConfig
	udpConn  *net.UDPConn
	quicConn quic.Connection
	stream   quic.Stream
	encoder  *json.Encoder
	decoder  *json.Decoder
}

func NewQuicClient(appCtx *appctx.AppContext, cfg QuicClientConfig) QuicClient {
	return &quicClient{
		appCtx: appCtx,
		cfg:    cfg,
	}
}

func (qc *quicClient) Name() string {
	return qc.cfg.Name
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
	appCtx       *appctx.AppContext
	cfg          QuicServerConfig
	udpConn      *net.UDPConn
	quicListener *quic.Listener
	Clients      *sync.Map
	handler      StreamHandler
}

func NewQuicServer(appCtx *appctx.AppContext, cfg QuicServerConfig, handler StreamHandler) QuicServer {
	return &quicServer{
		appCtx:  appCtx,
		cfg:     cfg,
		Clients: &sync.Map{},
		handler: handler,
	}
}

func (qs *quicServer) Name() string {
	return qs.cfg.Name
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

func (qs *quicServer) StartServer(ctx context.Context) error {
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
	if qs.quicListener != nil {
		if err := qs.quicListener.Close(); err != nil {
			return err
		}
	}

	if qs.udpConn != nil {
		if err := qs.udpConn.Close(); err != nil {
			return err
		}
	}

	return nil
}
