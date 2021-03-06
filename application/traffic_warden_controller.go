package application

import (
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func NewTrafficWardenManager(trafficWardenService services.TrafficWardenService) TrafficWardenManager {
	return TrafficWardenManager{
		trafficWardenService: trafficWardenService,
	}
}

type TrafficWardenManager struct {
	trafficWardenService services.TrafficWardenService
}

func (a TrafficWardenManager) CreateTrafficWarden(c *gin.Context) {
	//Valida o input
	var input input2.CreateTrafficWarden

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.trafficWardenService.CreateTrafficWarden(input)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Response": "Guarda de trânsito criado"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}

func (a TrafficWardenManager) UpdateTrafficWarden(c *gin.Context) {
	wardenIDString := c.Param("wardenID")
	wardenIDInt, _ := strconv.Atoi(wardenIDString)
	wardenID := uint(wardenIDInt)

	var input input2.UpdateTrafficWarden
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.trafficWardenService.UpdateTrafficWarden(input, wardenID)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Response": "Guarda de trânsito alterado"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}

func (a TrafficWardenManager) DeleteTrafficWardenByID(c *gin.Context) {
	wardenIDString := c.Param("trafficwardenID")
	wardenIDInt, _ := strconv.Atoi(wardenIDString)
	wardenID := uint(wardenIDInt)

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := a.trafficWardenService.DeleteTrafficWardenByID(input, wardenID)

	if err == nil {
		c.JSON(http.StatusOK, "Guarda deletado com sucesso!")
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}
