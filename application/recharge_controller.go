package application

import (
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/services"
	"github.com/RobinsonMarques/parking-system/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func CreateBearer() (string, error) {
	var Token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX25hbWUiOiJyb2JpbmhvbWFycXVlcy5ybTJAZ21haWwuY29tIiwic2NvcGUiOlsiYWxsIl0sImV4cCI6MTYyMjMwNzcwMSwianRpIjoicWVfdkZsVHV4cWhKWS0yUThiOVNsUUFBU1c4IiwiY2xpZW50X2lkIjoiUzNDeUtoT09nQTZMeWx0cSJ9.fvgQx7DyCebX6VyjGEd-fC3SS25gFBUCeZDoFZJUcLmtHSuMDz3FA-gL2OXCcDzAR9VQWeMpDZbIiTN6CQD6WvoNdJr1tuESCm1_wfD3q0Xjrni4jtEOgmCUIj4CKzaS-0Y2flF80D1SgID2DWm5OrtkXHnmeiDcvI7fxu_glpIjqcZlaHr-N-t7aEAslljKCOjRcp_0MWcA2CgTp9PbxIGcjRn3QCzLkKwSHy9IoiAR6j7aTp4gc9lt3UUE_HerVjDM0u1aMfVzLY4Ms3UGqpg38tQJ52HL64EnJvW5hkQlV0TDQi8eM6L6-a8oRBT3VCdwpUl0Dvjled1C1tckcQ"
	var Bearer = "Bearer" + Token
	Token, err := utils.CreateAccessToken(Bearer, Token)
	if err != nil {
		return "", err
	}
	Bearer = "Bearer" + Token
	return Bearer, nil
}

func NewRechargeManager(db *gorm.DB) (RechargeManager, error) {
	Bearer, err := CreateBearer()
	if err != nil {
		return RechargeManager{}, err
	}
	return RechargeManager{db: db, Bearer: Bearer}, nil
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
	rechargeService, err := services.NewRechargeService(a.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
	err = rechargeService.CreateRecharge(input, url)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Response": "Recarga criada"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}

func (a RechargeManager) GetRechargesStatus(c *gin.Context) {
	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rechargeService, err := services.NewRechargeService(a.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
	err = rechargeService.GetRechargeStatus(input)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Response": "Saldo alterado"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
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

	rechargeService, err := services.NewRechargeService(a.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
	err = rechargeService.DeleteRechargeByID(input, rechargeID)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Response": "Recarga deletada"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}
