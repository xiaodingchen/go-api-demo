package internal

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"test.local/internal/controllers"
	"test.local/internal/routes"
	"test.local/pkg/utils"
	"test.local/pkg/xredis"
)

func Init(g *gin.Engine) {
	controllers.Init()
	routes.InitRoutes(g)
	// 初始化redis
	_, err := xredis.NewClient(utils.RedisDefault)
	if err != nil{
		utils.Log().Error("init redis client err", zap.Error(err))
	}
}
