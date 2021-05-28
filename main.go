package main

import (
	"github.com/RobinsonMarques/parking-system/application"
	"github.com/RobinsonMarques/parking-system/database"
	"github.com/RobinsonMarques/parking-system/dependencies"
	"gorm.io/gorm"
	"log"
)

func main() {
	db, err := dependencies.CreateConnection()
	if err != nil {
		log.Fatal(err.Error())
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer sqlDB.Close()
	if err := migrate(db); err != nil {
		log.Fatal(err.Error())
	}
	app, err := application.NewApplication(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

func migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&database.Person{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&database.User{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&database.TrafficWarden{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&database.Admin{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&database.ParkingTicket{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&database.Vehicle{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&database.Billet{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&database.Recharge{}); err != nil {
		return err
	}
	return nil
}
