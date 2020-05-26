package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct {
}

func NewUser() *User {
	return &User{}
}

func (u *User) Index(ctx *gin.Context) {
	ctx.String(http.StatusOK, "/user/index")
}
