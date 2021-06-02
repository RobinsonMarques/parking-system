package application

import (
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/interfaces"
	"github.com/RobinsonMarques/parking-system/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func NewBilletManager(billetInterface interfaces.BilletInterface, utilInterface interfaces.UtilInterface) BilletManager {
	return BilletManager{
		billetInterface: billetInterface,
		utilInterface:   utilInterface,
	}
}

type BilletManager struct {
	billetInterface interfaces.BilletInterface
	utilInterface   interfaces.UtilInterface
}

func (a BilletManager) DeleteBilletByID(c *gin.Context) {
	billetService := services.NewBilletService(a.billetInterface, a.utilInterface)

	billetIDString := c.Param("billetID")
	billetIDInt, _ := strconv.Atoi(billetIDString)
	billetID := uint(billetIDInt)

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := billetService.DeleteBilletByID(input, billetID)

	if err == nil {
		c.JSON(http.StatusOK, "Boleto deletado com sucesso!")
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}
