package cache

import (
	"encoding/json"
	"log"
	"strings"
	"sync"
)

// 对本地缓存的 layer 信息做一个缓存
// 这样在查询的时候就不用系统调用了
// 不一定用得到，仅仅是一个简单的优化
// TODO: 可能需要定期清理缓存，LRU ？
var (
	layerInfoDB     = map[string]int{}
	layerInfoDBLock sync.RWMutex
)

type ManifestTypeStruct struct {
	MediaType string `json:"mediaType"`
}

type LayersManifest struct {
	Config struct {
		Digest    string `json:"digest"`
		MediaType string `json:"mediaType"`
		Size      int    `json:"size"`
	} `json:"config"`
	Layers []struct {
		Digest    string `json:"digest"`
		MediaType string `json:"mediaType"`
		Size      int    `json:"size"`
	} `json:"layers"`
	MediaType     string `json:"mediaType"`
	SchemaVersion int    `json:"schemaVersion"`
}

func ParseAndSetLayersInfo(rawData []byte) error {
	mts := &ManifestTypeStruct{}
	if err := json.Unmarshal(rawData, mts); err != nil {
		log.Println("----1", err)
		return err
	}
	if !strings.Contains(mts.MediaType, "vnd.docker.distribution.manifest.v2+json") {
		log.Println("----2", mts.MediaType)
		return nil
	}
	layersManifest := &LayersManifest{}
	if err := json.Unmarshal(rawData, layersManifest); err != nil {
		log.Println("-----3", err)
		return err
	}

	layerInfoDBLock.Lock()
	defer layerInfoDBLock.Unlock()
	layerInfoDB[layersManifest.Config.Digest[len("sha256:"):]] = layersManifest.Config.Size
	log.Println("[debug] layer cache update", layersManifest.Config.Digest, layersManifest.Config.Size)
	for _, layer := range layersManifest.Layers {
		layerInfoDB[layer.Digest[len("sha256:"):]] = layer.Size
		log.Println("[debug] layer cache update", layer.Digest, layer.Size)
	}
	return nil
}

func LayerInfo(layerid string) (int, bool, error) {
	layerInfoDBLock.RLock()
	defer layerInfoDBLock.RUnlock()
	size, ok := layerInfoDB[layerid]
	return size, ok, nil
}
