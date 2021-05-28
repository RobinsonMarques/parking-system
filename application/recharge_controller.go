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
	var Token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX25hbWUiOiJyb2JpbmhvbWFycXVlcy5ybTJAZ21haWwuY29tIiwic2NvcGUiOlsiYWxsIl0sImV4cCI6MTYyMjMxNjQ3NSwianRpIjoiVk9rckZUcmxOMTFzZ1Ffazh6TDRWTE9WRWVZIiwiY2xpZW50X2lkIjoiUzNDeUtoT09nQTZMeWx0cSJ9.VPPewOzPaxpHazBxOyA-58zXMI0xE_9R5-zTvBET2kZkbNenONiFz336pPJ0rxv8SbYBItFFW7o9-YMolPknn71gqTAtx0BuPL-la6K8sbK3GtGuavN2P4JK3LtukD0mv0Ehu-HZGC2wgVZIz0kqEXANMS4lfm202GpPp87-jDbhQdnOyVcyEGa3IQ7KWDUFB2TWEH815iToIgHR-1aDbDm9p0ItdNhR65BqomHYv_a7XyW4p40AGtgJ9c67tLhPOTaOMQXlxFOIkdldVYEI0LvmkGiafbJdMqtO0Zx-FNN4HK0p0nwQhFVeEEBt6x0BCghuj1L1JRtAdMR6J2YS3w"
	var Bearer = "Bearer" + Token
	Token, err := utils.CreateAccessToken(Bearer, Token)
	if err != nil {
		return "", err
	}
	Bearer = "Bearer" + Token
	return Bearer, nil
}

func NewRechargeController(db *gorm.DB) (RechargeController, error) {
	Bearer, err := CreateBearer()
	if err != nil {
		return RechargeController{}, err
	}
	return RechargeController{db: db, Bearer: Bearer}, nil
}

type RechargeController struct {
	db     *gorm.DB
	Bearer string
}

func (a RechargeController) CreateRecharge(c *gin.Context) {
	url := "https://sandbox.boletobancario.com/api-integration/charges"
	var input input2.CreateRecharge
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rechargeService, err := services.NewRechargeService(a.db)
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
	rechargeService, err := services.NewRechargeService(a.db)
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
