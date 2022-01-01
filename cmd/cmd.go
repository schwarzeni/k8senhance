package cmd

import (
	"github.com/schwarzeni/k8senhance/agent"
	"github.com/schwarzeni/k8senhance/crd"
	"github.com/schwarzeni/k8senhance/dns"
	"github.com/schwarzeni/k8senhance/netproxy"
	"github.com/schwarzeni/k8senhance/nodescheduler"
	"github.com/schwarzeni/k8senhance/svcstore"
	"github.com/spf13/cobra"
)

var (
	configPath string
)

var rootCmd = &cobra.Command{
	Use:   "k8senhance",
	Short: "",
	Long:  "",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "config file path")
}

func Execute() {
	rootCmd.AddCommand(agent.InitCMD(&configPath))
	rootCmd.AddCommand(crd.InitCMD(&configPath))
	rootCmd.AddCommand(dns.InitCMD(&configPath))
	rootCmd.AddCommand(nodescheduler.InitCMD(&configPath))
	rootCmd.AddCommand(svcstore.InitCMD(&configPath))
	rootCmd.AddCommand(netproxy.InitCMD(&configPath))
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
