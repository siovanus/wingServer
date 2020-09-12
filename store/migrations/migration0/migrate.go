package migration0

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/siovanus/wingServer/http/common"
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

type Price struct {
	Name  string `gorm:"primary_key"`
	Price uint64
}

type TrackHeight struct {
	Name   string `gorm:"primary_key"`
	Height uint32
}

type UserFlashPoolOverview struct {
	UserAddress      string `gorm:"primary_key"`
	SupplyBalance    uint64
	BorrowBalance    uint64
	InsuranceBalance uint64
	BorrowLimit      uint64
	NetApy           int64
	Info             string
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

	err = tx.AutoMigrate(Price{}).Error
	if err != nil {
		return errors.Wrap(err, "failed to auto migrate Price")
	}

	err = tx.AutoMigrate(TrackHeight{}).Error
	if err != nil {
		return errors.Wrap(err, "failed to auto migrate TrackHeight")
	}

	err = tx.AutoMigrate(UserFlashPoolOverview{}).Error
	if err != nil {
		return errors.Wrap(err, "failed to auto migrate UserFlashPoolOverview")
	}

	err = tx.AutoMigrate(common.Market{}).Error
	if err != nil {
		return errors.Wrap(err, "failed to auto migrate UserFlashPoolOverview")
	}

	return nil
}
