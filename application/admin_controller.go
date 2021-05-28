package application

import (
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func NewAdminController(db *gorm.DB) AdminController {
	return AdminController{db: db}
}

type AdminController struct {
	db *gorm.DB
}

func (a AdminController) CreateAdmin(c *gin.Context) {
	//Valida o input
	var input input2.CreateAdminInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	adminService := services.NewAdminService(a.db)
	err := adminService.CreateAdmin(input)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Response": "Admin criado"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}

func (a AdminController) UpdateAdmin(c *gin.Context) {
	adminIDString := c.Param("adminID")
	adminIDInt, _ := strconv.Atoi(adminIDString)
	adminID := uint(adminIDInt)

	var input input2.UpdateAdminInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	adminService := services.NewAdminService(a.db)
	err := adminService.UpdateAdmin(input, adminID)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Response": "Admin alterado"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}

func (a AdminController) DeleteAdminByID(c *gin.Context) {
	adminIDString := c.Param("adminID")
	adminIDInt, _ := strconv.Atoi(adminIDString)
	adminID := uint(adminIDInt)

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	adminService := services.NewAdminService(a.db)
	err := adminService.DeleteAdminByID(input, adminID)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Response": "Admin deletado"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}
