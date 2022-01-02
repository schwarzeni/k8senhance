package imagecachecontroller

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/schwarzeni/k8senhance/config"
	"github.com/schwarzeni/k8senhance/pkg/model"
	"github.com/spf13/cobra"
)

var (
	globalConfigPath *string
	globalConfig     *config.Config
)

var cmd = &cobra.Command{
	Use:   "icc",
	Short: "cloud controller for image cache service",
	Long:  "cloud controller for image cache service",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("start image cache controller")
		globalConfig = config.MustParse(*globalConfigPath)
		// TODO: 接收信息上报，定期移除过期的节点
		r := mux.NewRouter()
		r.HandleFunc("/healthz/{nodeID}", func(resp http.ResponseWriter, req *http.Request) {
			nodeID := mux.Vars(req)["nodeID"]
			log.Printf("[debug] get health ping from %s", nodeID)
			data, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Println("[err] read body", err)
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			imageCacheAgentHealthReqDTO := &model.ImageCacheAgentHealthReqDTO{}
			if err := json.Unmarshal(data, imageCacheAgentHealthReqDTO); err != nil {
				log.Println("[err] parse req body", err)
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			_ = setNodeRegionDB(imageCacheAgentHealthReqDTO)
			allNodesInRegion, _ := getAllNodesInRegion(imageCacheAgentHealthReqDTO.Region)
			imageCacheAgentHealthRespDTO := &model.ImageCacheAgentHealthRespDTO{
				Region: imageCacheAgentHealthReqDTO.Region,
				Nodes:  allNodesInRegion,
			}
			respData, err := json.Marshal(imageCacheAgentHealthRespDTO)
			if err != nil {
				log.Println("[err] encode resp data", err)
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, _ = io.Copy(resp, bytes.NewReader(respData))
		}).Methods(http.MethodPost)
		srv := &http.Server{
			Handler:      r,
			Addr:         globalConfig.Cloud.ImageCacheController.Addr,
			WriteTimeout: 10000 * time.Second,
			ReadTimeout:  10000 * time.Second,
		}
		log.Fatal(srv.ListenAndServe())
	},
}

func InitCMD(configPath *string) *cobra.Command {
	globalConfigPath = configPath
	return cmd
}
