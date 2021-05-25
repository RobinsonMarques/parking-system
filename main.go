package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/database"
	"github.com/RobinsonMarques/parking-system/dependencies"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/utils"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var Token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX25hbWUiOiJyb2JpbmhvbWFycXVlcy5ybTJAZ21haWwuY29tIiwic2NvcGUiOlsiYWxsIl0sImV4cCI6MTYyMjA0NzE5OCwianRpIjoiakpnUVlRdlVQLVNuUlhqWE5GMHpPeXVNTUI4IiwiY2xpZW50X2lkIjoiUzNDeUtoT09nQTZMeWx0cSJ9.VeYRDvLvOdKBf9aTzF_cvuHPxySZv48PM3s3UWKFKr4eWxA2VqQYN7rRmgw3OzW8rs2m1pwBAW4BN6EPHAvDN2HfFNvcS2ch0SMw3D0l--fZIt399tUns2aoPWxpmGOfg1io63PCbz97WmNEHn0mnlkQXrz2zkPiFuGrPowbUUXBMWjQWm189dPDcc5V7cAL43xsws2qjIpyjM4EmXNWSY5iOPBQbvsdrTc6AQ6ozFlqj3rowvjxnM1YmTfQVlezVvRxKkBWjfmnuG9NRMCFqHualzXyjZVF3yk42ufjSP0v9e011n-P92X1tYN97Up_WZ_ukO0dorqI6sJYcXRMKw"
var Bearer = "Bearer" + Token

func main() {
	db := dependencies.CreateConnection()
	db.AutoMigrate(&database.Person{})
	db.AutoMigrate(&database.User{})
	db.AutoMigrate(&database.TrafficWarden{})
	db.AutoMigrate(&database.Admin{})
	db.AutoMigrate(&database.ParkingTicket{})
	db.AutoMigrate(&database.Vehicle{})
	db.AutoMigrate(&database.Billet{})
	db.AutoMigrate(&database.Recharge{})

	Token = utils.CreateAccessToken(Bearer, Token)
	Bearer = "Bearer" + Token
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
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
							go utils.AlterVehicleStatus(resp2, input.ParkingTime, db)
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
		url := "https://sandbox.boletobancario.com/api-integration/charges"
		var input input2.CreateRecharge
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, db)
		user := crud.GetUserByEmail(input.LoginInput.Email, db)
		if resp == "user" {
			date := time.Now()
			var chargeString string = fmt.Sprintf(`{
"charge": {
            "description": "Recarga de crédito",
            "amount": %d,
            "paymentTypes": ["BOLETO"]
        },
        "billing": {
            "name": "%s",
            "document": "%s",
            "email": "%s",
            "notify": true
        }
}`, input.Value, user.Person.Name, user.Document, user.Person.Email)
			var jsonRequest = []byte(chargeString)
			//jsonRecharge, _ := json.Marshal(jsonStr)

			req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonRequest))
			req.Header.Add("X-Api-Version", "2")
			req.Header.Add("Authorization", Bearer)
			req.Header.Add("X-Resource-Token", "1AD89A918E8A9AD595BDD578188A496D6FC9A7743D79F9658CF4BC4C8E18FBCC")
			req.Header.Add("Content-Type", "application/json")

			client := &http.Client{}

			res, err := client.Do(req)

			if err != nil {
				log.Println("Error", err)
			}
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)

			if err != nil {
				log.Println("Error reading the response:", err)
			}

			response := input2.Response{}
			json.Unmarshal(body, &response)

			finalRecharge := database.Recharge{
				Date:         date.String(),
				Value:        input.Value,
				IsPaid:       false,
				PaymentType:  input.PaymentType,
				UserID:       user.ID,
				RechargeHash: response.Embedded.Charges[0].ID,
			}

			crud.CreateRecharge(finalRecharge, db)
			rechargeReturn := crud.GetRechargeByUserId(user.ID, db)
			len := len(rechargeReturn)
			billet := database.Billet{
				BilletLink: response.Embedded.Charges[0].Link,
				RechargeID: rechargeReturn[len-1].ID,
			}
			crud.CreateBillet(billet, db)
			c.JSON(http.StatusOK, gin.H{"Response": "Recarga criada"})
			//log.Println("Recharge Hash:", response.Embedded.Charges.ID)
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
			vehicle = crud.GetVehicleByLicensePlate(licensePlate, db)
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
