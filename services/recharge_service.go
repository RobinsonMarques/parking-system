package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/database"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/utils"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"time"
)

func CreateBearer() (string, error) {
	var Token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX25hbWUiOiJyb2JpbmhvbWFycXVlcy5ybTJAZ21haWwuY29tIiwic2NvcGUiOlsiYWxsIl0sImV4cCI6MTYyMjU2NzUyOCwianRpIjoiTExwM3Fkd1A2eFlQRmU5UjFPandUeDQ5YVc0IiwiY2xpZW50X2lkIjoiUzNDeUtoT09nQTZMeWx0cSJ9.dSXDLWgnicqig_nIoNCqrB_WKacLD89AuWLtx0bfVj2TrQDZ5GNNwmsnxF5koGKSCchcO05N_D8kOISE2-2006V2AgADDgGGkiEweNP7gSKVHKZ8n_0_oFjY7-D1J8L9OxZma4OUciSwc4ZsL0WS4YR_VA_OBx5H23re423IYN0fe7Ons-_a8yJSfzJPJmwV1n8MgH_0B0DoyCefURI8YR0UbuTAdAiuoUw5uSmn2Plt8nx_U10bj1ZcjK_pFGsf7xmXX5FznIghxabYMlI8uMDJ7VlIxKMhVjtsb67IU_kXNObLJsU2yeRnoBRMn04r-mcS86iiyda7J4COPJg5bw"
	var Bearer = "Bearer" + Token
	Token, err := utils.CreateAccessToken(Bearer, Token)
	if err != nil {
		return "", err
	}
	Bearer = "Bearer" + Token
	return Bearer, nil
}

func NewRechargeService(db *gorm.DB) (RechargeService, error) {
	Bearer, err := CreateBearer()
	if err != nil {
		return RechargeService{}, err
	}
	return RechargeService{db: db, Bearer: Bearer}, nil
}

type RechargeService struct {
	db     *gorm.DB
	Bearer string
}

func (r RechargeService) CreateRecharge(input input2.CreateRecharge, url string) error {
	resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, r.db)
	rechargeCrud := crud.NewRechargeCrud(r.db)
	crud := crud.NewCrud(r.db)
	user, err := crud.UserCrud.GetUserByEmail(input.LoginInput.Email)
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

		err = rechargeCrud.CreateRecharge(finalRecharge)
		if err != nil {
			return err
		}

		rechargeReturn, err := rechargeCrud.GetRechargeByUserId(user.ID)
		if err != nil {
			return err
		}
		leng := len(rechargeReturn)
		billet := database.Billet{
			BilletLink: response.Embedded.Charges[0].Link,
			RechargeID: rechargeReturn[leng-1].ID,
		}
		err = crud.BilletCrud.CreateBillet(billet)
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
	resp := utils.Login(input.Email, input.Password, r.db)
	rechargeCrud := crud.NewRechargeCrud(r.db)
	crud := crud.NewCrud(r.db)
	if resp == "user" {
		user, err := crud.UserCrud.GetUserByEmail(input.Email)
		if err != nil {
			return err
		}
		unpaidRecharges, err := rechargeCrud.GetUserUnpaidRechargesByID(user.ID)
		if err != nil {
			return err
		}

		for _, unpaidRecharge := range unpaidRecharges {
			status, err := utils.GetBilletStatus(unpaidRecharge.RechargeHash, r.Bearer)
			if err != nil {
				return err
			}

			if status == "CANCELLED" || status == "MANUAL_RECONCILIATION" || status == "FAILED" {
				err := rechargeCrud.DeleteRechargeByID(unpaidRecharge.ID, crud)
				if err != nil {
					return err
				}
			}

			if status == "PAID" {
				err := crud.UserCrud.UpdateBalance(user.Person.Email, unpaidRecharge.Value, crud)
				if err != nil {
					return err
				}
				err = rechargeCrud.UpdateIsPaid(unpaidRecharge.ID)
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
	resp := utils.Login(input.Email, input.Password, r.db)
	rechargeCrud := crud.NewRechargeCrud(r.db)
	crud := crud.NewCrud(r.db)
	if resp == "admin" {
		err := rechargeCrud.DeleteRechargeByID(rechargeID, crud)
		if err != nil {
			return err
		}
		return nil
	} else {
		err := errors.New(resp)
		return err
	}
}
