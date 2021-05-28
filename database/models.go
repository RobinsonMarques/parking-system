package database

import "gorm.io/gorm"

//Modelos das tabelas do BD

type Person struct {
	Name     string `gorm:"not null" json:"Name" binding:"required"`
	Email    string `gorm:"unique;not null" json:"Email" binding:"required"`
	Password string `gorm:"not null" json:"Password" binding:"required"`
}

type User struct {
	gorm.Model
	Person   Person     `gorm:"embedded"`
	Document string     `gorm:"unique;not null" json:"Document"`
	Balance  float64    `gorm:"default:0.0" json:"Balance"`
	Recharge []Recharge `json:"Recharges"`
	Vehicle  []Vehicle  `json:"Vehicles"`
}

type TrafficWarden struct {
	gorm.Model
	Person Person `gorm:"embedded"`
}

type Admin struct {
	gorm.Model
	Person Person `gorm:"embedded"`
}

type Vehicle struct {
	gorm.Model
	LicensePlate  string          `gorm:"unique;not null" json:"LicensePlate"`
	VehicleModel  string          `gorm:"not null" json:"VehicleModel"`
	VehicleType   string          `gorm:"not null" json:"VehicleType"`
	IsActive      bool            `gorm:"not null default:false" json:"IsActive"`
	IsParked      bool            `gorm:"not null default:false" json:"IsParked"`
	UserID        uint            `gorm:"not null" json:"UserID"`
	ParkingTicket []ParkingTicket `json:"ParkingTickets"`
}

type ParkingTicket struct {
	gorm.Model
	Location    string  `gorm:"not null" json:"Location"`
	ParkingTime int     `gorm:"not null" json:"ParkingTime"`
	StartTime   string  `gorm:"not null" json:"StartTime"`
	EndTime     string  `gorm:"not null" json:"EndTime"`
	Price       float64 `gorm:"not null" json:"Price"`
	VehicleID   uint    `gorm:"not null" json:"VehicleID"`
}

type Recharge struct {
	gorm.Model
	Date         string  `gorm:"not null" json:"Date"`
	Value        float64 `gorm:"not null" json:"Value"`
	IsPaid       bool    `gorm:"not null" json:"IsPaid"`
	PaymentType  string  `gorm:"not null" json:"PaymentType"`
	UserID       uint    `gorm:"not null" json:"UserID"`
	RechargeHash string  `gorm:"not null" json:"rechargeHash"`
	Billet       Billet  `json:"Billet"`
}

type Billet struct {
	gorm.Model
	BilletLink string `json:"BilletLink"`
	RechargeID uint   `json:"RechargeID"`
}
