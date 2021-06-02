package services

import (
	"errors"
	"github.com/RobinsonMarques/parking-system/database"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/interfaces"
	"github.com/RobinsonMarques/parking-system/utils"
)

func NewAdminService(adminInterface interfaces.AdminInterface, utilInterface interfaces.UtilInterface) AdminService {
	return AdminService{
		adminInterface: adminInterface,
		utilInterface:  utilInterface,
	}
}

type AdminService struct {
	adminInterface interfaces.AdminInterface
	utilInterface  interfaces.UtilInterface
}

func (a AdminService) CreateAdmin(input input2.CreateAdminInput) error {
	resp := a.utilInterface.Login(input.LoginInput.Email, input.LoginInput.Password)
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
		err = a.adminInterface.CreateAdmin(admin)
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
	resp := a.utilInterface.Login(input.LoginInput.Email, input.LoginInput.Password)
	if resp == "admin" {
		var err error
		input.Person.Password, err = utils.CreateHashPassword(input.Person.Password)
		if err != nil {
			return err
		}
		admin, err := a.adminInterface.GetAdminByID(adminID)
		if err == nil {
			admin.Person = input.Person
			a.adminInterface.UpdateAdmin(admin)
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
	resp := a.utilInterface.Login(input.Email, input.Password)
	if resp == "admin" {
		err := a.adminInterface.DeleteAdminByID(adminID)
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
