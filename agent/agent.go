package agent

import (
	"log"

	"github.com/schwarzeni/k8senhance/agent/netproxy"
	"github.com/schwarzeni/k8senhance/agent/nodemonitor"

	"github.com/schwarzeni/k8senhance/agent/imagecache"
	"github.com/schwarzeni/k8senhance/config"
	"github.com/spf13/cobra"
)

var (
	enableImageCacheFlag  bool
	enableNodeMonitorFlag bool
	enableNetProxyFlag    bool
	globalConfig          *config.Config
	globalConfigPath      *string
)
var cmd = &cobra.Command{
	Use:   "agent",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		globalConfig = config.MustParse(*globalConfigPath)
		log.Printf("config: %+v", *globalConfig)
		if enableImageCacheFlag {
			go func() {
				if err := imagecache.NewImageCache(globalConfig).Run(); err != nil {
					log.Fatal("start image cache", err)
				}
			}()
		}
		if enableNodeMonitorFlag {
			go func() {
				if err := nodemonitor.NewNodeMonitor(globalConfig).Run(); err != nil {
					log.Fatal("start node monitor", err)
				}
			}()
		}
		if enableNetProxyFlag {
			go func() {
				if err := netproxy.NewNetProxy(globalConfig).Run(); err != nil {
					log.Fatal("start net proxy", err)
				}
			}()
		}
		select {}
	},
}

func init() {
	cmd.Flags().BoolVar(&enableImageCacheFlag, "image-cache", false, "enable image cache feature")
	cmd.Flags().BoolVar(&enableNodeMonitorFlag, "node-monitor", false, "enable agent to collect system info")
	cmd.Flags().BoolVar(&enableNetProxyFlag, "net-proxy", false, "enable agent to do net proxy for container")
}

func InitCMD(configPath *string) *cobra.Command {
	globalConfigPath = configPath
	return cmd
}
