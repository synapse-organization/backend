package http

import (
	"barista/internal/modules"
	"barista/pkg/log"
	"barista/pkg/models"
	"barista/pkg/utils"
	"context"
	"net/http"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

var (
	TimeOut = 5 * time.Second
)

type User struct {
	Handler *modules.UserHandler
}

func (u User) SignUp(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	var data models.User

	err := c.ShouldBindJSON(&data)
	if err != nil {
		log.GetLog().Errorf("Unable to bind json. error: %v", err)
		c.JSON(400, gin.H{"error": "Unable to bind json"})
		return
	}

	err = u.Handler.SignUp(ctx, &data)
	if err != nil {
		if !utils.IsCommonError(err) {
			log.GetLog().Errorf("Unable to sign up. error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to sign up"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	token, err := u.Handler.Login(ctx, &user)
	if err != nil {
		errValue := err.Error()
		if !utils.IsCommonError(err) {
			log.GetLog().Errorf("Unable to login. error: %v", err)
			errValue = "Unable to login"
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": errValue})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"token":  token,
	})
	return
}

func (u User) GetUser(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
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
		log.GetLog().Errorf("Unable to process forget-password. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	return
}

func (u User) UserProfile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	userID, exists := c.Get("userID")
	if !exists {
		log.GetLog().Errorf("Unable to get token ID.")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get token ID"})
		return
	}

	user, err := u.Handler.UserProfile(ctx, fmt.Sprintf("%v", userID))
	if err != nil {
		log.GetLog().Errorf("Unable to get user profile. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "ok",
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"phone":      user.Phone,
		"sex":        user.Sex,
	})
	return
}

func (u User) EditProfile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	var newDetail models.User
	err := c.ShouldBindJSON(&newDetail)
	if err != nil {
		log.GetLog().Errorf("Unable to bind json. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.GetLog().Errorf("Unable to get token ID.")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get token ID"})
		return
	}

	err = u.Handler.EditProfile(ctx, &newDetail, fmt.Sprintf("%v", userID))
	if err != nil {
		log.GetLog().Errorf("Unable to edit profile. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
	return
}