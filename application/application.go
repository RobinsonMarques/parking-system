package application

import (
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewApplication(db *gorm.DB) (Application, error) {
	parkingTicketCrud := crud.NewParkingTicketCrud(db)
	vehicleCrud := crud.NewVehicleCrud(db)
	userCrud := crud.NewUserCrud(db)
	billetCrud := crud.NewBilletCrud(db)
	rechargeCrud := crud.NewRechargeCrud(db)
	trafficWardenCrud := crud.NewTrafficWardenCrud(db)
	adminCrud := crud.NewAdminCrud(db)
	utilCrud := services.NewUtilCrud(userCrud, adminCrud, trafficWardenCrud)
	adminService := services.NewAdminService(adminCrud, utilCrud)
	billetService := services.NewBilletService(billetCrud, utilCrud)
	userService := services.NewUserService(userCrud, vehicleCrud, rechargeCrud, billetCrud, utilCrud)
	trafficWardenService := services.NewTrafficWardenService(trafficWardenCrud, utilCrud)
	vehicleService := services.NewVehicleService(vehicleCrud, userCrud, parkingTicketCrud, utilCrud)
	parkingTicketService := services.NewParkingTicketService(parkingTicketCrud, vehicleCrud, userCrud, utilCrud)
	rechargeService, err := services.NewRechargeService(rechargeCrud, userCrud, utilCrud, billetCrud)
	if err != nil {
		return Application{}, err
	}
	UserManager := NewUserController(userService)
	AdminManager := NewAdminController(adminService)
	TrafficWardenManager := NewTrafficWardenManager(trafficWardenService)
	VehicleManager := NewVehicleManager(vehicleService)
	ParkingTicketManager := NewParkingTicketManager(parkingTicketService)
	RechargeManager, err := NewRechargeController(rechargeService)
	if err != nil {
		return Application{}, err
	}
	BilletManager := NewBilletManager(billetService)
	return Application{
		UserManager:          UserManager,
		AdminManager:         AdminManager,
		TrafficWardenManager: TrafficWardenManager,
		VehicleManager:       VehicleManager,
		ParkingTicketManager: ParkingTicketManager,
		RechargeManager:      RechargeManager,
		BilletManager:        BilletManager,
	}, nil
}

type Application struct {
	UserManager          UserController
	AdminManager         AdminController
	TrafficWardenManager TrafficWardenManager
	VehicleManager       VehicleManager
	ParkingTicketManager ParkingTicketManager
	RechargeManager      RechargeController
	BilletManager        BilletManager
}

func (a Application) Run() error {

	r := gin.Default()

	r.GET("/vehicles", a.VehicleManager.GetAllVehicles)
	r.POST("/users", a.UserManager.CreateUser)
	r.POST("/admins", a.AdminManager.CreateAdmin)
	r.POST("/trafficwarden", a.TrafficWardenManager.CreateTrafficWarden)
	r.POST("/vehicles", a.VehicleManager.CreateVehicle)
	r.POST("/parkingtickets", a.ParkingTicketManager.CreateParkingTicket)
	r.POST("/recharge", a.RechargeManager.CreateRecharge)
	r.GET("/vehicles/:licensePlate", a.VehicleManager.GetVehicleByLicensePlate)
	r.GET("/recharges", a.RechargeManager.GetRechargesStatus)
	r.GET("users/:document", a.UserManager.GetUserByDocument)
	r.PUT("/users/:userID", a.UserManager.UpdateUser)
	r.PUT("/admins/:adminID", a.AdminManager.UpdateAdmin)
	r.PUT("/trafficwarden/:wardenID", a.TrafficWardenManager.UpdateTrafficWarden)
	r.PUT("/vehicles", a.VehicleManager.UpdateVehicle)
	r.PUT("/vehicles/updateowner/:vehicleID", a.VehicleManager.UpdateVehicleOwner)
	r.DELETE("/users/:userID", a.UserManager.DeleteUserByID)
	r.DELETE("/trafficwarden/:trafficwardenID", a.TrafficWardenManager.DeleteTrafficWardenByID)
	r.DELETE("/admins/:adminID", a.AdminManager.DeleteAdminByID)
	r.DELETE("/parkingtickets/:ticketID", a.ParkingTicketManager.DeleteParkingTicketByID)
	r.DELETE("/recharge/:rechargeID", a.RechargeManager.DeleteRechargeByID)
	r.DELETE("/billets/:billetID", a.BilletManager.DeleteBilletByID)
	r.DELETE("/vehicles/:vehicleID", a.VehicleManager.DeleteVehicleByID)

	return r.Run() // listen and serve on 0.0.0.0:8080
}
