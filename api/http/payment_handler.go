package http

import (
	"barista/internal/modules"
	"barista/pkg/errors"
	"barista/pkg/log"
	"barista/pkg/models"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
)

type Payment struct {
	Handler *modules.PaymentHandler
}

func (h Payment) Transfer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()
	var req models.RequestTransfer

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.GetLog().Errorf("Unable to bind json. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRequest.Error().Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.GetLog().Errorf("Unable to get token ID.")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUnableToGetUser.Error()})
		return
	}

	err = h.Handler.Transfer(ctx, cast.ToInt32(userID), &req)
	if err != nil {
		log.GetLog().Errorf("Unable to create cafe. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInternalError.Error().Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})

}

func (h Payment) Deposit(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()
	var req models.RequestDeposit

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.GetLog().Errorf("Unable to bind json. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRequest.Error().Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.GetLog().Errorf("Unable to get token ID.")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUnableToGetUser.Error()})
		return
	}

	err = h.Handler.Deposit(ctx, cast.ToInt32(userID), &req)
	if err != nil {
		log.GetLog().Errorf("Unable to create cafe. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInternalError.Error().Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})

}

func (h Payment) Withdraw(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()
	var req models.RequestWithdraw

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.GetLog().Errorf("Unable to bind json. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRequest.Error().Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.GetLog().Errorf("Unable to get token ID.")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUnableToGetUser.Error()})
		return
	}

	err = h.Handler.Withdraw(ctx, cast.ToInt32(userID), &req)
	if err != nil {
		log.GetLog().Errorf("Unable to create cafe. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInternalError.Error().Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h Payment) Balance(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	userID, exists := c.Get("userID")
	if !exists {
		log.GetLog().Errorf("Unable to get token ID.")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUnableToGetUser.Error()})
		return
	}

	balance := h.Handler.Balance(ctx, cast.ToInt32(userID))

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}
