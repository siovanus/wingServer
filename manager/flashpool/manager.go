package flashpool

import (
	"fmt"
	sdk "github.com/ontio/ontology-go-sdk"
	ocommon "github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/http/common"
	"github.com/siovanus/wingServer/manager/governance"
	"time"
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
	oracleAddress   ocommon.Address
	sdk             *sdk.OntologySdk
}

func NewFlashPoolManager(contractAddress, oracleAddress ocommon.Address, sdk *sdk.OntologySdk) *FlashPoolManager {
	manager := &FlashPoolManager{
		contractAddress: contractAddress,
		oracleAddress:   oracleAddress,
		sdk:             sdk,
	}

	return manager
}

func (this *FlashPoolManager) AssetPrice(asset string) (uint64, error) {
	return this.assetPrice(asset)
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
		totalDistribution, err := this.getTotalDistribution(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolMarketDistribution, this.getInsuranceApy error: %s", err)
		}
		distributedDay := (uint64(time.Now().Unix()) - governance.GenesisTime) / governance.DaySecond
		distribution := &common.Distribution{
			Icon:            IconMap[AssetMap[address.ToHexString()]],
			Name:            AssetMap[address.ToHexString()],
			PerDay:          totalDistribution / distributedDay,
			SupplyAmount:    supplyAmount,
			BorrowAmount:    borrowAmount,
			InsuranceAmount: insuranceAmount,
			Total:           totalDistribution,
		}
		flashPoolMarketDistribution = append(flashPoolMarketDistribution, distribution)
	}
	return &common.FlashPoolMarketDistribution{FlashPoolMarketDistribution: flashPoolMarketDistribution}, nil
}

func (this *FlashPoolManager) PoolDistribution() (*common.Distribution, error) {
	allMarkets, err := this.getAllMarkets()
	if err != nil {
		return nil, fmt.Errorf("PoolDistribution, this.getAllMarkets error: %s", err)
	}
	distribution := new(common.Distribution)
	for _, address := range allMarkets {
		supplyAmount, err := this.getSupplyAmount(address)
		if err != nil {
			return nil, fmt.Errorf("PoolDistribution, this.getSupplyAmount error: %s", err)
		}
		borrowAmount, err := this.getBorrowAmount(address)
		if err != nil {
			return nil, fmt.Errorf("PoolDistribution, this.getSupplyAmount error: %s", err)
		}
		insuranceAmount, err := this.getInsuranceAmount(address)
		if err != nil {
			return nil, fmt.Errorf("PoolDistribution, this.getSupplyAmount error: %s", err)
		}
		totalDistribution, err := this.getTotalDistribution(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolMarketDistribution, this.getInsuranceApy error: %s", err)
		}
		price, err := this.assetPrice(AssetMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("PoolDistribution, this.assetPrice error: %s", err)
		}
		distribution.SupplyAmount += supplyAmount * price
		distribution.BorrowAmount += borrowAmount * price
		distribution.InsuranceAmount += insuranceAmount * price
		distribution.Total += totalDistribution
	}
	distributedDay := (uint64(time.Now().Unix()) - governance.GenesisTime) / governance.DaySecond
	distribution.PerDay = distribution.Total / distributedDay
	return distribution, nil
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

func (this *FlashPoolManager) UserFlashPoolOverview(accountStr string) (*common.UserFlashPoolOverview, error) {
	//account, err := ocommon.AddressFromBase58(accountStr)
	//if err != nil {
	//	return nil, fmt.Errorf("UserFlashPoolOverview, ocommon.AddressFromBase58 error: %s", err)
	//}
	//assetsIn, err := this.getAssetsIn(account)
	//if err != nil {
	//	return nil, fmt.Errorf("UserFlashPoolOverview, this.getAssetsIn error: %s", err)
	//}
	result := &common.UserFlashPoolOverview{
		CurrentSupply:    make([]*common.Supply, 0),
		CurrentBorrow:    make([]*common.Borrow, 0),
		CurrentInsurance: make([]*common.Insurance, 0),
		AllMarket:        make([]*common.UserMarket, 0),
	}
	//for _, address := range assetsIn {
	//	supplyAmount, err := this.getSupplyAmount(address)
	//	if err != nil {
	//		return nil, fmt.Errorf("UserFlashPoolOverview, this.getSupplyApy error: %s", err)
	//	}
	//	borrowAmount, err := this.getBorrowAmount(address)
	//	if err != nil {
	//		return nil, fmt.Errorf("UserFlashPoolOverview, this.getBorrowApy error: %s", err)
	//	}
	//	insuranceAmount, err := this.getInsuranceAmount(address)
	//	if err != nil {
	//		return nil, fmt.Errorf("UserFlashPoolOverview, this.getInsuranceApy error: %s", err)
	//	}
	//}
	return result, nil
}
