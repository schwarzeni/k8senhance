package netproxy

import (
	"log"

	"github.com/schwarzeni/k8senhance/config"
)

type NetProxy struct {
	gc *config.Config
}

func (ic *NetProxy) Run() error {
	log.Println("start net proxy service")
	return nil
}

func NewNetProxy(gc *config.Config) *NetProxy {
	return &NetProxy{gc: gc}
}
