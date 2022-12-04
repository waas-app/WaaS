package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/hjoshi123/WaaS/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "waas",
		Short: "waas is a command line tool for interacting with the Wireguard",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initConfig()
		},
	}
	cfgFile     string
	environment string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.waas.yaml)")
	rootCmd.PersistentFlags().StringVar(&environment, "environment", "development", "environment to run in")
	viper.BindPFlag("environment", rootCmd.PersistentFlags().Lookup("environment"))

	rootCmd.AddCommand(serve)
}

func Execute() error {
	return rootCmd.Execute()
}

func initConfig() {
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
	log.Println("Environment: ", viper.GetString("environment"))

	if environment == "" {
		viper.Set("environment", config.Development)
	}
}
