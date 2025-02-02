package cmd

import (
	"github.com/spf13/cobra"
	"infinitoon.dev/infinitoon/apps/cli/utils"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/shared/schema"
)

var ListCommand CommandFunc = func(appCtx *appctx.AppContext) *cobra.Command {
	validStatus := []schema.TunnelStatus{schema.TunnelStatusActive, schema.TunnelStatusInactive, schema.TunnelStatusError}
	list := &cobra.Command{
		Use:   "list",
		Short: "List all your tunnels",
		Long:  "List all your tunnels",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(utils.Banner)
			cmd.Println("List all your local http servers")
			status, err := cmd.Flags().GetString("status")
			if err != nil {
				cmd.Println("Error: ", err)
				return
			}

			// check if status is valid
			if status != "" {
				isValid := false
				for _, s := range validStatus {
					if string(s) == status {
						isValid = true
						break
					}
				}
				if !isValid {
					cmd.Println("Invalid status flag")
					return
				}
			}

		},
	}

	// add flag to filter based on status
	list.Flags().StringP("status", "s", "", "Filter based on status")
	return list
}
