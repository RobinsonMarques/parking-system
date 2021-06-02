package application

import (
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/interfaces"
	"github.com/RobinsonMarques/parking-system/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func NewTrafficWardenManager(trafficWardenInterface interfaces.TrafficWardenInterface, utilInterface interfaces.UtilInterface) TrafficWardenManager {
	return TrafficWardenManager{
		trafficWardenInterface: trafficWardenInterface,
		utilInterface:          utilInterface,
	}
}

type TrafficWardenManager struct {
	trafficWardenInterface interfaces.TrafficWardenInterface
	utilInterface          interfaces.UtilInterface
}

func (a TrafficWardenManager) CreateTrafficWarden(c *gin.Context) {

	trafficWardenService := services.NewTrafficWardenService(a.trafficWardenInterface, a.utilInterface)
	//Valida o input
	var input input2.CreateTrafficWarden

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := trafficWardenService.CreateTrafficWarden(input)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Response": "Guarda de trânsito criado"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}

func (a TrafficWardenManager) UpdateTrafficWarden(c *gin.Context) {
	trafficWardenService := services.NewTrafficWardenService(a.trafficWardenInterface, a.utilInterface)

	wardenIDString := c.Param("wardenID")
	wardenIDInt, _ := strconv.Atoi(wardenIDString)
	wardenID := uint(wardenIDInt)

	var input input2.UpdateTrafficWarden
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := trafficWardenService.UpdateTrafficWarden(input, wardenID)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Response": "Guarda de trânsito alterado"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}

func (a TrafficWardenManager) DeleteTrafficWardenByID(c *gin.Context) {
	trafficWardenService := services.NewTrafficWardenService(a.trafficWardenInterface, a.utilInterface)

	wardenIDString := c.Param("trafficwardenID")
	wardenIDInt, _ := strconv.Atoi(wardenIDString)
	wardenID := uint(wardenIDInt)

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := trafficWardenService.DeleteTrafficWardenByID(input, wardenID)

	if err == nil {
		c.JSON(http.StatusOK, "Guarda deletado com sucesso!")
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}
