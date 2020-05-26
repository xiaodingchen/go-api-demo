package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

const requestID = "requestId"

func Request() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := uuid.NewV4()
		c.Set(requestID, u.String())
	}
}
