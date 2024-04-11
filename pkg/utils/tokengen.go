package utils

import (
	"barista/pkg/repo"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5"
)

type SignedDetails struct {
	Uid        string
	First_Name string
	Last_Name  string
	Email      string
	jwt.StandardClaims
}

var SECRET_KEY = os.Getenv("SECRET_KEY")

func TokenGenerator(email, firstname, lastname, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Uid:        uid,
		First_Name: firstname,
		Last_Name:  lastname,
		Email:      email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(10)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(4)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, err
}

func ValidateToken(postgres *pgx.Conn, signedToken string) (claims *SignedDetails, msg string) {
	exists, err := repo.CheckTokenExistence(postgres, signedToken)
	if err != nil {
		return nil, err.Error()
	}

	if !exists {
		return nil, "Token doesn't exist"
	}

	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "The token is invalid"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "Token is already expired"
		return
	}

	return claims, msg
}

func UpdateAllTokens(postgres *pgx.Conn, signedToken, refreshToken, userID string) (newSignedToken, newSignedRefreshToken string, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Validate the refresh token
	refreshClaims, msg := ValidateToken(postgres, refreshToken)
	if msg != "" {
		return "", "", errors.New(msg)
	}

	// Check if the user ID in the refresh token matches the provided user ID
	if refreshClaims.Uid != userID {
		return "", "", errors.New("invalid user ID in refresh token")
	}

	// Generate new access token and refresh token
	newSignedToken, newSignedRefreshToken, err = TokenGenerator(refreshClaims.Email, refreshClaims.First_Name, refreshClaims.Last_Name, refreshClaims.Uid)
	if err != nil {
		return "", "", err
	}

	_, err = postgres.Exec(ctx,
		`UPDATE users
        SET token = $1, refresh_token = $2, updated_at = $3
        WHERE user_id = $4`,
		newSignedToken, newSignedRefreshToken, time.Now(), userID)
	if err != nil {
		panic(fmt.Sprintf("Unable to update tokens. error: %v", err))
	}

	return newSignedToken, newSignedRefreshToken, nil
}
