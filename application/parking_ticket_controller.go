package application

import (
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/interfaces"
	"github.com/RobinsonMarques/parking-system/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func NewParkingTicketManager(parkingTicketInterface interfaces.ParkingTicketInterface, vehicleInterface interfaces.VehicleInterface, userInterface interfaces.UserInterface, utilInterface interfaces.UtilInterface) ParkingTicketManager {
	return ParkingTicketManager{
		parkingTicketInterface: parkingTicketInterface,
		vehicleInterface:       vehicleInterface,
		userInterface:          userInterface,
		utilInterface:          utilInterface,
	}
}

type ParkingTicketManager struct {
	parkingTicketInterface interfaces.ParkingTicketInterface
	vehicleInterface       interfaces.VehicleInterface
	userInterface          interfaces.UserInterface
	utilInterface          interfaces.UtilInterface
}

func (a ParkingTicketManager) CreateParkingTicket(c *gin.Context) {
	parkingTicketService := services.NewParkingTicketService(a.parkingTicketInterface, a.vehicleInterface, a.userInterface, a.utilInterface)
	var input input2.CreateParkingTicket
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := parkingTicketService.CreateParkingTicket(input)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Response": "Ticket criado"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}

func (a ParkingTicketManager) DeleteParkingTicketByID(c *gin.Context) {
	parkingTicketService := services.NewParkingTicketService(a.parkingTicketInterface, a.vehicleInterface, a.userInterface, a.utilInterface)
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
