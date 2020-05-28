package internal

import (
	"github.com/gin-gonic/gin"
	"test.local/internal/controllers"
	"test.local/internal/routes"
	"test.local/pkg/utils"
	"test.local/pkg/xredis"
)

func InitApi(g *gin.Engine) error {
	// 初始化redis
	var err error
	err = initRedis()
	if err != nil{
		return err
	}
	// 初始化数据库
	err = initDB()
	if err != nil{
		return err
	}
	controllers.Init()
	routes.InitRoutes(g)
	return nil
}


func initRedis() error {
	_, err := xredis.NewClient(utils.RedisDefault)
	return err
}

func initDB()error{
	return nil
}
