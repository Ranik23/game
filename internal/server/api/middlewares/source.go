package middlewares

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)





func Ip() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		slog.Info("Request from", slog.String("IP", ip))
		c.Next()
	}
}