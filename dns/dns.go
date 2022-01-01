package dns

import (
	"fmt"
	"log"

	"github.com/miekg/dns"

	"github.com/schwarzeni/k8senhance/config"
	"github.com/spf13/cobra"
)

var (
	globalConfigPath *string
	globalConfig     *config.Config
)

var cmd = &cobra.Command{
	Use:   "dns",
	Short: "simple dns service",
	Long:  "simple dns service",
	Run: func(cmd *cobra.Command, args []string) {
		globalConfig = config.MustParse(*globalConfigPath)
		if len(globalConfig.Cloud.DNS.CloudEpIps) == 0 {
			panic("cloud ep ip is empty")
		}
		log.Println("dns service start, listening at port ", globalConfig.Cloud.DNS.Port)
		service()
	},
}

func InitCMD(configPath *string) *cobra.Command {
	globalConfigPath = configPath
	return cmd
}

func service() {
	dns.HandleFunc(".", handleDnsRequest)

	// start server
	server := &dns.Server{Addr: ":" + globalConfig.Cloud.DNS.Port, Net: "udp"}
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		//case dns.TypeA:
		default:
			// TODO: 这里默认请求格式都是正确的
			log.Printf("Query for %s\n", q.Name)
			// TODO: 目前仅支持单个 ep
			rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, globalConfig.Cloud.DNS.CloudEpIps[0]))
			if err == nil {
				m.Answer = append(m.Answer, rr)
			}
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	fmt.Println(m.Answer)
	w.WriteMsg(m)
}
