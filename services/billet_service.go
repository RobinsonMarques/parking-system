package services

import (
	"errors"
	"github.com/RobinsonMarques/parking-system/crud"
	input2 "github.com/RobinsonMarques/parking-system/input"
)

func NewBilletService(billetCrud crud.BilletCrud, utilCrud crud.UtilCrud) BilletService {
	return BilletService{
		billetCrud: billetCrud,
		utilCrud:   utilCrud,
	}
}

type BilletService struct {
	billetCrud crud.BilletCrud
	utilCrud   crud.UtilCrud
}

func (b BilletService) DeleteBilletByID(input input2.LoginInput, billetID uint, service BilletService) error {
	resp := service.utilCrud.Login(input.Email, input.Password)
	if resp == "admin" {
		err := service.billetCrud.DeleteBilletByID(billetID)
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
