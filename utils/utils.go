package utils

import (
	"encoding/json"
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/database"
	"github.com/RobinsonMarques/parking-system/input"
	"go/types"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func CreateHashPassword(password string) string {
	bytePassword := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, 8)

	if err != nil {
		log.Fatal(err)
	}
	return string(hashedPassword)
}

func ComparePassword(password string, userEmail string, userType string, db *gorm.DB) error {
	var err error

	if userType == "user" {
		hashedPassword := []byte(crud.GetPassword(userEmail, userType, db))
		bytePassword := []byte(password)
		err = bcrypt.CompareHashAndPassword(hashedPassword, bytePassword)
		return err
	} else if userType == "admin" {
		hashedPassword := []byte(crud.GetPassword(userEmail, userType, db))
		bytePassword := []byte(password)
		err = bcrypt.CompareHashAndPassword(hashedPassword, bytePassword)
		return err
	} else if userType == "trafficWarden" {
		hashedPassword := []byte(crud.GetPassword(userEmail, userType, db))
		bytePassword := []byte(password)
		err = bcrypt.CompareHashAndPassword(hashedPassword, bytePassword)
		return err
	} else {
		err = types.Error{Msg: "Tipo de usuário inválido"}
		return err
	}

}

func Login(email string, password string, db *gorm.DB) string {
	response := ""
	if crud.GetUserByEmail(email, db).Person.Name != "" {
		err := ComparePassword(password, email, "user", db)
		if err == nil {
			response = "user"
		} else {
			response = "Senha inválida!"
		}
	} else if crud.GetAdminByEmail(email, db).Person.Name != "" {
		err := ComparePassword(password, email, "admin", db)
		if err == nil {
			response = "admin"
		} else {
			response = "Senha inválida"
		}
	} else if crud.GetTrafficWardenByEmail(email, db).Person.Name != "" {
		err := ComparePassword(password, email, "trafficWarden", db)
		if err == nil {
			response = "trafficWarden"
		} else {
			response = "Senha inválida"
		}
	} else {
		response = "Usuário não cadastrado"
	}

	return response
}

//func AlterVehicleStatus(vehicle database.Vehicle, db *gorm.DB) {
//ticket := crud.GetLastParkingTicketFromVehicle(vehicle.ID, db)
//layout := "2006-01-02T15:04:05+07:00"

//endTime, _ := time.Parse(layout, ticket[0].EndTime)
//timeNow := time.Now()
//if timeNow.After(endTime) {
//crud.UpdateIsActive(vehicle.ID, false, db)
//crud.UpdateIsParked(vehicle.ID, false, db)
//}
//}

func AlterVehicleStatus(vehicle database.Vehicle, parkingTime int, db *gorm.DB) {
	duration := time.Duration(parkingTime)
	for true {
		time.Sleep(duration * time.Hour)
		crud.UpdateIsActive(vehicle.ID, false, db)
		crud.UpdateIsParked(vehicle.ID, false, db)
	}
}

func GetBilletStatus(rechargeID, Bearer string) string {
	url := "https://sandbox.boletobancario.com/api-integration/charges/" + rechargeID

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-Api-Version", "2")
	req.Header.Add("Authorization", Bearer)
	req.Header.Add("X-Resource-Token", "1AD89A918E8A9AD595BDD578188A496D6FC9A7743D79F9658CF4BC4C8E18FBCC")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Println("Error", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error reading the response:", err)
	}

	billet := input.Payment{}
	json.Unmarshal(body, &billet)

	//fmt.Println("Response status:", resp.Status)
	//fmt.Println("Response Headers:", resp.Header)
	//fmt.Println("Response Body:", string(body))
	return billet.Status
}

func CreateAccessToken(bearer, Token string) string {
	endpoint := "https://sandbox.boletobancario.com/authorization-server/oauth/token"
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	if validateToken(bearer) {
		return Token
	} else {

		req, _ := http.NewRequest("Post", endpoint, strings.NewReader(data.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
		req.Header.Add("Authorization", "Basic UzNDeUtoT09nQTZMeWx0cTouKjFEekY+QlM4UFR6em80MXRqTE9jfSRGaStmQWdIZA==")

		client := &http.Client{}

		resp, err := client.Do(req)

		if err != nil {
			log.Println("Error", err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Println("Error reading the response:", err)
		}

		token := input.AccessToken{}
		json.Unmarshal(body, &token)

		//fmt.Println("Response status:", resp.Status)
		//fmt.Println("Response Headers:", resp.Header)
		//fmt.Println("Response Body:", string(body))

		return token.AccessToken
	}
}

func validateToken(bearer string) bool {
	url := "https://sandbox.boletobancario.com/api-integration/digital-accounts"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-Api-Version", "2")
	req.Header.Add("Authorization", bearer)
	req.Header.Add("X-Resource-Token", "1AD89A918E8A9AD595BDD578188A496D6FC9A7743D79F9658CF4BC4C8E18FBCC")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Println("Error", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error reading the response:", err)
	}

	token := input.ValidateAccessToken{}
	json.Unmarshal(body, &token)

	//fmt.Println("Response status:", resp.Status)
	//fmt.Println("Response Headers:", resp.Header)
	//fmt.Println("Response Body:", string(body))

	if resp.Status != "200 OK" {
		return false
	} else {
		return true
	}
}
