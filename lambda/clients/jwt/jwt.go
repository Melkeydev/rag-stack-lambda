package jwt

import (
	"errors"
	"log"
	"time"

	ragDynamo "melkeydev/ragStackCDK/clients/dynamo"

	"github.com/golang-jwt/jwt"
)

type TokenValidator interface {
	GenerateAccessToken(username string) (string, error)
	GenerateRefreshToken(username string) (string, error)
	ValidateRefreshToken(refreshToken string) (string, error)
}

type JWTClient struct {
	db ragDynamo.UserStorageDB
}

func NewJWTClient(db ragDynamo.UserStorageDB) TokenValidator {
	return &JWTClient{
		db: db,
	}
}

type MyCustomClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (j *JWTClient) GenerateAccessToken(username string) (string, error) {
	// TODO: this should come from env
	mySigningKey := []byte("randomString")

	expirationTime := time.Now().Add(10 * time.Minute)

	// TODO: Issue should also come from env
	claims := MyCustomClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "test",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	acessString, err := accessToken.SignedString(mySigningKey)
	if err != nil {
		log.Printf("Failed to sign the token due to: %v", err)
		return "", err
	}

	return acessString, nil
}

func (j *JWTClient) GenerateRefreshToken(username string) (string, error) {
	expirationTimeRefresh := time.Now().Add(30 * 24 * time.Hour) // 30 days from now
	refreshClaims := MyCustomClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTimeRefresh.Unix(),
			Issuer:    "test",
		},
	}

	mySigningKey := []byte("randomString")
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString(mySigningKey)

	if err != nil {
		return "", err
	}

	return refreshString, nil
}

func (j *JWTClient) ValidateRefreshToken(refreshToken string) (string, error) {
	mySigningKey := []byte("randomString")
	token, err := jwt.ParseWithClaims(refreshToken, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok || !token.Valid {
		return "", errors.New("Invalid or expired refresh token")
	}

	valid := j.db.ValidateRefreshToken(claims.Username, refreshToken)
	if !valid {
		return "", errors.New("Invalid refresh token")
	}

	return claims.Username, nil
}
