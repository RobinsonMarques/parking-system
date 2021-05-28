package application

import (
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func NewVehicleManager(db *gorm.DB) VehicleManager {
	return VehicleManager{db: db}
}

type VehicleManager struct {
	db *gorm.DB
}

func (a VehicleManager) GetAllVehicles(c *gin.Context) {
	//valida o input
	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vehicleService := services.NewVehicleService(a.db)
	vehicles, err := vehicleService.GetAllVehicles(input)
	if err == nil {
		c.JSON(200, vehicles)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}

}

func (a VehicleManager) CreateVehicle(c *gin.Context) {
	var input input2.CreateVehicle
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vehicleService := services.NewVehicleService(a.db)
	err := vehicleService.CreateVehicle(input)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Response": "Veículo criado"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}

func (a VehicleManager) GetVehicleByLicensePlate(c *gin.Context) {
	licensePlate := c.Param("licensePlate")

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vehicleService := services.NewVehicleService(a.db)
	vehicle, err := vehicleService.GetVehicleByLicensePlate(input, licensePlate)
	if err == nil {
		c.JSON(200, vehicle)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}

func (a VehicleManager) UpdateVehicle(c *gin.Context) {
	var input input2.UpdateVehicle
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vehicleService := services.NewVehicleService(a.db)
	err := vehicleService.UpdateVehicle(input)
	if err == nil {
		c.JSON(http.StatusOK, "Usuário alterado com sucesso")
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}

func (a VehicleManager) UpdateVehicleOwner(c *gin.Context) {
	vehicleIDString := c.Param("vehicleID")
	vehicleIDInt, _ := strconv.Atoi(vehicleIDString)
	vehicleID := uint(vehicleIDInt)
	var input input2.UpdateVehicleOwner
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vehicleService := services.NewVehicleService(a.db)
	err := vehicleService.UpdateVehicleOwner(input, vehicleID)
	if err == nil {
		c.JSON(http.StatusOK, "Dono do veículo alterado com sucesso")
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}

func (a VehicleManager) DeleteVehicleByID(c *gin.Context) {
	vehicleIDString := c.Param("vehicleID")
	vehicleIDInt, _ := strconv.Atoi(vehicleIDString)
	vehicleID := uint(vehicleIDInt)

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vehicleService := services.NewVehicleService(a.db)
	err := vehicleService.DeleteVehicleByID(input, vehicleID)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Response": "Veículo deletado"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}
