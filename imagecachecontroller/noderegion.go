package imagecachecontroller

import (
	"sync"

	"github.com/schwarzeni/k8senhance/pkg/model"
)

var (
	// TODO: 优化此处的数据结构
	nodeRegionDB     = map[string]map[string]*model.ImageCacheAgentHealthReqDTO{}
	nodeRegionDBLock sync.RWMutex
)

func setNodeRegionDB(nodeinfo *model.ImageCacheAgentHealthReqDTO) error {
	nodeRegionDBLock.Lock()
	defer nodeRegionDBLock.Unlock()
	if _, ok := nodeRegionDB[nodeinfo.Region]; !ok {
		nodeRegionDB[nodeinfo.Region] = map[string]*model.ImageCacheAgentHealthReqDTO{}
	}
	nodeRegionDB[nodeinfo.Region][nodeinfo.NodeID] = nodeinfo
	return nil
}

func getAllNodesInRegion(region string) ([]*model.ImageCacheAgentHealthReqDTO, error) {
	nodeRegionDBLock.RLock()
	defer nodeRegionDBLock.RUnlock()
	var res []*model.ImageCacheAgentHealthReqDTO
	nodes := nodeRegionDB[region]
	for _, v := range nodes {
		res = append(res, v)
	}
	return res, nil
}
