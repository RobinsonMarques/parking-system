package crud

import (
	_ "gorm.io/driver/postgres"
	_ "gorm.io/gorm"
)

//func NewCrud(db *gorm.DB) Crud{
//	userCrud := NewUserCrud(db)
//	adminCrud := NewAdminCrud(db)
//	trafficWardenCrud := NewTrafficWardenCrud(db)
//	vehicleCrud := NewVehicleCrud(db)
//	rechargeCrud := NewRechargeCrud(db)
//	parkingTicketCrud := NewParkingTicketCrud(db)
//	billetCrud := NewBilletCrud(db)
//	utilCrud := NewUtilCrud(db)
//	return Crud{
//		UserCrud: userCrud,
//		AdminCrud: adminCrud,
//		TrafficWardenCrud: trafficWardenCrud,
//		VehicleCrud: vehicleCrud,
//		RechargeCrud: rechargeCrud,
//		ParkingTicketCrud: parkingTicketCrud,
//		BilletCrud: billetCrud,
//		UtilCrud: utilCrud,
//	}
//}
//type Crud struct {
//	UserCrud UserCrud
//	AdminCrud AdminCrud
//	TrafficWardenCrud TrafficWardenCrud
//	VehicleCrud VehicleCrud
//	RechargeCrud RechargeCrud
//	ParkingTicketCrud ParkingTicketCrud
//	BilletCrud BilletCrud
//	UtilCrud UtilCrud
//}
