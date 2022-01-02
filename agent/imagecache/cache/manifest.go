package cache

import "sync"

// 对镜像的 manifest 文件进行一个缓存
// 用于在断网的时候其它节点可以正常获取镜像
var (
	manifestDB     = map[string]map[string][]byte{}
	manifestDBLock sync.RWMutex
)

func SetImageManifest(imageName, imageVersion string, data []byte) error {
	manifestDBLock.Lock()
	defer manifestDBLock.Unlock()
	if _, ok := manifestDB[imageName]; !ok {
		manifestDB[imageName] = make(map[string][]byte)
	}
	manifestDB[imageName][imageVersion] = data
	return nil
}

func ImageManifest(imageName, imageVersion string) ([]byte, bool, error) {
	manifestDBLock.RLock()
	defer manifestDBLock.RUnlock()
	if _, ok := manifestDB[imageName]; !ok {
		return nil, false, nil
	}
	data, ok := manifestDB[imageName][imageVersion]
	return data, ok, nil
}
