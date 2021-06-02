package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/RobinsonMarques/parking-system/database"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/interfaces"
	"github.com/RobinsonMarques/parking-system/utils"
	"io/ioutil"
	"net/http"
	"time"
)

func CreateBearer() (string, error) {
	var Token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX25hbWUiOiJyb2JpbmhvbWFycXVlcy5ybTJAZ21haWwuY29tIiwic2NvcGUiOlsiYWxsIl0sImV4cCI6MTYyMjY2NjE2NywianRpIjoibEw0OUhkWHVPU3dqMWdFajZ3WVo0RUY1R004IiwiY2xpZW50X2lkIjoiUzNDeUtoT09nQTZMeWx0cSJ9.OHPlcsWNuUBPgzhpsemAkioX6wkHt4_vxxGuBA7SO0-TMjzqG_g9zN-Y72kUO8Qd47Zc8oX1bG9DgAci8XQ5jMqS1HY4MTZcz85ZQ5TlsV4K2QZlmjAl8d1DNmGmS9gTwakMuUW8EoyAWebZ8OM5l_VmNXvvDHEixo9OTJTJj9ZsO2WmlHOPxO0u5-TiKm9YHjtExAtkXTcF24jwzYW1J1sxHNt7r7pp7FSdfWhJWe2bou5KzfYlki_I5BObqc_-n6oL6RtxgfZ-jaZmH3F3lDjPevc_QPKEnH-NunkCMma8u1Z5JP9tv4DZPr075EgsnM-Xi0D-gtg1xKqLenGNYA"
	var Bearer = "Bearer" + Token
	Token, err := utils.CreateAccessToken(Bearer, Token)
	if err != nil {
		return "", err
	}
	Bearer = "Bearer" + Token
	return Bearer, nil
}

func NewRechargeService(rechargeInterface interfaces.RechargeInterface, userInterface interfaces.UserInterface, utilInterface interfaces.UtilInterface, billetInterface interfaces.BilletInterface) (RechargeService, error) {
	Bearer, err := CreateBearer()
	if err != nil {
		return RechargeService{}, err
	}
	return RechargeService{
		rechargeInterface: rechargeInterface,
		userInterface:     userInterface,
		utilInterface:     utilInterface,
		billetInterface:   billetInterface,
		Bearer:            Bearer}, nil
}

type RechargeService struct {
	rechargeInterface interfaces.RechargeInterface
	userInterface     interfaces.UserInterface
	utilInterface     interfaces.UtilInterface
	billetInterface   interfaces.BilletInterface
	Bearer            string
}

func (r RechargeService) CreateRecharge(input input2.CreateRecharge, url string) error {
	resp := r.utilInterface.Login(input.LoginInput.Email, input.LoginInput.Password)
	user, err := r.userInterface.GetUserByEmail(input.LoginInput.Email)
	if err != nil {
		return err
	}
	if resp == "user" {
		date := time.Now()
		var chargeString = fmt.Sprintf(`{
"charge": {
            "description": "Recarga de crédito",
            "amount": %.2f,
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

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonRequest))
		if err != nil {
			return err
		}
		req.Header.Add("X-Api-Version", "2")
		req.Header.Add("Authorization", r.Bearer)
		req.Header.Add("X-Resource-Token", "1AD89A918E8A9AD595BDD578188A496D6FC9A7743D79F9658CF4BC4C8E18FBCC")
		req.Header.Add("Content-Type", "application/json")

		client := &http.Client{}

		res, err := client.Do(req)
		if res.StatusCode != 200 {
			err := errors.New(res.Status)
			return err
		}
		if err != nil {
			return err
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			return err
		}

		response := input2.Response{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			return err
		}

		finalRecharge := database.Recharge{
			Date:         date.String(),
			Value:        input.Value,
			IsPaid:       false,
			PaymentType:  input.PaymentType,
			UserID:       user.ID,
			RechargeHash: response.Embedded.Charges[0].ID,
		}

		err = r.rechargeInterface.CreateRecharge(finalRecharge)
		if err != nil {
			return err
		}

		rechargeReturn, err := r.rechargeInterface.GetRechargeByUserId(user.ID)
		if err != nil {
			return err
		}
		leng := len(rechargeReturn)
		billet := database.Billet{
			BilletLink: response.Embedded.Charges[0].Link,
			RechargeID: rechargeReturn[leng-1].ID,
		}
		err = r.billetInterface.CreateBillet(billet)
		if err != nil {
			return err
		}
		return nil
	} else {
		err := errors.New(resp)
		return err
	}
}

func (r RechargeService) GetRechargeStatus(input input2.LoginInput) error {
	resp := r.utilInterface.Login(input.Email, input.Password)
	if resp == "user" {
		user, err := r.userInterface.GetUserByEmail(input.Email)
		if err != nil {
			return err
		}
		unpaidRecharges, err := r.rechargeInterface.GetUserUnpaidRechargesByID(user.ID)
		if err != nil {
			return err
		}

		for _, unpaidRecharge := range unpaidRecharges {
			status, err := utils.GetBilletStatus(unpaidRecharge.RechargeHash, r.Bearer)
			if err != nil {
				return err
			}

			if status == "CANCELLED" || status == "MANUAL_RECONCILIATION" || status == "FAILED" {
				err := r.rechargeInterface.DeleteRechargeByID(unpaidRecharge.ID)
				if err != nil {
					return err
				}
				err = r.billetInterface.DeleteBilletByRechargeID(unpaidRecharge.ID)
				if err != nil {
					return err
				}
			}

			if status == "PAID" {
				err := r.userInterface.UpdateBalance(user.Person.Email, unpaidRecharge.Value)
				if err != nil {
					return err
				}
				err = r.rechargeInterface.UpdateIsPaid(unpaidRecharge.ID)
				if err != nil {
					return err
				}
				return nil
			}
		}
	}
	err := errors.New(resp)
	return err
}

func (r RechargeService) DeleteRechargeByID(input input2.LoginInput, rechargeID uint) error {
	resp := r.utilInterface.Login(input.Email, input.Password)
	if resp == "admin" {
		err := r.rechargeInterface.DeleteRechargeByID(rechargeID)
		if err != nil {
			return err
		}
		err = r.billetInterface.DeleteBilletByRechargeID(rechargeID)
		if err != nil {
			return err
		}
		return nil
	} else {
		err := errors.New(resp)
		return err
	}
}
