package internal

import (
	"github.com/gin-gonic/gin"
	"test.local/internal/controllers"
	"test.local/internal/routes"
)

func Init(g *gin.Engine) {
	controllers.Init()
	routes.InitRoutes(g)
}
