package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/waas-app/WaaS/config"
	"github.com/waas-app/WaaS/util"
	"go.uber.org/zap"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var (
	rootCmd = &cobra.Command{
		Use:   "waas",
		Short: "waas is a command line tool for interacting with the Wireguard",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			file := new(os.File)
			var err error
			if config.Spec.OTLPEndpoint == "" {
				file, err = os.Create("traces.txt")
				if err != nil {
					util.Logger(cmd.Context()).Error("failed to create file", zap.Error(err))
					return err
				}
				defer file.Close()
			}

			ctx := context.Background()
			if cmd.Context() == nil {
				cmd.SetContext(ctx)
			}
			ctx, tCleanup, err := util.InitOTEL(ctx, "true", config.ServiceName, true, file)
			if err != nil {
				util.Logger(ctx).Error("failed to initialize opentelemetry", zap.Error(err))
				return err
			}
			defer tCleanup(ctx)
			return nil
		},
	}
	cfgFile string
)

func init() {
	cobra.OnInitialize(InitConfig)

	viper.AutomaticEnv()
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.waas.yaml)")
	rootCmd.PersistentFlags().StringVar(&config.Spec.OTLPEndpoint, "OTLP_ENDPOINT", "10.0.0.123:4317", "OTLP endpoint")
	rootCmd.PersistentFlags().StringVar(&config.Spec.Environment, "environment", "development", "environment to run in")
	rootCmd.PersistentFlags().StringVar(&config.Spec.AdminUserName, "WG_ADMIN_USERNAME", "admin", "admin username")
	rootCmd.PersistentFlags().StringVar(&config.Spec.AdminPassword, "WG_ADMIN_PASSWORD", "admin", "admin password")
	rootCmd.PersistentFlags().IntVar(&config.Spec.WG.Port, "WG_PORT", 51810, "port to run wireguard on")
	rootCmd.PersistentFlags().StringVar(&config.Spec.WG.Interface, "WG_INTERFACE", "wg0", "interface to run wireguard on")
	rootCmd.PersistentFlags().StringVar(&config.Spec.WG.PrivateKey, "WG_PRIVATE_KEY", "", "private key to run wireguard on")
	rootCmd.PersistentFlags().IntVar(&config.Spec.Port, "UI_PORT", 8000, "port to run wireguard on")
	rootCmd.PersistentFlags().StringVar(&config.Spec.ExternalHost, "EXTERNAL_HOST", "localhost", "external host to run wireguard on")
	rootCmd.PersistentFlags().StringVar(&config.Spec.Storage, "STORAGE", "postgresql://waas:WaasPassw0rd@postgresql:5432/waas", "storage to run wireguard on")
	rootCmd.PersistentFlags().BoolVar(&config.Spec.WG.Enabled, "WG_ENABLED", true, "enable wireguard")
	rootCmd.PersistentFlags().StringVar(&config.Spec.VPN.CIDR, "VPN_CIDR", "192.168.2.0/24", "cidr to run wireguard on")
	rootCmd.PersistentFlags().StringVar(&config.Spec.VPN.GatewayInterface, "VPN_GATEWAY_INTERFACE", "eth0", "gateway interface to run wireguard on")
	rootCmd.PersistentFlags().StringArrayVar(&config.Spec.VPN.AllowedIPs, "VPN_ALLOWED_IPS", []string{"0.0.0.0/0"}, "allowed ips to run wireguard on")
	rootCmd.PersistentFlags().BoolVar(&config.Spec.DNS.Enabled, "DNS_ENABLED", true, "dns to run wireguard on")
	rootCmd.PersistentFlags().StringArrayVar(&config.Spec.DNS.Upstream, "DNS_UPSTREAM", []string{"1.1.1.1"}, "upstream dns to run wireguard on")
	rootCmd.PersistentFlags().StringVar(&config.Spec.RootURL, "ROOT_URL", "http://localhost:3000", "root url to run wireguard on")
	rootCmd.PersistentFlags().StringVar(&config.Spec.SessionSecret, "SESSION_SECRET", "3bcf9f7cbc479b854f6877e917f82df03110db179d121f0c00bfd3afaa28f52eaff20af628b1e67caf9b7b39648e1c892df11036f9d2f2f767ede807d4c2779", "session secret")
	rootCmd.PersistentFlags().StringVar(&config.Spec.CookieDomain, "COOKIE_DOMAIN", "localhost", "cookie domain")
	rootCmd.PersistentFlags().StringVar(&config.Spec.Redis, "REDIS_URL", "redis://redis:6379", "redis url")

	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("environment", rootCmd.PersistentFlags().Lookup("environment"))
	viper.BindPFlag("admin_username", rootCmd.PersistentFlags().Lookup("WG_ADMIN_USERNAME"))
	viper.BindPFlag("admin_password", rootCmd.PersistentFlags().Lookup("WG_ADMIN_PASSWORD"))
	viper.BindPFlag("wg-port", rootCmd.PersistentFlags().Lookup("WG_PORT"))
	viper.BindPFlag("wg-interface", rootCmd.PersistentFlags().Lookup("WG_INTERFACE"))
	viper.BindPFlag("wg-privateKey", rootCmd.PersistentFlags().Lookup("WG_PRIVATE_KEY"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("UI_PORT"))
	viper.BindPFlag("externalHost", rootCmd.PersistentFlags().Lookup("EXTERNAL_HOST"))
	viper.BindPFlag("storage", rootCmd.PersistentFlags().Lookup("STORAGE"))
	viper.BindPFlag("wg-enabled", rootCmd.PersistentFlags().Lookup("WG_ENABLED"))
	viper.BindPFlag("vpn-cidr", rootCmd.PersistentFlags().Lookup("VPN_CIDR"))
	viper.BindPFlag("vpn-gatewayInterface", rootCmd.PersistentFlags().Lookup("VPN_GATEWAY_INTERFACE"))
	viper.BindPFlag("vpn-allowedIPs", rootCmd.PersistentFlags().Lookup("VPN_ALLOWED_IPS"))
	viper.BindPFlag("dns-enabled", rootCmd.PersistentFlags().Lookup("DNS_ENABLED"))
	viper.BindPFlag("dns-upstream", rootCmd.PersistentFlags().Lookup("DNS_UPSTREAM"))
	viper.BindPFlag("otlp_endpoint", rootCmd.PersistentFlags().Lookup("OTLP_ENDPOINT"))
	viper.BindPFlag("root_url", rootCmd.PersistentFlags().Lookup("ROOT_URL"))
	viper.BindPFlag("session_secret", rootCmd.PersistentFlags().Lookup("SESSION_SECRET"))
	viper.BindPFlag("cookie_domain", rootCmd.PersistentFlags().Lookup("COOKIE_DOMAIN"))
	viper.BindPFlag("redis_url", rootCmd.PersistentFlags().Lookup("REDIS_URL"))

	rootCmd.AddCommand(serve)
}

func Execute() error {
	// create new context from root command
	// file := new(os.File)
	// var err error
	// if config.Spec.OTLPEndpoint == "" {
	// 	file, err = os.Create("traces.txt")
	// 	if err != nil {
	// 		util.Logger(rootCmd.Context()).Error("failed to create file", zap.Error(err))
	// 		return err
	// 	}
	// 	defer file.Close()
	// }

	ctx := context.Background()
	if rootCmd.Context() == nil {
		rootCmd.SetContext(ctx)
	}
	// ctx, tCleanup, err := util.InitOTEL(ctx, "true", config.ServiceName, true, file)
	// if err != nil {
	// 	util.Logger(ctx).Error("failed to initialize opentelemetry", zap.Error(err))
	// 	return err
	// }
	// defer tCleanup(ctx)

	return rootCmd.ExecuteContext(ctx)
}

func InitConfig() {
	var wd string
	var err error
	ctx := context.Background()
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		wd, err = os.Getwd()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(wd)
		viper.SetConfigType("yml")
		viper.SetConfigName("waas")
	}

	log.Println("Working Directory: ", wd)

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	viper.Unmarshal(&config.Spec)

	if config.Spec.Environment == "" {
		viper.Set("environment", config.Development)
	}

	util.InitLogger()

	if config.Spec.WG.PrivateKey == "" {
		key, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			util.Logger(ctx).Fatal("failed to generate a server private key", zap.Error(err))
		}
		config.Spec.WG.PrivateKey = key.String()
	}
}
