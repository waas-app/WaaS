package cmd

import (
	"github.com/hjoshi123/WaaS/config"
	"github.com/spf13/cobra"
)

var (
	serve = &cobra.Command{
		Use:   "serve",
		Short: "Start the WaaS server",
		Long:  "Starting the WaaS server with config",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config.Logger(ctx).Info("Starting the WaaS server")
		},
	}
)
