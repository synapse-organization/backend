package modules

import (
	"barista/pkg/errors"
	"barista/pkg/log"
	"barista/pkg/models"
	"barista/pkg/repo"
	"barista/pkg/utils"
	"context"
	"fmt"
	"math/rand"

	"github.com/jackc/pgx/v5"
)

type UserHandler struct {
	UserRepo repo.UsersRepo
	Postgres *pgx.Conn
}

const (
    passwordLength = 8
)

func (u UserHandler) ForgetPassword(ctx context.Context, user *models.User) error {
	// Find the user based on email
	foundUser, err := u.UserRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		log.GetLog().Errorf("Email does not exist. error: %v", err)
		return err
	}

	// Generate a random password
    newPassword := utils.GenerateRandomPassword(passwordLength)

    // Update user's password with the new random password
    hashedPassword, err := utils.HashPassword(newPassword)
    if err != nil {
        log.GetLog().Errorf("Unable to hash password. error: %v", err)
        return err
    }
    foundUser.Password = hashedPassword

    // Save the updated password to the database
    if err := u.UserRepo.UpdatePassword(ctx, foundUser.ID, foundUser.Password); err != nil {
        log.GetLog().Errorf("Unable to update user's password. error: %v", err)
        return err
    }

	// Send email with the new password to the user
	emailBody := fmt.Sprintf(`Hello %s,<br><br>
	A new password has been requested for your Barista account associated with %s.<br><br>

	Here is your new password: <strong>%s</strong><br><br>

	You can use your random generated password to login to your account.<br><br>

	After logging in your account, you can reset your password in your profile section.<br><br>

	Yours,<br>
	The Synapse team`, foundUser.FirstName, foundUser.Email, newPassword)
	err = utils.SendEmail(foundUser.Email, "Barista account recovery", emailBody)
	if err != nil {
		log.GetLog().Errorf("Unable to send email. error: %v", err)
        return err
	}

	return nil
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
