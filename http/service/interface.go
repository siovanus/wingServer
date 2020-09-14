package service

import (
	ocommon "github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/http/common"
	"github.com/siovanus/wingServer/store"
)

type GovernanceManager interface {
	GovBannerOverview() (*common.GovBannerOverview, error)
	GovBanner() (*common.GovBanner, error)
}

type FlashPoolManager interface {
	AssetPrice(asset string) (string, error)
	FlashPoolMarketDistribution() (*common.FlashPoolMarketDistribution, error)
	PoolDistribution() (*common.Distribution, error)
	FlashPoolBanner() (*common.FlashPoolBanner, error)
	FlashPoolDetail() (*common.FlashPoolDetail, error)
	FlashPoolDetailForStore() (*store.FlashPoolDetail, error)
	FlashPoolMarketStore() error
	FlashPoolAllMarket() (*common.FlashPoolAllMarket, error)
	FlashPoolAllMarketForStore() (*common.FlashPoolAllMarket, error)
	UserFlashPoolOverview(account string) (*common.UserFlashPoolOverview, error)
	UserFlashPoolOverviewForStore(account string) (*common.UserFlashPoolOverview, error)
	GetAllMarkets() ([]ocommon.Address, error)
	GetInsuranceAddress(ocommon.Address) (ocommon.Address, error)
	ClaimWing(account string) (string, error)
}
