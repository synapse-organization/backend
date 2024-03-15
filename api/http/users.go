package http

import (
	"barista/internal/modules"
	"barista/pkg/log"
	"barista/pkg/models"
	"barista/pkg/utils"
	"github.com/gin-gonic/gin"
)

type User struct {
	Handler *modules.UserHandler
}

func (u User) GetSignIn(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"users": "sing in"})
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
