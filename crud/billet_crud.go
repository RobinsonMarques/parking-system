package crud

import (
	"github.com/RobinsonMarques/parking-system/database"
	"gorm.io/gorm"
)

func NewBilletCrud(db *gorm.DB) BilletCrud {
	return BilletCrud{db: db}
}

type BilletCrud struct {
	db *gorm.DB
}

func (b BilletCrud) CreateBillet(billet database.Billet) error {
	return b.db.Create(&billet).Error
}

func (b BilletCrud) GetBilletByRechargeId(rechargeID uint) (database.Billet, error) {
	var billet database.Billet
	err := b.db.Where("recharge_id = ?", rechargeID).Find(&billet).Error
	return billet, err
}

func (b BilletCrud) GetBilletsByRechargeID(rechargeID uint) ([]database.Billet, error) {
	var billetts []database.Billet

	err := b.db.Where("recharge_id = ?", rechargeID).Find(&billetts).Error
	return billetts, err
}

func (b BilletCrud) UpdateBilletLink(billetID uint, link string) error {
	err := b.db.Table("billets").Where("id = ?", billetID).Update("billet_link", link).Error
	return err
}

func (b BilletCrud) DeleteBilletByID(billetID uint) error {
	err := b.db.Table("billets").Where("id = ?", billetID).Delete(&database.Billet{}).Error
	return err
}

func (b BilletCrud) DeleteBilletByRechargeID(rechargeID uint) error {
	err := b.db.Table("billets").Where("recharge_id = ?", rechargeID).Delete(&database.Billet{}).Error
	return err
}
