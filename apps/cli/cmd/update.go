package cmd

import (
	"github.com/spf13/cobra"
	"infinitoon.dev/infinitoon/apps/cli/utils"
	appctx "infinitoon.dev/infinitoon/pkg/context"
)

var UpdateCommand CommandFunc = func(appCtx *appctx.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:    "update",
		Short:  "Update InfiniToon CLI",
		PreRun: utils.DefaultPreRun,
		Run: func(cmd *cobra.Command, args []string) {
			remoteVersion := "v0.0.1"
			cmd.Printf("Current version : %s\n", utils.Version)
			cmd.Printf("Remote version  : %s\n", remoteVersion)

			if utils.Version == remoteVersion {
				cmd.Println("You are already using the latest version of InfiniToon CLI")
				return
			}
			cmd.Println("Updating InfiniToon CLI")
		},
	}
}
