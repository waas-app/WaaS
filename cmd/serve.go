package cmd

import (
	"github.com/hjoshi123/WaaS/config"
	"github.com/hjoshi123/WaaS/util"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	serve = &cobra.Command{
		Use:   "serve",
		Short: "Start the WaaS server",
		Long:  "Starting the WaaS server with config",
		Run: func(cmd *cobra.Command, args []string) {
			RunServe(cmd, args)
		},
	}
)

func RunServe(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	util.Logger(ctx).Debug("Starting the WaaS server", zap.Any("config", config.Spec))
}
