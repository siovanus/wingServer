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

func (this *FlashPoolManager) FlashPoolBanner() (*FlashPoolBanner, error) {
	return this.flashPoolBanner()
}

func (this *FlashPoolManager) FlashPoolDetail() (*FlashPoolDetail, error) {
	return this.flashPoolDetail()
}

func (this *FlashPoolManager) FlashPoolAllMarket() (*FlashPoolAllMarket, error) {
	return this.flashPoolAllMarket()
}
