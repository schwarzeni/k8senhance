package service

import (
	"github.com/gin-gonic/gin"
	"github.com/schwarzeni/k8senhance/appserver/config"
	"github.com/schwarzeni/k8senhance/appserver/model"
)

var WebsvcPage PageService = func(config *config.ServiceConfig) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		menuConfig := model.TemplateNavStatus{
			Menu:    DefaultMenuConfig.Menu,
			User:    DefaultMenuConfig.User,
			CurrURL: ctx.Param("action"),
		}
		if err := config.T.Execute(ctx.Writer, menuConfig); err != nil {
			panic(err)
		}
	}
}
