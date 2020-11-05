package migration0

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/siovanus/wingServer/http/common"
)

type FlashPoolDetail struct {
	Timestamp      uint64 `gorm:"primary_key"`
	TotalSupply    string
	TotalBorrow    string
	TotalInsurance string
}

type FlashPoolMarket struct {
	ID             uint64
	Name           string
	Timestamp      uint64
	TotalSupply    string
	TotalBorrow    string
	TotalInsurance string
}

type Price struct {
	Name  string `gorm:"primary_key"`
	Price string
}

type TrackHeight struct {
	Name   string `gorm:"primary_key"`
	Height uint32
}

type UserAssetBalance struct {
	UserAddress  string `gorm:"primary_key"`
	AssetName    string `gorm:"primary_key"`
	AssetAddress string
	Icon         string
	FToken       string
	BorrowAmount string
	BorrowIndex  string
	Itoken       string
	IfCollateral bool
}

type WingApy struct {
	AssetName    string `gorm:"primary_key"`
	SupplyApy    string
	BorrowApy    string
	InsuranceApy string
}

type IFInfo struct {
	Name  string `gorm:"primary_key"`
	Total string
	Cap   string
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

	err = tx.AutoMigrate(UserAssetBalance{}).Error
	if err != nil {
		return errors.Wrap(err, "failed to auto migrate UserAssetBalance")
	}

	err = tx.AutoMigrate(common.Market{}).Error
	if err != nil {
		return errors.Wrap(err, "failed to auto migrate Market")
	}

	err = tx.AutoMigrate(WingApy{}).Error
	if err != nil {
		return errors.Wrap(err, "failed to auto migrate WingApy")
	}

	err = tx.AutoMigrate(IFInfo{}).Error
	if err != nil {
		return errors.Wrap(err, "failed to auto migrate IFInfo")
	}

	return nil
}
