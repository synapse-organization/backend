package utils

import (
	"log"
	"os"
	"time"
	"errors"
	jwt "github.com/golang-jwt/jwt"
)

type SignedDetails struct {
	Uid string
	First_Name string
	Last_Name string
	Email string
	jwt.StandardClaims
}


var SECRET_KEY = os.Getenv("SECRET_KEY")

func TokenGenerator(email, firstname, lastname, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Uid: uid,
		First_Name: firstname,
		Last_Name: lastname,
		Email: email,
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
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "the token is invalid"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is already expired"
		return
	}

	return claims, msg
}

func UpdateAllTokens(signedToken, refreshToken, userID string) (newSignedToken, newSignedRefreshToken string, err error) {
    // Validate the refresh token
    refreshClaims, err := ValidateRefreshToken(refreshToken)
    if err != nil {
        return "", "", err
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

    return newSignedToken, newSignedRefreshToken, nil
}

func ValidateRefreshToken(signedRefreshToken string) (*SignedDetails, error) {
    token, err := jwt.ParseWithClaims(signedRefreshToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(SECRET_KEY), nil
    })
    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*SignedDetails)
    if !ok || !token.Valid {
        return nil, errors.New("invalid refresh token")
    }

    return claims, nil
}
