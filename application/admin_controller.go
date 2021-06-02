package application

import (
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func NewAdminController(adminService services.AdminService) AdminController {
	return AdminController{
		adminService: adminService,
	}
}

type AdminController struct {
	adminService services.AdminService
}

func (a AdminController) CreateAdmin(c *gin.Context) {
	//Valida o input
	var input input2.CreateAdminInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.adminService.CreateAdmin(input)
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
	err := a.adminService.UpdateAdmin(input, adminID)
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
	err := a.adminService.DeleteAdminByID(input, adminID)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Response": "Admin deletado"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}
