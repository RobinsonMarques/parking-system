package services

import (
	"errors"
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/database"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/utils"
	"gorm.io/gorm"
	"time"
)

func NewParkingTicketService(db *gorm.DB) ParkingTicketService {
	return ParkingTicketService{db: db}
}

type ParkingTicketService struct {
	db *gorm.DB
}

func (p ParkingTicketService) CreateParkingTicket(input input2.CreateParkingTicket) error {
	resp := utils.Login(input.Login.Email, input.Login.Password, p.db)
	parkingTicketCrud := crud.NewParkingTicketCrud(p.db)
	crud := crud.NewCrud(p.db)
	resp2, err := crud.VehicleCrud.GetVehicleById(input.VehicleID)
	if err != nil {
		return err
	}
	user, err := crud.UserCrud.GetUserByEmail(input.Login.Email)
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
						err := parkingTicketCrud.CreateParkingTicket(ticket)
						if err != nil {
							return err
						}
						err2 := crud.VehicleCrud.UpdateIsParked(input.VehicleID, true)
						if err2 != nil {
							return nil
						}
						err2 = crud.UserCrud.UpdateBalance(input.Login.Email, -price, crud)
						if err2 != nil {
							return err2
						}
						err2 = crud.VehicleCrud.UpdateIsActive(input.VehicleID, true)
						if err2 != nil {
							return err2
						}
						go utils.AlterVehicleStatus(resp2, input.ParkingTime, p.db)
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

func (p ParkingTicketService) DeleteParkingTicketByID(input input2.LoginInput, ticketID uint) error {
	resp := utils.Login(input.Email, input.Password, p.db)
	parkingTicketCrud := crud.NewParkingTicketCrud(p.db)
	if resp == "admin" {
		err := parkingTicketCrud.DeleteParkingTicketByID(ticketID)
		if err != nil {
			return err
		}
		return nil
	} else {
		err := errors.New(resp)
		return err
	}
}
