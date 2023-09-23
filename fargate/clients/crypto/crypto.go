package crypto

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword(unHashedPassword string) ([]byte, error) {
	// Hash the password before storing it in the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(unHashedPassword), 10)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return []byte{}, err
	}

	return hashedPassword, nil
}

func ComparePasswords(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
