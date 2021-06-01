package services

import (
	"errors"
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/database"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/utils"
)

func NewTrafficWardenService(trafficWardenCrud crud.TrafficWardenCrud, utilCrud crud.UtilCrud) TrafficWardenService {
	return TrafficWardenService{
		trafficWardenCrud: trafficWardenCrud,
		utilCrud:          utilCrud,
	}
}

type TrafficWardenService struct {
	trafficWardenCrud crud.TrafficWardenCrud
	utilCrud          crud.UtilCrud
}

func (t TrafficWardenService) CreateTrafficWarden(input input2.CreateTrafficWarden, service TrafficWardenService) error {
	resp := service.utilCrud.Login(input.LoginInput.Email, input.LoginInput.Password)
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
		err = service.trafficWardenCrud.CreateTrafficWarden(warden)
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

func (t TrafficWardenService) UpdateTrafficWarden(input input2.UpdateTrafficWarden, wardenID uint, service TrafficWardenService) error {
	resp := service.utilCrud.Login(input.LoginInput.Email, input.LoginInput.Password)
	if resp == "trafficWarden" || resp == "admin" {
		trafficWarden, err := service.trafficWardenCrud.GetTrafficWardenByID(wardenID)
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
			err = service.trafficWardenCrud.UpdateTrafficWarden(trafficWarden)
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

func (t TrafficWardenService) DeleteTrafficWardenByID(input input2.LoginInput, wardenID uint, service TrafficWardenService) error {
	resp := service.utilCrud.Login(input.Email, input.Password)
	if resp == "admin" {
		err := service.trafficWardenCrud.DeleteTrafficWardenByID(wardenID)
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
