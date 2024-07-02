package http

import (
	"barista/internal/modules"
	"barista/pkg/errors"
	"barista/pkg/log"
	"barista/pkg/models"
	"barista/pkg/utils"
	"context"

	"fmt"
	"net/http"
	"time"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRequest.Error()})
		return
	}

	err = u.Handler.SignUp(ctx, &data)
	if err != nil {
		if !utils.IsCommonError(err) {
			log.GetLog().Errorf("Unable to sign up. error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInternalError.Error()})
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
		log.GetLog().WithError(err).Error("Unable to sign up")
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRequest.Error()})
		return
	}

	tokens, isCompleted, err := u.Handler.Login(ctx, &user)
	if err != nil {
		if !utils.IsCommonError(err) {
			log.GetLog().WithError(err).Error("Unable to sign up")
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInternalError.Error()})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "ok",
		"token":        tokens,
		"is_completed": isCompleted,
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRequest, "message": err})
		return
	}

	err = u.Handler.ForgetPassword(ctx, &user)
	if err != nil {
		log.GetLog().Errorf("Unable to process forget-password. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUnableToGetUser.Error()})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		log.GetLog().Errorf("Unable to get user role.")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUnableToGetUser.Error()})

	}

	user, err := u.Handler.UserProfile(ctx, fmt.Sprintf("%v", userID))
	if err != nil {
		log.GetLog().Errorf("Unable to get user profile. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if role.(int32) == 2 {
		c.JSON(http.StatusOK, gin.H{
			"status":       "ok",
			"first_name":   user.FirstName,
			"last_name":    user.LastName,
			"email":        user.Email,
			"phone":        user.Phone,
			"sex":          user.Sex,
			"bank_account": user.BankAccount,
			"national_id":  user.NationalID,
		})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.GetLog().Errorf("Unable to get token ID.")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUnableToGetUser.Error()})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		log.GetLog().Errorf("Unable to get user role.")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUnableToGetUser.Error()})

	}

	if role.(int32) != 2 {
		newDetail.NationalID = ""
		newDetail.BankAccount = ""
	}

	err = u.Handler.EditProfile(ctx, &newDetail, fmt.Sprintf("%v", userID))
	if err != nil {
		log.GetLog().Errorf("Unable to edit profile. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
	return
}

type RequestChangePassword struct {
	CurrentPassword string `json:"current_password"`
	Password        string `json:"password"`
	Password2       string `json:"password2"`
}

func (u User) ChangePassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	var data RequestChangePassword
	err := c.ShouldBindJSON(&data)
	if err != nil {
		log.GetLog().Errorf("Unable to bind json. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRequest.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.GetLog().Errorf("Unable to get token ID.")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUnableToGetUser.Error()})
		return
	}

	if data.Password != data.Password2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrPasswordNotMatch.Error()})
		return
	}

	if !utils.CheckPasswordValidity(data.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrPasswordIncorrect.Error()})
		return
	}

	err = u.Handler.ChangePassword(ctx, userID.(int32), data.Password, data.CurrentPassword)
	if err != nil {
		log.GetLog().Errorf("Unable to change password. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
	return

}

func (u User) Logout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	token := c.Request.Header["Authorization"]
	if len(token) == 0 || token[0] == "" {
		log.GetLog().Error("Token is empty")
		c.JSON(http.StatusBadRequest, gin.H{"message": errors.ErrDidntLogin.Error()})
		return
	}

	err := u.Handler.Logout(ctx, token[0])
	if err != nil {
		log.GetLog().Errorf("Unable to logout. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	return

}

type RequestManagerAgreement struct {
	NationalID  string `json:"national_id"`
	BankAccount string `json:"bank_account"`
}

func (u User) ManagerAgreement(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	var data RequestManagerAgreement
	err := c.ShouldBindJSON(&data)
	if err != nil {
		log.GetLog().Errorf("Unable to bind json. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRequest.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.GetLog().Errorf("Unable to get token ID.")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUnableToGetUser.Error()})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		log.GetLog().Errorf("Unable to get user role.")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUnableToGetUser.Error()})
		return
	}

	if role.(int32) != 2 {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrForbidden.Error()})
		return
	}

	err = u.Handler.ManagerAgreement(ctx, userID.(int32), data.NationalID, data.BankAccount)
	if err != nil {
		log.GetLog().Errorf("Unable to manager agreement. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	return
}

func (u User) UserReservations(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, TimeOut)
	defer cancel()

	userID, exists := c.Get("userID")
	if !exists {
		log.GetLog().Errorf("Unable to get token ID.")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUnableToGetUser.Error()})
		return
	}

	dayStr := c.Query("day")
	day, err := time.Parse("2006-01-02", dayStr)
	if err != nil {
		log.GetLog().Errorf("Unable to parse day. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid day format"})
		return
	}

	userReservations, err := u.Handler.UserReservations(ctx, userID.(int32), day)
	if err != nil {
		log.GetLog().Errorf("Unable to get user reservations. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reservations": userReservations})
	return
}