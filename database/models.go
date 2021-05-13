package database

//Modelos das tabelas do BD

type Person struct {
	PersonID uint   `gorm:"primary_key;AUTO_INCREMENT" json:"ID"`
	Name     string `gorm:"not null" json:"Name"`
	Email    string `gorm:"unique;not null" json:"Email"`
	Password string `gorm:"not null" json:"Password"`
}

type User struct {
	UserID   uint `gorm:"primary_key;AUTO_INCREMENT" json:"ID"`
	Person   Person
	Document string     `gorm:"unique;not null" json:"Document"`
	Balance  float64    `gorm:"default:0.0" json:"Balance"`
	Recharge []Recharge `json:"Recharges"`
	Vehicle  []Vehicle  `json:"Vehicles"`
}

type TrafficWarden struct {
	TrafficWardenID uint `gorm:"primary_key;AUTO_INCREMENT" json:"ID"`
	Person          Person
}

type Admin struct {
	AdminId uint `gorm:"primary_key;AUTO_INCREMENT" json:"ID"`
	Person  Person
}

type Vehicle struct {
	VehicleID     uint            `gorm:"primary_key;AUTO_INCREMENT" json:"ID"`
	LicensePlate  string          `gorm:"unique;not null" json:"LicensePlate"`
	VehicleModel  string          `gorm:"not null" json:"VehicleModel"`
	VehicleType   string          `gorm:"not null" json:"VehicleType"`
	IsActive      bool            `gorm:"not null" json:"IsActive"`
	IsParked      bool            `gorm:"not null" json:"IsParked"`
	UserID        uint            `gorm:"not null" json:"UserID"`
	ParkingTicket []ParkingTicket `json:"ParkingTickets"`
}

type ParkingTicket struct {
	ParkingTicketID uint    `gorm:"primary_key;AUTO_INCREMENT" json:"ID"`
	Location        string  `gorm:"not null" json:"Location"`
	ParkingTime     int     `gorm:"not null" json:"ParkingTime"`
	StartTime       string  `gorm:"not null" json:"StartTime"`
	EndTime         string  `gorm:"not null" json:"EndTime"`
	Price           float64 `gorm:"not null" json:"Price"`
	VehicleID       uint    `gorm:"not null" json:"VehicleID"`
}

type Recharge struct {
	RechargeID  uint    `gorm:"primary_key;AUTO_INCREMENT" json:"ID"`
	Date        string  `gorm:"not null" json:"Date"`
	Value       float64 `gorm:"not null" json:"Value"`
	IsPaid      bool    `gorm:"not null" json:"IsPaid"`
	PaymentType string  `gorm:"not null" json:"PaymentType"`
	UserID      uint    `gorm:"not null" json:"UserID"`
	Billet      Billet  `json:"Billet"`
}

type Billet struct {
	BilletID   uint   `gorm:"primary_key;AUTO_INCREMENT" json:"ID"`
	BilletLink string `json:"BilletLink"`
	RechargeID uint   `gorm:"not null" json:"RechargeID"`
}
