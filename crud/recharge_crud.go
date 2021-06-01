package crud

import (
	"github.com/RobinsonMarques/parking-system/database"
	"gorm.io/gorm"
)

func NewRechargeCrud(db *gorm.DB) RechargeCrud {
	return RechargeCrud{db: db}
}

type RechargeCrud struct {
	db *gorm.DB
}

func (r RechargeCrud) CreateRecharge(recharge database.Recharge) error {
	return r.db.Create(&recharge).Error
}

func (r RechargeCrud) GetUserUnpaidRechargesByID(userID uint) ([]database.Recharge, error) {
	var recharges []database.Recharge
	err := r.db.Where("is_paid = false AND user_id = ?", userID).Find(&recharges).Error
	return recharges, err
}

func (r RechargeCrud) GetRechargeByUserId(userID uint) ([]database.Recharge, error) {
	var recharges []database.Recharge
	err := r.db.Where("user_id = ?", userID).Find(&recharges).Error
	return recharges, err
}

func (r RechargeCrud) UpdateIsPaid(rechargeID uint) error {
	err := r.db.Table("recharges").Where("id = ?", rechargeID).Update("is_paid", true).Error
	return err
}

func (r RechargeCrud) DeleteRechargeByID(rechargeID uint) error {
	err := r.db.Table("recharges").Where("id = ?", rechargeID).Delete(&database.Recharge{}).Error
	return err
}

func (r RechargeCrud) DeleteRechargeByUserID(userID uint) error {
	err := r.db.Table("recharges").Where("user_id = ?", userID).Delete(&database.Recharge{}).Error
	if err != nil {
		return err
	}
	return nil
}
