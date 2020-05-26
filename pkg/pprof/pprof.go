package pprof

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/pprof"
)

const defeatPrefix = "/debug/pprof"

func getPrefix(prefixOptions ...string) string {
	prefix := defeatPrefix
	if len(prefixOptions) > 0 {
		prefix = prefixOptions[0]
	}

	return prefix
}

func Register(route *gin.Engine, prefixOptions ...string) {
	prefix := getPrefix(prefixOptions)
	router := route.Group(prefix)
	router.GET("/", pprofHandler(pprof.Index))
	router.GET("/cmdline", pprofHandler(pprof.Cmdline))
	router.GET("/profile", pprofHandler(pprof.Profile))
	router.GET("/symbol", pprofHandler(pprof.Symbol))
	router.POST("/symbol", pprofHandler(pprof.Symbol))
}

func pprofHandler(h http.HandlerFunc) gin.HandlerFunc {
	handler := http.HandlerFunc(h)
	return func(ctx *gin.Context) {
		handler.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
