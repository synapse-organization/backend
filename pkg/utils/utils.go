package utils

import (
	"barista/pkg/errors"
	"barista/pkg/log"
	"barista/pkg/models"
	"context"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"strings"
	"math/rand"
	"time"
)


func GenerateRandomPassword(length int) (string) {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    random := rand.New(rand.NewSource(time.Now().UnixNano()))
    password := make([]byte, length)
    for i := range password {
        password[i] = charset[random.Intn(len(charset))]
    }
    return string(password)
}

func SplitMethodPrefix(methodName string) (string, string) {
	for i, char := range methodName {
		if i > 0 && strings.ToUpper(string(char)) == string(char) {
			return methodName[:i], methodName[i:]
		}
	}
	return "", methodName
}

func NewPostgres(option models.Postgres) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), option.GetPostgresURL())
	if err != nil {
		log.GetLog().Fatalf("Unable to create connection pool. host: %v, error: %v", option.Host, err)
	}
	return conn
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CheckPasswordValidity(password string) bool {
	if len(password) < 8 {
		return false
	}
	return true
}

func CheckNameValidity(name string) bool {
	if len(name) < 3 {
		return false
	}
	return true
}

func CheckEmailValidity(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func CheckPhoneValidity(phone int64) bool {
	return phone > 0
}

func IsCommonError(t interface{}) bool {
	switch t.(type) {
	case errors.StringError:
		return false
	default:
		return true
	}
}
