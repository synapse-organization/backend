package middlewares

import (
	"barista/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type AuthMiddleware struct {
	IgnorePaths []string
	Postgres    *pgx.Conn
}

func (a AuthMiddleware) IsAuthorized(c *gin.Context) {

	token := c.Request.Header["Authorization"]
	if token[0] == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Token is empty"})
		c.Abort()
		return
	}

	claims, err := utils.ValidateToken(a.Postgres, token[0])
	if err != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		c.Abort()
		return
	}

	c.Set("email", claims.Email)
	c.Next()
}
