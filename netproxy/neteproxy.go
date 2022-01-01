package netproxy

import (
	"log"
	"net"

	"github.com/schwarzeni/k8senhance/config"
	"github.com/schwarzeni/k8senhance/pkg/pb"
	"google.golang.org/grpc"

	"github.com/spf13/cobra"
)

var (
	globalConfigPath *string
	globalConfig     *config.Config
)

var cmd = &cobra.Command{
	Use:   "netproxy",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("start netproxy service")
		globalConfig = config.MustParse(*globalConfigPath)
		lis, err := net.Listen("tcp", globalConfig.Cloud.Netproxy.GrpcAddr)
		if err != nil {
			panic(err)
		}
		grpcServer := grpc.NewServer()
		proxyChannelManager := NewProxyChannelManager()
		pb.RegisterProxyHttpServiceServer(grpcServer, &CloudEdgeServer{
			RequestCallBackMap:  make(map[string]ProxyRequestCallBack),
			proxyChannelManager: proxyChannelManager,
		})
		go RunHTTPProxyServer(globalConfig.Cloud.Netproxy.HTTPProxyAddr, proxyChannelManager)
		log.Println("[cloud start]")
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	},
}

func InitCMD(configPath *string) *cobra.Command {
	globalConfigPath = configPath
	return cmd
}
