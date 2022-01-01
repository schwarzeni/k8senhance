package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	NodeName string `yaml:"node_name"`
	Agent    struct {
		Imagecache struct {
			Addr           string `yaml:"addr"`
			RemoteRegistry string `yaml:"remote_registry"`
			CurrentIP      string `yaml:"current_ip"`
			CachePath      string `yaml:"cache_path"`
			ControllerAddr string `yaml:"controller_addr"`
		} `yaml:"imagecache"`
		Netproxy struct {
			CloudGrpcAddr string `yaml:"cloud_grpc_addr"`
			Addr          string `yaml:"addr"`
		} `yaml:"netproxy"`
		NodeMonitor struct {
			ControllerAddr string `yaml:"controller_addr"`
		} `yaml:"node_monitor"`
	} `yaml:"agent"`
	Cloud struct {
		Crd struct {
			K8Sconfig string `yaml:"k8sconfig"`
			StoreAddr string `yaml:"store_addr"`
		} `yaml:"crd"`
		DNS struct {
			CloudEpIps []string `yaml:"cloud_ep_ips"`
			Port       string   `yaml:"port"`
		} `yaml:"dns"`
		Netproxy struct {
			GrpcAddr      string `yaml:"grpc_addr"`
			StoreAddr     string `yaml:"store_addr"`
			HTTPProxyAddr string `yaml:"http_proxy_addr"`
		} `yaml:"netproxy"`
		NodeScheduler struct {
			Addr string `yaml:"addr"`
		} `yaml:"node_scheduler"`
		ServiceStore struct {
			Addr      string `yaml:"addr"`
			K8Sconfig string `yaml:"k8sconfig"`
		} `yaml:"service_store"`
	} `yaml:"cloud"`
}

func MustParse(configPath string) *Config {
	if len(configPath) == 0 {
		panic("you should set config file path")
	}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic("read file " + err.Error())
	}
	cf := &Config{}
	if err := yaml.Unmarshal(data, cf); err != nil {
		panic("parse config " + err.Error())
	}
	return cf
}
