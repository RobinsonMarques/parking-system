package application

import (
	"github.com/RobinsonMarques/parking-system/crud"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func NewParkingTicketManager(parkingTicketCrud crud.ParkingTicketCrud, vehicleCrud crud.VehicleCrud, userCrud crud.UserCrud, utilCrud crud.UtilCrud) ParkingTicketManager {
	return ParkingTicketManager{
		parkingTicketCrud: parkingTicketCrud,
		vehicleCrud:       vehicleCrud,
		userCrud:          userCrud,
		utilCrud:          utilCrud,
	}
}

type ParkingTicketManager struct {
	parkingTicketCrud crud.ParkingTicketCrud
	vehicleCrud       crud.VehicleCrud
	userCrud          crud.UserCrud
	utilCrud          crud.UtilCrud
}

func (a ParkingTicketManager) CreateParkingTicket(c *gin.Context) {
	parkingTicketService := services.NewParkingTicketService(a.parkingTicketCrud, a.vehicleCrud, a.userCrud, a.utilCrud)
	var input input2.CreateParkingTicket
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := parkingTicketService.CreateParkingTicket(input, parkingTicketService)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Response": "Ticket criado"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}

func (a ParkingTicketManager) DeleteParkingTicketByID(c *gin.Context) {
	parkingTicketService := services.NewParkingTicketService(a.parkingTicketCrud, a.vehicleCrud, a.userCrud, a.utilCrud)
	ticketIDString := c.Param("ticketID")
	ticketIDInt, _ := strconv.Atoi(ticketIDString)
	ticketID := uint(ticketIDInt)

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := parkingTicketService.DeleteParkingTicketByID(input, ticketID, parkingTicketService)
	if err == nil {
		c.JSON(http.StatusOK, "Usuário deletado com sucesso!")
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}
