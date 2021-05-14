package main

import (
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/dependencies"
	"github.com/gin-gonic/gin"
)

func main() {
	db := dependencies.CreateConnection()

	r := gin.Default()

	//Path /getallvehicles
	r.GET("/getallvehicles", func(c *gin.Context) {
		vehicles := crud.GetAllVehicles(db)
		c.JSON(200, vehicles)
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
