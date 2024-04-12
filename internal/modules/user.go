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
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type UserHandler struct {
	UserRepo repo.UsersRepo
	TokenRepo repo.TokensRepo
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

	err = u.TokenRepo.Create(ctx, token, refreshToken, user.ID, time.Now())
	if err != nil {
		log.GetLog().Errorf("Unable to create token. error: %v", err)
		return "", "", err
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
	if err != nil && !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		log.GetLog().Errorf("Unable to create user. error: %v", err)
		return err
	}
	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return errors.ErrEmailExists.Error()
	}

	encryptedEmail, err := utils.Encrypt(user.Email)
	if err == nil {
		log.GetLog().Errorf("Unable to encrypt email. error: %v", err)

		emailBody := fmt.Sprintf(`Hello %s,<br><br>
	To verify your email address, please click the link below:<br><br>
	
	<a href="http://localhost:8080/api/user/verify-email?c=%s&callback=http://localhost:5173">Verify Email</a><br><br>

	Yours,<br>
	The Synapse team`, user.FirstName, encryptedEmail)

		err = utils.SendEmail(user.Email, "Barista account verification", emailBody)
		if err != nil {
			log.GetLog().Errorf("Unable to send email. error: %v", err)
		}
	}

	return nil
}

func (u UserHandler) VerifyEmail(ctx context.Context, email string) error {
	decryptedEmail, err := utils.Decrypt(email)
	if err != nil {
		log.GetLog().Errorf("Unable to decrypt email. error: %v", err)
		return err
	}

	err = u.UserRepo.Verify(ctx, decryptedEmail)
	if err != nil {
		log.GetLog().Errorf("Unable to verify user. error: %v", err)
		return err
	}

	return nil
}
