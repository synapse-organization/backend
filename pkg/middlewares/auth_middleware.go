package middlewares

import (
	"barista/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	IgnorePaths []string
	Postgres    *pgxpool.Pool
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

	c.Set("userID", claims.Uid)
	c.Set("email", claims.Email)
	c.Set("tokenID", claims.TokenID)
	c.Set("claims", claims)
	c.Next()
}
