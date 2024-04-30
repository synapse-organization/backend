package http

import (
	"barista/internal/modules"
	"barista/pkg/log"
	"barista/pkg/models"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
	if err != nil {
		log.GetLog().Errorf("Unable to create cafe. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create cafe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
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

func (h Cafe) PublicCafeProfile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	cafeID := c.GetHeader(http.CanonicalHeaderKey("Cafe-ID"))
	if cafeID == "" {
		log.GetLog().Errorf("cafe id is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "cafe is is empty"})
		return
	}

	cafe, err := h.Handler.PublicCafeProfile(ctx, cafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to get public cafe profile. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"cafe":   cafe,
	})
	return
}

type RequestAddComment struct {
	CafeID  int32  `json:"cafe_id"`
	Comment string `json:"comment"`
}

func (h Cafe) AddComment(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()
	var req RequestAddComment

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.GetLog().Errorf("Unable to bind json. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to bind json."})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.GetLog().Errorf("Unable to get token ID.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to get token ID"})
		return
	}

	err = h.Handler.AddComment(ctx, req.CafeID, fmt.Sprintf("%v", userID), req.Comment)
	if err != nil {
		log.GetLog().Errorf("Unable to add comment. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add comment."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	return
}

func (h Cafe) GetComments(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	cafeID := c.GetHeader(http.CanonicalHeaderKey("X-Cafe-ID"))
	if cafeID == "" {
		log.GetLog().Errorf("cafe id is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "cafe id is empty"})
		return
	}

	cafe_id, err := strconv.Atoi(cafeID)
	if err != nil {
		log.GetLog().Errorf("Invalid cafe id. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cafe id"})
	}

	apiCallsCounter := c.GetHeader("X-Api-Calls-Counter")
	if apiCallsCounter == "" {
		log.GetLog().Errorf("API calls counter is missing")
		c.JSON(http.StatusBadRequest, gin.H{"error": "API calls counter is missing"})
		return
	}

	counter, err := strconv.Atoi(apiCallsCounter)
	if err != nil {
		log.GetLog().Errorf("Invalid API calls counter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid API calls counter"})
		return
	}

	comments, err := h.Handler.GetComments(ctx, int32(cafe_id), counter)
	if err != nil {
		log.GetLog().Errorf("Unable to get comments. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
	return
}

type RequestCreateEvent struct {
	CafeID      int32     `json:"cafe_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	ImageID     string    `json:"image_id"`
}

func (h Cafe) CreateEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()
	var req RequestCreateEvent

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.GetLog().Errorf("Unable to bind json. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to bind json"})
		return
	}

	err = h.Handler.CreateEvent(ctx, req.CafeID, req.Name, req.Description, req.StartTime, req.EndTime, req.ImageID)
	if err != nil {
		log.GetLog().Errorf("Unable to create event. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	return
}
