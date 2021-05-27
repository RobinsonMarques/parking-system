package crud

import (
	"github.com/RobinsonMarques/parking-system/database"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
	"time"
)

type Result struct {
	Data *gorm.DB
}

func CreateUser(user database.User, db *gorm.DB) Result {
	data := db.Create(&user)
	return Result{Data: data}
}

func CreateVehicle(vehicle database.Vehicle, db *gorm.DB) Result {

	data := db.Create(&vehicle)
	return Result{Data: data}
}

func CreateTrafficWarden(trafficWarden database.TrafficWarden, db *gorm.DB) Result {
	data := db.Create(&trafficWarden)
	return Result{Data: data}
}

func CreateAdmin(admin database.Admin, db *gorm.DB) Result {
	data := db.Create(&admin)
	return Result{Data: data}
}

func CreateParkingTicket(parkingTicket database.ParkingTicket, db *gorm.DB) Result {
	data := db.Create(&parkingTicket)
	return Result{Data: data}
}

func CreateRecharge(recharge database.Recharge, db *gorm.DB) Result {
	data := db.Create(&recharge)
	return Result{Data: data}
}

func CreateBillet(billet database.Billet, db *gorm.DB) Result {
	data := db.Create(&billet)
	return Result{Data: data}
}

func GetUserByEmail(email string, db *gorm.DB) database.User {
	var user database.User

	db.Where("email = ?", email).First(&user)
	return user
}

func GetUserUnpaidRechargesByID(userID uint, db *gorm.DB) []database.Recharge {
	var recharges []database.Recharge
	db.Where("is_paid = false").Find(&recharges)
	return recharges
}

func GetUserByID(id uint, db *gorm.DB) (database.User, error) {
	var user database.User
	err := db.Where("id = ?", id).First(&user).Error
	return user, err
}

func GetTrafficWardenByEmail(email string, db *gorm.DB) database.TrafficWarden {
	var trafficWarden database.TrafficWarden

	db.Where("Email = ?", email).First(&trafficWarden)
	return trafficWarden
}

func GetTrafficWardenByID(id uint, db *gorm.DB) (database.TrafficWarden, error) {
	var warden database.TrafficWarden
	err := db.Where("id = ?", id).First(&warden).Error
	return warden, err
}

func GetAdminByEmail(email string, db *gorm.DB) database.Admin {
	var admin database.Admin
	db.Where("Email = ?", email).First(&admin)
	return admin
}

func GetAdminByID(id uint, db *gorm.DB) (database.Admin, error) {
	var admin database.Admin
	err := db.Where("id = ?", id).First(&admin).Error
	return admin, err
}

func GetAllVehicles(db *gorm.DB) []database.Vehicle {
	var vehicle []database.Vehicle

	db.Find(&vehicle)
	return vehicle
}

func GetRechargeByUserId(userID uint, db *gorm.DB) ([]database.Recharge, error) {
	var recharges []database.Recharge
	err := db.Where("user_id = ?", userID).Find(&recharges).Error
	return recharges, err
}

func GetBilletByRechargeId(rechargeID uint, db *gorm.DB) (database.Billet, error) {
	var billet database.Billet
	err := db.Where("recharge_id = ?", rechargeID).Find(&billet).Error
	return billet, err
}

func GetVehiclesByUserId(userID uint, db *gorm.DB) ([]database.Vehicle, error) {
	var vehicles []database.Vehicle
	err := db.Where("user_id = ?", userID).Find(&vehicles).Error
	return vehicles, err
}

func GetUserByDocument(document string, db *gorm.DB) (database.User, error) {
	var user database.User

	err := db.Where("Document = ?", document).First(&user).Error
	return user, err
}

func GetVehicleByLicensePlate(licensePlate string, db *gorm.DB) (database.Vehicle, error) {
	var vehicle database.Vehicle

	err := db.Where("license_plate = ?", licensePlate).First(&vehicle).Error
	return vehicle, err
}

func GetVehicleById(id uint, db *gorm.DB) database.Vehicle {
	var vehicle database.Vehicle

	db.Where("id = ?", id).First(&vehicle)
	return vehicle
}

func GetLastParkingTicketFromVehicle(id uint, db *gorm.DB) ([]database.ParkingTicket, error) {
	var tickets []database.ParkingTicket

	err := db.Where("vehicle_id = ?", id).Last(&tickets).Error
	return tickets, err
}

func GetBilletsByRechargeID(rechargeID uint, db *gorm.DB) []database.Billet {
	var billetts []database.Billet

	db.Where("recharge_id = ?", rechargeID).Find(&billetts)
	return billetts
}

func GetBalance(email string, db *gorm.DB) float64 {
	user := GetUserByEmail(email, db)
	balance := user

	return balance.Balance
}

func GetPassword(email string, userType string, db *gorm.DB) string {
	if userType == "user" {
		user := GetUserByEmail(email, db)
		return user.Person.Password
	} else if userType == "admin" {
		admin := GetAdminByEmail(email, db)
		return admin.Person.Password
	} else if userType == "trafficWarden" {
		trafficWarden := GetTrafficWardenByEmail(email, db)
		return trafficWarden.Person.Password
	} else {
		return "Tipo de usuário inválido"
	}

}

func UpdateUser(user database.User, db *gorm.DB) {
	db.Table("users").Where("id = ?", user.ID).Update("name", user.Person.Name)
	db.Table("users").Where("id = ?", user.ID).Update("email", user.Person.Email)
	db.Table("users").Where("id = ?", user.ID).Update("document", user.Document)
	db.Table("users").Where("id = ?", user.ID).Update("password", user.Person.Password)

}

func UpdateAdmin(admin database.Admin, db *gorm.DB) {
	db.Table("admins").Where("id = ?", admin.ID).Update("name", admin.Person.Name)
	db.Table("admins").Where("id = ?", admin.ID).Update("email", admin.Person.Email)
	db.Table("admins").Where("id = ?", admin.ID).Update("password", admin.Person.Password)
}

func UpdateTrafficWarden(trafficWarden database.TrafficWarden, db *gorm.DB) {
	db.Table("traffic_wardens").Where("id = ?", trafficWarden.ID).Update("name", trafficWarden.Person.Name)
	db.Table("traffic_wardens").Where("id = ?", trafficWarden.ID).Update("email", trafficWarden.Person.Email)
	db.Table("traffic_wardens").Where("id = ?", trafficWarden.ID).Update("password", trafficWarden.Person.Password)
}

func UpdateVehicle(vehicle database.Vehicle, db *gorm.DB) {
	db.Table("vehicles").Where("id = ?", vehicle.ID).Update("license_plate", vehicle.LicensePlate)
	db.Table("vehicles").Where("id = ?", vehicle.ID).Update("vehicle_model", vehicle.VehicleModel)
	db.Table("vehicles").Where("id = ?", vehicle.ID).Update("vehicle_type", vehicle.VehicleType)
}

func UpdateVehicleOwner(vehicleID, newOwnerID uint, db *gorm.DB) {
	db.Table("vehicles").Where("id = ?", vehicleID).Update("user_id", newOwnerID)
}

func UpdateBalance(email string, extra float64, db *gorm.DB) {
	balance := GetBalance(email, db)
	db.Table("users").Where("email = ?", email).Update("balance", balance+extra)
}

func UpdateEndTime(ticketID uint, db *gorm.DB) {
	currentTime := time.Now()
	db.Table("parking_tickets").Where("id = ?", ticketID).Update("end_time", currentTime.String())
}

func UpdateIsPaid(rechargeID uint, db *gorm.DB) {
	db.Table("recharges").Where("id = ?", rechargeID).Update("is_paid", true)
}

func UpdateIsParked(vehicleID uint, value bool, db *gorm.DB) {
	db.Table("vehicles").Where("id = ?", vehicleID).Update("is_parked", value)
}

func UpdateIsActive(vehicleID uint, value bool, db *gorm.DB) {
	db.Table("vehicles").Where("id = ?", vehicleID).Update("is_active", value)
}

func UpdateBilletLink(billetID uint, link string, db *gorm.DB) {
	db.Table("billets").Where("id = ?", billetID).Update("billet_link", link)
}

func DeleteUserByID(userID uint, db *gorm.DB) {
	db.Table("users").Where("id = ?", userID).Delete(&database.User{})
	DeleteVehiclesByUserID(userID, db)
	DeleteRechargeByUserID(userID, db)
}

func DeleteTrafficWardenByID(trafficWardenID uint, db *gorm.DB) error{
	err := db.Table("traffic_wardens").Where("id = ?", trafficWardenID).Delete(&database.TrafficWarden{}).Error
	return  err
}

func DeleteAdminByID(adminID uint, db *gorm.DB) error{
	err := db.Table("admins").Where("id = ?", adminID).Delete(&database.Admin{}).Error
	return err
}

func DeleteVehicleByID(vehicleID uint, db *gorm.DB) {
	db.Table("vehicles").Where("id = ?", vehicleID).Delete(&database.Vehicle{})
	DeleteParkingTicketByVehicleID(vehicleID, db)
}

func DeleteParkingTicketByID(parkingTicketID uint, db *gorm.DB) {
	db.Table("parking_tickets").Where("id = ?", parkingTicketID).Delete(&database.ParkingTicket{})
}

func DeleteParkingTicketByVehicleID(vehicleId uint, db *gorm.DB) {
	db.Table("parking_tickets").Where("vehicle_id = ?", vehicleId).Delete(&database.ParkingTicket{})
}

func DeleteRechargeByID(rechargeID uint, db *gorm.DB) {
	db.Table("recharges").Where("id = ?", rechargeID).Delete(&database.Recharge{})
	DeleteBilletByRechargeID(rechargeID, db)
}

func DeleteBilletByID(billetID uint, db *gorm.DB) error{
	err := db.Table("billets").Where("id = ?", billetID).Delete(&database.Billet{}).Error
	return err
}

func DeleteVehiclesByUserID(userID uint, db *gorm.DB) {
	db.Table("vehicles").Where("user_id = ?", userID).Delete(&database.Vehicle{})
}

func DeleteBilletByRechargeID(rechargeID uint, db *gorm.DB) {
	db.Table("billets").Where("recharge_id = ?", rechargeID).Delete(&database.Billet{})
}

func DeleteRechargeByUserID(userID uint, db *gorm.DB) {
	recharges, _ := GetRechargeByUserId(userID, db)
	db.Table("recharges").Where("user_id = ?", userID).Delete(&database.Recharge{})

	for i := range recharges {
		DeleteBilletByRechargeID(recharges[i].ID, db)
	}

}
