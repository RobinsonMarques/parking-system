package application

import (
	"github.com/RobinsonMarques/parking-system/crud"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func NewTrafficWardenManager(db *gorm.DB) TrafficWardenManager {
	return TrafficWardenManager{db: db}
}

type TrafficWardenManager struct {
	db *gorm.DB
}

func (a TrafficWardenManager) CreateTrafficWarden(c *gin.Context) {
	trafficWardenCrud := crud.NewTrafficWardenCrud(a.db)
	utilCrud := crud.NewUtilCrud(a.db)
	trafficWardenService := services.NewTrafficWardenService(trafficWardenCrud, utilCrud)
	//Valida o input
	var input input2.CreateTrafficWarden

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := trafficWardenService.CreateTrafficWarden(input, trafficWardenService)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Response": "Guarda de trânsito criado"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}

func (a TrafficWardenManager) UpdateTrafficWarden(c *gin.Context) {
	trafficWardenCrud := crud.NewTrafficWardenCrud(a.db)
	utilCrud := crud.NewUtilCrud(a.db)
	trafficWardenService := services.NewTrafficWardenService(trafficWardenCrud, utilCrud)

	wardenIDString := c.Param("wardenID")
	wardenIDInt, _ := strconv.Atoi(wardenIDString)
	wardenID := uint(wardenIDInt)

	var input input2.UpdateTrafficWarden
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := trafficWardenService.UpdateTrafficWarden(input, wardenID, trafficWardenService)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Response": "Guarda de trânsito alterado"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}

func (a TrafficWardenManager) DeleteTrafficWardenByID(c *gin.Context) {
	trafficWardenCrud := crud.NewTrafficWardenCrud(a.db)
	utilCrud := crud.NewUtilCrud(a.db)
	trafficWardenService := services.NewTrafficWardenService(trafficWardenCrud, utilCrud)

	wardenIDString := c.Param("trafficwardenID")
	wardenIDInt, _ := strconv.Atoi(wardenIDString)
	wardenID := uint(wardenIDInt)

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := trafficWardenService.DeleteTrafficWardenByID(input, wardenID, trafficWardenService)

	if err == nil {
		c.JSON(http.StatusOK, "Guarda deletado com sucesso!")
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}
