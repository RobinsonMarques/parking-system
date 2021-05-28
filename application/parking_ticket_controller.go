package application

import (
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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
	parkingTicketService := services.NewParkingTicketService(a.db)
	err := parkingTicketService.CreateParkingTicket(input)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Response": "Ticket criado"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
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
	parkingTicketService := services.NewParkingTicketService(a.db)
	err := parkingTicketService.DeleteParkingTicketByID(input, ticketID)
	if err == nil {
		c.JSON(http.StatusOK, "Usuário deletado com sucesso!")
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}
