package job

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/schwarzeni/k8senhance/agent/imagecache/cache"
	"github.com/schwarzeni/k8senhance/config"
	"github.com/schwarzeni/k8senhance/pkg/model"
)

// 定期上报节点状态
type HealthJob struct {
	gc *config.Config
	// TODO: 可能存在并发访问的风险，可以使用 atomic 来实现
	online bool
}

func (hj *HealthJob) Online() bool {
	return hj.online
}

func (hj *HealthJob) Run() {
	targetAddr := hj.gc.Agent.Imagecache.ControllerAddr
	targetAPI := targetAddr + "/healthz/" + hj.gc.NodeName
	log.Println("health check job start, target at", targetAddr)

	nodeInfo := &model.ImageCacheAgentHealthReqDTO{
		NodeID:               hj.gc.NodeName,
		IP:                   hj.gc.Agent.Imagecache.CurrentIP,
		Region:               hj.gc.Agent.Imagecache.Region,
		ImageCacheServerPort: extractPort(hj.gc.Agent.Imagecache.Addr),
	}
	data, err := json.Marshal(nodeInfo)
	if err != nil {
		panic("failed to generate ImageCacheAgentHealthReqDTO json byte " + err.Error())
	}
	for {
		time.Sleep(time.Second * 2)
		resp, err := http.Post(targetAPI, "application/json", bytes.NewReader(data))
		if err != nil {
			log.Println("[health check] failed post, maybe it is offline", err)
			hj.online = false
			continue
		}
		hj.online = true
		respRawData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("[health check] failed to read resp body", err)
			resp.Body.Close()
			continue
		}
		respData := &model.ImageCacheAgentHealthRespDTO{}
		if err := json.Unmarshal(respRawData, respData); err != nil {
			log.Println("[health check] failed to parse resp data", err)
			resp.Body.Close()
			continue
		}
		_ = cache.SetNodeInfos(respData.Nodes)
		resp.Body.Close()
	}
}

func NewHealthJob(gc *config.Config) *HealthJob {
	return &HealthJob{gc: gc}
}

func extractPort(url string) string {
	// TODO: 解析 port，可能存在一些边界情况，未验证，默认输入值合法
	if len(url) > 0 && url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}
	res := strings.Split(url, ":")
	if len(res) < 2 {
		return "80"
	}
	return res[len(res)-1]
}
