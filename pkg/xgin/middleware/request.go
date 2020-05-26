package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

const requestID = "X-RequestId"

func Request() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := uuid.NewV4()
		c.Set(requestID, u.String())
		c.Writer.Header().Set(requestID, u.String())
	}
}
