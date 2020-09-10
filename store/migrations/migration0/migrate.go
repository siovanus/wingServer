package migration0

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type FlashPoolDetail struct {
	Timestamp      uint64 `gorm:"primary_key"`
	TotalSupply    uint64
	TotalBorrow    uint64
	TotalInsurance uint64
}

type FlashPoolMarket struct {
	ID             uint64
	Name           string
	Timestamp      uint64
	TotalSupply    uint64
	TotalBorrow    uint64
	TotalInsurance uint64
}

// Migrate runs the initial migration
func Migrate(tx *gorm.DB) error {
	err := tx.AutoMigrate(&FlashPoolDetail{}).Error
	if err != nil {
		return errors.Wrap(err, "failed to auto migrate FlashPoolDetail")
	}

	err = tx.AutoMigrate(&FlashPoolMarket{}).Error
	if err != nil {
		return errors.Wrap(err, "failed to auto migrate FlashPoolMarket")
	}

	return nil
}
