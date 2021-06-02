package application

import (
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/interfaces"
	"github.com/RobinsonMarques/parking-system/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func NewUserController(userInterface interfaces.UserInterface, vehicleInterface interfaces.VehicleInterface, rechargeInterface interfaces.RechargeInterface, billetInterface interfaces.BilletInterface, utilInterface interfaces.UtilInterface) UserController {
	return UserController{
		userInterface:     userInterface,
		vehicleInterface:  vehicleInterface,
		rechargeInterface: rechargeInterface,
		billetInterface:   billetInterface,
		utilInterface:     utilInterface,
	}
}

type UserController struct {
	userInterface     interfaces.UserInterface
	vehicleInterface  interfaces.VehicleInterface
	rechargeInterface interfaces.RechargeInterface
	billetInterface   interfaces.BilletInterface
	utilInterface     interfaces.UtilInterface
}

func (a UserController) CreateUser(c *gin.Context) {

	userService := services.NewUserService(a.userInterface, a.vehicleInterface, a.rechargeInterface, a.billetInterface, a.utilInterface)
	//Valida o input
	var input input2.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := userService.CreateUser(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Response": "Usuário criado"})
}

func (a UserController) GetUserByDocument(c *gin.Context) {
	userService := services.NewUserService(a.userInterface, a.vehicleInterface, a.rechargeInterface, a.billetInterface, a.utilInterface)

	document := c.Param("document")
	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := userService.GetUserByDocument(input, document)
	if err == nil {
		c.JSON(200, user)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}

func (a UserController) UpdateUser(c *gin.Context) {
	userService := services.NewUserService(a.userInterface, a.vehicleInterface, a.rechargeInterface, a.billetInterface, a.utilInterface)

	userIDString := c.Param("userID")
	userIDInt, _ := strconv.Atoi(userIDString)
	userID := uint(userIDInt)

	var input input2.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := userService.UpdateUser(input, userID)
	if err == nil {
		c.JSON(http.StatusOK, "Usuário alterado com sucesso!")
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}

func (a UserController) DeleteUserByID(c *gin.Context) {
	userService := services.NewUserService(a.userInterface, a.vehicleInterface, a.rechargeInterface, a.billetInterface, a.utilInterface)

	userIDString := c.Param("userID")
	userIDInt, _ := strconv.Atoi(userIDString)
	userID := uint(userIDInt)

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := userService.DeleteUserByID(input, userID)
	if err == nil {
		c.JSON(http.StatusOK, "Usuário deletado com sucesso!")
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}
