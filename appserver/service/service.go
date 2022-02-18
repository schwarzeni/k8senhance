package service

import (
	"github.com/gin-gonic/gin"
	"github.com/schwarzeni/k8senhance/appserver/config"
	"github.com/schwarzeni/k8senhance/appserver/model"
)

type PageService func(config *config.ServiceConfig) func(ctx *gin.Context)

var PageServices = map[string]PageService{
	"node":      NodePage,
	"container": ContainerPage,
	"scheduler": SchedulerPage,
	"websvc":    WebsvcPage,
}

var DefaultMenuConfig = model.TemplateNavStatus{
	Menu: []model.MenuItem{
		{
			ID:   "menu-node",
			Name: "节点信息",
			URL:  "node",
		},
		{
			ID:   "menu-websvc",
			Name: "网络服务",
			URL:  "websvc",
		},
		{
			ID:   "menu-container",
			Name: "镜像缓存使用率",
			URL:  "container",
		},
		{
			ID:   "menu-scheduler",
			Name: "调度算法规则",
			URL:  "scheduler",
		},
		{
			ID:   "menu-admin",
			Name: "人员管理",
			URL:  "admin",
		},
	},
	User: model.User{
		ID:   "1",
		Name: "admin",
	},
	CurrURL: "",
}
