package cache

import (
	"sync"

	"github.com/schwarzeni/k8senhance/pkg/metrics"
)

// 对当前节点的 metric 信息进行缓存
var (
	nodeMetric     *metrics.NodeMetric
	nodeMetricLock sync.RWMutex
)

func SetNodeMetric(data *metrics.NodeMetric) error {
	nodeMetricLock.Lock()
	defer nodeMetricLock.Unlock()
	nodeMetric = data
	return nil
}

func NodeMetric() (*metrics.NodeMetric, error) {
	nodeMetricLock.RLock()
	defer nodeMetricLock.RUnlock()
	return nodeMetric, nil
}
