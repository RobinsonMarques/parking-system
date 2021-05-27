package services

import (
	"errors"
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/database"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/utils"
	"gorm.io/gorm"
)

func NewTrafficWardenService (db *gorm.DB) TrafficWardenService{
	return TrafficWardenService{db: db}
}

type TrafficWardenService struct {
	db *gorm.DB
}

func (t TrafficWardenService) CreateTrafficWarden(input input2.CreateTrafficWarden) error {
	resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, t.db)

	if resp == "admin" {
		input.Person.Password = utils.CreateHashPassword(input.Person.Password)

		//Cria o traffic warden
		warden := database.TrafficWarden{
			Person: input.Person,
		}
		resp := crud.CreateTrafficWarden(warden, t.db)
		if resp.Data.Error != nil {
			return nil
		}else{
			return resp.Data.Error
		}
	} else {
		err := errors.New(resp)
		return err
	}
}

func (t TrafficWardenService) UpdateTrafficWarden( input input2.UpdateTrafficWarden, wardenID uint) error{
	resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, t.db)
	if resp == "trafficWarden" || resp == "admin" {
		trafficWarden, err := crud.GetTrafficWardenByID(wardenID, t.db)
		if err != nil{
			return err
		}
		if resp == "trafficWarden" && trafficWarden.Person.Email != input.LoginInput.Email {
			err := errors.New("usuário não possui permissão")
			return err
		} else {
			input.Person.Password = utils.CreateHashPassword(input.Person.Password)
			trafficWarden.Person = input.Person
			crud.UpdateTrafficWarden(trafficWarden, t.db)
			return nil
		}
	} else {
		err := errors.New(resp)
		return err
	}
}


func (t TrafficWardenService)DeleteTrafficWardenByID(input input2.LoginInput, wardenID uint) error{
	resp := utils.Login(input.Email, input.Password, t.db)

	if resp == "admin" {
		err := crud.DeleteTrafficWardenByID(wardenID, t.db)
		if err == nil {
			return nil
		}else{
			return err
		}
	} else {
		err := errors.New(resp)
		return err
	}
}
