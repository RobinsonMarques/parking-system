package services

import (
	"errors"
	"github.com/RobinsonMarques/parking-system/database"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/interfaces"
)

func NewVehicleService(VehicleInterface interfaces.VehicleInterface, UserInterface interfaces.UserInterface, ParkingTicketInterface interfaces.ParkingTicketInterface, UtilInterface interfaces.UtilInterface) VehicleService {
	return VehicleService{
		vehicleInterface:       VehicleInterface,
		userInterface:          UserInterface,
		parkingTicketInterface: ParkingTicketInterface,
		utilInterface:          UtilInterface,
	}
}

type VehicleService struct {
	vehicleInterface       interfaces.VehicleInterface
	userInterface          interfaces.UserInterface
	parkingTicketInterface interfaces.ParkingTicketInterface
	utilInterface          interfaces.UtilInterface
}

func (v VehicleService) GetAllVehicles(input input2.LoginInput) ([]database.Vehicle, error) {
	resp := v.utilInterface.Login(input.Email, input.Password)
	if resp == "trafficWarden" || resp == "admin" {
		vehicles, err := v.vehicleInterface.GetAllVehicles()
		if err != nil {
			return []database.Vehicle{}, err
		}
		return vehicles, nil
	} else {
		err := errors.New(resp)
		return []database.Vehicle{}, err
	}
}

func (v VehicleService) CreateVehicle(input input2.CreateVehicle) error {
	_, err := v.userInterface.GetUserByID(input.UserID)
	if err == nil {
		//Cria o veículo
		veiculo := database.Vehicle{
			LicensePlate:  input.LicensePlate,
			VehicleModel:  input.VehicleModel,
			VehicleType:   input.VehicleType,
			IsActive:      false,
			IsParked:      false,
			UserID:        input.UserID,
			ParkingTicket: nil,
		}
		err := v.vehicleInterface.CreateVehicle(veiculo)
		if err == nil {
			return nil
		} else {
			return err
		}
	} else {
		err := errors.New("usuário não encontrado")
		return err
	}
}

func (v VehicleService) GetVehicleByLicensePlate(input input2.LoginInput, licensePlate string) (database.Vehicle, error) {
	resp := v.utilInterface.Login(input.Email, input.Password)
	if resp == "trafficWarden" {
		vehicle, err := v.vehicleInterface.GetVehicleByLicensePlate(licensePlate)
		if err != nil {
			return vehicle, err
		}
		ticket, err := v.parkingTicketInterface.GetLastParkingTicketFromVehicle(vehicle.ID)
		if err != nil {
			return vehicle, err
		}
		vehicle.ParkingTicket = ticket
		return vehicle, err
	} else {
		err := errors.New(resp)
		vehicle := database.Vehicle{}
		return vehicle, err
	}
}

func (v VehicleService) UpdateVehicle(input input2.UpdateVehicle) error {
	resp := v.utilInterface.Login(input.LoginInput.Email, input.LoginInput.Password)
	user, err := v.userInterface.GetUserByEmail(input.LoginInput.Email)
	if err != nil {
		return err
	}
	vehicles, err := v.vehicleInterface.GetVehiclesByUserId(user.ID)
	if err != nil {
		return err
	}
	var resp2 bool
	for i := range vehicles {
		if vehicles[i].UserID == user.ID {
			resp2 = true
		}
	}

	if resp == "user" {
		if resp2 {
			vehicle, err := v.vehicleInterface.GetVehicleByLicensePlate(input.LicensePlate)
			if err != nil {
				return err
			}
			vehicle.VehicleModel = input.VehicleModel
			vehicle.VehicleType = input.VehicleType
			vehicle.LicensePlate = input.NewLicensePlate
			err = v.vehicleInterface.UpdateVehicle(vehicle)
			if err != nil {
				return err
			}
			return nil
		} else {
			err := errors.New("usuário não possui este veículo")
			return err
		}
	} else {
		err := errors.New(resp)
		return err
	}
}

func (v VehicleService) UpdateVehicleOwner(input input2.UpdateVehicleOwner, vehicleID uint) error {
	resp := v.utilInterface.Login(input.LoginInput.Email, input.LoginInput.Password)
	user, err := v.userInterface.GetUserByID(input.NewUserID)
	if err != nil {
		return err
	}
	if resp == "admin" {
		if user.Person.Name != "" {
			err := v.vehicleInterface.UpdateVehicleOwner(vehicleID, input.NewUserID)
			if err != nil {
				return err
			}
			return nil
		} else {
			err := errors.New("usuário inexistente")
			return err
		}
	} else {
		err := errors.New(resp)
		return err
	}
}

func (v VehicleService) DeleteVehicleByID(input input2.LoginInput, vehicleID uint) error {
	resp := v.utilInterface.Login(input.Email, input.Password)
	if resp == "user" || resp == "admin" {
		vehicle, err := v.vehicleInterface.GetVehicleById(vehicleID)
		if err != nil {
			return err
		}
		user, err := v.userInterface.GetUserByEmail(input.Email)
		if err != nil {
			return err
		}
		if resp == "user" && vehicle.UserID != user.ID {
			err := errors.New("usuário logado não possui permissão")
			return err
		} else {
			err := v.vehicleInterface.DeleteVehicleByID(vehicleID)
			if err != nil {
				return err
			}
			err = v.parkingTicketInterface.DeleteParkingTicketByVehicleID(vehicleID)
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
