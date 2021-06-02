package application

import (
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func NewUserController(userService services.UserService) UserController {
	return UserController{
		userService: userService,
	}
}

type UserController struct {
	userService services.UserService
}

func (a UserController) CreateUser(c *gin.Context) {
	//Valida o input
	var input input2.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.userService.CreateUser(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Response": "Usuário criado"})
}

func (a UserController) GetUserByDocument(c *gin.Context) {
	document := c.Param("document")
	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := a.userService.GetUserByDocument(input, document)
	if err == nil {
		c.JSON(200, user)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}

func (a UserController) UpdateUser(c *gin.Context) {
	userIDString := c.Param("userID")
	userIDInt, _ := strconv.Atoi(userIDString)
	userID := uint(userIDInt)

	var input input2.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.userService.UpdateUser(input, userID)
	if err == nil {
		c.JSON(http.StatusOK, "Usuário alterado com sucesso!")
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}

func (a UserController) DeleteUserByID(c *gin.Context) {
	userIDString := c.Param("userID")
	userIDInt, _ := strconv.Atoi(userIDString)
	userID := uint(userIDInt)

	var input input2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.userService.DeleteUserByID(input, userID)
	if err == nil {
		c.JSON(http.StatusOK, "Usuário deletado com sucesso!")
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
	}
}
