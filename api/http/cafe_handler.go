package http

import (
	"barista/internal/modules"
	"barista/pkg/errors"
	"barista/pkg/log"
	"barista/pkg/models"
	"context"
	"fmt"
	"net/http"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRequest.Error()})
		return
	}

	err = h.Handler.Create(ctx, &req)
	if err != nil {
		log.GetLog().Errorf("Unable to create cafe. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInternalError.Error()})
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
		log.GetLog().WithError(err).Error("Unable to bind json")
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRequest.Error()})
		return
	}

	cafes, err := h.Handler.SearchCafe(ctx, req.Name, req.Province, req.City, req.Category)
	if err != nil {
		log.GetLog().Errorf("Unable to search cafe. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInternalError.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cafes": cafes})

}

type RequestPublicCafe struct {
	CafeID int32 `json:"cafe_id"`
}

func (h Cafe) PublicCafeProfile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	var req RequestPublicCafe

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.GetLog().WithError(err).Error("Unable to bind json")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrBadRequest.Error()})
		return
	}

	cafe, err := h.Handler.PublicCafeProfile(ctx, req.CafeID)
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
		log.GetLog().WithError(err).Error("Unable to bind json")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrBadRequest.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.GetLog().Errorf("Unable to get token ID.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to get token ID"})
		return
	}

	AddedComment, err := h.Handler.AddComment(ctx, req.CafeID, fmt.Sprintf("%v", userID), req.Comment)
	if err != nil {
		log.GetLog().Errorf("Unable to add comment. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add comment."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comment": AddedComment})
	return
}

type RequestGetComments struct {
	CafeID  int32 `json:"cafe_id"`
	Counter int   `json:"counter"`
}

func (h Cafe) GetComments(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	var req RequestGetComments

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.GetLog().WithError(err).Error("Unable to bind json")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrBadRequest.Error()})
		return
	}

	comments, err := h.Handler.GetComments(ctx, req.CafeID, req.Counter)
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

func (h Cafe) AddMenuItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()
	var req models.MenuItem

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.GetLog().WithError(err).Error("Unable to bind json")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrBadRequest.Error()})
		return
	}

	err = h.Handler.AddMenuItem(ctx, &req)
	if err != nil {
		log.GetLog().Errorf("Unable to add menu item. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	return
}
