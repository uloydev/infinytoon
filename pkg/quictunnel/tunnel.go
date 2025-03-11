package quictunnel

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/logger"
	"infinitoon.dev/infinitoon/shared/packets"
)

type StreamHandler func(*appctx.AppContext, quic.Stream, *json.Encoder, packets.Message)
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
	appCtx  *appctx.AppContext
	log     *logger.Logger
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

func NewQuicTunnel(appCtx *appctx.AppContext) QuicTunnel {
	return &quicTunnel{
		appCtx:  appCtx,
		log:     appCtx.Get(appctx.LoggerKey).(*logger.Logger),
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
		go func(k QuicClientKey, client QuicClient) {
			qt.log.Info().Any("client", k).Msg("starting quic client")
			if err := client.Setup(context.Background()); err != nil {
				qt.log.Error().Err(err).Any("client", k).Msg("failed to start quic client")
				return
			}
			qt.log.Info().Any("client", k).Msg("quic client started")
		}(k, client)
	}

	for k, server := range qt.servers {
		go func(k QuicServerKey, server QuicServer) {
			qt.log.Info().Any("server", k).Msg("starting quic server")
			if err := server.StartServer(context.Background()); err != nil {
				qt.log.Error().Err(err).Any("server", k).Msg("failed to start quic server")
				return
			}
		}(k, server)
	}
}

func (qt *quicTunnel) Shutdown() {
	for k, client := range qt.clients {
		qt.log.Info().Any("client", k).Msg("shutting down quic client")
		if err := client.ShutdownClient(context.Background()); err != nil {
			qt.log.Error().Err(err).Any("client", k).Msg("failed to shutdown quic client")
			continue
		}
		qt.log.Info().Any("client", k).Msg("quic client shutdown")
	}

	for k, server := range qt.servers {
		qt.log.Info().Any("server", k).Msg("shutting down quic server")
		if err := server.ShutdownServer(context.Background()); err != nil {
			qt.log.Error().Err(err).Any("server", k).Msg("failed to shutdown quic server")
			continue
		}
		qt.log.Info().Any("server", k).Msg("quic server shutdown")
	}
}

type QuicClientConfig struct {
	Name       string `mapstructure:"name"`
	IP         string `mapstructure:"ip"`
	Port       int    `mapstructure:"port"`
	TLSConfing *tls.Config
}

type QuicServerConfig struct {
	Name       string `mapstructure:"name"`
	IP         string `mapstructure:"ip"`
	Port       int    `mapstructure:"port"`
	TLSConfing *tls.Config
}

type quicClient struct {
	appCtx   *appctx.AppContext
	log      *logger.Logger
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
		log:    appCtx.Get(appctx.LoggerKey).(*logger.Logger),
		cfg:    cfg,
	}
}

func (qc *quicClient) Name() string {
	return qc.cfg.Name
}

func (qc *quicClient) Setup(ctx context.Context) error {
	for {

		addr := &net.UDPAddr{
			IP:   net.ParseIP(qc.cfg.IP),
			Port: qc.cfg.Port,
		}

		quicConn, err := quic.DialAddr(ctx, addr.String(), qc.cfg.TLSConfing, &quic.Config{
			KeepAlivePeriod: 10 * time.Minute,
		})
		if err != nil {
			qc.log.Error().Err(err).Any("client", qc.cfg.Name).Msg("failed to dial quic address, retrying...")
			time.Sleep(5 * time.Second)
			continue
		}

		qc.quicConn = quicConn

		stream, err := quicConn.OpenStream()
		if err != nil {
			qc.log.Error().Err(err).Any("client", qc.cfg.Name).Msg("failed to open stream, retrying...")
			time.Sleep(5 * time.Second)
			continue
		}

		qc.stream = stream

		qc.encoder = json.NewEncoder(stream)
		qc.decoder = json.NewDecoder(stream)

		go qc.echo()

		return nil
	}
}

func (qc *quicClient) echo() {
	qc.sendEcho(time.Now())
	//  send echo every minute
	tick := time.NewTicker(time.Minute)
	for t := range tick.C {
		qc.sendEcho(t)
	}
}

func (qc *quicClient) sendEcho(t time.Time) {
	echo := packets.NewEchoPayload(qc.cfg.Name)
	echoRq, err := echo.EncodeRq()
	if err != nil {
		qc.log.Error().Err(err).Any("client", qc.cfg.Name).Msg("error encoding echo message")
		return
	}
	qc.log.Info().Any("client", qc.cfg.Name).Any("payload", echoRq).Any("time", t).Msg("sending echo message")
	res, err := qc.SendMessage(context.Background(), echoRq)
	if err != nil {
		qc.log.Error().Err(err).Any("client", qc.cfg.Name).Msg("error sending echo message")
		return
	}
	qc.log.Info().Any("client", qc.cfg.Name).Any("payload", res).Any("time", t).Msg("echo response received")
}

func (qc *quicClient) SendMessage(ctx context.Context, msg *packets.Message) (*packets.Message, error) {
	return sendMessage(ctx, qc.encoder, qc.decoder, msg)
}

func (qc *quicClient) Stream(ctx context.Context, handler StreamHandler) {
	for {
		stream, err := qc.quicConn.AcceptStream(ctx)
		if err != nil {
			handleStreamError(qc.log, qc.quicConn.RemoteAddr().String(), err)
			continue
		}

		encoder := json.NewEncoder(stream)
		decoder := json.NewDecoder(stream)

		var (
			req, res packets.Message
		)
		if err := decoder.Decode(&req); err != nil {
			qc.log.Error().Any("client", qc.cfg.Name).Any("streamID", stream.StreamID()).Err(err).Msg("error decoding message")
			res.Type = packets.ErrInvalidPayload
			res.ClientID = stream.StreamID().InitiatedBy().String()
			encoder.Encode(res)
		}
		go handler(qc.appCtx, stream, encoder, req)
	}
}

func (qc *quicClient) ShutdownClient(ctx context.Context) error {
	if qc.stream != nil {
		if err := qc.stream.Close(); err != nil {
			return err
		}
	}

	if qc.quicConn != nil {
		return qc.quicConn.CloseWithError(0, "close normal")
	}
	return nil
}

type quicServer struct {
	appCtx       *appctx.AppContext
	log          *logger.Logger
	cfg          QuicServerConfig
	udpConn      *net.UDPConn
	quicListener *quic.Listener
	Clients      *sync.Map
	handler      StreamHandler
}

func NewQuicServer(appCtx *appctx.AppContext, cfg QuicServerConfig, handler StreamHandler) QuicServer {
	return &quicServer{
		appCtx:  appCtx,
		log:     appCtx.Get(appctx.LoggerKey).(*logger.Logger),
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

	quicLs, err := quic.Listen(udpConn, qs.cfg.TLSConfing, &quic.Config{
		KeepAlivePeriod: 10 * time.Minute,
	})
	if err != nil {
		return err
	}
	qs.quicListener = quicLs

	qs.log.Info().Any("server", qs.cfg.Name).Any("address", quicLs.Addr()).Msg("quic server listening")
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
	qs.log.Info().Any("server", qs.cfg.Name).Any("client", conn.RemoteAddr()).Msg("new client connected to server")

	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			handleStreamError(qs.log, conn.RemoteAddr().String(), err)

			// remove client connection from map
			qs.Clients.Delete(conn.RemoteAddr().String())
			qs.log.Info().Any("server", qs.cfg.Name).Any("client", conn.RemoteAddr()).Msg("client disconnected")
			break
		}

		decoder := json.NewDecoder(stream)
		encoder := json.NewEncoder(stream)

		for {
			// check if stream is closed
			if _, err := stream.Read(nil); err != nil {
				handleStreamError(qs.log, conn.RemoteAddr().String(), err)
				break
			}

			var (
				req, res packets.Message
			)
			if err := decoder.Decode(&req); err != nil {

				// check if err timeout: no recent network activity
				if err.Error() == "timeout: no recent network activity" {
					qs.log.Info().Any("server", qs.cfg.Name).Any("client", conn.RemoteAddr()).Msg("client idle timeout")
					break
				}

				qs.log.Error().Any("server", qs.cfg.Name).Any("client", conn.RemoteAddr()).Any("streamID", stream.StreamID()).Err(err).Msg("error decoding message")
				res.Type = packets.ErrInvalidPayload
				res.ClientID = stream.StreamID().InitiatedBy().String()
				encoder.Encode(res)
			}
			go func() {
				qs.handler(qs.appCtx, stream, encoder, req)
			}()
		}
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
