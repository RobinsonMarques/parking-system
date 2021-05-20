package main

import (
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/database"
	"github.com/RobinsonMarques/parking-system/dependencies"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
			if resp2.UserID == user.ID {
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
				c.JSON(http.StatusBadRequest, gin.H{"Response": "Usuário não possui o veículo"})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"Response": resp})
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
	r.PUT("/users/:userID", func(c *gin.Context) {
		userIDString := c.Param("userID")
		userIDInt, _ := strconv.Atoi(userIDString)
		userID := uint(userIDInt)

		var input input2.UpdateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, db)

		if resp == "user" || resp == "admin" {
			user := crud.GetUserByID(userID, db)
			if resp == "user" && user.Person.Email != input.LoginInput.Email {
				c.JSON(http.StatusBadRequest, gin.H{"Response": "Usuário não possui permissão"})
			} else {
				input.Person.Password = utils.CreateHashPassword(input.Person.Password)
				user.Person = input.Person
				user.Document = input.Document
				crud.UpdateUser(user, db)
				c.JSON(http.StatusOK, gin.H{"Response": "Usuário alterado"})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		}
	})

	//Update admin
	r.PUT("/admins/:adminID", func(c *gin.Context) {
		adminIDString := c.Param("adminID")
		adminIDInt, _ := strconv.Atoi(adminIDString)
		adminID := uint(adminIDInt)

		var input input2.UpdateAdminInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, db)

		if resp == "admin" {
			input.Person.Password = utils.CreateHashPassword(input.Person.Password)
			admin := crud.GetAdminByID(adminID, db)
			admin.Person = input.Person
			crud.UpdateAdmin(admin, db)
			c.JSON(http.StatusOK, gin.H{"Response": "Admin alterado"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		}

	})

	//Update traffic warden
	r.PUT("/trafficwarden/:wardenID", func(c *gin.Context) {
		wardenIDString := c.Param("wardenID")
		wardenIDInt, _ := strconv.Atoi(wardenIDString)
		wardenID := uint(wardenIDInt)

		var input input2.UpdateTrafficWarden
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, db)
		if resp == "trafficWarden" || resp == "admin" {
			trafficWarden := crud.GetTrafficWardenByID(wardenID, db)
			if resp == "trafficWarden" && trafficWarden.Person.Email != input.LoginInput.Email {
				c.JSON(http.StatusBadRequest, gin.H{"Response": "Usuário não possui permissão"})
			} else {

				input.Person.Password = utils.CreateHashPassword(input.Person.Password)

				trafficWarden.Person = input.Person
				crud.UpdateTrafficWarden(trafficWarden, db)
				c.JSON(http.StatusOK, gin.H{"Response": "Guarda de trânsito alterado"})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		}

	})

	//Update vehicle
	r.PUT("/vehicles", func(c *gin.Context) {
		var input input2.UpdateVehicle
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, db)
		user := crud.GetUserByEmail(input.LoginInput.Email, db)
		vehicles := crud.GetVehiclesByUserId(user.ID, db)
		var resp2 bool
		for i := range vehicles {
			if vehicles[i].UserID == user.ID {
				resp2 = true
			}
		}

		if resp == "user" {
			if resp2 {
				vehicle := crud.GetVehicleByLicensePlate(input.LicensePlate, db)
				vehicle.VehicleModel = input.VehicleModel
				vehicle.VehicleType = input.VehicleType
				vehicle.LicensePlate = input.NewLicensePlate
				crud.UpdateVehicle(vehicle, db)
				c.JSON(http.StatusOK, gin.H{"Response": "Veículo alterado"})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Usuário não possui este veículo"})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		}

	})

	//Update vehicle owner
	r.PUT("/vehicles/updateowner/:vehicleID", func(c *gin.Context) {
		vehicleIDString := c.Param("vehicleID")
		vehicleIDInt, _ := strconv.Atoi(vehicleIDString)
		vehicleID := uint(vehicleIDInt)
		var input input2.UpdateVehicleOwner
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, db)
		user := crud.GetUserByID(input.NewUserID, db)
		if resp == "admin" {
			if user.Person.Name != "" {
				crud.UpdateVehicleOwner(vehicleID, input.NewUserID, db)
				c.JSON(http.StatusOK, gin.H{"Response": "Dono do veículo alterado"})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Usuário inexistente"})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		}
	})

	//Delete user by id
	r.DELETE("/users/:userID", func(c *gin.Context) {
		userIDString := c.Param("userID")
		userIDInt, _ := strconv.Atoi(userIDString)
		userID := uint(userIDInt)

		var input input2.LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp := utils.Login(input.Email, input.Password, db)

		if resp == "user" || resp == "admin" {
			user := crud.GetUserByEmail(input.Email, db)
			if resp == "user" && user.ID != userID {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Usuário logado não possui permissão"})
			} else {
				crud.DeleteUserByID(userID, db)
				c.JSON(http.StatusOK, gin.H{"Response": "Usuário deletado"})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		}
	})

	//Path delete traffic warden by id
	r.DELETE("/trafficwarden/:trafficwardenID", func(c *gin.Context) {
		wardenIDString := c.Param("trafficwardenID")
		wardenIDInt, _ := strconv.Atoi(wardenIDString)
		wardenID := uint(wardenIDInt)

		var input input2.LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp := utils.Login(input.Email, input.Password, db)

		if resp == "admin" {
			crud.DeleteTrafficWardenByID(wardenID, db)
			c.JSON(http.StatusOK, gin.H{"Response": "Guarda de Trânsito deletado"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		}
	})

	//Path delete admin by id
	r.DELETE("/admins/:adminID", func(c *gin.Context) {
		adminIDString := c.Param("adminID")
		adminIDInt, _ := strconv.Atoi(adminIDString)
		adminID := uint(adminIDInt)

		var input input2.LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp := utils.Login(input.Email, input.Password, db)

		if resp == "admin" {
			crud.DeleteAdminByID(adminID, db)
			c.JSON(http.StatusOK, gin.H{"Response": "Admin deletado"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		}
	})

	//Path delete parking ticket by id
	r.DELETE("/parkingtickets/:ticketID", func(c *gin.Context) {
		ticketIDString := c.Param("ticketID")
		ticketIDInt, _ := strconv.Atoi(ticketIDString)
		ticketID := uint(ticketIDInt)

		var input input2.LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp := utils.Login(input.Email, input.Password, db)

		if resp == "admin" {
			crud.DeleteParkingTicketByID(ticketID, db)
			c.JSON(http.StatusOK, gin.H{"Response": "Ticket deletado"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		}
	})

	//Path delete recharge by id
	r.DELETE("/recharge/:rechargeID", func(c *gin.Context) {
		rechargeIDString := c.Param("rechargeID")
		rechargeIDInt, _ := strconv.Atoi(rechargeIDString)
		rechargeID := uint(rechargeIDInt)

		var input input2.LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp := utils.Login(input.Email, input.Password, db)

		if resp == "admin" {
			crud.DeleteRechargeByID(rechargeID, db)
			c.JSON(http.StatusOK, gin.H{"Response": "Recarga deletada"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		}
	})

	//Path delete billet by id
	r.DELETE("/billets/:billetID", func(c *gin.Context) {
		billetIDString := c.Param("billetID")
		billetIDInt, _ := strconv.Atoi(billetIDString)
		billetID := uint(billetIDInt)

		var input input2.LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp := utils.Login(input.Email, input.Password, db)

		if resp == "admin" {
			crud.DeleteBilletByID(billetID, db)
			c.JSON(http.StatusOK, gin.H{"Response": "Boleto deletado"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		}
	})

	//Path delete vehicle by id
	r.DELETE("/vehicles/:vehicleID", func(c *gin.Context) {
		vehicleIDString := c.Param("vehicleID")
		vehicleIDInt, _ := strconv.Atoi(vehicleIDString)
		vehicleID := uint(vehicleIDInt)

		var input input2.LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp := utils.Login(input.Email, input.Password, db)

		if resp == "user" || resp == "admin" {
			vehicle := crud.GetVehicleById(vehicleID, db)
			user := crud.GetUserByEmail(input.Email, db)
			if resp == "user" && vehicle.UserID != user.ID {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Usuário logado não possui permissão"})
			} else {
				crud.DeleteVehicleByID(vehicleID, db)
				c.JSON(http.StatusOK, gin.H{"Response": "Veículo deletado"})
			}
		}
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
