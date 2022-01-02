package server

import "github.com/schwarzeni/k8senhance/pkg/metrics"

type QueryLayerRespDTO struct {
	HasLayer bool               `json:"has_layer"`
	Metric   metrics.NodeMetric `json:"metric"`
}
