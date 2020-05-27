package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"test.local/pkg/utils"
	"test.local/pkg/xtrace"
	"time"
)

type User struct {
}

func NewUser() *User {
	return &User{}
}

func (u *User) Index(ctx *gin.Context) {
	ctx.String(http.StatusOK, "/user/index")
	client := utils.Redis(utils.DefaultRedis)
	if client == nil{
		utils.Log().Error("redis client nil")
		return
	}
	client = xtrace.WrapperRedis(ctx, client)
	pingCmd := client.Ping()
	utils.Log().Debug(
		"redis ping",
		zap.Any("val", pingCmd.Val()),
		zap.Error(pingCmd.Err()),
	)
	key := "test:redis:key"
	cmd := client.Set(key, 1, 10 * time.Hour)
	utils.Log().Debug(
		"redis set",
		zap.Any("val", cmd.Val()),
		zap.Error(cmd.Err()),
	)

	stringCmd := client.Get(key)
	utils.Log().Debug(
		"redis get",
		zap.Any("val", stringCmd.Val()),
		zap.Error(stringCmd.Err()),
	)

	hashCmd := client.HGetAll(key)
	utils.Log().Debug(
		"redis hgetall",
		zap.Any("val", hashCmd.Val()),
		zap.Error(hashCmd.Err()),
	)
}
