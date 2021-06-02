package services

import (
	"errors"
	"github.com/RobinsonMarques/parking-system/interfaces"
	"go/types"
	"golang.org/x/crypto/bcrypt"
)

func NewUtilCrud(userInterface interfaces.UserInterface, adminInterface interfaces.AdminInterface, trafficWardenInterface interfaces.TrafficWardenInterface) UtilCrud {
	return UtilCrud{
		userInterface:          userInterface,
		adminInterface:         adminInterface,
		trafficWardenInterface: trafficWardenInterface,
	}
}

type UtilCrud struct {
	userInterface          interfaces.UserInterface
	adminInterface         interfaces.AdminInterface
	trafficWardenInterface interfaces.TrafficWardenInterface
}

func NewCrud(userInterface interfaces.UserInterface, adminInterface interfaces.AdminInterface, trafficWardenInterface interfaces.TrafficWardenInterface) Crud {
	return Crud{
		UserInterface:          userInterface,
		AdminInterface:         adminInterface,
		TrafficWardenInterface: trafficWardenInterface,
	}
}

type Crud struct {
	UserInterface          interfaces.UserInterface
	AdminInterface         interfaces.AdminInterface
	TrafficWardenInterface interfaces.TrafficWardenInterface
}

func (u UtilCrud) GetPassword(email string, userType string) (string, error) {
	if userType == "user" {
		user, err := u.userInterface.GetUserByEmail(email)
		if err != nil {
			return "", err
		}
		return user.Person.Password, err
	} else if userType == "admin" {
		admin, err := u.adminInterface.GetAdminByEmail(email)
		if err != nil {
			return "", err
		}
		return admin.Person.Password, err
	} else if userType == "trafficWarden" {
		trafficWarden, err := u.trafficWardenInterface.GetTrafficWardenByEmail(email)
		if err != nil {
			return "", err
		}
		return trafficWarden.Person.Password, err
	} else {
		err := errors.New("tipo de usuário inválido")
		return "", err
	}

}

func (u UtilCrud) Login(email string, password string) string {
	response := ""
	user, _ := u.userInterface.GetUserByEmail(email)
	admin, _ := u.adminInterface.GetAdminByEmail(email)
	warden, _ := u.trafficWardenInterface.GetTrafficWardenByEmail(email)
	if user.Person.Name != "" {
		err := u.ComparePassword(password, email, "user")
		if err == nil {
			response = "user"
		} else {
			response = "Senha inválida!"
		}
	} else if admin.Person.Name != "" {
		err := u.ComparePassword(password, email, "admin")
		if err == nil {
			response = "admin"
		} else {
			response = "Senha inválida"
		}
	} else if warden.Person.Name != "" {
		err := u.ComparePassword(password, email, "trafficWarden")
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

func (u UtilCrud) ComparePassword(password string, userEmail string, userType string) error {
	var err error
	if userType == "user" {
		userPassword, err := u.GetPassword(userEmail, userType)
		hashedPassword := []byte(userPassword)
		bytePassword := []byte(password)
		err = bcrypt.CompareHashAndPassword(hashedPassword, bytePassword)
		return err
	} else if userType == "admin" {
		adminPassword, err := u.GetPassword(userEmail, userType)
		if err != nil {
			return err
		}
		hashedPassword := []byte(adminPassword)
		bytePassword := []byte(password)
		err = bcrypt.CompareHashAndPassword(hashedPassword, bytePassword)
		return err
	} else if userType == "trafficWarden" {
		wardenPassword, err := u.GetPassword(userEmail, userType)
		if err != nil {
			return err
		}
		hashedPassword := []byte(wardenPassword)
		bytePassword := []byte(password)
		err = bcrypt.CompareHashAndPassword(hashedPassword, bytePassword)
		return err
	} else {
		err = types.Error{Msg: "Tipo de usuário inválido"}
		return err
	}
}
