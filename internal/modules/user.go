package modules

import (
	"barista/pkg/errors"
	"barista/pkg/log"
	"barista/pkg/models"
	"barista/pkg/repo"
	"barista/pkg/utils"
	"context"
	go_error "errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type UserHandler struct {
	UserRepo  repo.UsersRepo
	TokenRepo repo.TokensRepo
	Postgres  *pgxpool.Pool
}

const (
	passwordLength = 8
)

func (u UserHandler) SignUp(ctx context.Context, user *models.User) error {
	user.ID = rand.Int31()

	if !utils.CheckEmailValidity(user.Email) {
		return errors.ErrEmailInvalid.Error()
	}

	if !utils.CheckPhoneValidity(user.Phone) {
		return errors.ErrPhoneInvalid.Error()
	}

	if !utils.CheckPasswordValidity(user.Password) {
		return errors.ErrPasswordIncorrect.Error()
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

func (u UserHandler) Login(ctx context.Context, user *models.User) (map[string]string, error) {
	foundUsers, err := u.UserRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		log.GetLog().Errorf("Incorrect name or password. error: %v", err)
		return nil, err
	}
	var correctUsers []*models.User
	for _, foundUser := range foundUsers {
		if utils.CheckPasswordHash(user.Password, foundUser.Password) {
			correctUsers = append(correctUsers, foundUser)
		}
	}
	if len(correctUsers) == 0 {
		return nil, errors.ErrPasswordIncorrect.Error()
	}

	tokens := map[string]string{}
	for _, foundUser := range correctUsers {
		claims, token, err := utils.TokenGenerator(foundUser.ID, foundUser.Email, foundUser.FirstName, foundUser.LastName, int32(foundUser.Role))
		if err != nil {
			log.GetLog().Errorf("Unable to generate tokens. error: %v", err)
		}

		expiresAt := time.Unix(claims.ExpiresAt, 0)
		err = u.TokenRepo.Create(ctx, claims.TokenID, token, foundUser.ID, expiresAt)
		if err != nil {
			log.GetLog().Errorf("Unable to create token. error: %v", err)
			return nil, err
		}
		tokens[models.RoleToString(foundUser.Role)] = token
	}

	return tokens, nil
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

func (u UserHandler) ForgetPassword(ctx context.Context, user *models.User) error {
	// Find the user based on email
	foundUsers, err := u.UserRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		log.GetLog().Errorf("Email does not exist. error: %v", err)
		return err
	}

	// Generate a random password
	newPassword := utils.GenerateRandomStr(passwordLength)

	for _, foundUser := range foundUsers {
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
	}

	return nil
}

func (u UserHandler) UserProfile(ctx context.Context, userID string) (*models.User, error) {
	user_id, err := strconv.Atoi(userID)
	if err != nil {
		log.GetLog().Errorf("Unable to convert userID to int32. error: %v", err)
		return nil, err
	}

	user, err := u.UserRepo.GetByID(ctx, int32(user_id))
	if err != nil {
		log.GetLog().Errorf("Unable to get user by id. errror: %v", err)
		return nil, err
	}

	return user, nil
}

func (u UserHandler) EditProfile(ctx context.Context, newDetail *models.User, userID string) error {
	user_id, err := strconv.Atoi(userID)
	if err != nil {
		log.GetLog().Errorf("Unable to convert user id to int32. error: %v", err)
		return err
	}

	foundUser, err := u.UserRepo.GetByID(ctx, int32(user_id))
	if err != nil {
		log.GetLog().Errorf("Incorrect user id. error: %v", err)
		return err
	}

	if newDetail.FirstName != foundUser.FirstName && newDetail.FirstName != "" {
		if !utils.CheckNameValidity(newDetail.FirstName) {
			return errors.ErrFirstNameInvalid.Error()
		}

		err = u.UserRepo.UpdateFirstName(ctx, int32(user_id), newDetail.FirstName)
		if err != nil {
			log.GetLog().Errorf("Unable to update user's first name. error: %v", err)
			return err
		}
	}

	if newDetail.LastName != foundUser.LastName && newDetail.LastName != "" {
		if !utils.CheckNameValidity(newDetail.LastName) {
			return errors.ErrLastNameInvalid.Error()
		}

		err = u.UserRepo.UpdateLastName(ctx, int32(user_id), newDetail.LastName)
		if err != nil {
			log.GetLog().Errorf("Unable to update user's last name. error: %v", err)
		}
	}

	if newDetail.Sex != foundUser.Sex && newDetail.Sex != 0 {
		err = u.UserRepo.UpdateSex(ctx, int32(user_id), fmt.Sprintf("%v", newDetail.Sex))
		if err != nil {
			log.GetLog().Errorf("Unable to update user's sex. error: %v", err)
			return err
		}
	}

	if newDetail.Phone != foundUser.Phone && newDetail.Phone != 0 {
		if !utils.CheckPhoneValidity(newDetail.Phone) {
			return errors.ErrPhoneInvalid.Error()
		}

		err = u.UserRepo.UpdatePhone(ctx, int32(user_id), int32(newDetail.Phone))
		if err != nil {
			log.GetLog().Errorf("Unable to update user's phone number. error: %v", err)
			return err
		}
	}

	if newDetail.BankAccount != "" && newDetail.NationalID != "" {
		err = u.UserRepo.UpdateExtraInfo(ctx, int32(user_id), map[string]interface{}{
			"national_id":  newDetail.NationalID,
			"bank_account": newDetail.BankAccount,
		})
		if err != nil {
			log.GetLog().Errorf("Unable to update user's extra info. error: %v", err)
			return err
		}
	}

	return err
}

func (u UserHandler) ChangePassword(ctx context.Context, userID int32, password, currentPassword string) error {

	foundUser, err := u.UserRepo.GetByID(ctx, userID)
	if err != nil {
		log.GetLog().WithError(err).Error("Unable to get user by id")
		return err
	}

	if !utils.CheckPasswordHash(currentPassword, foundUser.Password) {
		return errors.ErrPasswordIncorrect.Error()
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.GetLog().Errorf("Unable to hash password. error: %v", err)
		return err
	}

	err = u.UserRepo.UpdatePassword(ctx, userID, hashedPassword)
	if err != nil {
		log.GetLog().Errorf("Unable to update user's password. error: %v", err)
		return err
	}

	return nil

}

func (u UserHandler) Logout(ctx context.Context, token string) error {
	claims, err := utils.ValidateToken(u.Postgres, token)
	if err != "" {
		log.GetLog().Errorf("Unable to validate token. error: %v", err)
		return go_error.New("unable to validate token")
	}

	err2 := u.TokenRepo.DeleteByID(ctx, claims.Uid)
	if err2 != nil {
		return err2
	}

	return nil
}

func (u UserHandler) ManagerAgreement(ctx context.Context, userID int32, nationalID, bankAccount string) error {
	err := u.UserRepo.UpdateExtraInfo(ctx, userID, map[string]interface{}{
		"national_id":  nationalID,
		"bank_account": bankAccount,
	})
	if err != nil {
		log.GetLog().Errorf("Unable to update user's extra info. error: %v", err)
		return err
	}
	return nil

}
