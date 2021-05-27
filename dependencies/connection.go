package dependencies

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//Conexão com o BD
func CreateConnection() (db *gorm.DB, err error) {
	//dsn := "host=110.97.88.192 user=park password=P@ssword dbname=park port=15432 sslmode=disable TimeZone=UTC-3"
	dsn := "host=localhost user=postgres password=Ak47#mp5 dbname=api port=5432 sslmode=disable TimeZone=UTC-3"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return db, err
}
