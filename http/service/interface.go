package service

import (
	"github.com/siovanus/wingServer/http/common"
	"github.com/siovanus/wingServer/manager/flashpool"
)

type GovernanceManager interface {
	GovBannerOverview() (*common.GovBannerOverview, error)
	GovBanner() (*common.GovBanner, error)
}

type FlashPoolManager interface {
	AssetPrice(asset string) (uint64, error)
	FlashPoolMarketDistribution() (*common.FlashPoolMarketDistribution, error)
	PoolDistribution() (*common.Distribution, error)
	FlashPoolBanner() (*flashpool.FlashPoolBanner, error)
	FlashPoolDetail() (*flashpool.FlashPoolDetail, error)
	FlashPoolAllMarket() (*flashpool.FlashPoolAllMarket, error)
	UserFlashPoolOverview(account string) (*common.UserFlashPoolOverview, error)
}
