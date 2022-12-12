package cmd

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/place1/wg-embed/pkg/wgembed"
	"github.com/spf13/cobra"
	"github.com/volatiletech/authboss/v3"
	"github.com/volatiletech/authboss/v3/remember"
	"github.com/waas-app/WaaS/config"
	"github.com/waas-app/WaaS/controller"
	"github.com/waas-app/WaaS/helpers/device"
	"github.com/waas-app/WaaS/infra"
	"github.com/waas-app/WaaS/infra/auth"
	"github.com/waas-app/WaaS/infra/middlewares"
	"github.com/waas-app/WaaS/ip"
	"github.com/waas-app/WaaS/util"
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

	wg wgembed.WireGuardInterface
)

func GetWgInterface() wgembed.WireGuardInterface {
	return wg
}

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

	wg = wgembed.NewNoOpInterface()

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

	dh := device.NewDeviceHelpers(wg)
	err = dh.RunSync(ctx)
	if err != nil {
		util.Logger(ctx).Fatal("Error running device sync", zap.Error(err))
		return err
	}

	auth.InitializeAuthBoss()

	router := mux.NewRouter()
	router.Use(middlewares.Logger)
	router.Use(middlewares.CheckUser)
	ab := auth.GetAuthBoss()
	router.Use(ab.LoadClientStateMiddleware, remember.Middleware(ab))

	router.Path("/ping").Methods(http.MethodGet).Handler(infra.CustomMux(controller.Ping))
	a := router.PathPrefix("/").Subrouter()
	a.Use(authboss.ModuleListMiddleware(ab))
	a.PathPrefix("/auth").Handler(ab.LoadClientStateMiddleware(http.StripPrefix("/auth", ab.Config.Core.Router)))

	site := router.PathPrefix("/").Subrouter()
	// site.Use(authboss.Middleware2(ab, authboss.RequireNone, authboss.RespondUnauthorized))
	site.PathPrefix("/api").Handler(controller.GRPCController(ctx, wg))

	w := router.PathPrefix("/").Subrouter()
	w.PathPrefix("/").Handler(controller.WebsiteRouter(ctx))

	address := fmt.Sprintf("0.0.0.0:%d", config.Spec.Port)
	srv := &http.Server{
		Addr:    address,
		Handler: router,
	}

	util.Logger(ctx).Info("Starting server on", zap.String("address", address))
	if err := srv.ListenAndServe(); err != nil {
		util.Logger(ctx).Fatal("Error starting server", zap.Error(err))
	}
	return nil
}
