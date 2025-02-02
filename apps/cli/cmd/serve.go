package cmd

import (
	"github.com/spf13/cobra"
	"infinitoon.dev/infinitoon/apps/cli/utils"
	appctx "infinitoon.dev/infinitoon/pkg/context"
)

var ServeCommand CommandFunc = func(appCtx *appctx.AppContext) *cobra.Command {
	serve := &cobra.Command{
		Use:    "serve",
		Short:  "Serve your local http server to Internet securely with InfiniToon",
		PreRun: utils.DefaultPreRun,
		Run: func(cmd *cobra.Command, args []string) {
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
