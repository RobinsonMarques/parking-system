package main

import (
	"fmt"
	"github.com/RobinsonMarques/parking-system/crud"
	"github.com/RobinsonMarques/parking-system/dependencies"
)

func main() {
	db := dependencies.CreateConnection()
	//db.AutoMigrate(&database.Person{})
	//db.AutoMigrate(&database.User{})
	//db.AutoMigrate(&database.TrafficWarden{})
	//db.AutoMigrate(&database.Admin{})
	//db.AutoMigrate(&database.ParkingTicket{})
	//db.AutoMigrate(&database.Vehicle{})
	//db.AutoMigrate(&database.Billet{})
	//db.AutoMigrate(&database.Recharge{})

	//Criar user
	//personUser := database.Person{
	//Model:    gorm.Model{},
	//Name:     "Jaaj",
	//Email:    "teste@teste.com",
	//Password: "123456",
	//}
	//user := database.User{
	//Model:    gorm.Model{},
	//Person:   personUser,
	//Document: "123.321.456-90",
	//Balance:  0,
	//Recharge: nil,
	//Vehicle:  nil,
	//}

	//crud.CreateUser(user, db)

	//Criar Admin
	//personAdmin := database.Person{
	//Model:    gorm.Model{},
	//Name:     "Admin",
	//Email:    "Admin@admin.com",
	//Password: "adm40028922",
	//}
	//admin := database.Admin{
	//Model:  gorm.Model{},
	//Person: personAdmin,
	//}

	//crud.CreateAdmin(admin, db)

	//Criar Guarda
	//personWarden := database.Person{
	//Model:    gorm.Model{},
	//Name:     "Guarda",
	//Email:    "guarda@guarda.com",
	//Password: "321456987",
	//}
	//warden := database.TrafficWarden{
	//Model:  gorm.Model{},
	//Person: personWarden,
	//}

	//crud.CreateTrafficWarden(warden, db)

	//Criar veiculo
	//veiculo := database.Vehicle{
	//Model:         gorm.Model{},
	//LicensePlate:  "abc-1234",
	//VehicleModel:  "Fusca",
	//VehicleType:   "carro",
	//IsActive:      false,
	//IsParked:      false,
	//UserID:        1,
	//ParkingTicket: nil,
	//}

	//crud.CreateVehicle(veiculo, db)

	//Criar ticket
	//currentTime := time.Now()
	//endTime := currentTime.Add(time.Hour)

	//ticket := database.ParkingTicket{
	//Model:       gorm.Model{},
	//Location:    "Ali ó",
	//ParkingTime: 1,
	//StartTime:   currentTime.String(),
	//EndTime:     endTime.String(),
	//Price:       1.99,
	//VehicleID:   1,
	//}

	//crud.CreateParkingTicket(ticket, db)

	//Criar recarga
	//currentTime = time.Now()

	//recarga := database.Recharge{
	//Model:       gorm.Model{},
	//Date:        currentTime.String(),
	//Value:       10,
	//IsPaid:      false,
	//PaymentType: "boleto",
	//UserID:      1,
	//Billet:      database.Billet{},
	//}

	//crud.CreateRecharge(recarga, db)

	//Criar boleto
	//boleto := database.Billet{
	//Model:      gorm.Model{},
	//BilletLink: "www.link.com",
	//RechargeID: 1,
	//}

	//crud.CreateBillet(boleto, db)

	//Consultar user por email
	fmt.Println(crud.GetUserByEmail("teste@teste.com", db))

	//Consultar guarda por email
	fmt.Println(crud.GetTrafficWardenByEmail("guarda@guarda.com", db))

	//Consultar admin por email
	fmt.Println(crud.GetAdminByEmail("Admin@admin.com", db))

	//Consultar todos os veículos
	fmt.Println(crud.GetAllVehicles(db))

	//Consultar veículos pelo id do usuário
	fmt.Println(crud.GetVehiclesByUserId(1, db))

	//Consultar user por documento
	fmt.Println(crud.GetUserByDocument("123.321.456-90", db))

	//Consultar veículo pela placa
	fmt.Println(crud.GetVehicleByLicensePlate("abc-1234", db))

	//Consultar veículo pelo id
	fmt.Println(crud.GetVehicleById(2, db))

	//Consultar ultimo ticket do veiculo
	fmt.Println(crud.GetLastParkingTicketFromVehicle(1, db))

	//Update no user
	crud.UpdateUser(1, "José", "josemaria@silva.com", "123.456-08", db)

}
