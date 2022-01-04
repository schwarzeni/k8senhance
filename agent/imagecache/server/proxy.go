// 用于代理 image 请求的所有 handler
package server

import (
	"net/http"
	"os"

	"github.com/schwarzeni/k8senhance/agent/imagecache/proxy"
)

// TODO: 未实现：缓存 manifest
func HandleProxy(server *Server) {
	cacheFolder := server.conf.Agent.Imagecache.CachePath
	_ = os.MkdirAll(cacheFolder, os.ModePerm)
	server.r.HandleFunc("/v2/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)
	server.r.PathPrefix("/v2/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//log.Println("[info]", r.RequestURI)
		urlInfo := proxy.MustParseURL(r.RequestURI)
		if urlInfo.SourceType == "manifests" {
			proxy.Manifest(server.conf, urlInfo, r, w)
			return
		}
		if urlInfo.SourceType == "blobs" {
			proxy.Layer(server.conf, urlInfo, r, w)
			return
		}
	}).Methods(http.MethodGet)
}
