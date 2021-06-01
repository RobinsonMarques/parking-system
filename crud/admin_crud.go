package crud

import (
	"github.com/RobinsonMarques/parking-system/database"
	"gorm.io/gorm"
)

func NewAdminCrud(db *gorm.DB) AdminCrud {
	return AdminCrud{db: db}
}

type AdminCrud struct {
	db *gorm.DB
}

func (a AdminCrud) CreateAdmin(admin database.Admin) error {
	return a.db.Create(&admin).Error
}

func (a AdminCrud) GetAdminByEmail(email string) (database.Admin, error) {
	var admin database.Admin
	err := a.db.Where("Email = ?", email).First(&admin).Error
	return admin, err
}

func (a AdminCrud) GetAdminByID(id uint) (database.Admin, error) {
	var admin database.Admin
	err := a.db.Where("id = ?", id).First(&admin).Error
	return admin, err
}

func (a AdminCrud) UpdateAdmin(admin database.Admin) error {
	err := a.db.Table("admins").Where("id = ?", admin.ID).Update("name", admin.Person.Name).Error
	if err != nil {
		return err
	}
	err = a.db.Table("admins").Where("id = ?", admin.ID).Update("email", admin.Person.Email).Error
	if err != nil {
		return err
	}
	err = a.db.Table("admins").Where("id = ?", admin.ID).Update("password", admin.Person.Password).Error
	if err != nil {
		return err
	}
	return nil
}

func (a AdminCrud) DeleteAdminByID(adminID uint) error {
	err := a.db.Table("admins").Where("id = ?", adminID).Delete(&database.Admin{}).Error
	return err
}
