package input

import "github.com/RobinsonMarques/parking-system/database"

type CreateUserInput struct {
	Person   database.Person
	Document string `json:"Document" binding:"required"`
}

type UpdateUserInput struct {
	Person     database.Person
	Document   string     `json:"Document"`
	LoginInput LoginInput `json:"Login" binding:"required"`
}

type UpdateAdminInput struct {
	Person     database.Person
	LoginInput LoginInput `json:"Login" binding:"required"`
}

type UpdateTrafficWarden struct {
	Person     database.Person
	LoginInput LoginInput `json:"Login" binding:"required"`
}

type UpdateVehicle struct {
	LicensePlate string     `gorm:"unique;not null" json:"LicensePlate"`
	VehicleModel string     `gorm:"not null" json:"VehicleModel"`
	VehicleType  string     `gorm:"not null" json:"VehicleType"`
	LoginInput   LoginInput `json:"Login" binding:"required"`
}

type CreateAdminInput struct {
	Person     database.Person
	LoginInput LoginInput `json:"Login" binding:"required"`
}

type CreateTrafficWarden struct {
	Person     database.Person
	LoginInput LoginInput `json:"Login" binding:"required"`
}

type CreateParkingTicket struct {
	Login       LoginInput
	Location    string `json:"Location"`
	ParkingTime int    `json:"ParkingTime"`
	VehicleID   uint   `json:"VehicleID"`
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

type CreateRecharge struct {
	Value       float64    `json:"Value" binding:"required"`
	PaymentType string     `json:"PaymentType" binding:"required"`
	LoginInput  LoginInput `json:"Login" binding:"required"`
}
