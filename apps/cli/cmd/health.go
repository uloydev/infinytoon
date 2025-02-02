package cmd

import (
	"github.com/spf13/cobra"
	"infinitoon.dev/infinitoon/apps/cli/utils"
	appctx "infinitoon.dev/infinitoon/pkg/context"
)

var HealthCommand CommandFunc = func(appCtx *appctx.AppContext) *cobra.Command {
	health := &cobra.Command{
		Use:   "health",
		Short: "Check health of your InfiniToon Tunnels",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(utils.Banner)
			cmd.Println("Checking health of your InfiniToon Tunnels")
		},
	}

	return health
}
