package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xiaodingchen/golibs/xpprof"
	"net/http"
	"test.local/internal/controllers"
)

func InitRoutes(g *gin.Engine) {
	g.Any("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})

	g.GET("/metrics", func(ctx *gin.Context) {
		promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
	})
	xpprof.GinRegister(g)
	// 处理业务路由
	g.GET("/user/index", controllers.Ctrl.User.Index)
}
