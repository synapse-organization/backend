package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthMiddleware struct {
	IgnorePaths []string
}

func (a AuthMiddleware) IsAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		// path = c.Request.URL.Path

		token := c.Request.Header["Authorization"]
		if token[0] == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Token is empty"})
			c.Abort()
			return
		}

		// TODO: Check with jwt and db
		if token[0] != "10" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
		return
	}
}
