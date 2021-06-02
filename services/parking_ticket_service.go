package services

import (
	"errors"
	"github.com/RobinsonMarques/parking-system/database"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/interfaces"
	"time"
)

func NewParkingTicketService(parkingTicketInterface interfaces.ParkingTicketInterface, vehicleInterface interfaces.VehicleInterface, userInterface interfaces.UserInterface, utilInterface interfaces.UtilInterface) ParkingTicketService {
	return ParkingTicketService{
		parkingTicketInterface: parkingTicketInterface,
		vehicleInterface:       vehicleInterface,
		userInterface:          userInterface,
		utilInterface:          utilInterface,
	}
}

type ParkingTicketService struct {
	parkingTicketInterface interfaces.ParkingTicketInterface
	vehicleInterface       interfaces.VehicleInterface
	userInterface          interfaces.UserInterface
	utilInterface          interfaces.UtilInterface
}

func (p ParkingTicketService) CreateParkingTicket(input input2.CreateParkingTicket) error {
	resp := p.utilInterface.Login(input.Login.Email, input.Login.Password)
	resp2, err := p.vehicleInterface.GetVehicleById(input.VehicleID)
	if err != nil {
		return err
	}
	user, err := p.userInterface.GetUserByEmail(input.Login.Email)
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
						err := p.parkingTicketInterface.CreateParkingTicket(ticket)
						if err != nil {
							return err
						}
						err2 := p.vehicleInterface.UpdateIsParked(input.VehicleID, true)
						if err2 != nil {
							return nil
						}
						err2 = p.userInterface.UpdateBalance(input.Login.Email, -price)
						if err2 != nil {
							return err2
						}
						err2 = p.vehicleInterface.UpdateIsActive(input.VehicleID, true)
						if err2 != nil {
							return err2
						}
						go p.vehicleInterface.AlterVehicleStatus(resp2, input.ParkingTime)
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
	resp := p.utilInterface.Login(input.Email, input.Password)
	if resp == "admin" {
		err := service.parkingTicketInterface.DeleteParkingTicketByID(ticketID)
		if err != nil {
			return err
		}
		return nil
	} else {
		err := errors.New(resp)
		return err
	}
}
