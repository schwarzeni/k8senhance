package imagecachecontroller

import (
	"sync"
	"time"

	"github.com/schwarzeni/k8senhance/pkg/model"
)

var (
	// TODO: 优化此处的数据结构
	nodeRegionDB     = map[string]map[string]NodeRegionInfo{}
	nodeRegionDBLock sync.RWMutex
	timeoutDuration  = time.Second * 10
)

type NodeRegionInfo struct {
	data       *model.ImageCacheAgentHealthReqDTO
	lastUpdate time.Time
}

func setNodeRegionDB(nodeinfo *model.ImageCacheAgentHealthReqDTO) error {
	nodeRegionDBLock.Lock()
	defer nodeRegionDBLock.Unlock()
	if _, ok := nodeRegionDB[nodeinfo.Region]; !ok {
		nodeRegionDB[nodeinfo.Region] = map[string]NodeRegionInfo{}
	}
	nodeRegionDB[nodeinfo.Region][nodeinfo.NodeID] = NodeRegionInfo{data: nodeinfo, lastUpdate: time.Now()}
	return nil
}

func getAllNodesInRegion(region string) ([]*model.ImageCacheAgentHealthReqDTO, error) {
	nodeRegionDBLock.RLock()
	defer nodeRegionDBLock.RUnlock()
	var res []*model.ImageCacheAgentHealthReqDTO
	nodes := nodeRegionDB[region]
	for key, v := range nodes {
		if time.Now().Sub(v.lastUpdate) < timeoutDuration {
			res = append(res, v.data)
		} else {
			delete(nodes, key)
		}
	}
	return res, nil
}
