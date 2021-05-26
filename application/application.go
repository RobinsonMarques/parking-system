package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/database"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func NewApplication(db *gorm.DB) Application {
	var Token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX25hbWUiOiJyb2JpbmhvbWFycXVlcy5ybTJAZ21haWwuY29tIiwic2NvcGUiOlsiYWxsIl0sImV4cCI6MTYyMjA0NzE5OCwianRpIjoiakpnUVlRdlVQLVNuUlhqWE5GMHpPeXVNTUI4IiwiY2xpZW50X2lkIjoiUzNDeUtoT09nQTZMeWx0cSJ9.VeYRDvLvOdKBf9aTzF_cvuHPxySZv48PM3s3UWKFKr4eWxA2VqQYN7rRmgw3OzW8rs2m1pwBAW4BN6EPHAvDN2HfFNvcS2ch0SMw3D0l--fZIt399tUns2aoPWxpmGOfg1io63PCbz97WmNEHn0mnlkQXrz2zkPiFuGrPowbUUXBMWjQWm189dPDcc5V7cAL43xsws2qjIpyjM4EmXNWSY5iOPBQbvsdrTc6AQ6ozFlqj3rowvjxnM1YmTfQVlezVvRxKkBWjfmnuG9NRMCFqHualzXyjZVF3yk42ufjSP0v9e011n-P92X1tYN97Up_WZ_ukO0dorqI6sJYcXRMKw"
	var Bearer = "Bearer" + Token
	Token = utils.CreateAccessToken(Bearer, Token)
	Bearer = "Bearer" + Token
	return Application{db: db, Bearer: Bearer}
}

type Application struct {
	db     *gorm.DB
	Bearer string
}

func (a Application) GetAllVehicles(c *gin.Context) {
	//valida o input
	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := utils.Login(input.Email, input.Password, a.db)

	if resp == "trafficWarden" || resp == "admin" {
		vehicles := crud.GetAllVehicles(a.db)
		c.JSON(200, vehicles)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
	}

}

func (a Application) CreateUser(c *gin.Context) {
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
	crud.CreateUser(user, a.db)
	c.JSON(http.StatusOK, gin.H{"Response": "Usuário criado"})
}

func (a Application) CreateAdmin(c *gin.Context) {
	//Valida o input
	var input input2.CreateAdminInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, a.db)

	if resp == "admin" {
		input.Person.Password = utils.CreateHashPassword(input.Person.Password)

		//Cria o admin
		admin := database.Admin{
			Person: input.Person,
		}
		crud.CreateAdmin(admin, a.db)
		c.JSON(http.StatusOK, gin.H{"Response": "Admin criado"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
	}
}

func (a Application) CreateTrafficWarden(c *gin.Context) {
	//Valida o input
	var input input2.CreateTrafficWarden

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, a.db)

	if resp == "admin" {
		input.Person.Password = utils.CreateHashPassword(input.Person.Password)

		//Cria o admin
		warden := database.TrafficWarden{
			Person: input.Person,
		}
		crud.CreateTrafficWarden(warden, a.db)
		c.JSON(http.StatusOK, gin.H{"Response": "Guarda criado"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
	}
}

func (a Application) CreateVehicle(c *gin.Context) {
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

func (a Application) CreateParkingTicket(c *gin.Context) {
	var input input2.CreateParkingTicket
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := utils.Login(input.Login.Email, input.Login.Password, a.db)
	resp2 := crud.GetVehicleById(input.VehicleID, a.db)
	user := crud.GetUserByEmail(input.Login.Email, a.db)
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
						crud.CreateParkingTicket(ticket, a.db)
						crud.UpdateIsParked(input.VehicleID, true, a.db)
						crud.UpdateBalance(input.Login.Email, -price, a.db)
						crud.UpdateIsActive(input.VehicleID, true, a.db)
						go utils.AlterVehicleStatus(resp2, input.ParkingTime, a.db)
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
}

func (a Application) CreateRecharge(c *gin.Context) {
	url := "https://sandbox.boletobancario.com/api-integration/charges"
	var input input2.CreateRecharge
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, a.db)
	user := crud.GetUserByEmail(input.LoginInput.Email, a.db)
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
		req.Header.Add("Authorization", a.Bearer)
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

		crud.CreateRecharge(finalRecharge, a.db)
		rechargeReturn := crud.GetRechargeByUserId(user.ID, a.db)
		len := len(rechargeReturn)
		billet := database.Billet{
			BilletLink: response.Embedded.Charges[0].Link,
			RechargeID: rechargeReturn[len-1].ID,
		}
		crud.CreateBillet(billet, a.db)
		c.JSON(http.StatusOK, gin.H{"Response": "Recarga criada"})
		//log.Println("Recharge Hash:", response.Embedded.Charges.ID)
	}
}

func (a Application) GetVehicleByLicensePlate(c *gin.Context) {
	licensePlate := c.Param("licensePlate")

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := utils.Login(input.Email, input.Password, a.db)

	if resp == "trafficWarden" {
		vehicle := crud.GetVehicleByLicensePlate(licensePlate, a.db)
		vehicle = crud.GetVehicleByLicensePlate(licensePlate, a.db)
		ticket := crud.GetLastParkingTicketFromVehicle(vehicle.ID, a.db)
		vehicle.ParkingTicket = ticket
		c.JSON(200, vehicle)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
	}
}

func (a Application) GetRechargesStatus(c *gin.Context) {
	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := utils.Login(input.Email, input.Password, a.db)

	if resp == "user" {
		user := crud.GetUserByEmail(input.Email, a.db)
		unpaidRecharges := crud.GetUserUnpaidRechargesByID(user.ID, a.db)

		for _, unpaidRecharge := range unpaidRecharges {
			status := utils.GetBilletStatus(unpaidRecharge.RechargeHash, a.Bearer)

			if status == "CANCELLED" || status == "MANUAL_RECONCILIATION" || status == "FAILED" {
				crud.DeleteRechargeByID(unpaidRecharge.ID, a.db)
			}

			if status == "PAID" {
				crud.UpdateBalance(user.Person.Email, float64(unpaidRecharge.Value), a.db)
				crud.UpdateIsPaid(unpaidRecharge.ID, a.db)
				c.JSON(200, "Saldo alterado com sucesso")
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
	}
}

func (a Application) GetUserByDocument(c *gin.Context) {
	document := c.Param("document")

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := utils.Login(input.Email, input.Password, a.db)

	if resp == "admin" {
		user := crud.GetUserByDocument(document, a.db)
		vehicles := crud.GetVehiclesByUserId(user.ID, a.db)
		recharges := crud.GetRechargeByUserId(user.ID, a.db)

		for i := range recharges {
			billet := crud.GetBilletByRechargeId(recharges[i].ID, a.db)
			recharges[i].Billet = billet
		}
		user.Vehicle = vehicles
		user.Recharge = recharges
		c.JSON(200, user)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
	}
}

func (a Application) UpdateUser(c *gin.Context) {
	userIDString := c.Param("userID")
	userIDInt, _ := strconv.Atoi(userIDString)
	userID := uint(userIDInt)

	var input input2.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, a.db)

	if resp == "user" || resp == "admin" {
		user := crud.GetUserByID(userID, a.db)
		if resp == "user" && user.Person.Email != input.LoginInput.Email {
			c.JSON(http.StatusBadRequest, gin.H{"Response": "Usuário não possui permissão"})
		} else {
			input.Person.Password = utils.CreateHashPassword(input.Person.Password)
			user.Person = input.Person
			user.Document = input.Document
			crud.UpdateUser(user, a.db)
			c.JSON(http.StatusOK, gin.H{"Response": "Usuário alterado"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
	}
}

func (a Application) UpdateAdmin(c *gin.Context) {
	adminIDString := c.Param("adminID")
	adminIDInt, _ := strconv.Atoi(adminIDString)
	adminID := uint(adminIDInt)

	var input input2.UpdateAdminInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, a.db)

	if resp == "admin" {
		input.Person.Password = utils.CreateHashPassword(input.Person.Password)
		admin := crud.GetAdminByID(adminID, a.db)
		admin.Person = input.Person
		crud.UpdateAdmin(admin, a.db)
		c.JSON(http.StatusOK, gin.H{"Response": "Admin alterado"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
	}
}

func (a Application) UpdateTrafficWarden(c *gin.Context) {
	wardenIDString := c.Param("wardenID")
	wardenIDInt, _ := strconv.Atoi(wardenIDString)
	wardenID := uint(wardenIDInt)

	var input input2.UpdateTrafficWarden
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, a.db)
	if resp == "trafficWarden" || resp == "admin" {
		trafficWarden := crud.GetTrafficWardenByID(wardenID, a.db)
		if resp == "trafficWarden" && trafficWarden.Person.Email != input.LoginInput.Email {
			c.JSON(http.StatusBadRequest, gin.H{"Response": "Usuário não possui permissão"})
		} else {

			input.Person.Password = utils.CreateHashPassword(input.Person.Password)

			trafficWarden.Person = input.Person
			crud.UpdateTrafficWarden(trafficWarden, a.db)
			c.JSON(http.StatusOK, gin.H{"Response": "Guarda de trânsito alterado"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
	}
}

func (a Application) UpdateVehicle(c *gin.Context) {
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

func (a Application) UpdateVehicleOwner(c *gin.Context) {
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

func (a Application) DeleteUserByID(c *gin.Context) {
	userIDString := c.Param("userID")
	userIDInt, _ := strconv.Atoi(userIDString)
	userID := uint(userIDInt)

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := utils.Login(input.Email, input.Password, a.db)

	if resp == "user" || resp == "admin" {
		user := crud.GetUserByEmail(input.Email, a.db)
		if resp == "user" && user.ID != userID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Usuário logado não possui permissão"})
		} else {
			crud.DeleteUserByID(userID, a.db)
			c.JSON(http.StatusOK, gin.H{"Response": "Usuário deletado"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
	}
}

func (a Application) DeleteTrafficWardenByID(c *gin.Context) {
	wardenIDString := c.Param("trafficwardenID")
	wardenIDInt, _ := strconv.Atoi(wardenIDString)
	wardenID := uint(wardenIDInt)

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := utils.Login(input.Email, input.Password, a.db)

	if resp == "admin" {
		crud.DeleteTrafficWardenByID(wardenID, a.db)
		c.JSON(http.StatusOK, gin.H{"Response": "Guarda de Trânsito deletado"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
	}
}

func (a Application) DeleteAdminByID(c *gin.Context) {
	adminIDString := c.Param("adminID")
	adminIDInt, _ := strconv.Atoi(adminIDString)
	adminID := uint(adminIDInt)

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := utils.Login(input.Email, input.Password, a.db)

	if resp == "admin" {
		crud.DeleteAdminByID(adminID, a.db)
		c.JSON(http.StatusOK, gin.H{"Response": "Admin deletado"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
	}
}

func (a Application) DeleteParkingTicketByID(c *gin.Context) {
	ticketIDString := c.Param("ticketID")
	ticketIDInt, _ := strconv.Atoi(ticketIDString)
	ticketID := uint(ticketIDInt)

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := utils.Login(input.Email, input.Password, a.db)

	if resp == "admin" {
		crud.DeleteParkingTicketByID(ticketID, a.db)
		c.JSON(http.StatusOK, gin.H{"Response": "Ticket deletado"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
	}
}

func (a Application) DeleteRechargeByID(c *gin.Context) {
	rechargeIDString := c.Param("rechargeID")
	rechargeIDInt, _ := strconv.Atoi(rechargeIDString)
	rechargeID := uint(rechargeIDInt)

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := utils.Login(input.Email, input.Password, a.db)

	if resp == "admin" {
		crud.DeleteRechargeByID(rechargeID, a.db)
		c.JSON(http.StatusOK, gin.H{"Response": "Recarga deletada"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
	}
}

func (a Application) DeleteBilletByID(c *gin.Context) {
	billetIDString := c.Param("billetID")
	billetIDInt, _ := strconv.Atoi(billetIDString)
	billetID := uint(billetIDInt)

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := utils.Login(input.Email, input.Password, a.db)

	if resp == "admin" {
		crud.DeleteBilletByID(billetID, a.db)
		c.JSON(http.StatusOK, gin.H{"Response": "Boleto deletado"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
	}
}

func (a Application) DeleteVehicleByID(c *gin.Context) {
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

func (a Application) Run() error {

	r := gin.Default()

	r.GET("/vehicles", a.GetAllVehicles)
	r.POST("/users", a.CreateUser)
	r.POST("/admins", a.CreateAdmin)
	r.POST("/trafficwarden", a.CreateTrafficWarden)
	r.POST("/vehicles", a.CreateVehicle)
	r.POST("/parkingtickets", a.CreateParkingTicket)
	r.POST("/recharge", a.CreateRecharge)
	r.GET("/vehicles/:licensePlate", a.GetVehicleByLicensePlate)
	r.GET("/recharges", a.GetRechargesStatus)
	r.GET("users/:document", a.GetUserByDocument)
	r.PUT("/users/:userID", a.UpdateUser)
	r.PUT("/admins/:adminID", a.UpdateAdmin)
	r.PUT("/trafficwarden/:wardenID", a.UpdateTrafficWarden)
	r.PUT("/vehicles", a.UpdateVehicle)
	r.PUT("/vehicles/updateowner/:vehicleID", a.UpdateVehicleOwner)
	r.DELETE("/users/:userID", a.DeleteUserByID)
	r.DELETE("/trafficwarden/:trafficwardenID", a.DeleteTrafficWardenByID)
	r.DELETE("/admins/:adminID", a.DeleteAdminByID)
	r.DELETE("/parkingtickets/:ticketID", a.DeleteParkingTicketByID)
	r.DELETE("/recharge/:rechargeID", a.DeleteRechargeByID)
	r.DELETE("/billets/:billetID", a.DeleteBilletByID)
	r.DELETE("/vehicles/:vehicleID", a.DeleteVehicleByID)

	return r.Run() // listen and serve on 0.0.0.0:8080
}
