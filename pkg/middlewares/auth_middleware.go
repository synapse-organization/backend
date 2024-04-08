package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"barista/pkg/utils"
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

		claims, err := utils.ValidateToken(token[0])
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		if token[0] != "10" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)
		c.Next()
		return
	}
}