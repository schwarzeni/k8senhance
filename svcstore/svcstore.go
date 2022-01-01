package svcstore

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/schwarzeni/k8senhance/pkg/model"
	"github.com/schwarzeni/k8senhance/svcstore/db"

	"github.com/schwarzeni/k8senhance/config"

	"github.com/spf13/cobra"
)

var gcpath *string

var cmd = &cobra.Command{
	Use:   "svcstore",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		config := config.MustParse(*gcpath)
		log.Printf("service store start, listening at port %s", config.Cloud.ServiceStore.Addr)
		db.InitK8SClient(config.Cloud.ServiceStore.K8Sconfig)
		router := gin.Default()
		router.GET("/api/v1/service", func(ctx *gin.Context) {
			serviceName := ctx.Query("name")
			if len(serviceName) == 0 {
				ctx.Status(http.StatusBadRequest)
				return
			}
			svc, err := db.GetService(serviceName)
			if err != nil {
				log.Println("[err] GET query db", err)
				ctx.Status(http.StatusInternalServerError)
				return
			}
			ctx.JSON(http.StatusOK, svc)
		})
		router.PUT("/api/v1/service", func(ctx *gin.Context) {
			var service model.ServiceInfo
			if err := ctx.BindJSON(&service); err != nil {
				log.Println("[err] PUT bind json", err)
				ctx.Status(http.StatusInternalServerError)
				return
			}
			if err := db.PutService(&service); err != nil {
				log.Println("[err] PUT save to db", err)
				ctx.Status(http.StatusInternalServerError)
				return
			}
			ctx.Status(http.StatusOK)
		})
		router.DELETE("/api/v1/service", func(ctx *gin.Context) {
			serviceName := ctx.Query("name")
			if len(serviceName) == 0 {
				ctx.Status(http.StatusBadRequest)
				return
			}
			if err := db.DelService(serviceName); err != nil {
				log.Println("[err] Delete data in db", err)
				ctx.Status(http.StatusInternalServerError)
				return
			}
			ctx.Status(http.StatusOK)
		})
		if err := router.Run(config.Cloud.ServiceStore.Addr); err != nil {
			panic(err)
		}
	},
}

func InitCMD(configPath *string) *cobra.Command {
	gcpath = configPath
	return cmd
}
