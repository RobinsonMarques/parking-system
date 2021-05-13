package dependencies

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

//Conexão com o BD
func CreateConnection() (db *gorm.DB) {
	//dsn := "host=110.97.88.192 user=park password=P@ssword dbname=park port=15432 sslmode=disable TimeZone=UTC-3"
	db, err := gorm.Open(sqlite.Open("parking.db"), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}
