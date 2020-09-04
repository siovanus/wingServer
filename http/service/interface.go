package service

import (
	"github.com/siovanus/wingServer/manager/flashpool"
	"github.com/siovanus/wingServer/manager/governance"
)

type GovernanceManager interface {
	GovBannerOverview() (*governance.GovBanner, error)
	GovBanner() (*governance.PoolBanner, error)
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
