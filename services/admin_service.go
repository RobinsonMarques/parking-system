package services

import (
	"errors"
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/database"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/utils"
	"gorm.io/gorm"
)

func NewAdminService(db *gorm.DB) AdminService {
	return AdminService{db: db}
}

type AdminService struct {
	db *gorm.DB
}

func (a AdminService) CreateAdmin(input input2.CreateAdminInput) error {
	resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, a.db)

	if resp == "admin" {
		var err error
		input.Person.Password, err = utils.CreateHashPassword(input.Person.Password)
		if err != nil {
			return err
		}

		//Cria o admin
		admin := database.Admin{
			Person: input.Person,
		}
		erro := crud.CreateAdmin(admin, a.db)
		if erro.Data.Error != nil {
			return erro.Data.Error
		}
		return nil
	} else {
		err := errors.New(resp)
		return err
	}
}

func (a AdminService) UpdateAdmin(input input2.UpdateAdminInput, adminID uint) error {
	resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, a.db)

	if resp == "admin" {
		var err error
		input.Person.Password, err = utils.CreateHashPassword(input.Person.Password)
		if err != nil {
			return err
		}
		admin, err := crud.GetAdminByID(adminID, a.db)
		if err == nil {
			admin.Person = input.Person
			crud.UpdateAdmin(admin, a.db)
			return nil
		} else {
			return err
		}
	} else {
		err := errors.New(resp)
		return err
	}
}

func (a AdminService) DeleteAdminByID(input input2.LoginInput, adminID uint) error {
	resp := utils.Login(input.Email, input.Password, a.db)

	if resp == "admin" {
		err := crud.DeleteAdminByID(adminID, a.db)
		if err == nil {
			return nil
		} else {
			return err
		}
	} else {
		err := errors.New(resp)
		return err
	}
}
