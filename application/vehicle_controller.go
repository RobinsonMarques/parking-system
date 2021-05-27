package application

import (
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/database"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/services"
	"github.com/RobinsonMarques/parking-system/utils"
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
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err})
	}

}

func (a VehicleManager) CreateVehicle(c *gin.Context) {
	var input input2.CreateVehicle
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if crud.GetUserByID(input.UserID, a.db).Person.Name != "" {
		//Cria o veículo
		veiculo := database.Vehicle{
			LicensePlate:  input.LicensePlate,
			VehicleModel:  input.VehicleModel,
			VehicleType:   input.VehicleType,
			IsActive:      false,
			IsParked:      false,
			UserID:        input.UserID,
			ParkingTicket: nil,
		}
		crud.CreateVehicle(veiculo, a.db)
		c.JSON(http.StatusOK, gin.H{"Response": "Veículo criado"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"Response": "Usuário não encontrado"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err})
	}
}

func (a VehicleManager) UpdateVehicle(c *gin.Context) {
	var input input2.UpdateVehicle
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, a.db)
	user := crud.GetUserByEmail(input.LoginInput.Email, a.db)
	vehicles := crud.GetVehiclesByUserId(user.ID, a.db)
	var resp2 bool
	for i := range vehicles {
		if vehicles[i].UserID == user.ID {
			resp2 = true
		}
	}

	if resp == "user" {
		if resp2 {
			vehicle := crud.GetVehicleByLicensePlate(input.LicensePlate, a.db)
			vehicle.VehicleModel = input.VehicleModel
			vehicle.VehicleType = input.VehicleType
			vehicle.LicensePlate = input.NewLicensePlate
			crud.UpdateVehicle(vehicle, a.db)
			c.JSON(http.StatusOK, gin.H{"Response": "Veículo alterado"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Usuário não possui este veículo"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
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

	resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, a.db)
	user := crud.GetUserByID(input.NewUserID, a.db)
	if resp == "admin" {
		if user.Person.Name != "" {
			crud.UpdateVehicleOwner(vehicleID, input.NewUserID, a.db)
			c.JSON(http.StatusOK, gin.H{"Response": "Dono do veículo alterado"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Usuário inexistente"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
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
	resp := utils.Login(input.Email, input.Password, a.db)

	if resp == "user" || resp == "admin" {
		vehicle := crud.GetVehicleById(vehicleID, a.db)
		user := crud.GetUserByEmail(input.Email, a.db)
		if resp == "user" && vehicle.UserID != user.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Usuário logado não possui permissão"})
		} else {
			crud.DeleteVehicleByID(vehicleID, a.db)
			c.JSON(http.StatusOK, gin.H{"Response": "Veículo deletado"})
		}
	}
}
