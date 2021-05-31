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
	userCrud := crud.NewUserCrud(u.db)
	err = userCrud.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (u UserService) GetUserByDocument(input input2.LoginInput, document string) (database.User, error) {
	resp := utils.Login(input.Email, input.Password, u.db)
	userCrud := crud.NewUserCrud(u.db)
	crud := crud.NewCrud(u.db)
	if resp == "admin" {
		user, err := userCrud.GetUserByDocument(document)
		if err != nil {
			return user, err
		}
		vehicles, _ := crud.VehicleCrud.GetVehiclesByUserId(user.ID)
		recharges, _ := crud.RechargeCrud.GetRechargeByUserId(user.ID)

		for i := range recharges {
			billet, _ := crud.BilletCrud.GetBilletByRechargeId(recharges[i].ID)
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
	userCrud := crud.NewUserCrud(u.db)
	if resp == "user" || resp == "admin" {
		user, err := userCrud.GetUserByID(userID)
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
			err = userCrud.UpdateUser(user)
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
	resp := utils.Login(input.Email, input.Password, u.db)
	userCrud := crud.NewUserCrud(u.db)
	crud := crud.NewCrud(u.db)
	if resp == "user" || resp == "admin" {
		user, err := userCrud.GetUserByEmail(input.Email)
		if err != nil {
			return err
		}
		if resp == "user" && user.ID != userID {
			err := errors.New("usuário não possui permissão")
			return err
		} else {
			err := userCrud.DeleteUserByID(userID, crud)
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
