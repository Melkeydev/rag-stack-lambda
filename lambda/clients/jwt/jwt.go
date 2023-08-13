package jwt

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

type MyCustomClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(username string) (string, error) {
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {

		log.Printf("Failed to sign the token due to: %v", err)
		return "", err
	}

	return ss, nil
}
