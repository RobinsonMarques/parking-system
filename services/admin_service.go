package services

import (
	"errors"
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/database"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/utils"
	"gorm.io/gorm"
)

func NewAdminService (db *gorm.DB) AdminService{
	return AdminService{db: db}
}

type AdminService struct {
	db *gorm.DB
}

func (a AdminService)CreateAdmin(input input2.CreateAdminInput) error {
	resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, a.db)

	if resp == "admin" {
		input.Person.Password = utils.CreateHashPassword(input.Person.Password)

		//Cria o admin
		admin := database.Admin{
			Person: input.Person,
		}
		err := crud.CreateAdmin(admin, a.db)
		if err.Data.Error != nil{
			return err.Data.Error
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
		input.Person.Password = utils.CreateHashPassword(input.Person.Password)
		admin, err := crud.GetAdminByID(adminID, a.db)
		if err == nil {
			admin.Person = input.Person
			crud.UpdateAdmin(admin, a.db)
			return nil
		}else{
			return err
		}
	}else {
		err := errors.New(resp)
		return err
	}
}

func (a AdminService)DeleteAdminByID(input input2.LoginInput, adminID uint) error {
	resp := utils.Login(input.Email, input.Password, a.db)

	if resp == "admin" {
		err := crud.DeleteAdminByID(adminID, a.db)
		if err == nil{
			return nil
		}else{
			return err
		}
	} else {
		err := errors.New(resp)
		return err
	}
}