package cmd

import (
	"github.com/spf13/cobra"
	"infinitoon.dev/infinitoon/apps/cli/utils"
	appctx "infinitoon.dev/infinitoon/pkg/context"
)

var AuthCommand CommandFunc = func(appCtx *appctx.AppContext) *cobra.Command {
	auth := &cobra.Command{
		Use:    "auth",
		Short:  "Authenticate your InfiniToon account",
		PreRun: utils.DefaultPreRun,
		Run: func(cmd *cobra.Command, args []string) {
			// check flags
			token, err := cmd.Flags().GetString("token")
			if err != nil {
				cmd.Println("Error: ", err)
				return
			}

			// check if already authenticated
			if oldToken, err := utils.ReadToken(); err == nil && oldToken != "" {
				replace, err := cmd.Flags().GetBool("replace")
				if err != nil {
					cmd.Println("Error: ", err)
					return
				}

				if !replace {
					cmd.Println("Already authenticated")
					cmd.Println("Use --replace or -r flag to replace existing token")
					return
				}
			}

			tokenValid := true

			cmd.Printf("Authenticating your InfiniToon account with token %s\n", token)

			if !tokenValid {
				cmd.Println("Invalid token")
				return
			}

			// save token to .inifinitoon file in user home directory
			if err := utils.SaveToken(token); err != nil {
				cmd.Println("Error: ", err)
				return
			}
		},
	}

	// replace flag to replace existing token
	auth.Flags().BoolP("replace", "r", false, "Replace existing token")

	// token flag required
	auth.Flags().StringP("token", "t", "", "Token of your InfiniToon account")
	auth.MarkFlagRequired("token")

	return auth
}
