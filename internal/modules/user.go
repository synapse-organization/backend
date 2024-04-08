package modules

import (
	"barista/pkg/errors"
	"barista/pkg/log"
	"barista/pkg/models"
	"barista/pkg/repo"
	"barista/pkg/utils"
	"context"
	"math/rand"

	"github.com/jackc/pgx/v5"
)

type UserHandler struct {
	UserRepo repo.UsersRepo
	Postgres *pgx.Conn
}

func (u UserHandler) Login(ctx context.Context, user *models.User) (string, string, error) {
	foundUser, err := u.UserRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		log.GetLog().Errorf("Incorrect name or password. error: %v", err)
		return "", "", err
	}

	if !utils.CheckPasswordHash(user.Password, foundUser.Password) {
		return "", "", errors.ErrPasswordIncorrect.Error()
	}

	token, refreshToken, err := utils.TokenGenerator(foundUser.Email, foundUser.FirstName, foundUser.LastName, string(foundUser.ID))
	if err != nil {
		log.GetLog().Errorf("Unable to generate tokens. error: %v", err)
	}

	utils.UpdateAllTokens(u.Postgres, token, refreshToken, string(foundUser.ID))
	if err != nil {
		log.GetLog().Errorf("Unable to update tokens. error: %v", err)
	}

	return token, refreshToken, nil
}

func (u UserHandler) SignUp(ctx context.Context, user *models.User) error {
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

	if !utils.CheckNameValidity(user.FirstName) {
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
