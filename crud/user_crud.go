package crud

import (
	"github.com/RobinsonMarques/parking-system/database"
	"gorm.io/gorm"
)

func NewUserCrud(db *gorm.DB) UserCrud {
	return UserCrud{db: db}
}

type UserCrud struct {
	db *gorm.DB
}

func (u UserCrud) CreateUser(user database.User) error {
	return u.db.Create(&user).Error
}

func (u UserCrud) GetUserByEmail(email string) (database.User, error) {
	var user database.User

	err := u.db.Where("email = ?", email).First(&user).Error
	return user, err
}

func (u UserCrud) GetUserByID(id uint) (database.User, error) {
	var user database.User
	err := u.db.Where("id = ?", id).First(&user).Error
	return user, err
}

func (u UserCrud) GetUserByDocument(document string) (database.User, error) {
	var user database.User

	err := u.db.Where("Document = ?", document).First(&user).Error
	return user, err
}

func (u UserCrud) GetBalance(email string, userCrud UserCrud) (float64, error) {
	user, err := userCrud.GetUserByEmail(email)
	balance := user

	return balance.Balance, err
}

func (u UserCrud) UpdateUser(user database.User) error {
	err := u.db.Table("users").Where("id = ?", user.ID).Update("name", user.Person.Name).Error
	if err != nil {
		return err
	}
	err = u.db.Table("users").Where("id = ?", user.ID).Update("email", user.Person.Email).Error
	if err != nil {
		return err
	}
	err = u.db.Table("users").Where("id = ?", user.ID).Update("document", user.Document).Error
	if err != nil {
		return err
	}
	err = u.db.Table("users").Where("id = ?", user.ID).Update("password", user.Person.Password).Error
	if err != nil {
		return err
	}
	return nil

}

func (u UserCrud) UpdateBalance(email string, extra float64) error {
	userCrud := NewUserCrud(u.db)
	balance, err := userCrud.GetBalance(email, userCrud)
	if err != nil {
		return err
	}
	err = u.db.Table("users").Where("email = ?", email).Update("balance", balance+extra).Error
	return err
}

func (u UserCrud) DeleteUserByID(userID uint) error {
	err := u.db.Table("users").Where("id = ?", userID).Delete(&database.User{}).Error
	if err != nil {
		return err
	}
	return nil
}
