package proxy

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/schwarzeni/k8senhance/agent/imagecache/cache"
	"github.com/schwarzeni/k8senhance/config"
)

func Manifest(config *config.Config, info *URLInfo, rawReq *http.Request, rawResp http.ResponseWriter) {
	_ = withDockerhubPullAuth(rawReq, info.ImageName)
	tmpResp, _ := doGetProxy(config.Agent.Imagecache.RemoteRegistry, rawReq)
	cache.SetHTTPHeaderCache(rawReq.RequestURI, tmpResp.Header)
	copyHeader(tmpResp.Header, rawResp.Header())
	rawData, _ := ioutil.ReadAll(tmpResp.Body)
	tmpResp.Body.Close()
	//log.Printf("[debug] %s, manifest %s: %s\n", imageName, sourceID, string(rawData))
	_ = cache.SetImageManifest(info.ImageName, info.SourceID, rawData)
	_ = cache.ParseAndSetLayersInfo(rawData)
	if _, err := bytes.NewReader(rawData).WriteTo(rawResp); err != nil {
		//panic(err)
	}
	return
}
