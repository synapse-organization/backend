package utils

import (
	"barista/pkg/errors"
	"barista/pkg/log"
	"barista/pkg/models"
	"context"
	"math/rand"
	"net/mail"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func GenerateRandomPassword(length int) string {
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

func NewPostgres(option models.Postgres) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), option.GetPostgresURL())
	if err != nil {
		log.GetLog().Fatalf("Unable to create connection pool. host: %v, error: %v", option.Host, err)
	}

	return pool
}

func ConnectDB(cfg models.Mongo) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.GetMongoURL()))
	if err != nil {
		log.GetLog().Fatal("Error: " + err.Error())
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.GetLog().Fatal("Error: " + err.Error())
	}
	log.GetLog().Println("Connected to MongoDB")
	return client
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

func CheckStartTime(startTime time.Time) bool {
	return startTime.After(time.Now())
}

func CheckEndTime(startTime time.Time, endTime time.Time) bool {
	return (endTime.After(time.Now()) && endTime.After(startTime))
}

func AppendIfNotExists(slice []string, str string) []string {
    for _, s := range slice {
        if s == str {
            return slice
        }
    }
    return append(slice, str)
}