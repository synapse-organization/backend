package modules

import (
	"barista/pkg/errors"
	"barista/pkg/log"
	"barista/pkg/models"
	"barista/pkg/repo"
	"barista/pkg/utils"
	"github.com/gin-gonic/gin"
	"math/rand"
)

type UserHandler struct {
	UserRepo repo.UsersRepo
}

func (u UserHandler) SignIn(ctx *gin.Context, user *models.User) error {
	// c.JSON(200, gin.H{"users": "singin"})

	return nil
}

func (u UserHandler) SignUp(ctx *gin.Context, user *models.User) error {
	user.ID = rand.Int31()

	if !utils.CheckEmailValidity(user.Email) {
		return errors.ErrEmailInvalid.Error()
	}

	if !utils.CheckPhoneValidity(user.Phone) {
		return errors.ErrPhoneInvalid.Error()
	}

	if !utils.CheckPasswordValidity(user.Password) {
		return errors.ErrPasswordInvalid.Error()
	}

	if !utils.CheckNameValidity(user.FistName) {
		return errors.ErrFirstNameInvalid.Error()
	}

	if !utils.CheckNameValidity(user.LastName) {
		return errors.ErrLastNameInvalid.Error()
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.GetLog().Errorf("Unable to hash password. error: %v", err)
		return err
	}
	user.Password = hashedPassword

	err = u.UserRepo.Create(ctx, user)
	if err != nil {
		log.GetLog().Errorf("Unable to create user. error: %v", err)
		return err
	}

	return nil
}
