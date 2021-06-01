package utils

import (
	"encoding/json"
	"errors"
	"github.com/RobinsonMarques/parking-system/input"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func CreateHashPassword(password string) (string, error) {
	bytePassword := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, 8)

	if err != nil {
		return "", err
	}
	return string(hashedPassword), err
}

//func ComparePassword(password string, userEmail string, userType string, db *gorm.DB) error {
//	var err error
//	crud := crud.NewCrud(db)
//	if userType == "user" {
//		userPassword, err := crud.UtilCrud.GetPassword(userEmail, userType, crud)
//		hashedPassword := []byte(userPassword)
//		bytePassword := []byte(password)
//		err = bcrypt.CompareHashAndPassword(hashedPassword, bytePassword)
//		return err
//	} else if userType == "admin" {
//		adminPassword, err := crud.UtilCrud.GetPassword(userEmail, userType, crud)
//		if err != nil {
//			return err
//		}
//		hashedPassword := []byte(adminPassword)
//		bytePassword := []byte(password)
//		err = bcrypt.CompareHashAndPassword(hashedPassword, bytePassword)
//		return err
//	} else if userType == "trafficWarden" {
//		wardenPassword, err := crud.UtilCrud.GetPassword(userEmail, userType, crud)
//		if err != nil {
//			return err
//		}
//		hashedPassword := []byte(wardenPassword)
//		bytePassword := []byte(password)
//		err = bcrypt.CompareHashAndPassword(hashedPassword, bytePassword)
//		return err
//	} else {
//		err = types.Error{Msg: "Tipo de usuário inválido"}
//		return err
//	}
//}
//
//func Login(email string, password string, db *gorm.DB) string {
//	response := ""
//	crud := crud.NewCrud(db)
//	user, _ := crud.UserCrud.GetUserByEmail(email)
//	admin, _ := crud.AdminCrud.GetAdminByEmail(email)
//	warden, _ := crud.TrafficWardenCrud.GetTrafficWardenByEmail(email)
//	if user.Person.Name != "" {
//		err := ComparePassword(password, email, "user", db)
//		if err == nil {
//			response = "user"
//		} else {
//			response = "Senha inválida!"
//		}
//	} else if admin.Person.Name != "" {
//		err := ComparePassword(password, email, "admin", db)
//		if err == nil {
//			response = "admin"
//		} else {
//			response = "Senha inválida"
//		}
//	} else if warden.Person.Name != "" {
//		err := ComparePassword(password, email, "trafficWarden", db)
//		if err == nil {
//			response = "trafficWarden"
//		} else {
//			response = "Senha inválida"
//		}
//	} else {
//		response = "Usuário não cadastrado"
//	}
//
//	return response
//}

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

func GetBilletStatus(rechargeID, Bearer string) (string, error) {
	url := "https://sandbox.boletobancario.com/api-integration/charges/" + rechargeID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("X-Api-Version", "2")
	req.Header.Add("Authorization", Bearer)
	req.Header.Add("X-Resource-Token", "1AD89A918E8A9AD595BDD578188A496D6FC9A7743D79F9658CF4BC4C8E18FBCC")

	client := &http.Client{}

	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		err := errors.New(resp.Status)
		return "", err
	}

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	billet := input.Payment{}
	err = json.Unmarshal(body, &billet)
	if err != nil {
		return "", err
	}

	return billet.Status, nil
}

func CreateAccessToken(bearer, Token string) (string, error) {
	endpoint := "https://sandbox.boletobancario.com/authorization-server/oauth/token"
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	val, err := validateToken(bearer)
	if err != nil {
		return "", err
	}
	if val {
		return Token, nil
	} else {

		req, err := http.NewRequest("Post", endpoint, strings.NewReader(data.Encode()))
		if err != nil {
			return "", err
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
		req.Header.Add("Authorization", "Basic UzNDeUtoT09nQTZMeWx0cTouKjFEekY+QlM4UFR6em80MXRqTE9jfSRGaStmQWdIZA==")

		client := &http.Client{}

		resp, err := client.Do(req)

		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return "", err
		}

		token := input.AccessToken{}
		err = json.Unmarshal(body, &token)
		if err != nil {
			return "", err
		}
		return token.AccessToken, nil
	}
}

func validateToken(bearer string) (bool, error) {
	url := "https://sandbox.boletobancario.com/api-integration/digital-accounts"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Add("X-Api-Version", "2")
	req.Header.Add("Authorization", bearer)
	req.Header.Add("X-Resource-Token", "1AD89A918E8A9AD595BDD578188A496D6FC9A7743D79F9658CF4BC4C8E18FBCC")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return false, nil
	} else {
		return true, nil
	}
}
