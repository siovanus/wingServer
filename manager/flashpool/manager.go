package flashpool

import (
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
)

type FlashPoolManager struct {
	contractAddress common.Address
	sdk             *sdk.OntologySdk
}

func NewFlashPoolManager(contractAddress common.Address, sdk *sdk.OntologySdk) *FlashPoolManager {
	manager := &FlashPoolManager{
		contractAddress: contractAddress,
		sdk:             sdk,
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

func (this *FlashPoolManager) UserFlashPoolOverview(address string) (*UserFlashPoolOverview, error) {
	return this.userFlashPoolOverview(address)
}
