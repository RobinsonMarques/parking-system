package application

import (
	"github.com/RobinsonMarques/parking-system/crud"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func NewAdminController(adminCrud crud.AdminCrud, utilCrud crud.UtilCrud) AdminController {
	return AdminController{
		adminCrud: adminCrud,
		utilCrud:  utilCrud,
	}
}

type AdminController struct {
	adminCrud crud.AdminCrud
	utilCrud  crud.UtilCrud
}

func (a AdminController) CreateAdmin(c *gin.Context) {
	adminService := services.NewAdminService(a.adminCrud, a.utilCrud)
	//Valida o input
	var input input2.CreateAdminInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := adminService.CreateAdmin(input, adminService)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Response": "Admin criado"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}

func (a AdminController) UpdateAdmin(c *gin.Context) {
	adminService := services.NewAdminService(a.adminCrud, a.utilCrud)

	adminIDString := c.Param("adminID")
	adminIDInt, _ := strconv.Atoi(adminIDString)
	adminID := uint(adminIDInt)

	var input input2.UpdateAdminInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := adminService.UpdateAdmin(input, adminID, adminService)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Response": "Admin alterado"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}

func (a AdminController) DeleteAdminByID(c *gin.Context) {
	adminService := services.NewAdminService(a.adminCrud, a.utilCrud)

	adminIDString := c.Param("adminID")
	adminIDInt, _ := strconv.Atoi(adminIDString)
	adminID := uint(adminIDInt)

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := adminService.DeleteAdminByID(input, adminID, adminService)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"Response": "Admin deletado"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}
