package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 处理跨域请求,支持options访问
func BackendTokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("authorization")
		if token == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// TODO: NEED More security!!!!
		if token != "adminToken" {
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			c.Next()
		}
	}
}
