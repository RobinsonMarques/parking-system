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

func CreateBearer() string {
	var Token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX25hbWUiOiJyb2JpbmhvbWFycXVlcy5ybTJAZ21haWwuY29tIiwic2NvcGUiOlsiYWxsIl0sImV4cCI6MTYyMjIwNTQ0NywianRpIjoiWVpPbXVNbW96b2RaZTFxNXV2TjRqbnZ3c1kwIiwiY2xpZW50X2lkIjoiUzNDeUtoT09nQTZMeWx0cSJ9.ivjMTl_msqTIvt5gihELyaGOMpY5ogCV4AMQM7C5jEbj6Cqy1k4Ej2V5sGScparAAyKCsMVp3RvHea96ZLtobj2_ojrGyc7FlcsoBGlsRsA7n6o36nJwvu8iu8sLJAWz9Zn_RupkQm1I6ffTJaQXNNygx4B4mngxftmvcQBsZLhDlDMGeuH7XzUwo0WS578P89hmpJaMbpvVb0pyvR-QVwZPB5378s3Qam3BpK0sF5ReFSYhjlRtqx3sTGkMM5E3HNsvD_MSnn6ZkZMDTwJ1o8mmh4CbWCdG8FfMBKUjj20yG2P69NdE7L1PrJaAvJ2x7VPYCQ3aNcFyGddyi0n5Fw"
	var Bearer = "Bearer" + Token
	Token = utils.CreateAccessToken(Bearer, Token)
	Bearer = "Bearer" + Token
	return Bearer
}

func NewRechargeManager(db *gorm.DB) RechargeManager {
	Bearer := CreateBearer()
	return RechargeManager{db: db, Bearer: Bearer}
}

type RechargeManager struct {
	db     *gorm.DB
	Bearer string
}

func (a RechargeManager) CreateRecharge(c *gin.Context) {
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

func (a RechargeManager) GetRechargesStatus(c *gin.Context) {
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

func (a RechargeManager) DeleteRechargeByID(c *gin.Context) {
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
