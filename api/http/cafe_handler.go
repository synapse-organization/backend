package http

import (
	"barista/internal/modules"
	"barista/pkg/log"
	"barista/pkg/models"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Cafe struct {
	Handler *modules.CafeHandler
}

func (h Cafe) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()
	var req models.Cafe

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.GetLog().Errorf("Unable to bind json. error: %v", err)
		c.JSON(400, gin.H{"error": "Unable to bind json"})
		return
	}

	err = h.Handler.Create(ctx, &req)
}

func (h Cafe) GetCafe(c *gin.Context) {

}

type RequestSearchCafe struct {
	Name     string `json:"name"`
	Province string `json:"province"`
	City     string `json:"city"`
	Category string `json:"category"`
}

func (h Cafe) SearchCafe(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()
	var req RequestSearchCafe

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.GetLog().Errorf("Unable to bind json. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to bind json"})
		return
	}

	cafes, err := h.Handler.SearchCafe(ctx, req.Name, req.Province, req.City, req.Category)
	if err != nil {
		log.GetLog().Errorf("Unable to search cafe. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to search cafe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cafes": cafes})

}
