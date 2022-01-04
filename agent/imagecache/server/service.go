// 用于对其他 agent 提供服务的 handler
package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/schwarzeni/k8senhance/agent/imagecache/proxy"

	"github.com/gorilla/mux"
	"github.com/schwarzeni/k8senhance/agent/imagecache/cache"
)

func HandleService(server *Server) {
	r := server.r
	cacheFolder := server.conf.Agent.Imagecache.CachePath

	r.HandleFunc("/agentapi/v1/layerquery/{layerid}", func(resp http.ResponseWriter, req *http.Request) {
		currMetric, _ := cache.NodeMetric()
		respData := proxy.QueryLayerRespDTO{
			Metric: *currMetric,
		}
		// TODO: 这里获取 layer 可以使用 cache
		layerid := mux.Vars(req)["layerid"]
		if _, err := os.Stat(path.Join(cacheFolder, layerid)); !os.IsNotExist(err) {
			respData.HasLayer = true
		}
		respDataBytes, _ := json.Marshal(&respData)
		_, _ = resp.Write(respDataBytes)
	}).Methods(http.MethodGet)

	r.HandleFunc("/agentapi/v1/layerdl/{layerid}", func(resp http.ResponseWriter, req *http.Request) {
		layerid := mux.Vars(req)["layerid"]
		if _, err := os.Stat(path.Join(cacheFolder, layerid)); os.IsNotExist(err) {
			resp.WriteHeader(http.StatusNotFound)
			return
		}
		f, err := os.Open(path.Join(cacheFolder, layerid))
		if err != nil {
			log.Printf("[err] failed to open %s: %v\n", path.Join(cacheFolder, layerid), err)
		}
		rawURL := req.URL.Query().Get("rawurl")
		if header, ok := cache.HTTPHeaderCache(rawURL); ok {
			for k, vv := range header {
				for _, v := range vv {
					resp.Header().Add(k, v)
				}
			}
		}
		defer f.Close()
		_, _ = io.Copy(resp, f)
	}).Methods(http.MethodGet)

	r.HandleFunc("/agentapi/v1/manifest", func(w http.ResponseWriter, r *http.Request) {
		imageName := r.URL.Query().Get("name")
		imageVersion := r.URL.Query().Get("version")
		data, ok, err := cache.ImageManifest(imageName, imageVersion)
		if err != nil {
			log.Printf("[err] get image manifest cache %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		rawURL := r.URL.Query().Get("rawurl")
		if header, ok := cache.HTTPHeaderCache(rawURL); ok {
			for k, vv := range header {
				for _, v := range vv {
					w.Header().Add(k, v)
				}
			}
		}
		_, _ = io.Copy(w, bytes.NewReader(data))
	}).Methods(http.MethodGet)
}
