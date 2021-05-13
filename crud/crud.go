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

func GetUserByEmail(email string, db *gorm.DB) []database.User {
	var user []database.User

	db.Where("Email = ?", email).First(&user)
	return user
}

func GetTrafficWardenByEmail(email string, db *gorm.DB) []database.TrafficWarden {
	var trafficWarden []database.TrafficWarden

	db.Where("Email = ?", email).First(&trafficWarden)
	return trafficWarden
}

func GetAdminByEmail(email string, db *gorm.DB) []database.Admin {
	var admin []database.Admin
	db.Where("Email = ?", email).First(&admin)
	return admin
}

func GetAllVehicles(db *gorm.DB) []database.Vehicle {
	var vehicle []database.Vehicle

	db.Find(&vehicle)
	return vehicle
}

func GetVehiclesByUserId(user database.User, db *gorm.DB) Result {
	var vehicles []database.Vehicle
	data := db.Where("VehicleID <> ?", user.Vehicle).Find(&vehicles)
	return Result{Data: data}
}

func GetUserByDocument(document string, db *gorm.DB) Result {
	var user []database.User

	data := db.Where("Document = ?", document).First(&user)
	return Result{Data: data}
}

func GetVehicleByLicensePlate(licensePlate string, db *gorm.DB) Result {
	var vehicle []database.Vehicle

	data := db.Where("LicensePlate = ?", licensePlate).First(&vehicle)
	return Result{Data: data}
}

func GetVehicleById(id uint, db *gorm.DB) Result {
	var vehicle []database.Vehicle

	data := db.Where("VehicleID = ?", id).First(&vehicle)
	return Result{Data: data}
}

func GetLastParkingTicketFromVehicle(vehicle database.Vehicle, db *gorm.DB) Result {
	var tickets []database.ParkingTicket

	data := db.Where("ParkingTicketID <> ?", vehicle.ParkingTicket).Last(&tickets)
	return Result{Data: data}
}

func UpdateUser(id uint, name string, email string, document string, db *gorm.DB) {
	db.Table("users").Where("id = ?", id).Update("name", name)
	db.Table("users").Where("id = ?", id).Update("email", email)
	db.Table("users").Where("id = ?", id).Update("document", document)

}

func UpdateAdmin(id uint, name string, email string, db *gorm.DB) {
	//admin := database.Admin{AdminId: id}
	//db.Model(&admin).Update("Name", name)
	//db.Model(&admin).Update("Email", email)
}

func UpdateTrafficWarden(id uint, name string, email string, db *gorm.DB) {
	//trafficWarden := database.TrafficWarden{TrafficWardenID: id}
	//db.Model(&trafficWarden).Update("Name", name)
	//db.Model(&trafficWarden).Update("Email", email)
}

func UpdateVehicle(vehicleID uint, licensePlate string, vehicleModel string, vehicleType string, db *gorm.DB) {
	//vehicle := database.Vehicle{VehicleID: vehicleID}
	//db.Model(&vehicle).Update("LicensePlate", licensePlate)
	//db.Model(&vehicle).Update("VehicleModel", vehicleModel)
	//db.Model(&vehicle).Update("VehicleType", vehicleType)
}

func UpdateVehicleOwner(userID uint, db *gorm.DB) {
	db.Table("Vehicle").Where("UserId = ?", userID).Update("UserID", userID)
}

func UpdateBalance(userID uint, extra string, db *gorm.DB) {
	db.Table("User").Where("UserID = ?", userID).Update("Balance", gorm.Expr("Balance +", extra))
}

func UpdateEndTime(ticketID uint, db *gorm.DB) {
	currentTime := time.Now()
	db.Table("ParkingTicket").Where("ParkingTicketID = ?", ticketID).Update("EndTime", currentTime.String())
}

func UpdateIsPaid(rechargeID uint, db *gorm.DB) {
	db.Table("Recharge").Where("RechargeID = ?", rechargeID).Update("IsPaid", true)
}

func UpdateBilletLink(billetID uint, link string, db *gorm.DB) {
	db.Table("Billet").Where("BilletID = ?", billetID).Update("BilletLink", link)
}

func DeleteUserByID(userID uint, db *gorm.DB) {
	db.Table("User").Where("UserID = ?", userID).Delete(&database.User{})
}

func DeleteTrafficWardenByID(trafficWardenID uint, db *gorm.DB) {
	db.Table("TrafficWarden").Where("TrafficWardenID = ?", trafficWardenID).Delete(&database.TrafficWarden{})
}

func DeleteAdminByID(adminID uint, db *gorm.DB) {
	db.Table("Admin").Where("AdminID = ?", adminID).Delete(&database.Admin{})
}

func DeleteVehicleByID(vehicleID uint, db *gorm.DB) {
	db.Table("Vehicle").Where("VehicleID = ?", vehicleID).Delete(&database.Vehicle{})
}

func DeleteParkingTicketByID(parkingTicketID uint, db *gorm.DB) {
	db.Table("ParkingTicket").Where("ParkingTicketID = ?", parkingTicketID).Delete(&database.ParkingTicket{})
}

func DeleteRechargeByID(rechargeID uint, db *gorm.DB) {
	db.Table("Recharge").Where("RechargeID = ?", rechargeID).Delete(&database.Recharge{})
}

func DeleteBilletByID(billetID uint, db *gorm.DB) {
	db.Table("Billet").Where("BilletID = ?", billetID).Delete(&database.Billet{})
}
