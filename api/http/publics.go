package http

import (
	"barista/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
)

type PublicHandler struct {
}

func (h PublicHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func (h PublicHandler) GetCities(c *gin.Context) {
	id := c.Query("id")
	val, ok := models.Cities[cast.ToInt(id)]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cities": val,
	})
}
