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
			c.JSON(http.StatusBadRequest, gin.H{"Response": "Usuário não encontrado"})
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
						crud.UpdateIsActive(input.VehicleID, true, db)
						c.JSON(http.StatusOK, gin.H{"Response": "Ticket criado"})
					} else {
						c.JSON(http.StatusBadRequest, gin.H{"Response": "Saldo insuficiente"})
					}
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"Response": "Veículo já estacionado"})
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"Response": "Veículo não encontrado"})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"Response": "Usuário não encontrado"})
		}

	})

	//Path create recharge
	r.POST("/recharge", func(c *gin.Context) {
		var input input2.CreateRecharge
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, db)
		user := crud.GetUserByEmail(input.LoginInput.Email, db)
		if resp == "user" {
			date := time.Now()
			recharge := database.Recharge{
				Date:        date.String(),
				Value:       input.Value,
				IsPaid:      false,
				PaymentType: input.PaymentType,
				UserID:      user.ID,
			}
			crud.CreateRecharge(recharge, db)
			user := crud.GetUserByEmail(input.LoginInput.Email, db)
			rechargeReturn := crud.GetRechargeByUserId(user.ID, db)
			len := len(rechargeReturn)
			billet := database.Billet{
				BilletLink: "link@link.com",
				RechargeID: rechargeReturn[len-1].ID,
			}
			crud.CreateBillet(billet, db)
			c.JSON(http.StatusOK, gin.H{"Response": "Recarga criada"})
		}

	})

	//Path get vehicle by license plate
	r.GET("/vehicles/:licensePlate", func(c *gin.Context) {
		licensePlate := c.Param("licensePlate")

		var input input2.LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := utils.Login(input.Email, input.Password, db)

		if resp == "trafficWarden" {
			vehicle := crud.GetVehicleByLicensePlate(licensePlate, db)
			ticket := crud.GetLastParkingTicketFromVehicle(vehicle.ID, db)
			vehicle.ParkingTicket = ticket
			c.JSON(200, vehicle)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		}

	})

	//Path get user by document
	r.GET("users/:document", func(c *gin.Context) {
		document := c.Param("document")

		var input input2.LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := utils.Login(input.Email, input.Password, db)

		if resp == "admin" {
			user := crud.GetUserByDocument(document, db)
			vehicles := crud.GetVehiclesByUserId(user.ID, db)
			recharges := crud.GetRechargeByUserId(user.ID, db)

			for i := range recharges {
				billet := crud.GetBilletByRechargeId(recharges[i].ID, db)
				recharges[i].Billet = billet
			}
			user.Vehicle = vehicles
			user.Recharge = recharges
			c.JSON(200, user)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		}
	})

	//Update user
	r.PUT("/users", func(c *gin.Context) {
		var input input2.UpdateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, db)

		if resp == "user" || resp == "admin" {
			input.Person.Password = utils.CreateHashPassword(input.Person.Password)
			user := crud.GetUserByEmail(input.LoginInput.Email, db)
			user.Person = input.Person
			user.Document = input.Document
			crud.UpdateUser(user, db)
			c.JSON(http.StatusOK, gin.H{"Response": "Usuário alterado"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		}
	})

	//Update admin
	r.PUT("/admins", func(c *gin.Context) {
		var input input2.UpdateAdminInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, db)

		if resp == "admin" {
			input.Person.Password = utils.CreateHashPassword(input.Person.Password)
			admin := crud.GetAdminByEmail(input.LoginInput.Email, db)
			admin.Person = input.Person
			crud.UpdateAdmin(admin, db)
			c.JSON(http.StatusOK, gin.H{"Response": "Admin alterado"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		}

	})

	//Update traffic warden
	r.PUT("/trafficwarden", func(c *gin.Context) {
		var input input2.UpdateTrafficWarden
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, db)
		if resp == "trafficWarden" || resp == "admin" {
			input.Person.Password = utils.CreateHashPassword(input.Person.Password)
			trafficWarden := crud.GetTrafficWardenByEmail(input.LoginInput.Email, db)
			trafficWarden.Person = input.Person
			crud.UpdateTrafficWarden(trafficWarden, db)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		}

	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
