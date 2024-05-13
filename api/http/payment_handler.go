package http

import (
	"barista/internal/modules"
	"barista/pkg/errors"
	"barista/pkg/log"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
)

type Payment struct {
	Handler *modules.PaymentHandler
}

type RequestTransfer struct {
	To     int32 `json:"to"`
	Amount int64 `json:"amount"`
}

func (h Payment) Transfer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()
	var req RequestTransfer

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

type RequestDeposit struct {
	Amount int64 `json:"amount"`
}

func (h Payment) Deposit(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()
	var req RequestDeposit

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

type RequestWithdraw struct {
	To     int32 `json:"to"`
	Amount int64 `json:"amount"`
}

func (h Payment) Withdraw(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()
	var req RequestWithdraw

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
