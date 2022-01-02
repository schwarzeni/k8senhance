package cache

import "sync"

// 对本地缓存的 layer 信息做一个缓存
// 这样在查询的时候就不用系统调用了
// 不一定用得到，仅仅是一个简单的优化
// TODO: 可能需要定期清理缓存，LRU ？
var (
	layerInfoDB     = map[string]struct{}{}
	layerInfoDBLock sync.RWMutex
)

func SetLayerInfo(layerid string) error {
	layerInfoDBLock.Lock()
	defer layerInfoDBLock.Unlock()
	layerInfoDB[layerid] = struct{}{}
	return nil
}

func HasLayer(layerid string) (bool, error) {
	layerInfoDBLock.RLock()
	defer layerInfoDBLock.RUnlock()
	_, ok := layerInfoDB[layerid]
	return ok, nil
}
