package middlewares

import (
	"barista/pkg/log"
	"barista/pkg/utils"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	IgnorePaths []string
	Postgres    *pgxpool.Pool
}

func (a AuthMiddleware) IsAuthorized(c *gin.Context) {

	token := c.Request.Header["Authorization"]
	if len(token) == 0 || token[0] == "" {
		log.GetLog().Errorf("Token is empty")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Token is empty"})
		c.Abort()
		return
	}

	claims, err := utils.ValidateToken(a.Postgres, token[0])
	if err != "" {
		log.GetLog().Errorf("Token is invalid. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		c.Abort()
		return
	}

	c.Set("userID", claims.Uid)
	c.Set("email", claims.Email)
	c.Set("role", claims.Role)
	c.Set("tokenID", claims.TokenID)
	c.Set("claims", claims)
	c.Next()
}
