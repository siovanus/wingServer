package flashpool

import (
	"fmt"
	"github.com/siovanus/wingServer/config"
	"github.com/siovanus/wingServer/store"
	"math/big"
	"sort"
	"time"

	sdk "github.com/ontio/ontology-go-sdk"
	ocommon "github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/http/common"
	"github.com/siovanus/wingServer/manager/governance"
)

const (
	BlockPerYear      = 60 * 60 * 24 * 365 * 2 / 3
	PercentageDecimal = 10000
)

var WingDecimal = new(big.Int).SetUint64(1000000)
var FrontDecimal = new(big.Int).SetUint64(10000000000000000)
var FrontPercentageDecimal = new(big.Int).SetUint64(100000000000000)

type FlashPoolManager struct {
	cfg             *config.Config
	contractAddress ocommon.Address
	oracleAddress   ocommon.Address
	sdk             *sdk.OntologySdk
	store           *store.Client
}

func NewFlashPoolManager(contractAddress, oracleAddress ocommon.Address, sdk *sdk.OntologySdk,
	store *store.Client, cfg *config.Config) *FlashPoolManager {
	manager := &FlashPoolManager{
		cfg:             cfg,
		contractAddress: contractAddress,
		oracleAddress:   oracleAddress,
		sdk:             sdk,
		store:           store,
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
			Icon: this.cfg.IconMap[this.cfg.AssetMap[address.ToHexString()]],
			Name: this.cfg.AssetMap[address.ToHexString()],
			// totalDistribution / distributedDay
			PerDay:          new(big.Int).Div(new(big.Int).Div(totalDistribution, new(big.Int).SetUint64(distributedDay)), FrontDecimal).Uint64(),
			SupplyAmount:    new(big.Int).Div(supplyAmount, FrontDecimal).Uint64(),
			BorrowAmount:    new(big.Int).Div(borrowAmount, FrontDecimal).Uint64(),
			InsuranceAmount: new(big.Int).Div(insuranceAmount, FrontDecimal).Uint64(),
			Total:           new(big.Int).Div(totalDistribution, WingDecimal).Uint64(),
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
		price, err := this.assetPrice(this.cfg.OracleMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("PoolDistribution, this.assetPrice error: %s", err)
		}
		// supplyAmount * price
		distribution.SupplyAmount += new(big.Int).Div(new(big.Int).Mul(supplyAmount, new(big.Int).SetUint64(price)), FrontDecimal).Uint64()
		// borrowAmount * price
		distribution.BorrowAmount += new(big.Int).Div(new(big.Int).Mul(borrowAmount, new(big.Int).SetUint64(price)), FrontDecimal).Uint64()
		// insuranceAmount * price
		distribution.InsuranceAmount += new(big.Int).Div(new(big.Int).Mul(insuranceAmount, new(big.Int).SetUint64(price)), FrontDecimal).Uint64()
		distribution.Total += new(big.Int).Div(totalDistribution, WingDecimal).Uint64()
	}
	distribution.Name = "Flash"
	distribution.Icon = this.cfg.IconMap[distribution.Name]
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
		total += new(big.Int).Div(totalDistribution, WingDecimal).Uint64()
	}
	today := governance.DailyDistibute[index]
	var share uint64 = 0
	if total == 0 {
		share = 0
	}

	return &common.FlashPoolBanner{
		Today: today,
		Share: share,
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
		price, err := this.assetPrice(this.cfg.OracleMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("FlashPoolDetail, this.assetPrice error: %s", err)
		}
		// supplyAmount * price
		// borrowAmount * price
		// insuranceAmount * price
		supplyDollar := new(big.Int).Div(new(big.Int).Mul(supplyAmount, new(big.Int).SetUint64(price)), FrontDecimal).Uint64()
		borrowDollar := new(big.Int).Div(new(big.Int).Mul(borrowAmount, new(big.Int).SetUint64(price)), FrontDecimal).Uint64()
		insuranceDollar := new(big.Int).Div(new(big.Int).Mul(insuranceAmount, new(big.Int).SetUint64(price)), FrontDecimal).Uint64()
		flashPoolDetail.TotalSupply += supplyDollar
		flashPoolDetail.TotalBorrow += borrowDollar
		flashPoolDetail.TotalInsurance += insuranceDollar

		name := this.cfg.AssetMap[address.ToHexString()]
		flashPoolDetail.SupplyMarketRank = append(flashPoolDetail.SupplyMarketRank, &common.MarketFund{
			Icon: this.cfg.IconMap[name],
			Name: name,
			Fund: supplyDollar,
		})
		flashPoolDetail.BorrowMarketRank = append(flashPoolDetail.BorrowMarketRank, &common.MarketFund{
			Icon: this.cfg.IconMap[name],
			Name: name,
			Fund: borrowDollar,
		})
		flashPoolDetail.InsuranceMarketRank = append(flashPoolDetail.InsuranceMarketRank, &common.MarketFund{
			Icon: this.cfg.IconMap[name],
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
	preFlashPoolDetailStore, err := this.store.LoadLatestFlashPoolDetail()
	if err != nil {
		return nil, fmt.Errorf("FlashPoolDetail, this.store.LoadLastestFlashPoolDetail error: %s", err)
	}
	flashPoolDetail.SupplyVolumeDaily = int64(flashPoolDetail.TotalSupply) - int64(preFlashPoolDetailStore.TotalSupply)
	flashPoolDetail.BorrowVolumeDaily = int64(flashPoolDetail.TotalBorrow) - int64(preFlashPoolDetailStore.TotalBorrow)
	flashPoolDetail.InsuranceVolumeDaily = int64(flashPoolDetail.TotalInsurance) - int64(preFlashPoolDetailStore.TotalInsurance)

	if flashPoolDetail.TotalSupply != 0 {
		flashPoolDetail.TotalSupplyRate = flashPoolDetail.SupplyVolumeDaily * 100 / int64(flashPoolDetail.TotalSupply)
	}
	if flashPoolDetail.TotalBorrow != 0 {
		flashPoolDetail.TotalBorrowRate = flashPoolDetail.BorrowVolumeDaily * 100 / int64(flashPoolDetail.TotalBorrow)
	}
	if flashPoolDetail.TotalInsurance != 0 {
		flashPoolDetail.TotalInsuranceRate = flashPoolDetail.InsuranceVolumeDaily * 100 / int64(flashPoolDetail.TotalInsurance)
	}
	return flashPoolDetail, nil
}

func (this *FlashPoolManager) FlashPoolDetailForStore() (*store.FlashPoolDetail, error) {
	allMarkets, err := this.getAllMarkets()
	if err != nil {
		return nil, fmt.Errorf("FlashPoolDetailForStore, this.getAllMarkets error: %s", err)
	}
	flashPoolDetail := new(store.FlashPoolDetail)
	for _, address := range allMarkets {
		supplyAmount, err := this.getSupplyAmount(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolDetailForStore, this.getSupplyAmount error: %s", err)
		}
		borrowAmount, err := this.getBorrowAmount(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolDetailForStore, this.getSupplyAmount error: %s", err)
		}
		insuranceAmount, err := this.getInsuranceAmount(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolDetailForStore, this.getSupplyAmount error: %s", err)
		}
		price, err := this.assetPrice(this.cfg.OracleMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("FlashPoolDetailForStore, this.assetPrice error: %s", err)
		}
		// supplyAmount * price
		// borrowAmount * price
		// insuranceAmount * price
		supplyDollar := new(big.Int).Div(new(big.Int).Mul(supplyAmount, new(big.Int).SetUint64(price)), FrontDecimal).Uint64()
		borrowDollar := new(big.Int).Div(new(big.Int).Mul(borrowAmount, new(big.Int).SetUint64(price)), FrontDecimal).Uint64()
		insuranceDollar := new(big.Int).Div(new(big.Int).Mul(insuranceAmount, new(big.Int).SetUint64(price)), FrontDecimal).Uint64()
		flashPoolDetail.TotalSupply += supplyDollar
		flashPoolDetail.TotalBorrow += borrowDollar
		flashPoolDetail.TotalInsurance += insuranceDollar
	}
	flashPoolDetail.Timestamp = uint64(time.Now().Unix())
	return flashPoolDetail, nil
}

func (this *FlashPoolManager) FlashPoolMarketStore() error {
	allMarkets, err := this.getAllMarkets()
	if err != nil {
		return fmt.Errorf("FlashPoolMarketStore, this.getAllMarkets error: %s", err)
	}
	timestamp := uint64(time.Now().Unix())
	for _, address := range allMarkets {
		flashPoolMarket := new(store.FlashPoolMarket)
		supplyAmount, err := this.getSupplyAmount(address)
		if err != nil {
			return fmt.Errorf("FlashPoolMarketStore, this.getSupplyAmount error: %s", err)
		}
		borrowAmount, err := this.getBorrowAmount(address)
		if err != nil {
			return fmt.Errorf("FlashPoolMarketStore, this.getSupplyAmount error: %s", err)
		}
		insuranceAmount, err := this.getInsuranceAmount(address)
		if err != nil {
			return fmt.Errorf("FlashPoolMarketStore, this.getSupplyAmount error: %s", err)
		}
		name := this.cfg.AssetMap[address.ToHexString()]
		price, err := this.assetPrice(this.cfg.OracleMap[address.ToHexString()])
		if err != nil {
			return fmt.Errorf("FlashPoolMarketStore, this.assetPrice error: %s", err)
		}
		// supplyAmount * price
		// borrowAmount * price
		// insuranceAmount * price
		supplyDollar := new(big.Int).Div(new(big.Int).Mul(supplyAmount, new(big.Int).SetUint64(price)), FrontDecimal).Uint64()
		borrowDollar := new(big.Int).Div(new(big.Int).Mul(borrowAmount, new(big.Int).SetUint64(price)), FrontDecimal).Uint64()
		insuranceDollar := new(big.Int).Div(new(big.Int).Mul(insuranceAmount, new(big.Int).SetUint64(price)), FrontDecimal).Uint64()
		flashPoolMarket.Name = name
		flashPoolMarket.TotalSupply = supplyDollar
		flashPoolMarket.TotalBorrow = borrowDollar
		flashPoolMarket.TotalInsurance = insuranceDollar
		flashPoolMarket.Timestamp = timestamp
		err = this.store.SaveFlashPoolMarket(flashPoolMarket)
		if err != nil {
			return fmt.Errorf("FlashPoolMarketStore, this.store.SaveFlashPoolMarket error: %s", err)
		}
	}
	return nil
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
		price, err := this.assetPrice(this.cfg.OracleMap[address.ToHexString()])
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
		market.Name = this.cfg.AssetMap[address.ToHexString()]
		market.Icon = this.cfg.IconMap[market.Name]

		// supplyAmount * price
		// borrowAmount * price
		// insuranceAmount * price
		market.TotalSupply = new(big.Int).Div(new(big.Int).Mul(supplyAmount, new(big.Int).SetUint64(price)), FrontDecimal).Uint64()
		market.TotalBorrow = new(big.Int).Div(new(big.Int).Mul(borrowAmount, new(big.Int).SetUint64(price)), FrontDecimal).Uint64()
		market.TotalInsurance = new(big.Int).Div(new(big.Int).Mul(insuranceAmount, new(big.Int).SetUint64(price)), FrontDecimal).Uint64()

		market.SupplyApy = supplyApy
		market.BorrowApy = borrowApy
		market.InsuranceApy = insuranceApy

		latestFlashPoolMarket, err := this.store.LoadLatestFlashPoolMarket(market.Name)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolAllMarket, this.store.LoadLatestFlashPoolMarket error: %s", err)
		}
		if market.TotalSupply != 0 {
			market.TotalSupplyRate = (market.TotalSupply - latestFlashPoolMarket.TotalSupply) * 100 / market.TotalSupply
		}
		if market.TotalBorrow != 0 {
			market.TotalBorrowRate = (market.TotalBorrow - latestFlashPoolMarket.TotalBorrow) * 100 / market.TotalBorrow
		}
		if market.TotalInsurance != 0 {
			market.TotalInsuranceRate = (market.TotalInsurance - latestFlashPoolMarket.TotalInsurance) * 100 / market.TotalInsurance
		}
		flashPoolAllMarket.FlashPoolAllMarket = append(flashPoolAllMarket.FlashPoolAllMarket, market)
	}
	return flashPoolAllMarket, nil
}

func (this *FlashPoolManager) UserFlashPoolOverview(accountStr string) (*common.UserFlashPoolOverview, error) {
	account, err := ocommon.AddressFromBase58(accountStr)
	if err != nil {
		return nil, fmt.Errorf("UserFlashPoolOverview, ocommon.AddressFromBase58 error: %s", err)
	}
	allMarkets, err := this.getAllMarkets()
	if err != nil {
		return nil, fmt.Errorf("UserFlashPoolOverview, this.getAllMarkets error: %s", err)
	}
	assetsIn, _ := this.getAssetsIn(account)
	userFlashPoolOverview := &common.UserFlashPoolOverview{
		CurrentSupply:    make([]*common.Supply, 0),
		CurrentBorrow:    make([]*common.Borrow, 0),
		CurrentInsurance: make([]*common.Insurance, 0),
		AllMarket:        make([]*common.UserMarket, 0),
	}

	var assetInSupplyDollar uint64 = 0
	var totalBorrowBalance uint64 = 0
	for _, address := range allMarkets {
		borrowAmount, err := this.getBorrowAmountByAccount(address, account)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverview, this.getSupplyAmountByAccount error: %s", err)
		}
		price, err := this.assetPrice(this.cfg.OracleMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverview, this.assetPrice error: %s", err)
		}
		// borrowAmount * price
		totalBorrowBalance += new(big.Int).Div(new(big.Int).Mul(borrowAmount, new(big.Int).SetUint64(price)), FrontDecimal).Uint64()
	}
	userFlashPoolOverview.BorrowBalance = totalBorrowBalance

	for _, address := range allMarkets {
		supplyAmount, err := this.getSupplyAmountByAccount(address, account)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverview, this.getSupplyAmountByAccount error: %s", err)
		}
		borrowAmount, err := this.getBorrowAmountByAccount(address, account)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverview, this.getSupplyAmountByAccount error: %s", err)
		}
		insuranceAmount, err := this.getInsuranceAmountByAccount(address, account)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverview, this.getSupplyAmountByAccount error: %s", err)
		}
		price, err := this.assetPrice(this.cfg.OracleMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverview, this.assetPrice error: %s", err)
		}
		// supplyAmount * price
		// borrowAmount * price
		// insuranceAmount * price
		supplyAmountU64 := new(big.Int).Div(new(big.Int).Mul(supplyAmount, new(big.Int).SetUint64(price)), FrontDecimal).Uint64()
		borrowAmountU64 := new(big.Int).Div(new(big.Int).Mul(borrowAmount, new(big.Int).SetUint64(price)), FrontDecimal).Uint64()
		insuranceAmountU64 := new(big.Int).Div(new(big.Int).Mul(insuranceAmount, new(big.Int).SetUint64(price)), FrontDecimal).Uint64()
		userFlashPoolOverview.SupplyBalance += supplyAmountU64
		userFlashPoolOverview.InsuranceBalance += insuranceAmountU64
		supplyApy, err := this.getSupplyApy(address)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverview, this.getSupplyApy error: %s", err)
		}
		borrowApy, err := this.getBorrowApy(address)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverview, this.getBorrowApy error: %s", err)
		}
		insuranceApy, err := this.getInsuranceApy(address)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverview, this.getInsuranceApy error: %s", err)
		}
		userFlashPoolOverview.NetApy += int64(supplyAmountU64*supplyApy+insuranceAmountU64*insuranceApy) -
			int64(borrowAmountU64*borrowApy)

		isAssetIn := false
		for _, a := range assetsIn {
			if address == a {
				isAssetIn = true
				break
			}
		}
		if isAssetIn {
			assetInSupplyDollar += supplyAmountU64
		}

		if supplyAmount.Uint64() != 0 {
			supply := &common.Supply{
				Name:          this.cfg.AssetMap[address.ToHexString()],
				Icon:          this.cfg.IconMap[this.cfg.AssetMap[address.ToHexString()]],
				SupplyDollar:  supplyAmountU64,
				SupplyBalance: new(big.Int).Div(supplyAmount, FrontDecimal).Uint64(),
				Apy:           supplyApy,
				IfCollateral:  isAssetIn,
			}
			userFlashPoolOverview.CurrentSupply = append(userFlashPoolOverview.CurrentSupply, supply)
		}
		if borrowAmount.Uint64() != 0 {
			borrow := &common.Borrow{
				Name:          this.cfg.AssetMap[address.ToHexString()],
				Icon:          this.cfg.IconMap[this.cfg.AssetMap[address.ToHexString()]],
				BorrowDollar:  borrowAmountU64,
				BorrowBalance: new(big.Int).Div(borrowAmount, FrontDecimal).Uint64(),
				Apy:           borrowApy,
				Limit:         borrowAmountU64 * PercentageDecimal / totalBorrowBalance,
			}
			userFlashPoolOverview.CurrentBorrow = append(userFlashPoolOverview.CurrentBorrow, borrow)
		}
		if insuranceAmount.Uint64() != 0 {
			insurance := &common.Insurance{
				Name:             this.cfg.AssetMap[address.ToHexString()],
				Icon:             this.cfg.IconMap[this.cfg.AssetMap[address.ToHexString()]],
				InsuranceDollar:  insuranceAmountU64,
				InsuranceBalance: new(big.Int).Div(insuranceAmount, FrontDecimal).Uint64(),
				Apy:              insuranceApy,
			}
			userFlashPoolOverview.CurrentInsurance = append(userFlashPoolOverview.CurrentInsurance, insurance)
		}

		totalBorrowAmount, err := this.getSupplyAmount(address)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverview, this.getSupplyAmount error: %s", err)
		}
		totalInsuranceAmount, err := this.getSupplyAmount(address)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverview, this.getSupplyAmount error: %s", err)
		}
		if supplyAmount.Uint64() == 0 && borrowAmount.Uint64() == 0 && insuranceAmount.Uint64() == 0 {
			userMarket := &common.UserMarket{
				Name:            this.cfg.AssetMap[address.ToHexString()],
				Icon:            this.cfg.IconMap[this.cfg.AssetMap[address.ToHexString()]],
				SupplyApy:       supplyApy,
				BorrowApy:       borrowApy,
				BorrowLiquidity: new(big.Int).Div(totalBorrowAmount, FrontDecimal).Uint64(),
				InsuranceApy:    insuranceApy,
				InsuranceAmount: new(big.Int).Div(totalInsuranceAmount, FrontDecimal).Uint64(),
			}
			userFlashPoolOverview.AllMarket = append(userFlashPoolOverview.AllMarket, userMarket)
		}
	}

	marketMeta, err := this.getMarketMeta()
	if err == nil {
		// assetInSupplyDollar(already multiply 100) * CollateralFactor / FrontDecimal
		userFlashPoolOverview.BorrowLimit = new(big.Int).Div(new(big.Int).Mul(new(big.Int).SetUint64(assetInSupplyDollar),
			marketMeta.CollateralFactor.ToBigInt()), FrontDecimal).Uint64()
	}
	if userFlashPoolOverview.SupplyBalance+userFlashPoolOverview.BorrowBalance+userFlashPoolOverview.InsuranceBalance != 0 {
		userFlashPoolOverview.NetApy = userFlashPoolOverview.NetApy / int64(userFlashPoolOverview.SupplyBalance+
			userFlashPoolOverview.BorrowBalance+userFlashPoolOverview.InsuranceBalance)
	}

	return userFlashPoolOverview, nil
}
