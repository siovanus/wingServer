package flashpool

import (
	"fmt"
	sdk "github.com/ontio/ontology-go-sdk"
	ocommon "github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/http/common"
)

const BlockPerYear = 60 * 60 * 24 * 365 * 2 / 3

var IconMap = map[string]string{
	"oETH":  "http://106.75.209.209/icon/eth_icon.svg",
	"oDAI":  "http://106.75.209.209/icon/asset_dai_icon.svg",
	"Flash": "http://106.75.209.209/icon/flash_icon.svg",
	"IF":    "http://106.75.209.209/icon/if_icon.svg",
}

var AssetMap = map[string]string{
	"a": "oWBTC",
	"b": "oETH",
	"c": "oDAI",
	"d": "ONT",
}

type FlashPoolManager struct {
	contractAddress ocommon.Address
	sdk             *sdk.OntologySdk
}

func NewFlashPoolManager(contractAddress ocommon.Address, sdk *sdk.OntologySdk) *FlashPoolManager {
	manager := &FlashPoolManager{
		contractAddress: contractAddress,
		sdk:             sdk,
	}

	return manager
}

func (this *FlashPoolManager) FlashPoolMarketDistribution() (*common.FlashPoolMarketDistribution, error) {
	allMarkets, err := this.getAllMarkets()
	if err != nil {
		return nil, fmt.Errorf("FlashPoolMarketDistribution, this.getAllMarkets error: %s", err)
	}
	flashPoolMarketDistribution := make([]*common.Distribution, 0)
	for _, address := range allMarkets {
		supplyAmount, err := this.getSupplyAmount(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolMarketDistribution, this.getSupplyApy error: %s", err)
		}
		borrowAmount, err := this.getBorrowAmount(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolMarketDistribution, this.getBorrowApy error: %s", err)
		}
		insuranceAmount, err := this.getInsuranceAmount(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolMarketDistribution, this.getInsuranceApy error: %s", err)
		}
		distribution := &common.Distribution{
			Icon: IconMap[AssetMap[address.ToHexString()]],
			Name: AssetMap[address.ToHexString()],
			//TODO: Per Day
			SupplyAmount:    supplyAmount,
			BorrowAmount:    borrowAmount,
			InsuranceAmount: insuranceAmount,
			//TODO: total distribution
		}
		flashPoolMarketDistribution = append(flashPoolMarketDistribution, distribution)
	}
	return &common.FlashPoolMarketDistribution{FlashPoolMarketDistribution: flashPoolMarketDistribution}, nil
}

func (this *FlashPoolManager) PoolDistribution() (*common.PoolDistribution, error) {
	allMarkets, err := this.getAllMarkets()
	if err != nil {
		return nil, fmt.Errorf("PoolDistribution, this.getAllMarkets error: %s", err)
	}
	poolDistribution := make([]*common.Distribution, 0)
	for _, address := range allPools {
		supplyAmount, err := this.getSupplyAmount(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolMarketDistribution, this.getSupplyApy error: %s", err)
		}
		borrowAmount, err := this.getBorrowAmount(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolMarketDistribution, this.getBorrowApy error: %s", err)
		}
		insuranceAmount, err := this.getInsuranceAmount(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolMarketDistribution, this.getInsuranceApy error: %s", err)
		}
		distribution := &common.Distribution{
			Icon: IconMap[AssetMap[address.ToHexString()]],
			Name: AssetMap[address.ToHexString()],
			//TODO: Per Day
			SupplyAmount:    supplyAmount,
			BorrowAmount:    borrowAmount,
			InsuranceAmount: insuranceAmount,
			//TODO: total distribution
		}
		flashPoolMarketDistribution = append(flashPoolMarketDistribution, distribution)
	}
	return &common.FlashPoolMarketDistribution{FlashPoolMarketDistribution: flashPoolMarketDistribution}, nil
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
