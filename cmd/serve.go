package cmd

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/hjoshi123/WaaS/config"
	"github.com/hjoshi123/WaaS/ip"
	"github.com/hjoshi123/WaaS/util"
	"github.com/place1/wg-embed/pkg/wgembed"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	serve = &cobra.Command{
		Use:   "serve",
		Short: "Start the WaaS server",
		Long:  "Starting the WaaS server with config",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunServe(cmd, args)
		},
	}
)

func RunServe(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	ctx, span := util.Tracer.Start(ctx, "ServeWG")
	defer span.End()
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
		}
	}()
	util.Logger(ctx).Debug("Starting the WaaS server", zap.Any("config", config.Spec))
	serverIP := ip.GetWireGuardServerIP(config.Spec.VPN.CIDR)
	config.Spec.VPN.AllowedIPs = append(config.Spec.VPN.AllowedIPs, fmt.Sprintf("%s/32", serverIP.IP.String()))

	wg := wgembed.NewNoOpInterface()

	// this needs to be run only if wireguard needs to be run.
	if config.Spec.WG.Enabled {
		wgimpl, err := wgembed.New(config.Spec.WG.Interface)
		if err != nil {
			util.Logger(ctx).Error("Error creating WireGuard interface", zap.Error(err))
			return err
		}
		defer wgimpl.Close()
		wg = wgimpl

		util.Logger(ctx).Info("Starting WireGuard server on 0.0.0.0:", zap.Int("port", config.Spec.WG.Port), zap.String("interface", config.Spec.WG.Interface))

		wgconfig := &wgembed.ConfigFile{
			Interface: wgembed.IfaceConfig{
				PrivateKey: config.Spec.WG.PrivateKey,
				Address:    serverIP.String(),
				ListenPort: &config.Spec.WG.Port,
			},
		}

		if err := wg.LoadConfig(wgconfig); err != nil {
			util.Logger(ctx).Error("Error loading WireGuard config", zap.Error(err))
			return err
		}

		if err := ip.ConfigureIPTables(ctx, config.Spec.WG.Interface, config.Spec.VPN.GatewayInterface, config.Spec.VPN.CIDR, config.Spec.VPN.AllowedIPs); err != nil {
			util.Logger(ctx).Error("Error configuring IPTables", zap.Error(err))
			return err
		}
	}

	if config.Spec.DNS.Enabled {
		dns, err := ip.New(ctx, config.Spec.DNS.Upstream)
		if err != nil {
			util.Logger(ctx).Error("Error creating DNS server", zap.Error(err))
			return err
		}

		defer dns.Close()
	}

	router := mux.NewRouter()

	return nil
}
