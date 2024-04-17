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
		errValue := err.Error()
		if !utils.IsCommonError(err) {
			log.GetLog().Errorf("Unable to sign up. error: %v", err)
			errValue = "Unable to sign up"
		}

		c.JSON(500, gin.H{"error": errValue})
		return
	}

	c.JSON(200, gin.H{"status": "ok"})
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
		errValue := err.Error()
		if !utils.IsCommonError(err) {
			log.GetLog().Errorf("Unable to login. error: %v", err)
			errValue = "Unable to login"
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": errValue})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "ok",
		"token":        token,
		"refreshToken": refreshToken,
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

	incoming_token := c.Request.Header["Authorization"][0]

	user, err := u.Handler.UserProfile(ctx, incoming_token)
	if err != nil {
		log.GetLog().Errorf("Unable to get user profile. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"first_name": user.FirstName,
		"last_name": user.LastName,
		"email": user.Email,
		"phone": user.Phone,
		"sex": user.Sex,
	})
	return
}