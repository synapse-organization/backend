package utils

import (
	"barista/pkg/errors"
	"barista/pkg/log"
	"barista/pkg/models"
	"context"
	"math/rand"
	"net"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func GenerateRandomStr(length int) string {
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

func CheckCafeTimeValidity(cafeTime int8) bool {
	return cafeTime >= 0 && cafeTime <= 23
}

func CheckCapacityValidity(capacity int32) bool {
	return capacity > 0
}

func CheckPriceValidity(price float64) bool {
	return price >= 0
}

func CheckReservability(preReserve bool, newReserve bool, preCapacity int32, newCapacity int32, attendees int32) (bool, bool, error) {
	updateCapacity := false
	updateReserve := false
	if !preReserve {
		if newCapacity > preCapacity {
			updateCapacity = true
			updateReserve = true
		} else {
			return updateCapacity, updateReserve, errors.ErrCapacityInvalid.Error()
		}
	} else {
		if newCapacity < attendees {
			return updateCapacity, updateReserve, errors.ErrCapacityInvalid.Error()
		} else if newCapacity == attendees {
			updateCapacity = true
			updateReserve = true
		} else {
			updateCapacity = true
		}
	}

	return updateCapacity, updateReserve, nil
}

func getBaseURL() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	addrs, err := net.LookupIP(hostname)
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil && !ipv4.IsLoopback() {
			if ipv4.String() == "127.0.0.1" {
				return "http://localhost:8080", nil
			}
			return "http://" + ipv4.String() + ":8080", nil
		}
	}

	return "http://localhost:8080", nil
}
