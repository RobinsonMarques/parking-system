package interfaces

import (
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/database"
)

type UserInterface interface {
	CreateUser(user database.User) error
	GetUserByEmail(email string) (database.User, error)
	GetUserByID(id uint) (database.User, error)
	GetUserByDocument(document string) (database.User, error)
	GetBalance(email string, userCrud crud.UserCrud) (float64, error)
	UpdateUser(user database.User) error
	UpdateBalance(email string, extra float64) error
	DeleteUserByID(userID uint) error
}

type VehicleInterface interface {
	CreateVehicle(vehicle database.Vehicle) error
	GetAllVehicles() ([]database.Vehicle, error)
	GetVehiclesByUserId(userID uint) ([]database.Vehicle, error)
	GetVehicleByLicensePlate(licensePlate string) (database.Vehicle, error)
	GetVehicleById(id uint) (database.Vehicle, error)
	UpdateVehicle(vehicle database.Vehicle) error
	UpdateVehicleOwner(vehicleID, newOwnerID uint) error
	UpdateIsParked(vehicleID uint, value bool) error
	UpdateIsActive(vehicleID uint, value bool) error
	AlterVehicleStatus(vehicle database.Vehicle, parkingTime int) error
	DeleteVehicleByID(vehicleID uint) error
	DeleteVehiclesByUserID(userID uint) error
}

type RechargeInterface interface {
	CreateRecharge(recharge database.Recharge) error
	GetUserUnpaidRechargesByID(userID uint) ([]database.Recharge, error)
	GetRechargeByUserId(userID uint) ([]database.Recharge, error)
	UpdateIsPaid(rechargeID uint) error
	DeleteRechargeByID(rechargeID uint) error
	DeleteRechargeByUserID(userID uint) error
}

type BilletInterface interface {
	CreateBillet(billet database.Billet) error
	GetBilletByRechargeId(rechargeID uint) (database.Billet, error)
	GetBilletsByRechargeID(rechargeID uint) ([]database.Billet, error)
	UpdateBilletLink(billetID uint, link string) error
	DeleteBilletByID(billetID uint) error
	DeleteBilletByRechargeID(rechargeID uint) error
}

type UtilInterface interface {
	Login(email string, password string) string
}

type AdminInterface interface {
	CreateAdmin(admin database.Admin) error
	GetAdminByEmail(email string) (database.Admin, error)
	GetAdminByID(id uint) (database.Admin, error)
	UpdateAdmin(admin database.Admin) error
	DeleteAdminByID(adminID uint) error
}

type TrafficWardenInterface interface {
	CreateTrafficWarden(trafficWarden database.TrafficWarden) error
	GetTrafficWardenByEmail(email string) (database.TrafficWarden, error)
	GetTrafficWardenByID(id uint) (database.TrafficWarden, error)
	UpdateTrafficWarden(trafficWarden database.TrafficWarden) error
	DeleteTrafficWardenByID(trafficWardenID uint) error
}

type ParkingTicketInterface interface {
	CreateParkingTicket(parkingTicket database.ParkingTicket) error
	GetLastParkingTicketFromVehicle(id uint) ([]database.ParkingTicket, error)
	UpdateEndTime(ticketID uint) error
	DeleteParkingTicketByID(parkingTicketID uint) error
	DeleteParkingTicketByVehicleID(vehicleId uint) error
}
