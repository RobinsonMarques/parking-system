package input

import "github.com/RobinsonMarques/parking-system/database"

type CreateUserInput struct {
	Person   database.Person
	Document string  `json:"Document" binding:"required"`
	Balance  float64 `json:"Balance"`
}

type CreateAdminInput struct {
	Person     database.Person
	LoginInput LoginInput `json:"Login" binding:"required"`
}

type CreateTrafficWarden struct {
	Person     database.Person
	LoginInput LoginInput `json:"Login" binding:"required"`
}

type CreateVehicle struct {
	LicensePlate string `json:"LicensePlate" binding:"required"`
	VehicleModel string `json:"VehicleModel" binding:"required"`
	VehicleType  string `json:"VehicleType" binding:"required"`
	UserID       uint   `json:"UserID" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"Email" binding:"required"`
	Password string `json:"Password" binding:"required"`
}
