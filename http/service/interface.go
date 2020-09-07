package service

import (
	"github.com/siovanus/wingServer/http/common"
	"github.com/siovanus/wingServer/store"
)

type GovernanceManager interface {
	GovBannerOverview() (*common.GovBannerOverview, error)
	GovBanner() (*common.GovBanner, error)
}

type FlashPoolManager interface {
	AssetPrice(asset string) (uint64, error)
	FlashPoolMarketDistribution() (*common.FlashPoolMarketDistribution, error)
	PoolDistribution() (*common.Distribution, error)
	FlashPoolBanner() (*common.FlashPoolBanner, error)
	FlashPoolDetail() (*common.FlashPoolDetail, error)
	FlashPoolDetailForStore() (*store.FlashPoolDetail, error)
	FlashPoolMarketStore() error
	FlashPoolAllMarket() (*common.FlashPoolAllMarket, error)
	UserFlashPoolOverview(account string) (*common.UserFlashPoolOverview, error)
}
