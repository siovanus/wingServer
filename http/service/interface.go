package service

import (
	"github.com/siovanus/wingServer/manager/flashpool"
	"github.com/siovanus/wingServer/manager/governance"
)

type GovernanceManager interface {
	GovBanner() (*governance.GovBanner, error)
}

type FlashPoolManager interface {
	MarketDistribution() (*flashpool.MarketDistribution, error)
	PoolDistribution() (*flashpool.PoolDistribution, error)
	PoolBanner() (*flashpool.PoolBanner, error)
}
