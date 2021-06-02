package crud

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

func GetPassword(email string, userType string, crud Crud) (string, error) {
	if userType == "user" {
		user, err := crud.UserInterface.GetUserByEmail(email)
		if err != nil {
			return "", err
		}
		return user.Person.Password, err
	} else if userType == "admin" {
		admin, err := crud.AdminInterface.GetAdminByEmail(email)
		if err != nil {
			return "", err
		}
		return admin.Person.Password, err
	} else if userType == "trafficWarden" {
		trafficWarden, err := crud.TrafficWardenInterface.GetTrafficWardenByEmail(email)
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
	crud := NewCrud(u.userInterface, u.adminInterface, u.trafficWardenInterface)
	user, _ := u.userInterface.GetUserByEmail(email)
	admin, _ := u.adminInterface.GetAdminByEmail(email)
	warden, _ := u.trafficWardenInterface.GetTrafficWardenByEmail(email)
	if user.Person.Name != "" {
		err := ComparePassword(password, email, "user", crud)
		if err == nil {
			response = "user"
		} else {
			response = "Senha inválida!"
		}
	} else if admin.Person.Name != "" {
		err := ComparePassword(password, email, "admin", crud)
		if err == nil {
			response = "admin"
		} else {
			response = "Senha inválida"
		}
	} else if warden.Person.Name != "" {
		err := ComparePassword(password, email, "trafficWarden", crud)
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

func ComparePassword(password string, userEmail string, userType string, crud Crud) error {
	var err error
	if userType == "user" {
		userPassword, err := GetPassword(userEmail, userType, crud)
		hashedPassword := []byte(userPassword)
		bytePassword := []byte(password)
		err = bcrypt.CompareHashAndPassword(hashedPassword, bytePassword)
		return err
	} else if userType == "admin" {
		adminPassword, err := GetPassword(userEmail, userType, crud)
		if err != nil {
			return err
		}
		hashedPassword := []byte(adminPassword)
		bytePassword := []byte(password)
		err = bcrypt.CompareHashAndPassword(hashedPassword, bytePassword)
		return err
	} else if userType == "trafficWarden" {
		wardenPassword, err := GetPassword(userEmail, userType, crud)
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
