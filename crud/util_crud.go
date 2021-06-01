package crud

import (
	"errors"
	"go/types"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func NewUtilCrud(db *gorm.DB) UtilCrud {
	return UtilCrud{db: db}
}

type UtilCrud struct {
	db *gorm.DB
}

func NewCrud(userCrud UserCrud, adminCrud AdminCrud, trafficWardenCrud TrafficWardenCrud) Crud {
	return Crud{
		UserCrud:          userCrud,
		AdminCrud:         adminCrud,
		TrafficWardenCrud: trafficWardenCrud,
	}
}

type Crud struct {
	UserCrud          UserCrud
	AdminCrud         AdminCrud
	TrafficWardenCrud TrafficWardenCrud
}

func GetPassword(email string, userType string, crud Crud) (string, error) {
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
	userCrud := NewUserCrud(u.db)
	adminCrud := NewAdminCrud(u.db)
	trafficWardenCrud := NewTrafficWardenCrud(u.db)
	crud := NewCrud(userCrud, adminCrud, trafficWardenCrud)
	user, _ := userCrud.GetUserByEmail(email)
	admin, _ := adminCrud.GetAdminByEmail(email)
	warden, _ := trafficWardenCrud.GetTrafficWardenByEmail(email)
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
