package job

import (
	"log"
	"time"

	"github.com/schwarzeni/k8senhance/agent/imagecache/cache"
	"github.com/schwarzeni/k8senhance/config"
	"github.com/schwarzeni/k8senhance/pkg/metrics"
)

// 定期收集存储系统信息
type MetricJob struct {
	gc *config.Config
}

func (mj *MetricJob) Run() {
	log.Println("start metric job")
	c := metrics.NewDefaultCollector()
	for {
		time.Sleep(time.Second * 2)
		nodeMetric := &metrics.NodeMetric{}
		if err := c.Collect(nodeMetric); err != nil {
			log.Println("[metric job] failed collect metric:", err)
			continue
		}
		_ = cache.SetNodeMetric(nodeMetric)
	}
}

func NewMetricJob(gc *config.Config) *MetricJob {
	return &MetricJob{gc: gc}
}
