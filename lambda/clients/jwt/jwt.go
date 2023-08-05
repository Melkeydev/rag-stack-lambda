package jwt

import (
	"log"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(username string) (string, error) {
	// TODO: this should come from env
	mySigningKey := []byte("randomString")

	type MyCustomClaims struct {
		Username string `json:"username"`
		jwt.StandardClaims
	}

	// TODO: Issue should also come from env
	claims := MyCustomClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: 15000,
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
