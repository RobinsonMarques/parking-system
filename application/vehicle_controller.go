package application

import (
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func NewVehicleManager(vehicleService services.VehicleService) VehicleManager {
	return VehicleManager{
		vehicleService: vehicleService,
	}
}

type VehicleManager struct {
	vehicleService services.VehicleService
}

func (a VehicleManager) GetAllVehicles(c *gin.Context) {
	//valida o input
	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vehicles, err := a.vehicleService.GetAllVehicles(input)
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
	err := a.vehicleService.CreateVehicle(input)
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
	vehicle, err := a.vehicleService.GetVehicleByLicensePlate(input, licensePlate)
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
	err := a.vehicleService.UpdateVehicle(input)
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
	err := a.vehicleService.UpdateVehicleOwner(input, vehicleID)
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
	err := a.vehicleService.DeleteVehicleByID(input, vehicleID)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Response": "Veículo deletado"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}
