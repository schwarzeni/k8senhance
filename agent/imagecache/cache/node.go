package cache

import (
	"sync"

	"github.com/schwarzeni/k8senhance/pkg/model"
)

// 用于缓存同一 region 中其它节点的信息
// 用于在和中心节点断开连接的时候依然可以获取数据
var (
	nodesInfoDB     []*model.ImageCacheAgentHealthReqDTO
	nodesInfoDBLock sync.RWMutex
)

func SetNodeInfos(info []*model.ImageCacheAgentHealthReqDTO) error {
	nodesInfoDBLock.Lock()
	defer nodesInfoDBLock.Unlock()
	nodesInfoDB = info
	return nil
}

func IterateNodeInfo(fn func(info *model.ImageCacheAgentHealthReqDTO)) {
	nodesInfoDBLock.RLock()
	defer nodesInfoDBLock.RUnlock()
	for _, v := range nodesInfoDB {
		fn(v)
	}
}
