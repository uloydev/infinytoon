package cmd

import (
	"github.com/spf13/cobra"
	"infinitoon.dev/infinitoon/apps/cli/utils"
	appctx "infinitoon.dev/infinitoon/pkg/context"
)

var AuthCommand CommandFunc = func(appCtx *appctx.AppContext) *cobra.Command {
	auth := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate your InfiniToon account",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(utils.Banner)
			// check flags
			token, err := cmd.Flags().GetString("token")
			if err != nil {
				cmd.Println("Error: ", err)
				return
			}
			cmd.Printf("Authenticating your InfiniToon account with token %s\n", token)
		},
	}

	// token flag required
	auth.Flags().StringP("token", "t", "", "Token of your InfiniToon account")
	auth.MarkFlagRequired("token")

	return auth
}
