package services

import (
	"errors"
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/database"
	input2 "github.com/RobinsonMarques/parking-system/input"
)

func NewVehicleService(VehicleCrud crud.VehicleCrud, UserCrud crud.UserCrud, ParkingTicketCrud crud.ParkingTicketCrud, UtilCrud crud.UtilCrud) VehicleService {
	return VehicleService{
		vehicleCrud:       VehicleCrud,
		userCrud:          UserCrud,
		parkingTicketCrud: ParkingTicketCrud,
		util:              UtilCrud,
	}
}

type VehicleService struct {
	vehicleCrud       crud.VehicleCrud
	userCrud          crud.UserCrud
	parkingTicketCrud crud.ParkingTicketCrud
	util              crud.UtilCrud
}

func (v VehicleService) GetAllVehicles(input input2.LoginInput, service VehicleService) ([]database.Vehicle, error) {
	resp := service.util.Login(input.Email, input.Password)
	if resp == "trafficWarden" || resp == "admin" {
		vehicles, err := service.vehicleCrud.GetAllVehicles()
		if err != nil {
			return []database.Vehicle{}, err
		}
		return vehicles, nil
	} else {
		err := errors.New(resp)
		return []database.Vehicle{}, err
	}
}

func (v VehicleService) CreateVehicle(input input2.CreateVehicle, service VehicleService) error {
	_, err := service.userCrud.GetUserByID(input.UserID)
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
		err := service.vehicleCrud.CreateVehicle(veiculo)
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

func (v VehicleService) GetVehicleByLicensePlate(input input2.LoginInput, licensePlate string, service VehicleService) (database.Vehicle, error) {
	resp := service.util.Login(input.Email, input.Password)
	if resp == "trafficWarden" {
		vehicle, err := service.vehicleCrud.GetVehicleByLicensePlate(licensePlate)
		if err != nil {
			return vehicle, err
		}
		ticket, err := service.parkingTicketCrud.GetLastParkingTicketFromVehicle(vehicle.ID)
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

func (v VehicleService) UpdateVehicle(input input2.UpdateVehicle, service VehicleService) error {
	resp := service.util.Login(input.LoginInput.Email, input.LoginInput.Password)
	user, err := service.userCrud.GetUserByEmail(input.LoginInput.Email)
	if err != nil {
		return err
	}
	vehicles, err := service.vehicleCrud.GetVehiclesByUserId(user.ID)
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
			vehicle, err := service.vehicleCrud.GetVehicleByLicensePlate(input.LicensePlate)
			if err != nil {
				return err
			}
			vehicle.VehicleModel = input.VehicleModel
			vehicle.VehicleType = input.VehicleType
			vehicle.LicensePlate = input.NewLicensePlate
			err = service.vehicleCrud.UpdateVehicle(vehicle)
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

func (v VehicleService) UpdateVehicleOwner(input input2.UpdateVehicleOwner, vehicleID uint, service VehicleService) error {
	resp := service.util.Login(input.LoginInput.Email, input.LoginInput.Password)
	user, err := service.userCrud.GetUserByID(input.NewUserID)
	if err != nil {
		return err
	}
	if resp == "admin" {
		if user.Person.Name != "" {
			err := service.vehicleCrud.UpdateVehicleOwner(vehicleID, input.NewUserID)
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

func (v VehicleService) DeleteVehicleByID(input input2.LoginInput, vehicleID uint, service VehicleService) error {
	resp := service.util.Login(input.Email, input.Password)
	if resp == "user" || resp == "admin" {
		vehicle, err := service.vehicleCrud.GetVehicleById(vehicleID)
		if err != nil {
			return err
		}
		user, err := service.userCrud.GetUserByEmail(input.Email)
		if err != nil {
			return err
		}
		if resp == "user" && vehicle.UserID != user.ID {
			err := errors.New("usuário logado não possui permissão")
			return err
		} else {
			err := service.vehicleCrud.DeleteVehicleByID(vehicleID)
			if err != nil {
				return err
			}
			err = service.parkingTicketCrud.DeleteParkingTicketByVehicleID(vehicleID)
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
