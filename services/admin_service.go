package services

import (
	"errors"
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/database"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/utils"
)

func NewAdminService(adminCrud crud.AdminCrud, utilCrud crud.UtilCrud) AdminService {
	return AdminService{
		adminCrud: adminCrud,
		utilCrud:  utilCrud,
	}
}

type AdminService struct {
	adminCrud crud.AdminCrud
	utilCrud  crud.UtilCrud
}

func (a AdminService) CreateAdmin(input input2.CreateAdminInput, service AdminService) error {
	resp := service.utilCrud.Login(input.LoginInput.Email, input.LoginInput.Password)
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
		err = service.adminCrud.CreateAdmin(admin)
		if err != nil {
			return err
		}
		return nil
	} else {
		err := errors.New(resp)
		return err
	}
}

func (a AdminService) UpdateAdmin(input input2.UpdateAdminInput, adminID uint, service AdminService) error {
	resp := service.utilCrud.Login(input.LoginInput.Email, input.LoginInput.Password)
	if resp == "admin" {
		var err error
		input.Person.Password, err = utils.CreateHashPassword(input.Person.Password)
		if err != nil {
			return err
		}
		admin, err := service.adminCrud.GetAdminByID(adminID)
		if err == nil {
			admin.Person = input.Person
			service.adminCrud.UpdateAdmin(admin)
			return nil
		} else {
			return err
		}
	} else {
		err := errors.New(resp)
		return err
	}
}

func (a AdminService) DeleteAdminByID(input input2.LoginInput, adminID uint, service AdminService) error {
	resp := service.utilCrud.Login(input.Email, input.Password)
	if resp == "admin" {
		err := service.adminCrud.DeleteAdminByID(adminID)
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
