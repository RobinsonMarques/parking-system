package application

import (
	"github.com/RobinsonMarques/parking-system/crud"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/services"
	"github.com/RobinsonMarques/parking-system/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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
	userCrud := crud.NewUserCrud(a.db)
	rechargeCrud := crud.NewRechargeCrud(a.db)
	utilCrud := crud.NewUtilCrud(a.db)
	billetCrud := crud.NewBilletCrud(a.db)
	rechargeService, err := services.NewRechargeService(rechargeCrud, userCrud, utilCrud, billetCrud)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	} else {
		err = rechargeService.CreateRecharge(input, url, rechargeService)
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
	userCrud := crud.NewUserCrud(a.db)
	rechargeCrud := crud.NewRechargeCrud(a.db)
	utilCrud := crud.NewUtilCrud(a.db)
	billetCrud := crud.NewBilletCrud(a.db)
	rechargeService, err := services.NewRechargeService(rechargeCrud, userCrud, utilCrud, billetCrud)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	} else {
		err = rechargeService.GetRechargeStatus(input, rechargeService)
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
	userCrud := crud.NewUserCrud(a.db)
	rechargeCrud := crud.NewRechargeCrud(a.db)
	utilCrud := crud.NewUtilCrud(a.db)
	billetCrud := crud.NewBilletCrud(a.db)
	rechargeService, err := services.NewRechargeService(rechargeCrud, userCrud, utilCrud, billetCrud)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	} else {
		err = rechargeService.DeleteRechargeByID(input, rechargeID, rechargeService)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"Response": "Recarga deletada"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
		}
	}
}
