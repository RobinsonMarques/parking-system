package application

import (
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func NewBilletManager(db *gorm.DB) BilletManager {
	return BilletManager{db: db}
}

type BilletManager struct {
	db *gorm.DB
}

func (a BilletManager) DeleteBilletByID(c *gin.Context) {
	billetIDString := c.Param("billetID")
	billetIDInt, _ := strconv.Atoi(billetIDString)
	billetID := uint(billetIDInt)

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	billetService := services.NewBilletService(a.db)
	err := billetService.DeleteBilletByID(input, billetID)

	if err == nil {
		c.JSON(http.StatusOK, "Boleto deletado com sucesso!")
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err})
	}
}
