package cmd

import (
	"github.com/spf13/cobra"
	"infinitoon.dev/infinitoon/apps/cli/utils"
	appctx "infinitoon.dev/infinitoon/pkg/context"
)

var VersionCommand CommandFunc = func(appCtx *appctx.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:    "version",
		Short:  "Print the version number of InfiniToon CLI",
		PreRun: utils.DefaultPreRun,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("InfiniToon CLI ", utils.Version)
		},
	}
}
