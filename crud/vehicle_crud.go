package crud

import (
	"github.com/RobinsonMarques/parking-system/database"
	"gorm.io/gorm"
)

func NewVehicleCrud(db *gorm.DB) VehicleCrud {
	return VehicleCrud{db: db}
}

type VehicleCrud struct {
	db *gorm.DB
}

func (v VehicleCrud) CreateVehicle(vehicle database.Vehicle) error {
	return v.db.Create(&vehicle).Error
}

func (v VehicleCrud) GetAllVehicles() ([]database.Vehicle, error) {
	var vehicle []database.Vehicle

	err := v.db.Find(&vehicle).Error
	return vehicle, err
}

func (v VehicleCrud) GetVehiclesByUserId(userID uint) ([]database.Vehicle, error) {
	var vehicles []database.Vehicle
	err := v.db.Where("user_id = ?", userID).Find(&vehicles).Error
	return vehicles, err
}

func (v VehicleCrud) GetVehicleByLicensePlate(licensePlate string) (database.Vehicle, error) {
	var vehicle database.Vehicle

	err := v.db.Where("license_plate = ?", licensePlate).First(&vehicle).Error
	return vehicle, err
}

func (v VehicleCrud) GetVehicleById(id uint) (database.Vehicle, error) {
	var vehicle database.Vehicle

	err := v.db.Where("id = ?", id).First(&vehicle).Error
	return vehicle, err
}

func (v VehicleCrud) UpdateVehicle(vehicle database.Vehicle) error {
	err := v.db.Table("vehicles").Where("id = ?", vehicle.ID).Update("license_plate", vehicle.LicensePlate).Error
	if err != nil {
		return err
	}
	err = v.db.Table("vehicles").Where("id = ?", vehicle.ID).Update("vehicle_model", vehicle.VehicleModel).Error
	if err != nil {
		return err
	}
	err = v.db.Table("vehicles").Where("id = ?", vehicle.ID).Update("vehicle_type", vehicle.VehicleType).Error
	if err != nil {
		return err
	}
	return nil
}

func (v VehicleCrud) UpdateVehicleOwner(vehicleID, newOwnerID uint) error {
	err := v.db.Table("vehicles").Where("id = ?", vehicleID).Update("user_id", newOwnerID).Error
	return err
}

func (v VehicleCrud) UpdateIsParked(vehicleID uint, value bool) error {
	err := v.db.Table("vehicles").Where("id = ?", vehicleID).Update("is_parked", value).Error
	return err
}

func (v VehicleCrud) UpdateIsActive(vehicleID uint, value bool) error {
	err := v.db.Table("vehicles").Where("id = ?", vehicleID).Update("is_active", value).Error
	return err
}

func (v VehicleCrud) DeleteVehicleByID(vehicleID uint, crud Crud) error {
	err := v.db.Table("vehicles").Where("id = ?", vehicleID).Delete(&database.Vehicle{}).Error
	crud.ParkingTicketCrud.DeleteParkingTicketByVehicleID(vehicleID)
	return err
}

func (v VehicleCrud) DeleteVehiclesByUserID(userID uint) error {
	err := v.db.Table("vehicles").Where("user_id = ?", userID).Delete(&database.Vehicle{}).Error
	return err
}
