package application

import (
	"github.com/RobinsonMarques/parking-system/crud"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func NewBilletManager(billetCrud crud.BilletCrud, utilCrud crud.UtilCrud) BilletManager {
	return BilletManager{
		billetCrud: billetCrud,
		utilCrud:   utilCrud,
	}
}

type BilletManager struct {
	billetCrud crud.BilletCrud
	utilCrud   crud.UtilCrud
}

func (a BilletManager) DeleteBilletByID(c *gin.Context) {
	billetService := services.NewBilletService(a.billetCrud, a.utilCrud)

	billetIDString := c.Param("billetID")
	billetIDInt, _ := strconv.Atoi(billetIDString)
	billetID := uint(billetIDInt)

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := billetService.DeleteBilletByID(input, billetID, billetService)

	if err == nil {
		c.JSON(http.StatusOK, "Boleto deletado com sucesso!")
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}
