package services

import (
	"errors"
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/database"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/utils"
	"gorm.io/gorm"
)

func NewUserService(db *gorm.DB) UserService {
	return UserService{db: db}
}

type UserService struct {
	db *gorm.DB
}

func (u UserService) CreateUser(input input2.CreateUserInput) error {
	input.Person.Password = utils.CreateHashPassword(input.Person.Password)
	//Cria o user
	user := database.User{
		Person:   input.Person,
		Document: input.Document,
		Balance:  0,
		Recharge: nil,
		Vehicle:  nil,
	}
	resp := crud.CreateUser(user, u.db)
	if resp.Data.Error != nil {
		return resp.Data.Error
	}
	return nil
}

func (u UserService) GetUserByDocument(input input2.LoginInput, document string) (database.User, error) {
	resp := utils.Login(input.Email, input.Password, u.db)

	if resp == "admin" {
		user, err := crud.GetUserByDocument(document, u.db)
		if err != nil {
			return user, err
		}
		vehicles, _ := crud.GetVehiclesByUserId(user.ID, u.db)
		recharges, _ := crud.GetRechargeByUserId(user.ID, u.db)

		for i := range recharges {
			billet, _ := crud.GetBilletByRechargeId(recharges[i].ID, u.db)
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
	resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, u.db)

	if resp == "user" || resp == "admin" {
		user, err := crud.GetUserByID(userID, u.db)
		if err != nil {
			return err
		}
		if resp == "user" && user.Person.Email != input.LoginInput.Email {
			err := errors.New("usuário não possui permissão")
			return err
		} else {
			input.Person.Password = utils.CreateHashPassword(input.Person.Password)
			user.Person = input.Person
			user.Document = input.Document
			crud.UpdateUser(user, u.db)
			return nil
		}
	} else {
		err := errors.New(resp)
		return err
	}
}

func (u UserService) DeleteUserByID(input input2.LoginInput, userID uint) error {
	resp := utils.Login(input.Email, input.Password, u.db)

	if resp == "user" || resp == "admin" {
		user, err := crud.GetUserByEmail(input.Email, u.db)
		if err != nil {
			return err
		}
		if resp == "user" && user.ID != userID {
			err := errors.New("usuário não possui permissão")
			return err
		} else {
			crud.DeleteUserByID(userID, u.db)
			return nil
		}
	} else {
		err := errors.New(resp)
		return err
	}
}
