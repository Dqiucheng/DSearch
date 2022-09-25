package router

import (
	"DSearch/controller"
	"DSearch/core"
	"DSearch/logger"
	"DSearch/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RunRouter 启动路由
func RunRouter(router *gin.Engine) {
	// 注册全局中间件
	globalMiddleware(router)

	// 注册路由
	register(router)
}

// globalMiddleware 注册全局中间件
func globalMiddleware(router *gin.Engine) {
	router.Use(core.Handle(middleware.RecoveryJSON)) // 自定义全局避免恐慌造成退出
	router.Use(core.Handle(middleware.Cors))         // 跨域处理

	logger.SysLog("全局中间件注册成功").Info("Server日志")
}

// register 注册路由
func register(router *gin.Engine) {
	// Html路由注册
	htmlRouter(router)

	// Api路由注册
	apiRouter(router)

	logger.SysLog("路由注册成功").Info("Server日志")
}

// htmlRouter Html路由注册
func htmlRouter(router *gin.Engine) {
	// 加载模板文件
	router.LoadHTMLGlob("template/view/**/*")
	// 加载静态文件
	router.Static("/static", "template/static")
	router.StaticFile("/favicon.ico", "template/favicon.ico")

	// 404 未知路由处理
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "error/404", gin.H{
			"errmsg": "抱歉，你访问的页面不存在",
		})
	})

	// 未知调用方式
	router.NoMethod(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "error/404", gin.H{
			"errmsg": "Not method",
		})
	})
}

// apiRouter Api路由注册
func apiRouter(router *gin.Engine) {
	apiGroup := router.Group("/api")
	apiGroup.Use(core.Handle(middleware.RequestLimit))
	{
		api := new(controller.Api)
		apiGroup.Any("/", core.Handle(api.Test))
	}
}
