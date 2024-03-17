package http

import (
	"barista/internal/modules"
	"barista/pkg/log"
	"barista/pkg/models"
	"barista/pkg/utils"
	"github.com/gin-gonic/gin"
	"time"
	"context"
	"net/http"
)

type User struct {
	Handler *modules.UserHandler
}

func (u User) PostLogin(ctx *gin.Context) {
	var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	
	var user models.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		log.GetLog().Errorf("Unable to bind json. error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = u.Handler.Login(ctx, &user)
	defer cancel()
	if err != nil {
		errValue := err.Error()
		if !utils.IsCommonError(err) {
			log.GetLog().Errorf("Unable to login. error: %v", err)
			errValue = "Unable to login"
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errValue})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"users": "login"})
}

func (u User) PostSignUp(ctx *gin.Context) {
	var data models.User

	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		log.GetLog().Errorf("Unable to bind json. error: %v", err)
		ctx.JSON(400, gin.H{"error": "Unable to bind json"})
		return
	}

	err = u.Handler.SignUp(ctx, &data)
	if err != nil {
		errValue := err.Error()
		if !utils.IsCommonError(err) {
			log.GetLog().Errorf("Unable to sign up. error: %v", err)
			errValue = "Unable to sign up"
		}

		ctx.JSON(500, gin.H{"error": errValue})
		return
	}

	ctx.JSON(200, gin.H{"users": "signup"})
}
