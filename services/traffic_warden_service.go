package services

import (
	"errors"
	"github.com/RobinsonMarques/parking-system/database"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/interfaces"
	"github.com/RobinsonMarques/parking-system/utils"
)

func NewTrafficWardenService(trafficWardenInterface interfaces.TrafficWardenInterface, utilInterface interfaces.UtilInterface) TrafficWardenService {
	return TrafficWardenService{
		trafficWardenInterface: trafficWardenInterface,
		utilInterface:          utilInterface,
	}
}

type TrafficWardenService struct {
	trafficWardenInterface interfaces.TrafficWardenInterface
	utilInterface          interfaces.UtilInterface
}

func (t TrafficWardenService) CreateTrafficWarden(input input2.CreateTrafficWarden) error {
	resp := t.utilInterface.Login(input.LoginInput.Email, input.LoginInput.Password)
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
		err = t.trafficWardenInterface.CreateTrafficWarden(warden)
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
	resp := t.utilInterface.Login(input.LoginInput.Email, input.LoginInput.Password)
	if resp == "trafficWarden" || resp == "admin" {
		trafficWarden, err := t.trafficWardenInterface.GetTrafficWardenByID(wardenID)
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
			err = t.trafficWardenInterface.UpdateTrafficWarden(trafficWarden)
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

func (t TrafficWardenService) DeleteTrafficWardenByID(input input2.LoginInput, wardenID uint) error {
	resp := t.utilInterface.Login(input.Email, input.Password)
	if resp == "admin" {
		err := t.trafficWardenInterface.DeleteTrafficWardenByID(wardenID)
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
