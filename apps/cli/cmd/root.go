package cmd

import (
	"github.com/spf13/cobra"
	"infinitoon.dev/infinitoon/apps/cli/utils"
	appctx "infinitoon.dev/infinitoon/pkg/context"
)

type CommandFunc func(appCtx *appctx.AppContext) *cobra.Command

// register commands here
var CommandList = []CommandFunc{
	VersionCommand,
	ServeCommand,
	ListCommand,
	AuthCommand,
	HealthCommand,
	UpdateCommand,
}

func Execute(appCtx *appctx.AppContext) {
	rootCmd := &cobra.Command{
		Use:    "toon",
		PreRun: utils.DefaultPreRun,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("This is a CLI tool for InfiniToon, Serve your local http server to Internet with InfiniToon.")
			cmd.Help()
		},
	}

	for _, cmd := range CommandList {
		rootCmd.AddCommand(cmd(appCtx))
	}

	rootCmd.Execute()
}
