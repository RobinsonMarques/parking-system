package services

import (
	"errors"
	input2 "github.com/RobinsonMarques/parking-system/input"
	"github.com/RobinsonMarques/parking-system/interfaces"
)

func NewBilletService(billetInterface interfaces.BilletInterface, utilInterface interfaces.UtilInterface) BilletService {
	return BilletService{
		billetInterface: billetInterface,
		utilInterface:   utilInterface,
	}
}

type BilletService struct {
	billetInterface interfaces.BilletInterface
	utilInterface   interfaces.UtilInterface
}

func (b BilletService) DeleteBilletByID(input input2.LoginInput, billetID uint) error {
	resp := b.utilInterface.Login(input.Email, input.Password)
	if resp == "admin" {
		err := b.billetInterface.DeleteBilletByID(billetID)
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
