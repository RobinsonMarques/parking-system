package main

import (
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/database"
	"github.com/RobinsonMarques/parking-system/dependencies"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	db := dependencies.CreateConnection()

	r := gin.Default()

	//Path get all vehicles
	r.GET("/vehicles", func(c *gin.Context) {
		//valida o input
		var input input2.LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := utils.Login(input.Email, input.Password, db)

		if resp == "trafficWarden" || resp == "admin" {
			vehicles := crud.GetAllVehicles(db)
			c.JSON(200, vehicles)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		}

	})

	//Path create user
	r.POST("/users", func(c *gin.Context) {

		//Valida o input
		var input input2.CreateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		input.Person.Password = utils.CreateHashPassword(input.Person.Password)

		//Cria o user

		user := database.User{
			Person:   input.Person,
			Document: input.Document,
			Balance:  0,
			Recharge: nil,
			Vehicle:  nil,
		}
		crud.CreateUser(user, db)
		c.JSON(http.StatusOK, gin.H{"Response": "Usuário criado"})

	})

	//Path create admin
	r.POST("/admins", func(c *gin.Context) {
		//Valida o input
		var input input2.CreateAdminInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := utils.Login(input.LoginInput.Email, input.LoginInput.Password, db)

		if resp == "admin" {
			input.Person.Password = utils.CreateHashPassword(input.Person.Password)

			//Cria o admin
			admin := database.Admin{
				Person: input.Person,
			}
			crud.CreateAdmin(admin, db)
			c.JSON(http.StatusOK, gin.H{"Response": "Admin criado"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		}
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
