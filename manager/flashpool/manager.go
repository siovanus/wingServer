package flashpool

import (
	"fmt"
	"sort"
	"time"

	sdk "github.com/ontio/ontology-go-sdk"
	ocommon "github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/http/common"
	"github.com/siovanus/wingServer/manager/governance"
)

const (
	BlockPerYear         = 60 * 60 * 24 * 365 * 2 / 3
	PercentageMultiplier = 10000
)

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
	"e": "oUSDT",
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
			return nil, fmt.Errorf("FlashPoolMarketDistribution, this.getSupplyAmount error: %s", err)
		}
		borrowAmount, err := this.getBorrowAmount(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolMarketDistribution, this.getBorrowAmount error: %s", err)
		}
		insuranceAmount, err := this.getInsuranceAmount(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolMarketDistribution, this.getInsuranceAmount error: %s", err)
		}
		totalDistribution, err := this.getTotalDistribution(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolMarketDistribution, this.getTotalDistribution error: %s", err)
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
			return nil, fmt.Errorf("PoolDistribution, this.getTotalDistribution error: %s", err)
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

func (this *FlashPoolManager) FlashPoolBanner() (*common.FlashPoolBanner, error) {
	distributed := uint64(time.Now().Unix()) - governance.GenesisTime
	index := distributed/governance.YearSecond + 1

	allMarkets, err := this.getAllMarkets()
	if err != nil {
		return nil, fmt.Errorf("FlashPoolBanner, this.getAllMarkets error: %s", err)
	}
	var total uint64 = 0
	for _, address := range allMarkets {
		totalDistribution, err := this.getTotalDistribution(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolBanner, this.getTotalDistribution error: %s", err)
		}
		total += totalDistribution
	}
	today := governance.DailyDistibute[index]

	return &common.FlashPoolBanner{
		Today: today,
		Share: today * PercentageMultiplier / total,
		Total: total,
	}, nil
}

func (this *FlashPoolManager) FlashPoolDetail() (*common.FlashPoolDetail, error) {
	allMarkets, err := this.getAllMarkets()
	if err != nil {
		return nil, fmt.Errorf("FlashPoolDetail, this.getAllMarkets error: %s", err)
	}
	flashPoolDetail := &common.FlashPoolDetail{
		SupplyMarketRank:    make([]*common.MarketFund, 0),
		BorrowMarketRank:    make([]*common.MarketFund, 0),
		InsuranceMarketRank: make([]*common.MarketFund, 0),
	}
	for _, address := range allMarkets {
		supplyAmount, err := this.getSupplyAmount(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolDetail, this.getSupplyAmount error: %s", err)
		}
		borrowAmount, err := this.getBorrowAmount(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolDetail, this.getSupplyAmount error: %s", err)
		}
		insuranceAmount, err := this.getInsuranceAmount(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolDetail, this.getSupplyAmount error: %s", err)
		}
		price, err := this.assetPrice(AssetMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("FlashPoolDetail, this.assetPrice error: %s", err)
		}
		supplyDollar := supplyAmount * price
		borrowDollar := borrowAmount * price
		insuranceDollar := insuranceAmount * price
		flashPoolDetail.TotalSupply += supplyDollar
		flashPoolDetail.TotalBorrow += borrowDollar
		flashPoolDetail.TotalInsurance += insuranceDollar

		name := AssetMap[address.ToHexString()]
		flashPoolDetail.SupplyMarketRank = append(flashPoolDetail.SupplyMarketRank, &common.MarketFund{
			Icon: IconMap[name],
			Name: name,
			Fund: supplyDollar,
		})
		flashPoolDetail.BorrowMarketRank = append(flashPoolDetail.BorrowMarketRank, &common.MarketFund{
			Icon: IconMap[name],
			Name: name,
			Fund: borrowDollar,
		})
		flashPoolDetail.InsuranceMarketRank = append(flashPoolDetail.InsuranceMarketRank, &common.MarketFund{
			Icon: IconMap[name],
			Name: name,
			Fund: insuranceDollar,
		})
	}
	sort.SliceStable(flashPoolDetail.SupplyMarketRank, func(i, j int) bool {
		return flashPoolDetail.SupplyMarketRank[i].Fund > flashPoolDetail.SupplyMarketRank[j].Fund
	})
	sort.SliceStable(flashPoolDetail.BorrowMarketRank, func(i, j int) bool {
		return flashPoolDetail.BorrowMarketRank[i].Fund > flashPoolDetail.BorrowMarketRank[j].Fund
	})
	sort.SliceStable(flashPoolDetail.InsuranceMarketRank, func(i, j int) bool {
		return flashPoolDetail.InsuranceMarketRank[i].Fund > flashPoolDetail.InsuranceMarketRank[j].Fund
	})
	//TODO: rate
	return flashPoolDetail, nil
}

func (this *FlashPoolManager) FlashPoolAllMarket() (*common.FlashPoolAllMarket, error) {
	allMarkets, err := this.getAllMarkets()
	if err != nil {
		return nil, fmt.Errorf("FlashPoolAllMarket, this.getAllMarkets error: %s", err)
	}
	flashPoolAllMarket := &common.FlashPoolAllMarket{
		FlashPoolAllMarket: make([]*common.Market, 0),
	}
	for _, address := range allMarkets {
		supplyAmount, err := this.getSupplyAmount(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolAllMarket, this.getSupplyAmount error: %s", err)
		}
		borrowAmount, err := this.getBorrowAmount(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolAllMarket, this.getSupplyAmount error: %s", err)
		}
		insuranceAmount, err := this.getInsuranceAmount(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolAllMarket, this.getSupplyAmount error: %s", err)
		}
		price, err := this.assetPrice(AssetMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("FlashPoolAllMarket, this.assetPrice error: %s", err)
		}

		supplyApy, err := this.getSupplyApy(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolAllMarket, this.getSupplyApy error: %s", err)
		}
		borrowApy, err := this.getBorrowApy(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolAllMarket, this.getBorrowApy error: %s", err)
		}
		insuranceApy, err := this.getInsuranceApy(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolAllMarket, this.getInsuranceApy error: %s", err)
		}

		market := new(common.Market)
		market.Name = AssetMap[address.ToHexString()]
		market.Icon = IconMap[market.Name]
		market.TotalSupply = supplyAmount * price
		market.TotalBorrow = borrowAmount * price
		market.TotalInsurance = insuranceAmount * price
		market.SupplyApy = supplyApy
		market.BorrowApy = borrowApy
		market.InsuranceApy = insuranceApy
		//TODO: rate
		flashPoolAllMarket.FlashPoolAllMarket = append(flashPoolAllMarket.FlashPoolAllMarket, market)
	}
	return flashPoolAllMarket, nil
}

func (this *FlashPoolManager) UserFlashPoolOverview(accountStr string) (*common.UserFlashPoolOverview, error) {
	//TODO
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
