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

func (v VehicleService) GetAllVehicles(input input2.LoginInput) ([]database.Vehicle, error){
	resp := utils.Login(input.Email, input.Password, v.db)

	if resp == "trafficWarden" || resp == "admin" {
		vehicles := crud.GetAllVehicles(v.db)
		return vehicles, nil
	} else {
		err := errors.New(resp)
		var vehicles []database.Vehicle
		return vehicles, err
	}
}

func (v VehicleService)GetVehicleByLicensePlate(input input2.LoginInput, licensePlate string) (database.Vehicle, error){
	resp := utils.Login(input.Email, input.Password, v.db)

	if resp == "trafficWarden" {
		vehicle, err := crud.GetVehicleByLicensePlate(licensePlate, v.db)
		if err != nil{
			return vehicle, err
		}
		ticket, err := crud.GetLastParkingTicketFromVehicle(vehicle.ID, v.db)
		if err != nil{
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
