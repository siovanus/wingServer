package flashpool

import (
	"github.com/ontio/ontology/common"
)

type FlashPoolManager struct {
	contractAddress common.Address
}

func NewFlashPoolManager(contractAddress common.Address) *FlashPoolManager {
	manager := &FlashPoolManager{
		contractAddress,
	}

	return manager
}

func (this *FlashPoolManager) MarketDistribution() (*MarketDistribution, error) {
	return this.marketDistribution()
}

func (this *FlashPoolManager) PoolDistribution() (*PoolDistribution, error) {
	return this.poolDistribution()
}
