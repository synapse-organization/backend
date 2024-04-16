package http

import (
	"barista/internal/modules"
	"barista/pkg/log"
	"barista/pkg/models"
	"barista/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
	"time"
)

var (
	TimeOut = 5 * time.Second
)

type User struct {
	Handler *modules.UserHandler
}

func (u User) ForgetPassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.GetLog().Errorf("Unable to bind json. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = u.Handler.ForgetPassword(ctx, &user)
	if err != nil {
		log.GetLog().Errorf("Unable to process forgetpassword. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	return
}

func (u User) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.GetLog().Errorf("Unable to bind json. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	token, refreshToken, err := u.Handler.Login(ctx, &user)
	if err != nil {
		if !utils.IsCommonError(err) {
			log.GetLog().Errorf("Unable to sign up. error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to sign up"})
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "ok",
		"token":        token,
		"refreshToken": refreshToken,
	})
	return
}

func (u User) SignUp(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	var data models.User

	err := c.ShouldBindJSON(&data)
	if err != nil {
		log.GetLog().Errorf("Unable to bind json. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to bind json"})
		return
	}

	err = u.Handler.SignUp(ctx, &data)
	if err != nil {
		if !utils.IsCommonError(err) {
			log.GetLog().Errorf("Unable to sign up. error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to sign up"})
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "ok"})
	return
}

func (u User) VerifyEmail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	email := c.Query("c")
	callback := c.Query("callback")

	err := u.Handler.VerifyEmail(ctx, cast.ToString(email))
	if err != nil {
		log.GetLog().Errorf("Unable to verify email. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.Redirect(http.StatusFound, callback)
	return
}

func (u User) GetUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
