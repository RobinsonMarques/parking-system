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
	var Token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX25hbWUiOiJyb2JpbmhvbWFycXVlcy5ybTJAZ21haWwuY29tIiwic2NvcGUiOlsiYWxsIl0sImV4cCI6MTYyMjMxNjQ3NSwianRpIjoiVk9rckZUcmxOMTFzZ1Ffazh6TDRWTE9WRWVZIiwiY2xpZW50X2lkIjoiUzNDeUtoT09nQTZMeWx0cSJ9.VPPewOzPaxpHazBxOyA-58zXMI0xE_9R5-zTvBET2kZkbNenONiFz336pPJ0rxv8SbYBItFFW7o9-YMolPknn71gqTAtx0BuPL-la6K8sbK3GtGuavN2P4JK3LtukD0mv0Ehu-HZGC2wgVZIz0kqEXANMS4lfm202GpPp87-jDbhQdnOyVcyEGa3IQ7KWDUFB2TWEH815iToIgHR-1aDbDm9p0ItdNhR65BqomHYv_a7XyW4p40AGtgJ9c67tLhPOTaOMQXlxFOIkdldVYEI0LvmkGiafbJdMqtO0Zx-FNN4HK0p0nwQhFVeEEBt6x0BCghuj1L1JRtAdMR6J2YS3w"
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
	user, err := crud.GetUserByEmail(input.LoginInput.Email, r.db)
	if err != nil {
		return err
	}
	if resp == "user" {
		date := time.Now()
		var chargeString = fmt.Sprintf(`{
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

		respo := crud.CreateRecharge(finalRecharge, r.db)
		if respo.Data.Error != nil {
			return respo.Data.Error
		}

		rechargeReturn, err := crud.GetRechargeByUserId(user.ID, r.db)
		if err != nil {
			return err
		}
		leng := len(rechargeReturn)
		billet := database.Billet{
			BilletLink: response.Embedded.Charges[0].Link,
			RechargeID: rechargeReturn[leng-1].ID,
		}
		crud.CreateBillet(billet, r.db)
		return nil
	} else {
		err := errors.New(resp)
		return err
	}
}

func (r RechargeService) GetRechargeStatus(input input2.LoginInput) error {
	resp := utils.Login(input.Email, input.Password, r.db)

	if resp == "user" {
		user, err := crud.GetUserByEmail(input.Email, r.db)
		if err != nil {
			return err
		}
		unpaidRecharges, err := crud.GetUserUnpaidRechargesByID(user.ID, r.db)
		if err != nil {
			return err
		}

		for _, unpaidRecharge := range unpaidRecharges {
			status, err := utils.GetBilletStatus(unpaidRecharge.RechargeHash, r.Bearer)
			if err != nil {
				return err
			}

			if status == "CANCELLED" || status == "MANUAL_RECONCILIATION" || status == "FAILED" {
				err := crud.DeleteRechargeByID(unpaidRecharge.ID, r.db)
				if err != nil {
					return err
				}
			}

			if status == "PAID" {
				err := crud.UpdateBalance(user.Person.Email, float64(unpaidRecharge.Value), r.db)
				if err != nil {
					return err
				}
				err = crud.UpdateIsPaid(unpaidRecharge.ID, r.db)
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
	if resp == "admin" {
		err := crud.DeleteRechargeByID(rechargeID, r.db)
		if err != nil {
			return err
		}
		return nil
	} else {
		err := errors.New(resp)
		return err
	}
}
