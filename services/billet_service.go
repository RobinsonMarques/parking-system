package services

import (
	"errors"
	"github.com/RobinsonMarques/parking-system/crud"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/utils"
	"gorm.io/gorm"
)

func NewBilletService(db *gorm.DB) BilletService {
	return BilletService{db: db}
}

type BilletService struct {
	db *gorm.DB
}

func (b BilletService) DeleteBilletByID(input input2.LoginInput, billetID uint) error {
	resp := utils.Login(input.Email, input.Password, b.db)
	billetCrud := crud.NewBilletCrud(b.db)
	if resp == "admin" {
		err := billetCrud.DeleteBilletByID(billetID)
		if err == nil {
			return nil
		} else {
			return err
		}
	} else {
		err := errors.New(resp)
		return err
	}
}
