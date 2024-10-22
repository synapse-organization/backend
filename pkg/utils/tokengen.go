package utils

import (
	"barista/pkg/log"
	"barista/pkg/repo"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	jwt "github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type SignedDetails struct {
	Uid        int32
	First_Name string
	Last_Name  string
	Email      string
	TokenID    int32 `json:"tid"`
	Role       int32
	jwt.StandardClaims
}

var SECRET_KEY = os.Getenv("SECRET_KEY")

func TokenGenerator(uid int32, email, firstname, lastname string, role int32) (*SignedDetails, string, error) {
	tokenID := uuid.New().ID()

	claims := &SignedDetails{
		Uid:        uid,
		First_Name: firstname,
		Last_Name:  lastname,
		Email:      email,
		TokenID:    int32(tokenID),
		Role:       role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24*3)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return nil, "", err
	}

	return claims, token, err
}

func ValidateToken(postgres *pgxpool.Pool, signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		log.GetLog().Errorf("Token is incorrect. error: %v", err)
		msg = "Token is incorrect"
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		log.GetLog().Errorf("The token is invalid.")
		msg = "The token is invalid"
		return
	}

	exists, err := repo.CheckTokenExistence(postgres, claims.TokenID)
	if err != nil {
		return nil, err.Error()
	}

	if !exists {
		return nil, "Token doesn't exist"
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "Token is already expired"
		repo.DeleteByID(postgres, claims.TokenID)
		return
	}

	return claims, msg
}

// func UpdateAllTokens(postgres *pgxpool.Pool, signedToken, refreshToken, userID string) (newSignedToken, newSignedRefreshToken string, err error) {
// 	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	// Validate the refresh token
// 	refreshClaims, msg := ValidateToken(postgres, refreshToken)
// 	if msg != "" {
// 		return "", "", errors.New(msg)
// 	}

// 	// Check if the user ID in the refresh token matches the provided user ID
// 	if refreshClaims.Uid != userID {
// 		return "", "", errors.New("invalid user ID in refresh token")
// 	}

// 	// Generate new access token and refresh token
// 	newSignedToken, newSignedRefreshToken, err = TokenGenerator(refreshClaims.Email, refreshClaims.First_Name, refreshClaims.Last_Name, refreshClaims.Uid)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	_, err = postgres.Exec(ctx,
// 		`UPDATE users
//         SET token = $1, refresh_token = $2, updated_at = $3
//         WHERE user_id = $4`,
// 		newSignedToken, newSignedRefreshToken, time.Now(), userID)
// 	if err != nil {
// 		panic(fmt.Sprintf("Unable to update tokens. error: %v", err))
// 	}

// 	return newSignedToken, newSignedRefreshToken, nil
// }
