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
	MarketDistribution() (*flashpool.MarketDistribution, error)
	PoolDistribution() (*flashpool.PoolDistribution, error)
	FlashPoolBanner() (*flashpool.FlashPoolBanner, error)
	FlashPoolDetail() (*flashpool.FlashPoolDetail, error)
	FlashPoolAllMarket() (*flashpool.FlashPoolAllMarket, error)
	UserFlashPoolOverview(address string) (*flashpool.UserFlashPoolOverview, error)
}

type OracleManager interface {
	AssetPrice(asset string) (uint64, error)
}
