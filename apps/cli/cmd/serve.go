package cmd

import (
	"context"
	"encoding/json"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/spf13/cobra"
	"infinitoon.dev/infinitoon/apps/cli/config"
	"infinitoon.dev/infinitoon/apps/cli/utils"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/http_forwarder"
	"infinitoon.dev/infinitoon/pkg/logger"
	"infinitoon.dev/infinitoon/pkg/quictunnel"
	"infinitoon.dev/infinitoon/shared/packets"
)

var ServeCommand CommandFunc = func(appCtx *appctx.AppContext) *cobra.Command {
	serve := &cobra.Command{
		Use:    "serve",
		Short:  "Serve your local http server to Internet securely with InfiniToon",
		PreRun: utils.DefaultPreRun,
		Run: func(cmd *cobra.Command, args []string) {
			cfg := appCtx.Get(appctx.ConfigKey).(*config.Config)
			forwarder := appCtx.Get(appctx.HttpForwarderKey).(http_forwarder.HttpForwarder)
			log := appCtx.Get(appctx.LoggerKey).(*logger.Logger)

			// check flags
			host, err := cmd.Flags().GetString("host")
			if err != nil {
				cmd.Println("Error: ", err)
				return
			}
			port, err := cmd.Flags().GetString("port")
			if err != nil {
				cmd.Println("Error: ", err)
				return
			}
			subdomain, err := cmd.Flags().GetString("subdomain")
			if err != nil {
				cmd.Println("Error: ", err)
				return
			}

			cmd.Printf("Serving your local http server to Internet securely with InfiniToon on %s:%s with subdomain %s\n", host, port, subdomain)

			quicClient := quictunnel.NewQuicClient(appCtx, cfg.TunnelClient)
			servePayload := packets.NewServeRq(host, port, cfg.TunnelClient.Name)

			// send serve request
			serveRq, err := servePayload.Encode()
			if err != nil {
				cmd.Println("Error: ", err)
				return
			}

			// TRY send serve request to server if failed retry until timeout 10s
			ctx, cancelCtx := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
			defer cancelCtx()

			if err := quicClient.Setup(ctx); err != nil {
				cmd.Println("Error: ", err)
				return

			}

			serveRs, err := quicClient.SendMessage(ctx, serveRq)
			if err != nil {
				cmd.Println("Error: ", err)
				return
			}

			if serveRs.Type != packets.ServeResponse {
				cmd.Println("Error: Invalid response type : " + serveRs.Type)
				return
			}

			quicClient.Stream(context.Background(), func(ac *appctx.AppContext, c quic.Connection, s quic.Stream, e *json.Encoder, m packets.Message) {

				log.Debug().Msgf("Received message: %s", m.Type)
				switch m.Type {
				case packets.HttpRequest:
					httpPayload := &packets.HttpRq{}
					err := httpPayload.Decode(&m)
					if err != nil {
						cmd.Println("Error: ", err)
						return
					}

					log.Info().Any("payload", httpPayload).Msgf("Received http request: %s %s", httpPayload.Method, httpPayload.Path)

					res, err := forwarder.Do(httpPayload)
					if err != nil {
						cmd.Println("Error: ", err)
						return
					}

					rs, err := res.Encode()
					if err != nil {
						cmd.Println("Error: ", err)
						return
					}

					log.Info().Any("payload", res).Msgf("Sending http response: %d", res.StatusCode)

					err = e.Encode(rs)
					if err != nil {
						cmd.Println("Error: ", err)
						return
					}
				}
			})
		},
	}

	// host flag required
	serve.Flags().StringP("host", "H", "", "Host of your local http server")
	serve.MarkFlagRequired("host")

	// port flag required
	serve.Flags().StringP("port", "p", "", "Port of your local http server")
	serve.MarkFlagRequired("port")

	// subdomain flag required
	serve.Flags().StringP("subdomain", "s", "", "Subdomain of your InfiniToon URL")
	serve.MarkFlagRequired("subdomain")

	return serve
}
