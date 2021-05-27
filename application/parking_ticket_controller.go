package application

import (
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/database"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

func NewParkingTicketManager(db *gorm.DB) ParkingTicketManager {
	return ParkingTicketManager{db: db}
}

type ParkingTicketManager struct {
	db *gorm.DB
}

func (a ParkingTicketManager) CreateParkingTicket(c *gin.Context) {
	var input input2.CreateParkingTicket
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := utils.Login(input.Login.Email, input.Login.Password, a.db)
	resp2 := crud.GetVehicleById(input.VehicleID, a.db)
	user := crud.GetUserByEmail(input.Login.Email, a.db)
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
						crud.CreateParkingTicket(ticket, a.db)
						crud.UpdateIsParked(input.VehicleID, true, a.db)
						crud.UpdateBalance(input.Login.Email, -price, a.db)
						crud.UpdateIsActive(input.VehicleID, true, a.db)
						go utils.AlterVehicleStatus(resp2, input.ParkingTime, a.db)
						c.JSON(http.StatusOK, gin.H{"Response": "Ticket criado"})
					} else {
						c.JSON(http.StatusBadRequest, gin.H{"Response": "Saldo insuficiente"})
					}
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"Response": "Veículo já estacionado"})
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"Response": "Veículo não encontrado"})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"Response": "Usuário não possui o veículo"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"Response": resp})
	}
}

func (a ParkingTicketManager) DeleteParkingTicketByID(c *gin.Context) {
	ticketIDString := c.Param("ticketID")
	ticketIDInt, _ := strconv.Atoi(ticketIDString)
	ticketID := uint(ticketIDInt)

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := utils.Login(input.Email, input.Password, a.db)

	if resp == "admin" {
		crud.DeleteParkingTicketByID(ticketID, a.db)
		c.JSON(http.StatusOK, gin.H{"Response": "Ticket deletado"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
	}
}
