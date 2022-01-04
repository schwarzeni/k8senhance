package proxy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/schwarzeni/k8senhance/agent/imagecache/cache"
	"github.com/schwarzeni/k8senhance/config"
	"github.com/schwarzeni/k8senhance/pkg/metrics"
	"github.com/schwarzeni/k8senhance/pkg/model"
)

func Layer(config *config.Config, info *URLInfo, rawReq *http.Request, rawResp http.ResponseWriter) {
	// 双写，在本地缓存一份
	layerID := info.SourceID[len("sha256:"):]
	cacheFile := path.Join(config.Agent.Imagecache.CachePath, layerID)
	f, err := os.OpenFile(cacheFile, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		log.Printf("[err] create file %s: %v\n", cacheFile, err)
	}
	defer f.Close()
	writer := io.MultiWriter(rawResp, f)
	responses := queryForLayer(config.NodeName, layerID)
	if idx := prioritizeNodes(layerID, responses); idx >= 0 {
		targetNode := responses[idx]
		resp, err := layerdl(targetNode.Metric.NodeInfo.IP, layerID, rawReq.RequestURI)
		if err == nil {
			defer resp.Body.Close()
			log.Printf("[info] dl layer %s from %s\n", layerID, targetNode.Metric.NodeInfo.IP)
			//_ = withDockerhubPullAuth(r, imageName)
			// need get the header, 否则无效，不清楚哪些 header 起作用
			//tmpResp, _ := doGetProxy(remoteRegister, r)
			//for k, vv := range tmpResp.Header {
			//	for _, v := range vv {
			//		w.Header().Add(k, v)
			//		//log.Println("[debug] blobs header: ", k, v)
			//	}
			//}
			//tmpResp.Body.Close()
			copyHeader(resp.Header, rawResp.Header())
			if _, err = bufio.NewReader(resp.Body).WriteTo(writer); err != nil {
				//panic(err)
			}
			return
		}
		log.Printf("[info] dl layer %s from %s has err: %v, using dockerhub instead\n", layerID, targetNode.Metric.NodeInfo.IP, err)
	} else {
		log.Printf("[info] prioritize node result -1, using dockerhub instead")
	}
	if err := withDockerhubPullAuth(rawReq, info.ImageName); err != nil {
		log.Println("[err] fetch token for ", info.ImageName, err)
		rawResp.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	resp, err := doGetProxy(config.Agent.Imagecache.RemoteRegistry, rawReq)
	if err != nil {
		log.Println("[err]", rawReq.RequestURI, err)
		rawResp.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	// !!! header 需要 copy
	copyHeader(resp.Header, rawResp.Header())
	cache.SetHTTPHeaderCache(rawReq.RequestURI, resp.Header)
	if _, err = bufio.NewReader(resp.Body).WriteTo(writer); err != nil {
		//panic(err)
	}
}

type QueryLayerRespDTO struct {
	HasLayer bool               `json:"has_layer"`
	Metric   metrics.NodeMetric `json:"metric"`
}

func queryForLayer(currNode string, layerID string) (responses []*QueryLayerRespDTO) {
	timeout := time.Millisecond * 100
	reschan := make(chan *QueryLayerRespDTO)
	client := http.Client{Timeout: timeout}
	count := 0
	cache.IterateNodeInfo(func(node *model.ImageCacheAgentHealthReqDTO) {
		if node.NodeID == currNode {
			return
		}
		count++
		go func(node *model.ImageCacheAgentHealthReqDTO) {
			apiURL := fmt.Sprintf("http://%s:%s/agentapi/v1/layerquery/%s", node.IP, node.ImageCacheServerPort, layerID)
			resp, err := client.Get(apiURL)
			if err != nil {
				log.Println("[info] err, access url ", apiURL, err)
				// just ignore
				reschan <- nil
				return
			}
			defer resp.Body.Close()
			respData := QueryLayerRespDTO{}
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("[info] err, read resp from ", apiURL, err)
				// just ignore it
				reschan <- nil
				return
			}
			_ = json.Unmarshal(bodyBytes, &respData)
			// TODO: 重构一下节点数据是由那个结构体记录的，from cache.metric or response or cache.nodeinfos
			respData.Metric.NodeInfo = metrics.NodeInfo{
				ID:              node.NodeID,
				IP:              node.IP,
				CacheServerPort: node.ImageCacheServerPort,
				Region:          node.Region,
			}
			reschan <- &respData
		}(node)
	})
	for i := 0; i < count; i++ {
		res := <-reschan
		if res != nil && res.HasLayer {
			responses = append(responses, res)
		}
	}
	return responses
}

func layerdl(targetNode string, layerID string, originURL string) (*http.Response, error) {
	return http.Get("http://" + targetNode + ":8888/agentapi/v1/layerdl/" + layerID + "?rawurl=" + originURL)
}
