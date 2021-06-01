package application

import (
	"github.com/RobinsonMarques/parking-system/crud"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func NewUserController(userCrud crud.UserCrud, vehicleCrud crud.VehicleCrud, rechargeCrud crud.RechargeCrud, billetCrud crud.BilletCrud, utilCrud crud.UtilCrud) UserController {
	return UserController{
		userCrud:     userCrud,
		vehicleCrud:  vehicleCrud,
		rechargeCrud: rechargeCrud,
		billetCrud:   billetCrud,
		utilCrud:     utilCrud,
	}
}

type UserController struct {
	userCrud     crud.UserCrud
	vehicleCrud  crud.VehicleCrud
	rechargeCrud crud.RechargeCrud
	billetCrud   crud.BilletCrud
	utilCrud     crud.UtilCrud
}

func (a UserController) CreateUser(c *gin.Context) {

	userService := services.NewUserService(a.userCrud, a.vehicleCrud, a.rechargeCrud, a.billetCrud, a.utilCrud)
	//Valida o input
	var input input2.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := userService.CreateUser(input, userService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Response": "Usuário criado"})
}

func (a UserController) GetUserByDocument(c *gin.Context) {
	userService := services.NewUserService(a.userCrud, a.vehicleCrud, a.rechargeCrud, a.billetCrud, a.utilCrud)

	document := c.Param("document")
	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := userService.GetUserByDocument(input, document, userService)
	if err == nil {
		c.JSON(200, user)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}

func (a UserController) UpdateUser(c *gin.Context) {
	userService := services.NewUserService(a.userCrud, a.vehicleCrud, a.rechargeCrud, a.billetCrud, a.utilCrud)

	userIDString := c.Param("userID")
	userIDInt, _ := strconv.Atoi(userIDString)
	userID := uint(userIDInt)

	var input input2.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := userService.UpdateUser(input, userID, userService)
	if err == nil {
		c.JSON(http.StatusOK, "Usuário alterado com sucesso!")
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}

func (a UserController) DeleteUserByID(c *gin.Context) {
	userService := services.NewUserService(a.userCrud, a.vehicleCrud, a.rechargeCrud, a.billetCrud, a.utilCrud)

	userIDString := c.Param("userID")
	userIDInt, _ := strconv.Atoi(userIDString)
	userID := uint(userIDInt)

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := userService.DeleteUserByID(input, userID, userService)
	if err == nil {
		c.JSON(http.StatusOK, "Usuário deletado com sucesso!")
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}
