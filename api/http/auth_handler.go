package http

import (
	"barista/internal/modules"
	"barista/pkg/log"
	"barista/pkg/models"
	"barista/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var (
	TimeOut = 5 * time.Second
)

type Auth struct {
	Handler *modules.UserHandler
}

func (u Auth) PostLogin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.GetLog().Errorf("Unable to bind json. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = u.Handler.Login(ctx, &user)
	if err != nil {
		errValue := err.Error()
		if !utils.IsCommonError(err) {
			log.GetLog().Errorf("Unable to login. error: %v", err)
			errValue = "Unable to login"
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": errValue})
		return
	}

	// TODO: return jwt token
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	return
}

func (u Auth) PostSignUp(c *gin.Context) {
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