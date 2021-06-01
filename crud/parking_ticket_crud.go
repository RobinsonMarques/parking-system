package crud

import (
	"github.com/RobinsonMarques/parking-system/database"
	"gorm.io/gorm"
	"time"
)

func NewParkingTicketCrud(db *gorm.DB) ParkingTicketCrud {
	return ParkingTicketCrud{db: db}
}

type ParkingTicketCrud struct {
	db *gorm.DB
}

func (p ParkingTicketCrud) CreateParkingTicket(parkingTicket database.ParkingTicket) error {
	return p.db.Create(&parkingTicket).Error
}

func (p ParkingTicketCrud) GetLastParkingTicketFromVehicle(id uint) ([]database.ParkingTicket, error) {
	var tickets []database.ParkingTicket

	err := p.db.Where("vehicle_id = ?", id).Last(&tickets).Error
	return tickets, err
}

func (p ParkingTicketCrud) UpdateEndTime(ticketID uint) error {
	currentTime := time.Now()
	err := p.db.Table("parking_tickets").Where("id = ?", ticketID).Update("end_time", currentTime.String()).Error
	return err
}

func (p ParkingTicketCrud) DeleteParkingTicketByID(parkingTicketID uint) error {
	err := p.db.Table("parking_tickets").Where("id = ?", parkingTicketID).Delete(&database.ParkingTicket{}).Error
	return err
}

func (p ParkingTicketCrud) DeleteParkingTicketByVehicleID(vehicleId uint) error {
	err := p.db.Table("parking_tickets").Where("vehicle_id = ?", vehicleId).Delete(&database.ParkingTicket{}).Error
	return err
}
