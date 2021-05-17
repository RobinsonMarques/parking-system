package utils

import (
	"github.com/RobinsonMarques/parking-system/crud"
	"go/types"
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

func ComparePassword(password string, userEmail string, userType string, db *gorm.DB) error {
	var err error

	if userType == "user" {
		hashedPassword := []byte(crud.GetPassword(userEmail, userType, db))
		bytePassword := []byte(password)
		err = bcrypt.CompareHashAndPassword(hashedPassword, bytePassword)
		return err
	} else if userType == "admin" {
		hashedPassword := []byte(crud.GetPassword(userEmail, userType, db))
		bytePassword := []byte(password)
		err = bcrypt.CompareHashAndPassword(hashedPassword, bytePassword)
		return err
	} else if userType == "trafficWarden" {
		hashedPassword := []byte(crud.GetPassword(userEmail, userType, db))
		bytePassword := []byte(password)
		err = bcrypt.CompareHashAndPassword(hashedPassword, bytePassword)
		return err
	} else {
		err = types.Error{Msg: "Tipo de usuário inválido"}
		return err
	}

}

func Login(email string, password string, db *gorm.DB) string {
	response := ""
	if crud.GetUserByEmail(email, db).Person.Name != "" {
		err := ComparePassword(password, email, "user", db)
		if err == nil {
			response = "user"
		} else {
			response = "Senha inválida!"
		}
	} else if crud.GetAdminByEmail(email, db).Person.Name != "" {
		err := ComparePassword(password, email, "admin", db)
		if err == nil {
			response = "admin"
		} else {
			response = "Senha inválida"
		}
	} else if crud.GetTrafficWardenByEmail(email, db).Person.Name != "" {
		err := ComparePassword(password, email, "trafficWarden", db)
		if err == nil {
			response = "trafficWarden"
		} else {
			response = "Senha inválida"
		}
	} else {
		response = "Usuário não cadastrado"
	}

	return response
}
