package http

import (
	"barista/internal/modules"
	"barista/pkg/errors"
	"barista/pkg/log"
	"barista/pkg/models"
	"context"
	"fmt"
	"net/http"
	"strconv"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRequest.Error().Error()})
		return
	}

	err = h.Handler.Create(ctx, &req)
	if err != nil {
		log.GetLog().Errorf("Unable to create cafe. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInternalError.Error().Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRequest.Error().Error()})
		return
	}

	cafes, err := h.Handler.SearchCafe(ctx, req.Name, req.Province, req.City, req.Category)
	if err != nil {
		log.GetLog().Errorf("Unable to search cafe. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInternalError.Error().Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cafes": cafes})

}

func (h Cafe) PublicCafeProfile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	cafeID := c.Query("cafe_id")
	cafe_id, err := strconv.Atoi(cafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to convert userID to int32. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRequest.Error().Error()})
		return
	}

	cafe, err := h.Handler.PublicCafeProfile(ctx, int32(cafe_id))
	if err != nil {
		log.GetLog().Errorf("Unable to get public cafe profile. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrBadRequest.Error().Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.GetLog().Errorf("Unable to get token ID.")
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRequest.Error().Error()})
		return
	}

	AddedComment, err := h.Handler.AddComment(ctx, req.CafeID, fmt.Sprintf("%v", userID), req.Comment)
	if err != nil {
		log.GetLog().Errorf("Unable to add comment. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comment": AddedComment})
	return
}

func (h Cafe) GetComments(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	CafeID := c.Query("cafe_id")
	Counter := c.Query("counter")

	cafe_id, err := strconv.Atoi(CafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to convert userID to int32. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRequest.Error().Error()})
		return
	}

	counter, err := strconv.Atoi(Counter)
	if err != nil {
		log.GetLog().Errorf("Unable to convert userID to int32. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRequest.Error().Error()})
		return
	}

	comments, err := h.Handler.GetComments(ctx, int32(cafe_id), counter)
	if err != nil {
		log.GetLog().Errorf("Unable to get comments. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
	return
}

func (h Cafe) CreateEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()
	var req models.Event

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.GetLog().Errorf("Unable to bind json. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to bind json"})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		log.GetLog().Errorf("Unable to get user role.")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUnableToGetUser.Error().Error()})
		return
	}

	if role.(int32) != 2 {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrForbidden.Error().Error()})
		return
	}

	err = h.Handler.CreateEvent(ctx, req)
	if err != nil {
		log.GetLog().Errorf("Unable to create event. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrBadRequest.Error().Error()})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		log.GetLog().Errorf("Unable to get user role")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUnableToGetUser.Error().Error()})
		return
	}

	if role.(int32) != 2 {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrForbidden.Error().Error()})
		return
	}

	item, err := h.Handler.AddMenuItem(ctx, &req)
	if err != nil {
		log.GetLog().Errorf("Unable to add menu item. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"item": item,
	})
	return
}

func (h Cafe) GetMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	CafeID := c.Query("cafe_id")

	cafe_id, err := strconv.Atoi(CafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to convert cafe id to int32. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRequest.Error().Error()})
		return
	}

	categories, menu, cafeName, cafeImage, err := h.Handler.GetMenu(ctx, int32(cafe_id))
	if err != nil {
		log.GetLog().Errorf("Unable to get menu. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
		"menu":       menu,
		"cafe_name":  cafeName,
		"cafe_image": cafeImage,
	})
	return
}

func (h Cafe) EditMenuItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	var req models.MenuItem

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.GetLog().WithError(err).Error("Unable to bind json")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrBadRequest.Error()})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		log.GetLog().Errorf("Unable to get user role")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUnableToGetUser.Error().Error()})
		return
	}

	if role.(int32) != 2 {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrForbidden.Error().Error()})
		return
	}

	err = h.Handler.EditMenuItem(ctx, req)
	if err != nil {
		log.GetLog().Errorf("Unable to edit menu item. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	return
}

func (h Cafe) DeleteMenuItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	Item := c.Query("item")
	itemID, err := strconv.Atoi(Item)
	if err != nil {
		log.GetLog().Errorf("Invalid item type. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRequest.Error().Error()})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		log.GetLog().Errorf("Unable to get user role")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUnableToGetUser.Error().Error()})
		return
	}

	if role.(int32) != 2 {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrForbidden.Error().Error()})
		return
	}

	err = h.Handler.DeleteMenuItem(ctx, int32(itemID))
	if err != nil {
		log.GetLog().Errorf("Unable to edit menu item. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	return
}

func (h Cafe) Home(c *gin.Context) {
	cafe, comments, err := h.Handler.Home(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cafes": cafe, "comments": comments})
}

type RequestReserveEvent struct {
	EventID int32 `json:"event_id"`
}

func (h Cafe) ReserveEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()
	var data RequestReserveEvent

	err := c.ShouldBindJSON(&data)
	if err != nil {
		log.GetLog().WithError(err).Error("Unable to bind json")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrBadRequest.Error().Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.GetLog().Errorf("Unable to get user id.")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUnableToGetUser.Error().Error()})
		return
	}

	err = h.Handler.ReserveEvent(ctx, data.EventID, userID.(int32))
	if err != nil {
		log.GetLog().Errorf("Unable to add menu item. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	return
}

func (h Cafe) PrivateCafe(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	cafeID := c.Query("cafe_id")
	cafe_id, err := strconv.Atoi(cafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to convert userID to int32. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRequest.Error().Error()})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		log.GetLog().Errorf("Unable to get user role")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUnableToGetUser.Error().Error()})
		return
	}

	if role.(int32) != 2 {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrForbidden.Error().Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.GetLog().Errorf("Unable to get user id")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUnableToGetUser.Error().Error()})
		return
	}

	cafe, err := h.Handler.CafeRepo.GetByID(ctx, int32(cafe_id))
	if err != nil {
		log.GetLog().Errorf("Unable to get user by id. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInternalError.Error().Error()})
		return
	}

	if cafe.OwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrForbidden.Error().Error()})
		return
	}

	privateCafe, err := h.Handler.PrivateCafe(ctx, *cafe)
	if err != nil {
		log.GetLog().Errorf("Unable to get public cafe profile. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cafe": privateCafe,
	})
	return
}

func (h Cafe) EditCafe(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	var data modules.RequestEditCafe

	err := c.ShouldBindJSON(&data)
	if err != nil {
		log.GetLog().WithError(err).Error("Unable to bind json")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrBadRequest.Error().Error()})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		log.GetLog().Errorf("Unable to get user role")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUnableToGetUser.Error().Error()})
		return
	}

	if role.(int32) != 2 {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrForbidden.Error().Error()})
		return
	}

	err = h.Handler.EditCafe(ctx, data)
	if err != nil {
		log.GetLog().Errorf("Unable to edit cafe. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	return
}
