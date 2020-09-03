package service

import "github.com/siovanus/wingServer/manager/flashpool"

type GovernanceManager interface {
	QueryData()
}

type FlashPoolManager interface {
	MarketDistribution() (*flashpool.MarketDistribution, error)
	PoolDistribution() (*flashpool.PoolDistribution, error)
}
