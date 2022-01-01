package nodemonitor

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/schwarzeni/k8senhance/config"
	"github.com/schwarzeni/k8senhance/pkg/metrics"
)

type NodeMonitor struct {
	gc *config.Config
}

func (ic *NodeMonitor) Run() error {
	log.Println("start node monitor service")
	c := metrics.NewDefaultCollector()
	httpClient := http.DefaultClient
	metricURL := ic.gc.Agent.NodeMonitor.ControllerAddr + "/api/v1/agenthealth/" + ic.gc.NodeName
	for {
		time.Sleep(time.Second)
		nodeMetric := &metrics.NodeMetric{}
		if err := c.Collect(nodeMetric); err != nil {
			log.Println("[err] collect metric:", err)
			continue
		}
		nodeMetric.NodeInfo = metrics.NodeInfo{ID: ic.gc.NodeName}
		binaryData, err := json.Marshal(nodeMetric)
		if err != nil {
			log.Println("[err] marshal json data:", err)
			continue
		}
		nodeMetric.Timestamp = time.Now()
		req, err := http.NewRequest(http.MethodPut, metricURL, bytes.NewBuffer(binaryData))
		if err != nil {
			log.Println("[err] gen http request:", err)
			continue
		}
		resp, err := httpClient.Do(req)
		if err != nil {
			log.Println(err)
			continue
		}
		_ = resp.Body.Close()
	}
}

func NewNodeMonitor(gc *config.Config) *NodeMonitor {
	return &NodeMonitor{gc: gc}
}
