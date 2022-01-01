package nodescheduler

import (
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/schwarzeni/k8senhance/config"
	"github.com/schwarzeni/k8senhance/pkg/metrics"
	"github.com/spf13/cobra"
)

var (
	gc     *config.Config
	gcpath *string
)

var cmd = &cobra.Command{
	Use:   "nodescheduler",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		gc = config.MustParse(*gcpath)
		addr := gc.Cloud.NodeScheduler.Addr
		log.Println("nodescheduler start, listening at addr ", addr)
		if err := service(); err != nil {
			panic(err)
		}
	},
}

func InitCMD(configPath *string) *cobra.Command {
	gcpath = configPath
	return cmd
}

func service() error {
	r := gin.New()
	ch := make(chan *metrics.NodeMetric)
	r.PUT("/api/v1/agenthealth/:nodeid", func(c *gin.Context) {
		nodeID := c.Param("nodeid")
		rawMetric := &metrics.NodeMetric{}
		if err := c.BindJSON(rawMetric); err != nil {
			log.Println("[err] parse json:", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		log.Println("got metric upload", nodeID, *rawMetric)
		ch <- rawMetric
	})
	r.POST("/api/v1/k8sextension/prioritize", func(c *gin.Context) {
		log.Println("[debug] access priority")
		priorityFunc(c)
	})
	r.PUT("/api/v1/processor/:id/:weight", func(c *gin.Context) {
		processorID := c.Param("id")
		newWeight := c.Param("weight")

		w, err := strconv.Atoi(newWeight)
		id, err2 := strconv.Atoi(processorID)
		_, ok := ProcessorMap[ProcessorType(id)]
		if err != nil || err2 != nil || w < 0 || w > math.MaxInt32 || !ok {
			c.Status(http.StatusBadRequest)
			return
		}
		ProcessorMap[ProcessorType(id)].ExtraWeight(int32(w))
	})
	go func() {
		for metric := range ch {
			processdata(metric)
		}
	}()
	return r.Run(":8080")
}
