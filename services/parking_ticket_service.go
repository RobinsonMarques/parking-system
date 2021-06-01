package services

import (
	"errors"
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/database"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"time"
)

func NewParkingTicketService(parkingTicketCrud crud.ParkingTicketCrud, vehicleCrud crud.VehicleCrud, userCrud crud.UserCrud, utilCrud crud.UtilCrud) ParkingTicketService {
	return ParkingTicketService{
		parkingTicketCrud: parkingTicketCrud,
		vehicleCrud:       vehicleCrud,
		userCrud:          userCrud,
		utilCrud:          utilCrud,
	}
}

type ParkingTicketService struct {
	parkingTicketCrud crud.ParkingTicketCrud
	vehicleCrud       crud.VehicleCrud
	userCrud          crud.UserCrud
	utilCrud          crud.UtilCrud
}

func (p ParkingTicketService) CreateParkingTicket(input input2.CreateParkingTicket, service ParkingTicketService) error {
	resp := service.utilCrud.Login(input.Login.Email, input.Login.Password)
	resp2, err := service.vehicleCrud.GetVehicleById(input.VehicleID)
	if err != nil {
		return err
	}
	user, err := service.userCrud.GetUserByEmail(input.Login.Email)
	if err != nil {
		return err
	}
	if resp == "user" {
		if resp2.UserID == user.ID {
			if resp2.LicensePlate != "" {
				if !resp2.IsParked {
					price := float64(input.ParkingTime) * 1.75
					currentTime := time.Now()
					endTime := currentTime.Add(time.Hour * time.Duration(input.ParkingTime))
					if user.Balance > price {
						ticket := database.ParkingTicket{
							Location:    input.Location,
							ParkingTime: input.ParkingTime,
							StartTime:   currentTime.String(),
							EndTime:     endTime.String(),
							Price:       price,
							VehicleID:   input.VehicleID,
						}
						err := service.parkingTicketCrud.CreateParkingTicket(ticket)
						if err != nil {
							return err
						}
						err2 := service.vehicleCrud.UpdateIsParked(input.VehicleID, true)
						if err2 != nil {
							return nil
						}
						err2 = service.userCrud.UpdateBalance(input.Login.Email, -price)
						if err2 != nil {
							return err2
						}
						err2 = service.vehicleCrud.UpdateIsActive(input.VehicleID, true)
						if err2 != nil {
							return err2
						}
						go service.vehicleCrud.AlterVehicleStatus(resp2, input.ParkingTime)
						return nil
					} else {
						err := errors.New("saldo insuficiente")
						return err
					}
				} else {
					err := errors.New("veículo já estacionado")
					return err
				}
			} else {
				err := errors.New("veículo não encontrado")
				return err
			}
		} else {
			err := errors.New("usuário não possui o veículo")
			return err
		}
	} else {
		err := errors.New(resp)
		return err
	}
}

func (p ParkingTicketService) DeleteParkingTicketByID(input input2.LoginInput, ticketID uint, service ParkingTicketService) error {
	resp := service.utilCrud.Login(input.Email, input.Password)
	if resp == "admin" {
		err := service.parkingTicketCrud.DeleteParkingTicketByID(ticketID)
		if err != nil {
			return err
		}
		return nil
	} else {
		err := errors.New(resp)
		return err
	}
}
