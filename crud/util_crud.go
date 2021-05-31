package crud

import (
	"errors"
	"gorm.io/gorm"
)

func NewUtilCrud(db *gorm.DB) UtilCrud {
	return UtilCrud{db: db}
}

type UtilCrud struct {
	db *gorm.DB
}

func (u UtilCrud) GetPassword(email string, userType string, crud Crud) (string, error) {
	if userType == "user" {
		user, err := crud.UserCrud.GetUserByEmail(email)
		if err != nil {
			return "", err
		}
		return user.Person.Password, err
	} else if userType == "admin" {
		admin, err := crud.AdminCrud.GetAdminByEmail(email)
		if err != nil {
			return "", err
		}
		return admin.Person.Password, err
	} else if userType == "trafficWarden" {
		trafficWarden, err := crud.TrafficWardenCrud.GetTrafficWardenByEmail(email)
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
	crud := crud.NewCrud(db)
	user, _ := crud.UserCrud.GetUserByEmail(email)
	admin, _ := crud.AdminCrud.GetAdminByEmail(email)
	warden, _ := crud.TrafficWardenCrud.GetTrafficWardenByEmail(email)
	if user.Person.Name != "" {
		err := ComparePassword(password, email, "user", db)
		if err == nil {
			response = "user"
		} else {
			response = "Senha inválida!"
		}
	} else if admin.Person.Name != "" {
		err := ComparePassword(password, email, "admin", db)
		if err == nil {
			response = "admin"
		} else {
			response = "Senha inválida"
		}
	} else if warden.Person.Name != "" {
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
