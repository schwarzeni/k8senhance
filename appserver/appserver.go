//package appserver
package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/schwarzeni/k8senhance/appserver/config"
	"github.com/schwarzeni/k8senhance/appserver/service"
)

var distFilePath = "/Users/nizhenyang/Desktop/论文 workspace/code/github.com/schwarzeni/k8senhance/appserver/dist"
var templateFilePath = "/Users/nizhenyang/Desktop/论文 workspace/code/github.com/schwarzeni/k8senhance/appserver/template"

func runServer() error {
	router := gin.Default()
	router.StaticFS("/dist", http.Dir(distFilePath))

	servicePageNames := []string{"container", "node", "scheduler", "websvc"}
	servicePageTemplates, err := initPageTemplate(servicePageNames)
	if err != nil {
		panic(err)
	}
	servicePages := make(map[string]func(ctx *gin.Context))
	for _, servicePageName := range servicePageNames {
		servicePages[servicePageName] = service.PageServices[servicePageName](&config.ServiceConfig{T: servicePageTemplates[servicePageName]})
	}
	router.GET("/service/:action", func(ctx *gin.Context) {
		targetPageService, ok := servicePages[ctx.Param("action")]
		if !ok {
			ctx.Status(http.StatusNotFound)
			return
		}
		targetPageService(ctx)
	})
	return router.Run(":8080")
}

func initPageTemplate(serviceNames []string) (map[string]*template.Template, error) {
	data := make(map[string]*template.Template)
	for _, pageFile := range serviceNames {
		templ, err := template.ParseFiles(
			path.Join(templateFilePath, pageFile+".html"),
			path.Join(templateFilePath, "nav.html"),
			path.Join(templateFilePath, "head.html"))
		if err != nil {
			return nil, fmt.Errorf("init %s failed: %v", pageFile, err)
		}
		data[pageFile] = templ
	}
	return data, nil
}

func main() {
	if err := runServer(); err != nil {
		panic(err)
	}
}
