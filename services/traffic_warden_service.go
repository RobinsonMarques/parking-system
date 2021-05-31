package services

import (
	"errors"
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/database"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/utils"
	"gorm.io/gorm"
)

func NewTrafficWardenService(db *gorm.DB) TrafficWardenService {
	return TrafficWardenService{db: db}
}

type TrafficWardenService struct {
	db *gorm.DB
}

func (t TrafficWardenService) CreateTrafficWarden(input input2.CreateTrafficWarden) error {
	resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, t.db)
	trafficWardenCrud := crud.NewTrafficWardenCrud(t.db)
	if resp == "admin" {
		var err error
		input.Person.Password, err = utils.CreateHashPassword(input.Person.Password)
		if err != nil {
			return err
		}

		//Cria o traffic warden
		warden := database.TrafficWarden{
			Person: input.Person,
		}
		err = trafficWardenCrud.CreateTrafficWarden(warden)
		if err != nil {
			return nil
		} else {
			return err
		}
	} else {
		err := errors.New(resp)
		return err
	}
}

func (t TrafficWardenService) UpdateTrafficWarden(input input2.UpdateTrafficWarden, wardenID uint) error {
	resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, t.db)
	trafficWardenCrud := crud.NewTrafficWardenCrud(t.db)
	if resp == "trafficWarden" || resp == "admin" {
		trafficWarden, err := trafficWardenCrud.GetTrafficWardenByID(wardenID)
		if err != nil {
			return err
		}
		if resp == "trafficWarden" && trafficWarden.Person.Email != input.LoginInput.Email {
			err := errors.New("usuário não possui permissão")
			return err
		} else {
			var err error
			input.Person.Password, err = utils.CreateHashPassword(input.Person.Password)
			if err != nil {
				return err
			}
			trafficWarden.Person = input.Person
			trafficWardenCrud.UpdateTrafficWarden(trafficWarden)
			return nil
		}
	} else {
		err := errors.New(resp)
		return err
	}
}

func (t TrafficWardenService) DeleteTrafficWardenByID(input input2.LoginInput, wardenID uint) error {
	resp := utils.Login(input.Email, input.Password, t.db)
	trafficWardenCrud := crud.NewTrafficWardenCrud(t.db)
	if resp == "admin" {
		err := trafficWardenCrud.DeleteTrafficWardenByID(wardenID)
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
