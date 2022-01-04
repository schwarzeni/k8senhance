package proxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/schwarzeni/k8senhance/agent/imagecache/cache"
	"github.com/schwarzeni/k8senhance/agent/imagecache/job"
	"github.com/schwarzeni/k8senhance/config"
	"github.com/schwarzeni/k8senhance/pkg/model"
)

func Manifest(config *config.Config, info *URLInfo, rawReq *http.Request, rawResp http.ResponseWriter) {
	if !job.GetHealthJobInstance().Online() {
		//if !job.GetHealthJobInstance().Online() || config.NodeName == "edge1node1" { // 仅用作测试
		if fetchManifestInRegion(config, info, rawReq, rawResp) {
			return
		}
		log.Printf("unable to get manifest for %+v, try fetch from dockerhub", *info)
	}

	_ = withDockerhubPullAuth(rawReq, info.ImageName)
	proxyResp, _ := doGetProxy(config.Agent.Imagecache.RemoteRegistry, rawReq)
	copyManifestResponse(info, proxyResp, rawReq, rawResp)
	_ = proxyResp.Body.Close()
	return
}

func fetchManifestInRegion(config *config.Config, info *URLInfo, rawReq *http.Request, rawResp http.ResponseWriter) bool {
	count := 0
	resChan := make(chan *fetchManifestResult)
	cache.IterateNodeInfo(func(node *model.ImageCacheAgentHealthReqDTO) {
		if node.NodeID == config.NodeName {
			return
		}
		count++
		go func(node *model.ImageCacheAgentHealthReqDTO) {
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%s/agentapi/v1/manifest", node.IP, node.ImageCacheServerPort), nil)
			q := req.URL.Query()
			q.Add("name", info.ImageName)
			q.Add("version", info.SourceID)
			req.URL.RawQuery = q.Encode()
			resp, err := http.DefaultClient.Do(req)
			res := &fetchManifestResult{}
			if err != nil || resp.StatusCode == http.StatusNotFound {
				//log.Printf("[debug] [manifest cache] fetch manifest %s %s from %s in region %s, maybe error (%v) or 404 (%v)", info.ImageName, info.SourceID, node.NodeID, node.Region, err, resp)
				res.hasManifest = false
			} else {
				res.hasManifest = true
				res.resp = resp
			}
			// TODO: 这里可能存在res.resp 的 body 没有 close 的情况，暂时不管
			resChan <- res
		}(node)
	})
	for i := 0; i < count; i++ {
		res := <-resChan
		if !res.hasManifest {
			continue
		}
		log.Printf("[debug] [manifest cache] success use cache for %s %s", info.ImageName, info.SourceID)
		copyManifestResponse(info, res.resp, rawReq, rawResp)
		_ = res.resp.Body.Close()
		return true
	}
	log.Printf("[debug] [manifest cache] failed use cache for %s %s", info.ImageName, info.SourceID)
	return false
}

func copyManifestResponse(info *URLInfo, proxyResp *http.Response, rawReq *http.Request, rawResp http.ResponseWriter) {
	cache.SetHTTPHeaderCache(rawReq.RequestURI, proxyResp.Header)
	copyHeader(proxyResp.Header, rawResp.Header())
	rawData, _ := ioutil.ReadAll(proxyResp.Body)
	//log.Printf("[debug] %s, manifest %s: %s\n", imageName, sourceID, string(rawData))
	_ = cache.SetImageManifest(info.ImageName, info.SourceID, rawData)
	_ = cache.ParseAndSetLayersInfo(rawData)
	if _, err := bytes.NewReader(rawData).WriteTo(rawResp); err != nil {
		//panic(err)
	}
}

type fetchManifestResult struct {
	resp        *http.Response
	hasManifest bool
}
