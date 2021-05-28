package crud

import (
	"errors"
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

func GetUserByEmail(email string, db *gorm.DB) (database.User, error) {
	var user database.User

	err := db.Where("email = ?", email).First(&user).Error
	return user, err
}

func GetUserUnpaidRechargesByID(userID uint, db *gorm.DB) ([]database.Recharge, error) {
	var recharges []database.Recharge
	err := db.Where("is_paid = false AND user_id = ?", userID).Find(&recharges).Error
	return recharges, err
}

func GetUserByID(id uint, db *gorm.DB) (database.User, error) {
	var user database.User
	err := db.Where("id = ?", id).First(&user).Error
	return user, err
}

func GetTrafficWardenByEmail(email string, db *gorm.DB) (database.TrafficWarden, error) {
	var trafficWarden database.TrafficWarden

	err := db.Where("Email = ?", email).First(&trafficWarden).Error
	return trafficWarden, err
}

func GetTrafficWardenByID(id uint, db *gorm.DB) (database.TrafficWarden, error) {
	var warden database.TrafficWarden
	err := db.Where("id = ?", id).First(&warden).Error
	return warden, err
}

func GetAdminByEmail(email string, db *gorm.DB) (database.Admin, error) {
	var admin database.Admin
	err := db.Where("Email = ?", email).First(&admin).Error
	return admin, err
}

func GetAdminByID(id uint, db *gorm.DB) (database.Admin, error) {
	var admin database.Admin
	err := db.Where("id = ?", id).First(&admin).Error
	return admin, err
}

func GetAllVehicles(db *gorm.DB) ([]database.Vehicle, error) {
	var vehicle []database.Vehicle

	err := db.Find(&vehicle).Error
	return vehicle, err
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

func GetVehicleById(id uint, db *gorm.DB) (database.Vehicle, error) {
	var vehicle database.Vehicle

	err := db.Where("id = ?", id).First(&vehicle).Error
	return vehicle, err
}

func GetLastParkingTicketFromVehicle(id uint, db *gorm.DB) ([]database.ParkingTicket, error) {
	var tickets []database.ParkingTicket

	err := db.Where("vehicle_id = ?", id).Last(&tickets).Error
	return tickets, err
}

func GetBilletsByRechargeID(rechargeID uint, db *gorm.DB) ([]database.Billet, error) {
	var billetts []database.Billet

	err := db.Where("recharge_id = ?", rechargeID).Find(&billetts).Error
	return billetts, err
}

func GetBalance(email string, db *gorm.DB) (float64, error) {
	user, err := GetUserByEmail(email, db)
	balance := user

	return balance.Balance, err
}

func GetPassword(email string, userType string, db *gorm.DB) (string, error) {
	if userType == "user" {
		user, err := GetUserByEmail(email, db)
		return user.Person.Password, err
	} else if userType == "admin" {
		admin, err := GetAdminByEmail(email, db)
		return admin.Person.Password, err
	} else if userType == "trafficWarden" {
		trafficWarden, err := GetTrafficWardenByEmail(email, db)
		return trafficWarden.Person.Password, err
	} else {
		err := errors.New("tipo de usuário inválido")
		return "", err
	}

}

func UpdateUser(user database.User, db *gorm.DB) error {
	err := db.Table("users").Where("id = ?", user.ID).Update("name", user.Person.Name).Error
	if err != nil {
		return err
	}
	err = db.Table("users").Where("id = ?", user.ID).Update("email", user.Person.Email).Error
	if err != nil {
		return err
	}
	err = db.Table("users").Where("id = ?", user.ID).Update("document", user.Document).Error
	if err != nil {
		return err
	}
	err = db.Table("users").Where("id = ?", user.ID).Update("password", user.Person.Password).Error
	if err != nil {
		return err
	}
	return nil

}

func UpdateAdmin(admin database.Admin, db *gorm.DB) error {
	err := db.Table("admins").Where("id = ?", admin.ID).Update("name", admin.Person.Name).Error
	if err != nil {
		return err
	}
	err = db.Table("admins").Where("id = ?", admin.ID).Update("email", admin.Person.Email).Error
	if err != nil {
		return err
	}
	err = db.Table("admins").Where("id = ?", admin.ID).Update("password", admin.Person.Password).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateTrafficWarden(trafficWarden database.TrafficWarden, db *gorm.DB) error {
	err := db.Table("traffic_wardens").Where("id = ?", trafficWarden.ID).Update("name", trafficWarden.Person.Name).Error
	if err != nil {
		return err
	}
	err = db.Table("traffic_wardens").Where("id = ?", trafficWarden.ID).Update("email", trafficWarden.Person.Email).Error
	if err != nil {
		return err
	}
	err = db.Table("traffic_wardens").Where("id = ?", trafficWarden.ID).Update("password", trafficWarden.Person.Password).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateVehicle(vehicle database.Vehicle, db *gorm.DB) error {
	err := db.Table("vehicles").Where("id = ?", vehicle.ID).Update("license_plate", vehicle.LicensePlate).Error
	if err != nil {
		return err
	}
	err = db.Table("vehicles").Where("id = ?", vehicle.ID).Update("vehicle_model", vehicle.VehicleModel).Error
	if err != nil {
		return err
	}
	err = db.Table("vehicles").Where("id = ?", vehicle.ID).Update("vehicle_type", vehicle.VehicleType).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateVehicleOwner(vehicleID, newOwnerID uint, db *gorm.DB) error {
	err := db.Table("vehicles").Where("id = ?", vehicleID).Update("user_id", newOwnerID).Error
	return err
}

func UpdateBalance(email string, extra float64, db *gorm.DB) error {
	balance, err := GetBalance(email, db)
	if err != nil {
		return err
	}
	err = db.Table("users").Where("email = ?", email).Update("balance", balance+extra).Error
	return err
}

func UpdateEndTime(ticketID uint, db *gorm.DB) error {
	currentTime := time.Now()
	err := db.Table("parking_tickets").Where("id = ?", ticketID).Update("end_time", currentTime.String()).Error
	return err
}

func UpdateIsPaid(rechargeID uint, db *gorm.DB) error {
	err := db.Table("recharges").Where("id = ?", rechargeID).Update("is_paid", true).Error
	return err
}

func UpdateIsParked(vehicleID uint, value bool, db *gorm.DB) error {
	err := db.Table("vehicles").Where("id = ?", vehicleID).Update("is_parked", value).Error
	return err
}

func UpdateIsActive(vehicleID uint, value bool, db *gorm.DB) error {
	err := db.Table("vehicles").Where("id = ?", vehicleID).Update("is_active", value).Error
	return err
}

func UpdateBilletLink(billetID uint, link string, db *gorm.DB) error {
	err := db.Table("billets").Where("id = ?", billetID).Update("billet_link", link).Error
	return err
}

func DeleteUserByID(userID uint, db *gorm.DB) error {
	err := db.Table("users").Where("id = ?", userID).Delete(&database.User{}).Error
	if err != nil {
		return err
	}
	err = DeleteVehiclesByUserID(userID, db)
	if err != nil {
		return err
	}
	err = DeleteRechargeByUserID(userID, db)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTrafficWardenByID(trafficWardenID uint, db *gorm.DB) error {
	err := db.Table("traffic_wardens").Where("id = ?", trafficWardenID).Delete(&database.TrafficWarden{}).Error
	return err
}

func DeleteAdminByID(adminID uint, db *gorm.DB) error {
	err := db.Table("admins").Where("id = ?", adminID).Delete(&database.Admin{}).Error
	return err
}

func DeleteVehicleByID(vehicleID uint, db *gorm.DB) error {
	err := db.Table("vehicles").Where("id = ?", vehicleID).Delete(&database.Vehicle{}).Error
	DeleteParkingTicketByVehicleID(vehicleID, db)
	return err
}

func DeleteParkingTicketByID(parkingTicketID uint, db *gorm.DB) error {
	err := db.Table("parking_tickets").Where("id = ?", parkingTicketID).Delete(&database.ParkingTicket{}).Error
	return err
}

func DeleteParkingTicketByVehicleID(vehicleId uint, db *gorm.DB) {
	db.Table("parking_tickets").Where("vehicle_id = ?", vehicleId).Delete(&database.ParkingTicket{})
}

func DeleteRechargeByID(rechargeID uint, db *gorm.DB) error {
	err := db.Table("recharges").Where("id = ?", rechargeID).Delete(&database.Recharge{}).Error
	if err != nil {
		return err
	}
	err = DeleteBilletByRechargeID(rechargeID, db)
	return err
}

func DeleteBilletByID(billetID uint, db *gorm.DB) error {
	err := db.Table("billets").Where("id = ?", billetID).Delete(&database.Billet{}).Error
	return err
}

func DeleteVehiclesByUserID(userID uint, db *gorm.DB) error {
	err := db.Table("vehicles").Where("user_id = ?", userID).Delete(&database.Vehicle{}).Error
	return err
}

func DeleteBilletByRechargeID(rechargeID uint, db *gorm.DB) error {
	err := db.Table("billets").Where("recharge_id = ?", rechargeID).Delete(&database.Billet{}).Error
	return err
}

func DeleteRechargeByUserID(userID uint, db *gorm.DB) error {
	recharges, err := GetRechargeByUserId(userID, db)
	if err != nil {
		return err
	}
	err = db.Table("recharges").Where("user_id = ?", userID).Delete(&database.Recharge{}).Error
	if err != nil {
		return err
	}

	for i := range recharges {
		err := DeleteBilletByRechargeID(recharges[i].ID, db)
		if err != nil {
			return err
		}
	}
	return nil
}
