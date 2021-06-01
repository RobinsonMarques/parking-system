package crud

import (
	"github.com/RobinsonMarques/parking-system/database"
	"gorm.io/gorm"
)

func NewTrafficWardenCrud(db *gorm.DB) TrafficWardenCrud {
	return TrafficWardenCrud{db: db}
}

type TrafficWardenCrud struct {
	db *gorm.DB
}

func (t TrafficWardenCrud) CreateTrafficWarden(trafficWarden database.TrafficWarden) error {
	return t.db.Create(&trafficWarden).Error
}

func (t TrafficWardenCrud) GetTrafficWardenByEmail(email string) (database.TrafficWarden, error) {
	var trafficWarden database.TrafficWarden

	err := t.db.Where("Email = ?", email).First(&trafficWarden).Error
	return trafficWarden, err
}

func (t TrafficWardenCrud) GetTrafficWardenByID(id uint) (database.TrafficWarden, error) {
	var trafficWarden database.TrafficWarden

	err := t.db.Where("id = ?", id).First(&trafficWarden).Error
	return trafficWarden, err
}

func (t TrafficWardenCrud) UpdateTrafficWarden(trafficWarden database.TrafficWarden) error {
	err := t.db.Table("traffic_wardens").Where("id = ?", trafficWarden.ID).Update("name", trafficWarden.Person.Name).Error
	if err != nil {
		return err
	}
	err = t.db.Table("traffic_wardens").Where("id = ?", trafficWarden.ID).Update("email", trafficWarden.Person.Email).Error
	if err != nil {
		return err
	}
	err = t.db.Table("traffic_wardens").Where("id = ?", trafficWarden.ID).Update("password", trafficWarden.Person.Password).Error
	if err != nil {
		return err
	}
	return nil
}

func (t TrafficWardenCrud) DeleteTrafficWardenByID(trafficWardenID uint) error {
	err := t.db.Table("traffic_wardens").Where("id = ?", trafficWardenID).Delete(&database.TrafficWarden{}).Error
	return err
}
