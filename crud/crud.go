package crud

import (
	"github.com/RobinsonMarques/parking-system/database"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
	"time"
)

type Result struct {
	Data interface{}
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

func GetUserByID(id uint, db *gorm.DB) database.User {
	var user database.User
	db.Where("id = ?", id).First(&user)
	return user
}

func GetTrafficWardenByEmail(email string, db *gorm.DB) database.TrafficWarden {
	var trafficWarden database.TrafficWarden

	db.Where("Email = ?", email).First(&trafficWarden)
	return trafficWarden
}

func GetAdminByEmail(email string, db *gorm.DB) database.Admin {
	var admin database.Admin
	db.Where("Email = ?", email).First(&admin)
	return admin
}

func GetAllVehicles(db *gorm.DB) []database.Vehicle {
	var vehicle []database.Vehicle

	db.Find(&vehicle)
	return vehicle
}

func GetVehiclesByUserId(userID uint, db *gorm.DB) []database.Vehicle {
	var vehicles []database.Vehicle
	db.Where("user_id = ?", userID).Find(&vehicles)
	return vehicles
}

func GetUserByDocument(document string, db *gorm.DB) database.User {
	var user database.User

	db.Where("Document = ?", document).First(&user)
	return user
}

func GetVehicleByLicensePlate(licensePlate string, db *gorm.DB) database.Vehicle {
	var vehicle database.Vehicle

	db.Where("license_plate = ?", licensePlate).First(&vehicle)
	return vehicle
}

func GetVehicleById(id uint, db *gorm.DB) database.Vehicle {
	var vehicle database.Vehicle

	db.Where("id = ?", id).First(&vehicle)
	return vehicle
}

func GetLastParkingTicketFromVehicle(id uint, db *gorm.DB) []database.ParkingTicket {
	var tickets []database.ParkingTicket

	db.Where("vehicle_id = ?", id).Last(&tickets)
	return tickets
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

func UpdateUser(id uint, name string, email string, document string, db *gorm.DB) {
	db.Table("users").Where("id = ?", id).Update("name", name)
	db.Table("users").Where("id = ?", id).Update("email", email)
	db.Table("users").Where("id = ?", id).Update("document", document)

}

func UpdateAdmin(id uint, name string, email string, db *gorm.DB) {
	db.Table("admins").Where("id = ?", id).Update("name", name)
	db.Table("admins").Where("id = ?", id).Update("email", email)
}

func UpdateTrafficWarden(id uint, name string, email string, db *gorm.DB) {
	db.Table("traffic_wardens").Where("id = ?", id).Update("name", name)
	db.Table("traffic_wardens").Where("id = ?", id).Update("Email", email)
}

func UpdateVehicle(id uint, licensePlate string, vehicleModel string, vehicleType string, db *gorm.DB) {
	db.Table("vehicles").Where("id = ?", id).Update("license_plate", licensePlate)
	db.Table("vehicles").Where("id = ?", id).Update("vehicle_model", vehicleModel)
	db.Table("vehicles").Where("id = ?", id).Update("vehicle_type", vehicleType)
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

func DeleteTrafficWardenByID(trafficWardenID uint, db *gorm.DB) {
	db.Table("traffic_wardens").Where("id = ?", trafficWardenID).Delete(&database.TrafficWarden{})
}

func DeleteAdminByID(adminID uint, db *gorm.DB) {
	db.Table("admins").Where("id = ?", adminID).Delete(&database.Admin{})
}

func DeleteVehicleByID(vehicleID uint, db *gorm.DB) {
	db.Table("vehicles").Where("id = ?", vehicleID).Delete(&database.Vehicle{})
}

func DeleteParkingTicketByID(parkingTicketID uint, db *gorm.DB) {
	db.Table("parking_tickets").Where("id = ?", parkingTicketID).Delete(&database.ParkingTicket{})
}

func DeleteRechargeByID(rechargeID uint, db *gorm.DB) {
	db.Table("recharges").Where("id = ?", rechargeID).Delete(&database.Recharge{})
	DeleteBilletByRechargeID(rechargeID, db)
}

func DeleteBilletByID(billetID uint, db *gorm.DB) {
	db.Table("billets").Where("id = ?", billetID).Delete(&database.Billet{})
}

func DeleteVehiclesByUserID(userID uint, db *gorm.DB) {
	db.Table("vehicles").Where("user_id = ?", userID).Delete(&database.Vehicle{})
}

func DeleteBilletByRechargeID(rechargeID uint, db *gorm.DB) {
	db.Table("billets").Where("recharge_id = ?", rechargeID).Delete(&database.Billet{})
}

func DeleteRechargeByUserID(userID uint, db *gorm.DB) {
	db.Table("recharges").Where("user_id = ?", userID).Delete(&database.Recharge{})
}
