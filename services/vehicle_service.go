package services

import (
	"errors"
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/database"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/utils"
	"gorm.io/gorm"
)

func NewVehicleService(db *gorm.DB) VehicleService {
	return VehicleService{db: db}
}

type VehicleService struct {
	db *gorm.DB
}

func (v VehicleService) GetAllVehicles(input input2.LoginInput) ([]database.Vehicle, error) {
	resp := utils.Login(input.Email, input.Password, v.db)

	if resp == "trafficWarden" || resp == "admin" {
		vehicles, err := crud.GetAllVehicles(v.db)
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
	_, err := crud.GetUserByID(input.UserID, v.db)
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
		resp := crud.CreateVehicle(veiculo, v.db)
		if resp.Data.Error == nil {
			return nil
		} else {
			return resp.Data.Error
		}
	} else {
		err := errors.New("usuário não encontrado")
		return err
	}
}

func (v VehicleService) GetVehicleByLicensePlate(input input2.LoginInput, licensePlate string) (database.Vehicle, error) {
	resp := utils.Login(input.Email, input.Password, v.db)

	if resp == "trafficWarden" {
		vehicle, err := crud.GetVehicleByLicensePlate(licensePlate, v.db)
		if err != nil {
			return vehicle, err
		}
		ticket, err := crud.GetLastParkingTicketFromVehicle(vehicle.ID, v.db)
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
	resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, v.db)
	user, err := crud.GetUserByEmail(input.LoginInput.Email, v.db)
	if err != nil {
		return err
	}
	vehicles, err := crud.GetVehiclesByUserId(user.ID, v.db)
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
			vehicle, err := crud.GetVehicleByLicensePlate(input.LicensePlate, v.db)
			if err != nil {
				return err
			}
			vehicle.VehicleModel = input.VehicleModel
			vehicle.VehicleType = input.VehicleType
			vehicle.LicensePlate = input.NewLicensePlate
			err = crud.UpdateVehicle(vehicle, v.db)
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
	resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, v.db)
	user, err := crud.GetUserByID(input.NewUserID, v.db)
	if err != nil {
		return err
	}
	if resp == "admin" {
		if user.Person.Name != "" {
			crud.UpdateVehicleOwner(vehicleID, input.NewUserID, v.db)
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
	resp := utils.Login(input.Email, input.Password, v.db)

	if resp == "user" || resp == "admin" {
		vehicle, err := crud.GetVehicleById(vehicleID, v.db)
		if err != nil {
			return err
		}
		user, err := crud.GetUserByEmail(input.Email, v.db)
		if err != nil {
			return err
		}
		if resp == "user" && vehicle.UserID != user.ID {
			err := errors.New("usuário logado não possui permissão")
			return err
		} else {
			err := crud.DeleteVehicleByID(vehicleID, v.db)
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
