package services

import (
	"errors"
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/database"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/utils"
)

func NewUserService(userCrud crud.UserCrud, vehicleCrud crud.VehicleCrud, rechargeCrud crud.RechargeCrud, billetCrud crud.BilletCrud, utilCrud crud.UtilCrud) UserService {
	return UserService{
		userCrud:     userCrud,
		vehicleCrud:  vehicleCrud,
		rechargeCrud: rechargeCrud,
		billetCrud:   billetCrud,
		util:         utilCrud,
	}
}

type UserService struct {
	userCrud     crud.UserCrud
	vehicleCrud  crud.VehicleCrud
	rechargeCrud crud.RechargeCrud
	billetCrud   crud.BilletCrud
	util         crud.UtilCrud
}

func (u UserService) CreateUser(input input2.CreateUserInput, service UserService) error {
	var err error
	input.Person.Password, err = utils.CreateHashPassword(input.Person.Password)
	if err != nil {
		return err
	}
	//Cria o user
	user := database.User{
		Person:   input.Person,
		Document: input.Document,
		Balance:  0,
		Recharge: nil,
		Vehicle:  nil,
	}
	err = service.userCrud.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (u UserService) GetUserByDocument(input input2.LoginInput, document string, service UserService) (database.User, error) {
	resp := service.util.Login(input.Email, input.Password)
	if resp == "admin" {
		user, err := service.userCrud.GetUserByDocument(document)
		if err != nil {
			return user, err
		}
		vehicles, _ := service.vehicleCrud.GetVehiclesByUserId(user.ID)
		recharges, _ := service.rechargeCrud.GetRechargeByUserId(user.ID)

		for i := range recharges {
			billet, _ := service.billetCrud.GetBilletByRechargeId(recharges[i].ID)
			recharges[i].Billet = billet
		}
		user.Vehicle = vehicles
		user.Recharge = recharges
		return user, err
	}
	err := errors.New(resp)
	user := database.User{}
	return user, err
}

func (u UserService) UpdateUser(input input2.UpdateUserInput, userID uint, service UserService) error {
	resp := service.util.Login(input.LoginInput.Email, input.LoginInput.Password)
	if resp == "user" || resp == "admin" {
		user, err := service.userCrud.GetUserByID(userID)
		if err != nil {
			return err
		}
		if resp == "user" && user.Person.Email != input.LoginInput.Email {
			err := errors.New("usuário não possui permissão")
			return err
		} else {
			var err error
			input.Person.Password, err = utils.CreateHashPassword(input.Person.Password)
			if err != nil {
				return err
			}
			user.Person = input.Person
			user.Document = input.Document
			err = service.userCrud.UpdateUser(user)
			if err != nil {
				return err
			}
			return nil
		}
	} else {
		err := errors.New(resp)
		return err
	}
}

func (u UserService) DeleteUserByID(input input2.LoginInput, userID uint, service UserService) error {
	resp := service.util.Login(input.Email, input.Password)
	if resp == "user" || resp == "admin" {
		user, err := service.userCrud.GetUserByEmail(input.Email)
		if err != nil {
			return err
		}
		if resp == "user" && user.ID != userID {
			err := errors.New("usuário não possui permissão")
			return err
		} else {
			err := service.userCrud.DeleteUserByID(userID)
			if err != nil {
				return err
			}
			err = service.vehicleCrud.DeleteVehiclesByUserID(userID)
			if err != nil {
				return err
			}
			err = service.rechargeCrud.DeleteRechargeByUserID(userID)
			if err != nil {
				return err
			}
			recharges, err := service.rechargeCrud.GetRechargeByUserId(userID)
			if err != nil {
				return err
			}
			for i := range recharges {
				err := service.billetCrud.DeleteBilletByRechargeID(recharges[i].ID)
				if err != nil {
					return err
				}
			}
			return nil
		}
	} else {
		err := errors.New(resp)
		return err
	}
}
