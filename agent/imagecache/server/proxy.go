// 用于代理 image 请求的所有 handler
package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/schwarzeni/k8senhance/agent/imagecache/cache"
	"github.com/schwarzeni/k8senhance/pkg/metrics"
	"github.com/schwarzeni/k8senhance/pkg/model"
)

// TODO: 未实现：缓存 manifest
func HandleProxy(server *Server) {
	r := server.r
	remoteRegister := server.conf.Agent.Imagecache.RemoteRegistry
	cacheFolder := server.conf.Agent.Imagecache.CachePath
	_ = os.MkdirAll(cacheFolder, os.ModePerm)
	r.HandleFunc("/v2/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)
	r.PathPrefix("/v2/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//log.Println("[info]", r.RequestURI)
		urlItems := strings.Split(r.RequestURI, "/")
		imageName := urlItems[2] + "/" + urlItems[3]
		sourceType := urlItems[4]
		sourceID := urlItems[5]
		_ = sourceID
		_ = sourceType

		if sourceType == "manifests" {
			_ = withDockerhubPullAuth(r, imageName)
			tmpResp, _ := doGetProxy(remoteRegister, r)
			for k, vv := range tmpResp.Header {
				for _, v := range vv {
					w.Header().Add(k, v)
					log.Println("[debug] manifest header: ", k, v)
				}
			}
			rawData, _ := ioutil.ReadAll(tmpResp.Body)
			tmpResp.Body.Close()
			//log.Printf("[debug] %s, manifest %s: %s\n", imageName, sourceID, string(rawData))
			_ = cache.SetImageManifest(imageName, sourceID, rawData)
			_ = cache.ParseAndSetLayersInfo(rawData)
			if _, err := bytes.NewReader(rawData).WriteTo(w); err != nil {
				//panic(err)
			}
			return
		}

		var writer io.Writer = w
		// 双写，在本地缓存一份
		if sourceType == "blobs" {
			layerID := sourceID[len("sha256:"):]
			cacheFile := path.Join(cacheFolder, layerID)
			f, err := os.OpenFile(cacheFile, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0600)
			if err != nil {
				log.Printf("[err] create file %s: %v\n", cacheFile, err)
			}
			defer f.Close()
			writer = io.MultiWriter(w, f)
		}
		if sourceType == "blobs" {
			layerID := sourceID[len("sha256:"):]
			responses := queryForLayer(server.conf.NodeName, layerID)
			if idx := prioritizeNodes(layerID, responses); idx >= 0 {
				targetNode := responses[idx]
				resp, err := layerdl(targetNode.Metric.NodeInfo.IP, layerID)
				if err == nil {
					defer resp.Body.Close()
					log.Printf("[info] dl layer %s from %s\n", layerID, targetNode.Metric.NodeInfo.IP)
					_ = withDockerhubPullAuth(r, imageName)
					// TODO: need get the header, 否则无效，不清楚哪些 header 起作用
					tmpResp, _ := doGetProxy(remoteRegister, r)
					for k, vv := range tmpResp.Header {
						for _, v := range vv {
							w.Header().Add(k, v)
							log.Println("[debug] blobs header: ", k, v)
						}
					}
					tmpResp.Body.Close()
					if _, err = bufio.NewReader(resp.Body).WriteTo(writer); err != nil {
						//panic(err)
					}
					return
				}
				log.Printf("[info] dl layer %s from %s has err: %v, using dockerhub instead\n", layerID, targetNode.Metric.NodeInfo.IP, err)
			} else {
				log.Printf("[info] prioritize node result -1, using dockerhub instead")
			}
		}

		if err := withDockerhubPullAuth(r, imageName); err != nil {
			log.Println("[err] fetch token for ", imageName, err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		resp, err := doGetProxy(remoteRegister, r)
		if err != nil {
			log.Println("[err]", r.RequestURI, err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()
		// !!! header 需要 copy
		for k, vv := range resp.Header {
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}
		if _, err = bufio.NewReader(resp.Body).WriteTo(writer); err != nil {
			//panic(err)
		}
	}).Methods(http.MethodGet)
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

func layerdl(targetNode string, layerID string) (*http.Response, error) {
	return http.Get("http://" + targetNode + ":8888/agentapi/v1/layerdl/" + layerID)
}

var (
	authCache     = map[string]string{}
	authCacheLock sync.RWMutex
)

func withDockerhubPullAuth(req *http.Request, imageName string) (err error) {
	// TODO: 后期缓存可以采用 cache + singleflight 优化
	authCacheLock.Lock()
	defer authCacheLock.Unlock()
	token, ok := authCache[imageName]
	if !ok {
		log.Println("[debug] request token for ", imageName)
		resp, err := http.Get(fmt.Sprintf("https://auth.docker.io/token?service=registry.docker.io&scope=repository:%s:pull", imageName))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		respObj := dockerhubAuthObj{}
		respData, _ := ioutil.ReadAll(resp.Body)
		if err = json.Unmarshal(respData, &respObj); err != nil {
			return err
		}
		token = respObj.Token
		authCache[imageName] = respObj.Token
	}
	req.Header.Add("Authorization", "Bearer "+token)
	return nil
}

type dockerhubAuthObj struct {
	Token string `json:"token"`
}

func doGetProxy(remoteAddr string, rawReq *http.Request) (*http.Response, error) {
	targetURL := rawReq.URL.String()
	// TODO: set this proxy to be configtable
	//os.Setenv("HTTP_PROXY", "http://10.211.55.2:7890")
	//os.Setenv("HTTPS_PROXY", "http://10.211.55.2:7890")
	req, err := http.NewRequest(http.MethodGet, remoteAddr+targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("newRequest for %s: %v", targetURL, err)
	}
	for k, vv := range rawReq.Header {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}
	client := http.DefaultClient
	client.Timeout = time.Second * 100000
	return client.Do(req)
}
