package utils

import (
	"github.com/RobinsonMarques/parking-system/crud"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
)

func CreateHashPassword(password string) string {
	bytePassword := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, 8)

	if err != nil {
		log.Fatal(err)
	}
	return string(hashedPassword)
}

func ComparePassword(password string, userEmail string, db *gorm.DB) error {
	hashedPassword := []byte(crud.GetPassword(userEmail, db))
	bytePassword := []byte(password)

	err := bcrypt.CompareHashAndPassword(hashedPassword, bytePassword)

	return err
}
