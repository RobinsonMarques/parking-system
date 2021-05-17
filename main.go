package main

import (
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/database"
	"github.com/RobinsonMarques/parking-system/dependencies"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	db := dependencies.CreateConnection()

	r := gin.Default()

	//Path get all vehicles
	r.GET("/vehicles", func(c *gin.Context) {
		//valida o input
		var input input2.LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := utils.Login(input.Email, input.Password, db)

		if resp == "trafficWarden" || resp == "admin" {
			vehicles := crud.GetAllVehicles(db)
			c.JSON(200, vehicles)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		}

	})

	//Path create user
	r.POST("/users", func(c *gin.Context) {

		//Valida o input
		var input input2.CreateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		input.Person.Password = utils.CreateHashPassword(input.Person.Password)

		//Cria o user

		user := database.User{
			Person:   input.Person,
			Document: input.Document,
			Balance:  0,
			Recharge: nil,
			Vehicle:  nil,
		}
		crud.CreateUser(user, db)
		c.JSON(http.StatusOK, gin.H{"Response": "Usuário criado"})

	})

	//Path create admin
	r.POST("/admins", func(c *gin.Context) {
		//Valida o input
		var input input2.CreateAdminInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, db)

		if resp == "admin" {
			input.Person.Password = utils.CreateHashPassword(input.Person.Password)

			//Cria o admin
			admin := database.Admin{
				Person: input.Person,
			}
			crud.CreateAdmin(admin, db)
			c.JSON(http.StatusOK, gin.H{"Response": "Admin criado"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		}
	})

	//Path create traffic warden
	r.POST("/trafficwarden", func(c *gin.Context) {
		//Valida o input
		var input input2.CreateTrafficWarden

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, db)

		if resp == "admin" {
			input.Person.Password = utils.CreateHashPassword(input.Person.Password)

			//Cria o admin
			warden := database.TrafficWarden{
				Person: input.Person,
			}
			crud.CreateTrafficWarden(warden, db)
			c.JSON(http.StatusOK, gin.H{"Response": "Guarda criado"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		}
	})

	//Path create vehicle
	r.POST("/vehicles", func(c *gin.Context) {
		var input input2.CreateVehicle
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if crud.GetUserByID(input.UserID, db).Person.Name != "" {
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
			crud.CreateVehicle(veiculo, db)
			c.JSON(http.StatusOK, gin.H{"Response": "Veículo criado"})
		} else {
			c.JSON(http.StatusOK, gin.H{"Response": "Usuário não encontrado"})
		}

	})

	//Path create parking ticket
	r.POST("/parkingtickets", func(c *gin.Context) {
		var input input2.CreateParkingTicket
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := utils.Login(input.Login.Email, input.Login.Password, db)
		resp2 := crud.GetVehicleById(input.VehicleID, db)
		user := crud.GetUserByEmail(input.Login.Email, db)
		if resp == "user" {
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
						crud.CreateParkingTicket(ticket, db)
						crud.UpdateIsParked(input.VehicleID, true, db)
						crud.UpdateBalance(input.Login.Email, -price, db)
						c.JSON(http.StatusOK, gin.H{"Response": "Ticket criado"})
					} else {
						c.JSON(http.StatusOK, gin.H{"Response": "Saldo insuficiente"})
					}
				} else {
					c.JSON(http.StatusOK, gin.H{"Response": "Veículo já estacionado"})
				}
			} else {
				c.JSON(http.StatusOK, gin.H{"Response": "Veículo não encontrado"})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"Response": "Usuário não encontrado"})
		}

	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
