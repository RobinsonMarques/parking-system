package services

import (
	"errors"
	"github.com/RobinsonMarques/parking-system/database"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/interfaces"
	"github.com/RobinsonMarques/parking-system/utils"
)

func NewUserService(userInterface interfaces.UserInterface, vehicleInterface interfaces.VehicleInterface, rechargeInterface interfaces.RechargeInterface, billetInterface interfaces.BilletInterface, utilInterface interfaces.UtilInterface) UserService {
	return UserService{
		userInterface:     userInterface,
		vehicleInterface:  vehicleInterface,
		rechargeInterface: rechargeInterface,
		billetInterface:   billetInterface,
		utilInterface:     utilInterface,
	}
}

type UserService struct {
	userInterface     interfaces.UserInterface
	vehicleInterface  interfaces.VehicleInterface
	rechargeInterface interfaces.RechargeInterface
	billetInterface   interfaces.BilletInterface
	utilInterface     interfaces.UtilInterface
}

func (u UserService) CreateUser(input input2.CreateUserInput) error {
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
	err = u.userInterface.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (u UserService) GetUserByDocument(input input2.LoginInput, document string) (database.User, error) {
	resp := u.utilInterface.Login(input.Email, input.Password)
	if resp == "admin" {
		user, err := u.userInterface.GetUserByDocument(document)
		if err != nil {
			return user, err
		}
		vehicles, _ := u.vehicleInterface.GetVehiclesByUserId(user.ID)
		recharges, _ := u.rechargeInterface.GetRechargeByUserId(user.ID)

		for i := range recharges {
			billet, _ := u.billetInterface.GetBilletByRechargeId(recharges[i].ID)
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

func (u UserService) UpdateUser(input input2.UpdateUserInput, userID uint) error {
	resp := u.utilInterface.Login(input.LoginInput.Email, input.LoginInput.Password)
	if resp == "user" || resp == "admin" {
		user, err := u.userInterface.GetUserByID(userID)
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
			err = u.userInterface.UpdateUser(user)
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

func (u UserService) DeleteUserByID(input input2.LoginInput, userID uint) error {
	resp := u.utilInterface.Login(input.Email, input.Password)
	if resp == "user" || resp == "admin" {
		user, err := u.userInterface.GetUserByEmail(input.Email)
		if err != nil {
			return err
		}
		if resp == "user" && user.ID != userID {
			err := errors.New("usuário não possui permissão")
			return err
		} else {
			err := u.userInterface.DeleteUserByID(userID)
			if err != nil {
				return err
			}
			err = u.vehicleInterface.DeleteVehiclesByUserID(userID)
			if err != nil {
				return err
			}
			err = u.rechargeInterface.DeleteRechargeByUserID(userID)
			if err != nil {
				return err
			}
			recharges, err := u.rechargeInterface.GetRechargeByUserId(userID)
			if err != nil {
				return err
			}
			for i := range recharges {
				err := u.billetInterface.DeleteBilletByRechargeID(recharges[i].ID)
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
