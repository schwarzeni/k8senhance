package server

import (
	"github.com/schwarzeni/k8senhance/agent/imagecache/cache"
	"github.com/schwarzeni/k8senhance/pkg/metrics"
	"github.com/schwarzeni/k8senhance/pkg/processor"
)

var metricProcessor = processor.NewScoreProcessor()

// prioritizeNodes 选出合适的节点
// 如果不存在，则返回 -1
func prioritizeNodes(layer string, nodes []*QueryLayerRespDTO) int {
	// TODO: extract max size to setting
	//log.Println("[debug] priority nodes")
	//for _, node := range nodes {
	//	fmt.Printf("%+v", *node)
	//}
	size, ok, _ := cache.LayerInfo(layer)
	//log.Println("[debug] priority size", layer, size)
	if ok && size < 1024*500 {
		return -1
	}
	maxScore := 0.0
	maxScoreIdx := -1
	for idx, node := range nodes {
		// TODO: maybe use cache to improve performance
		score, _ := metricProcessor.Score(&metrics.NodeFullMetric{RawMetric: node.Metric})
		if score > maxScore {
			maxScore = score
			maxScoreIdx = idx
		}
		//log.Println("[debug] priority score", node, score)
	}
	return maxScoreIdx
}
