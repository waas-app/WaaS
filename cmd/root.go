package cmd

import (
	"fmt"
	"os"

	"github.com/hjoshi123/WaaS/config"
	"github.com/hjoshi123/WaaS/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "waas",
		Short: "waas is a command line tool for interacting with the Wireguard",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			InitConfig()
		},
	}
	cfgFile string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.waas.yaml)")
	rootCmd.PersistentFlags().StringVar(&config.Spec.Environment, "environment", "", "environment to run in")
	rootCmd.PersistentFlags().StringVar(&config.Spec.AdminUserName, "WG_ADMIN_USERNAME", "admin", "admin username")
	rootCmd.PersistentFlags().StringVar(&config.Spec.AdminPassword, "WG_ADMIN_PASSWORD", "admin", "admin password")
	rootCmd.PersistentFlags().IntVar(&config.Spec.WG.Port, "WG_PORT", 51820, "port to run wireguard on")
	rootCmd.PersistentFlags().StringVar(&config.Spec.WG.Interface, "WG_INTERFACE", "wg0", "interface to run wireguard on")
	rootCmd.PersistentFlags().StringVar(&config.Spec.WG.PrivateKey, "WG_PRIVATE_KEY", "", "private key to run wireguard on")
	rootCmd.PersistentFlags().IntVar(&config.Spec.Port, "UI_PORT", 8080, "port to run wireguard on")
	rootCmd.PersistentFlags().StringVar(&config.Spec.ExternalHost, "EXTERNAL_HOST", "localhost", "external host to run wireguard on")
	rootCmd.PersistentFlags().StringVar(&config.Spec.Storage, "STORAGE", "memory", "storage to run wireguard on")
	rootCmd.PersistentFlags().BoolVar(&config.Spec.WG.Enabled, "WG_ENABLED", true, "enable wireguard")
	rootCmd.PersistentFlags().StringVar(&config.Spec.VPN.CIDR, "VPN_CIDR", "192.168.2.0/24", "cidr to run wireguard on")
	rootCmd.PersistentFlags().StringVar(&config.Spec.VPN.GatewayInterface, "VPN_GATEWAY_INTERFACE", "eth0", "gateway interface to run wireguard on")
	rootCmd.PersistentFlags().StringArrayVar(&config.Spec.VPN.AllowedIPs, "VPN_ALLOWED_IPS", []string{"0.0.0.0/0"}, "allowed ips to run wireguard on")
	rootCmd.PersistentFlags().BoolVar(&config.Spec.DNS.Enabled, "DNS_ENABLED", true, "dns to run wireguard on")
	rootCmd.PersistentFlags().StringArrayVar(&config.Spec.DNS.Upstream, "DNS_UPSTREAM", []string{"1.1.1.1"}, "upstream dns to run wireguard on")

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

	rootCmd.AddCommand(serve)
}

func Execute() error {
	return rootCmd.Execute()
}

func InitConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".waas")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	viper.Unmarshal(&config.Spec)

	if config.Spec.Environment == "" {
		viper.Set("environment", config.Development)
	}

	util.InitLogger()
}
