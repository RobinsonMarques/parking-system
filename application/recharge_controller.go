package application

import (
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/interfaces"
	"github.com/RobinsonMarques/parking-system/services"
	"github.com/RobinsonMarques/parking-system/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

func NewRechargeController(rechargeInterface interfaces.RechargeInterface, userInterface interfaces.UserInterface, utilInterface interfaces.UtilInterface, billetInterface interfaces.BilletInterface) (RechargeController, error) {
	Bearer, err := CreateBearer()
	if err != nil {
		return RechargeController{}, err
	}
	return RechargeController{
		rechargeInterface: rechargeInterface,
		userInterface:     userInterface,
		utilInterface:     utilInterface,
		billetInterface:   billetInterface,
		Bearer:            Bearer}, nil
}

type RechargeController struct {
	rechargeInterface interfaces.RechargeInterface
	userInterface     interfaces.UserInterface
	utilInterface     interfaces.UtilInterface
	billetInterface   interfaces.BilletInterface
	Bearer            string
}

func (a RechargeController) CreateRecharge(c *gin.Context) {
	url := "https://sandbox.boletobancario.com/api-integration/charges"
	var input input2.CreateRecharge
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rechargeService, err := services.NewRechargeService(a.rechargeInterface, a.userInterface, a.utilInterface, a.billetInterface)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	} else {
		err = rechargeService.CreateRecharge(input, url)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"Response": "Recarga criada"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
		}
	}

}

func (a RechargeController) GetRechargesStatus(c *gin.Context) {
	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rechargeService, err := services.NewRechargeService(a.rechargeInterface, a.userInterface, a.utilInterface, a.billetInterface)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	} else {
		err = rechargeService.GetRechargeStatus(input)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"Response": "Saldo alterado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
		}
	}

}

func (a RechargeController) DeleteRechargeByID(c *gin.Context) {
	rechargeIDString := c.Param("rechargeID")
	rechargeIDInt, _ := strconv.Atoi(rechargeIDString)
	rechargeID := uint(rechargeIDInt)

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rechargeService, err := services.NewRechargeService(a.rechargeInterface, a.userInterface, a.utilInterface, a.billetInterface)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	} else {
		err = rechargeService.DeleteRechargeByID(input, rechargeID)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"Response": "Recarga deletada"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
		}
	}
}
