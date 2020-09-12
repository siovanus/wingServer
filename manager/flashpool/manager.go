package flashpool

import (
	"fmt"
	"math/big"
	"sort"
	"time"

	sdk "github.com/ontio/ontology-go-sdk"
	ocommon "github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/config"
	"github.com/siovanus/wingServer/http/common"
	"github.com/siovanus/wingServer/manager/governance"
	"github.com/siovanus/wingServer/store"
)

const (
	BlockPerYear      = 60 * 60 * 24 * 365 * 2 / 3
	PercentageDecimal = 10000
)

var PriceDecimal = new(big.Int).SetUint64(1000000000)

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

func (this *FlashPoolManager) AssetStoredPrice(asset string) (uint64, error) {
	if asset == "USDT" {
		return PriceDecimal.Uint64(), nil
	}
	price, err := this.store.LoadPrice(asset)
	if err != nil {
		return 0, fmt.Errorf("AssetStoredPrice, this.store.LoadPrice error: %s", err)
	}
	return price.Price, nil
}

func (this *FlashPoolManager) FlashPoolMarketDistribution() (*common.FlashPoolMarketDistribution, error) {
	allMarkets, err := this.GetAllMarkets()
	if err != nil {
		return nil, fmt.Errorf("FlashPoolMarketDistribution, this.GetAllMarkets error: %s", err)
	}
	flashPoolMarketDistribution := make([]*common.Distribution, 0)
	for _, address := range allMarkets {
		market, err := this.store.LoadFlashMarket(this.cfg.AssetMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("FlashPoolDetail, this.store.LoadFlashMarket error: %s", err)
		}
		supplyAmount := new(big.Int).SetUint64(market.TotalSupply)
		borrowAmount := new(big.Int).SetUint64(market.TotalBorrow)
		insuranceAmount := new(big.Int).SetUint64(market.TotalInsurance)

		totalDistribution, err := this.getTotalDistribution(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolMarketDistribution, this.getTotalDistribution error: %s", err)
		}
		distributedDay := (uint64(time.Now().Unix()) - governance.GenesisTime) / governance.DaySecond
		distribution := &common.Distribution{
			Icon: this.cfg.IconMap[this.cfg.AssetMap[address.ToHexString()]],
			Name: this.cfg.AssetMap[address.ToHexString()],
			// totalDistribution / distributedDay
			PerDay:          new(big.Int).Div(totalDistribution, new(big.Int).SetUint64(distributedDay)).Uint64(),
			SupplyAmount:    supplyAmount.Uint64(),
			BorrowAmount:    borrowAmount.Uint64(),
			InsuranceAmount: insuranceAmount.Uint64(),
			Total:           totalDistribution.Uint64(),
		}
		flashPoolMarketDistribution = append(flashPoolMarketDistribution, distribution)
	}
	return &common.FlashPoolMarketDistribution{FlashPoolMarketDistribution: flashPoolMarketDistribution}, nil
}

func (this *FlashPoolManager) PoolDistribution() (*common.Distribution, error) {
	allMarkets, err := this.GetAllMarkets()
	if err != nil {
		return nil, fmt.Errorf("PoolDistribution, this.GetAllMarkets error: %s", err)
	}
	distribution := new(common.Distribution)
	for _, address := range allMarkets {
		market, err := this.store.LoadFlashMarket(this.cfg.AssetMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("PoolDistribution, this.store.LoadFlashMarket error: %s", err)
		}
		supplyAmount := new(big.Int).SetUint64(market.TotalSupply)
		borrowAmount := new(big.Int).SetUint64(market.TotalBorrow)
		insuranceAmount := new(big.Int).SetUint64(market.TotalInsurance)

		totalDistribution, err := this.getTotalDistribution(address)
		if err != nil {
			return nil, fmt.Errorf("PoolDistribution, this.getTotalDistribution error: %s", err)
		}
		price, err := this.AssetStoredPrice(this.cfg.OracleMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("PoolDistribution, this.AssetStoredPrice error: %s", err)
		}
		// supplyAmount * price
		distribution.SupplyAmount += new(big.Int).Mul(supplyAmount, new(big.Int).SetUint64(price)).Uint64()
		// borrowAmount * price
		distribution.BorrowAmount += new(big.Int).Mul(borrowAmount, new(big.Int).SetUint64(price)).Uint64()
		// insuranceAmount * price
		distribution.InsuranceAmount += new(big.Int).Mul(insuranceAmount, new(big.Int).SetUint64(price)).Uint64()
		distribution.Total += totalDistribution.Uint64()
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

	allMarkets, err := this.GetAllMarkets()
	if err != nil {
		return nil, fmt.Errorf("FlashPoolBanner, this.GetAllMarkets error: %s", err)
	}
	var total uint64 = 0
	for _, address := range allMarkets {
		totalDistribution, err := this.getTotalDistribution(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolBanner, this.getTotalDistribution error: %s", err)
		}
		total += totalDistribution.Uint64()
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
	allMarkets, err := this.GetAllMarkets()
	if err != nil {
		return nil, fmt.Errorf("FlashPoolDetail, this.GetAllMarkets error: %s", err)
	}
	flashPoolDetail := &common.FlashPoolDetail{
		SupplyMarketRank:    make([]*common.MarketFund, 0),
		BorrowMarketRank:    make([]*common.MarketFund, 0),
		InsuranceMarketRank: make([]*common.MarketFund, 0),
	}
	for _, address := range allMarkets {
		market, err := this.store.LoadFlashMarket(this.cfg.AssetMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("FlashPoolDetail, this.store.LoadFlashMarket error: %s", err)
		}
		supplyAmount := new(big.Int).SetUint64(market.TotalSupply)
		borrowAmount := new(big.Int).SetUint64(market.TotalBorrow)
		insuranceAmount := new(big.Int).SetUint64(market.TotalInsurance)

		price, err := this.AssetStoredPrice(this.cfg.OracleMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("FlashPoolDetail, this.AssetStoredPrice error: %s", err)
		}
		// supplyAmount * price
		// borrowAmount * price
		// insuranceAmount * price
		supplyDollar := new(big.Int).Div(new(big.Int).Mul(supplyAmount, new(big.Int).SetUint64(price)), PriceDecimal).Uint64()
		borrowDollar := new(big.Int).Div(new(big.Int).Mul(borrowAmount, new(big.Int).SetUint64(price)), PriceDecimal).Uint64()
		insuranceDollar := new(big.Int).Div(new(big.Int).Mul(insuranceAmount, new(big.Int).SetUint64(price)), PriceDecimal).Uint64()
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
	allMarkets, err := this.GetAllMarkets()
	if err != nil {
		return nil, fmt.Errorf("FlashPoolDetailForStore, this.GetAllMarkets error: %s", err)
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
		price, err := this.AssetPrice(this.cfg.OracleMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("FlashPoolDetailForStore, this.AssetStoredPrice error: %s", err)
		}
		// supplyAmount * price
		// borrowAmount * price
		// insuranceAmount * price
		supplyDollar := new(big.Int).Div(new(big.Int).Mul(supplyAmount, new(big.Int).SetUint64(price)), PriceDecimal).Uint64()
		borrowDollar := new(big.Int).Div(new(big.Int).Mul(borrowAmount, new(big.Int).SetUint64(price)), PriceDecimal).Uint64()
		insuranceDollar := new(big.Int).Div(new(big.Int).Mul(insuranceAmount, new(big.Int).SetUint64(price)), PriceDecimal).Uint64()
		flashPoolDetail.TotalSupply += supplyDollar
		flashPoolDetail.TotalBorrow += borrowDollar
		flashPoolDetail.TotalInsurance += insuranceDollar
	}
	flashPoolDetail.Timestamp = uint64(time.Now().Unix())
	return flashPoolDetail, nil
}

func (this *FlashPoolManager) FlashPoolMarketStore() error {
	allMarkets, err := this.GetAllMarkets()
	if err != nil {
		return fmt.Errorf("FlashPoolMarketStore, this.GetAllMarkets error: %s", err)
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
		price, err := this.AssetStoredPrice(this.cfg.OracleMap[address.ToHexString()])
		if err != nil {
			return fmt.Errorf("FlashPoolMarketStore, this.AssetStoredPrice error: %s", err)
		}
		// supplyAmount * price
		// borrowAmount * price
		// insuranceAmount * price
		supplyDollar := new(big.Int).Div(new(big.Int).Mul(supplyAmount, new(big.Int).SetUint64(price)), PriceDecimal).Uint64()
		borrowDollar := new(big.Int).Div(new(big.Int).Mul(borrowAmount, new(big.Int).SetUint64(price)), PriceDecimal).Uint64()
		insuranceDollar := new(big.Int).Div(new(big.Int).Mul(insuranceAmount, new(big.Int).SetUint64(price)), PriceDecimal).Uint64()
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
	allMarkets, err := this.GetAllMarkets()
	if err != nil {
		return nil, fmt.Errorf("FlashPoolAllMarket, this.GetAllMarkets error: %s", err)
	}
	flashPoolAllMarket := &common.FlashPoolAllMarket{
		FlashPoolAllMarket: make([]*common.Market, 0),
	}
	for _, address := range allMarkets {
		market, err := this.store.LoadFlashMarket(this.cfg.AssetMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("FlashPoolAllMarket, this.store.LoadFlashMarket error: %s", err)
		}
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
		flashPoolAllMarket.FlashPoolAllMarket = append(flashPoolAllMarket.FlashPoolAllMarket, &market)
	}
	return flashPoolAllMarket, nil
}

func (this *FlashPoolManager) FlashPoolAllMarketForStore() (*common.FlashPoolAllMarket, error) {
	allMarkets, err := this.GetAllMarkets()
	if err != nil {
		return nil, fmt.Errorf("FlashPoolAllMarketForStore, this.GetAllMarkets error: %s", err)
	}
	flashPoolAllMarket := &common.FlashPoolAllMarket{
		FlashPoolAllMarket: make([]*common.Market, 0),
	}
	for _, address := range allMarkets {
		supplyAmount, err := this.getSupplyAmount(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolAllMarketForStore, this.getSupplyAmount error: %s", err)
		}
		borrowAmount, err := this.getBorrowAmount(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolAllMarketForStore, this.getSupplyAmount error: %s", err)
		}
		insuranceAmount, err := this.getInsuranceAmount(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolAllMarketForStore, this.getSupplyAmount error: %s", err)
		}
		price, err := this.AssetStoredPrice(this.cfg.OracleMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("FlashPoolAllMarketForStore, this.AssetStoredPrice error: %s", err)
		}

		supplyApy, err := this.getSupplyApy(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolAllMarketForStore, this.getSupplyApy error: %s", err)
		}
		borrowApy, err := this.getBorrowApy(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolAllMarketForStore, this.getBorrowApy error: %s", err)
		}
		insuranceApy, err := this.getInsuranceApy(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolAllMarketForStore, this.getInsuranceApy error: %s", err)
		}

		market := new(common.Market)
		market.Name = this.cfg.AssetMap[address.ToHexString()]
		market.Icon = this.cfg.IconMap[market.Name]

		// supplyAmount * price
		// borrowAmount * price
		// insuranceAmount * price
		market.TotalSupply = new(big.Int).Div(new(big.Int).Mul(supplyAmount, new(big.Int).SetUint64(price)), PriceDecimal).Uint64()
		market.TotalBorrow = new(big.Int).Div(new(big.Int).Mul(borrowAmount, new(big.Int).SetUint64(price)), PriceDecimal).Uint64()
		market.TotalInsurance = new(big.Int).Div(new(big.Int).Mul(insuranceAmount, new(big.Int).SetUint64(price)), PriceDecimal).Uint64()

		market.SupplyApy = supplyApy
		market.BorrowApy = borrowApy
		market.InsuranceApy = insuranceApy
		flashPoolAllMarket.FlashPoolAllMarket = append(flashPoolAllMarket.FlashPoolAllMarket, market)
	}
	return flashPoolAllMarket, nil
}

func (this *FlashPoolManager) UserFlashPoolOverview(accountStr string) (*common.UserFlashPoolOverview, error) {
	userFlashPoolOverview, err := this.store.LoadUserFlashPoolOverview(accountStr)
	if err != nil {
		userFlashPoolOverview := &common.UserFlashPoolOverview{
			CurrentSupply:    make([]*common.Supply, 0),
			CurrentBorrow:    make([]*common.Borrow, 0),
			CurrentInsurance: make([]*common.Insurance, 0),
			AllMarket:        make([]*common.UserMarket, 0),
		}
		allMarkets, err := this.GetAllMarkets()
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverview, this.GetAllMarkets error: %s", err)
		}
		for _, address := range allMarkets {
			market, err := this.store.LoadFlashMarket(this.cfg.AssetMap[address.ToHexString()])
			if err != nil {
				return nil, fmt.Errorf("FlashPoolAllMarket, this.store.LoadFlashMarket error: %s", err)
			}
			supplyApy := market.SupplyApy
			borrowApy := market.BorrowApy
			insuranceApy := market.InsuranceApy
			totalBorrowAmount := market.TotalBorrow
			totalInsuranceAmount := market.TotalInsurance

			userMarket := &common.UserMarket{
				Name:            this.cfg.AssetMap[address.ToHexString()],
				Icon:            this.cfg.IconMap[this.cfg.AssetMap[address.ToHexString()]],
				SupplyApy:       supplyApy,
				BorrowApy:       borrowApy,
				BorrowLiquidity: totalBorrowAmount,
				InsuranceApy:    insuranceApy,
				InsuranceAmount: totalInsuranceAmount,
			}
			userFlashPoolOverview.AllMarket = append(userFlashPoolOverview.AllMarket, userMarket)
		}

		return userFlashPoolOverview, nil
	}
	return userFlashPoolOverview, nil
}

func (this *FlashPoolManager) UserFlashPoolOverviewForStore(accountStr string) (*common.UserFlashPoolOverview, error) {
	account, err := ocommon.AddressFromBase58(accountStr)
	if err != nil {
		return nil, fmt.Errorf("UserFlashPoolOverviewForStore, ocommon.AddressFromBase58 error: %s", err)
	}
	allMarkets, err := this.GetAllMarkets()
	if err != nil {
		return nil, fmt.Errorf("UserFlashPoolOverviewForStore, this.GetAllMarkets error: %s", err)
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
			return nil, fmt.Errorf("UserFlashPoolOverviewForStore, this.getSupplyAmountByAccount error: %s", err)
		}
		price, err := this.AssetStoredPrice(this.cfg.OracleMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverviewForStore, this.AssetStoredPrice error: %s", err)
		}
		// borrowAmount * price
		totalBorrowBalance += new(big.Int).Div(new(big.Int).Mul(borrowAmount, new(big.Int).SetUint64(price)), PriceDecimal).Uint64()
	}
	userFlashPoolOverview.BorrowBalance = totalBorrowBalance
	netApy := new(big.Int).SetUint64(0)

	for _, address := range allMarkets {
		supplyAmount, err := this.getSupplyAmountByAccount(address, account)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverviewForStore, this.getSupplyAmountByAccount error: %s", err)
		}
		borrowAmount, err := this.getBorrowAmountByAccount(address, account)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverviewForStore, this.getSupplyAmountByAccount error: %s", err)
		}
		insuranceAmount, err := this.getInsuranceAmountByAccount(address, account)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverviewForStore, this.getSupplyAmountByAccount error: %s", err)
		}
		price, err := this.AssetStoredPrice(this.cfg.OracleMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverviewForStore, this.AssetStoredPrice error: %s", err)
		}
		marketMeta, err := this.getMarketMeta(address)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverviewForStore, this.getMarketMeta error: %s", err)
		}
		wingAccrued, err := this.getWingAccrued(address)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverviewForStore, this.getWingAccrued error: %s", err)
		}
		// supplyAmount * price
		// borrowAmount * price
		// insuranceAmount * price
		supplyAmountU64 := new(big.Int).Div(new(big.Int).Mul(supplyAmount, new(big.Int).SetUint64(price)), PriceDecimal).Uint64()
		borrowAmountU64 := new(big.Int).Div(new(big.Int).Mul(borrowAmount, new(big.Int).SetUint64(price)), PriceDecimal).Uint64()
		insuranceAmountU64 := new(big.Int).Div(new(big.Int).Mul(insuranceAmount, new(big.Int).SetUint64(price)), PriceDecimal).Uint64()
		userFlashPoolOverview.SupplyBalance += supplyAmountU64
		userFlashPoolOverview.InsuranceBalance += insuranceAmountU64
		userFlashPoolOverview.WingAccrued += wingAccrued.Uint64()
		supplyApy, err := this.getSupplyApy(address)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverviewForStore, this.getSupplyApy error: %s", err)
		}
		borrowApy, err := this.getBorrowApy(address)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverviewForStore, this.getBorrowApy error: %s", err)
		}
		insuranceApy, err := this.getInsuranceApy(address)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverviewForStore, this.getInsuranceApy error: %s", err)
		}
		a := new(big.Int).Mul(new(big.Int).SetUint64(supplyAmountU64), new(big.Int).SetUint64(supplyApy))
		b := new(big.Int).Mul(new(big.Int).SetUint64(insuranceAmountU64), new(big.Int).SetUint64(insuranceApy))
		c := new(big.Int).Mul(new(big.Int).SetUint64(borrowAmountU64), new(big.Int).SetUint64(borrowApy))
		netApy = new(big.Int).Add(netApy, new(big.Int).Sub(new(big.Int).Add(a, b), c))

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
				Name:             this.cfg.AssetMap[address.ToHexString()],
				Icon:             this.cfg.IconMap[this.cfg.AssetMap[address.ToHexString()]],
				SupplyDollar:     supplyAmountU64,
				SupplyBalance:    supplyAmount.Uint64(),
				Apy:              supplyApy,
				CollateralFactor: marketMeta.CollateralFactorMantissa.Uint64(),
				IfCollateral:     isAssetIn,
			}
			userFlashPoolOverview.CurrentSupply = append(userFlashPoolOverview.CurrentSupply, supply)
		}
		if borrowAmount.Uint64() != 0 {
			borrow := &common.Borrow{
				Name:             this.cfg.AssetMap[address.ToHexString()],
				Icon:             this.cfg.IconMap[this.cfg.AssetMap[address.ToHexString()]],
				BorrowDollar:     borrowAmountU64,
				BorrowBalance:    borrowAmount.Uint64(),
				Apy:              borrowApy,
				Limit:            borrowAmountU64 * PercentageDecimal / totalBorrowBalance,
				CollateralFactor: marketMeta.CollateralFactorMantissa.Uint64(),
			}
			userFlashPoolOverview.CurrentBorrow = append(userFlashPoolOverview.CurrentBorrow, borrow)
		}
		if insuranceAmount.Uint64() != 0 {
			insurance := &common.Insurance{
				Name:             this.cfg.AssetMap[address.ToHexString()],
				Icon:             this.cfg.IconMap[this.cfg.AssetMap[address.ToHexString()]],
				InsuranceDollar:  insuranceAmountU64,
				InsuranceBalance: insuranceAmount.Uint64(),
				Apy:              insuranceApy,
				CollateralFactor: marketMeta.CollateralFactorMantissa.Uint64(),
			}
			userFlashPoolOverview.CurrentInsurance = append(userFlashPoolOverview.CurrentInsurance, insurance)
		}

		totalBorrowAmount, err := this.getBorrowAmount(address)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverviewForStore, this.getSupplyAmount error: %s", err)
		}
		totalInsuranceAmount, err := this.getInsuranceAmount(address)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverviewForStore, this.getInsuranceAmount error: %s", err)
		}
		if supplyAmount.Uint64() == 0 && borrowAmount.Uint64() == 0 && insuranceAmount.Uint64() == 0 {
			userMarket := &common.UserMarket{
				Name:             this.cfg.AssetMap[address.ToHexString()],
				Icon:             this.cfg.IconMap[this.cfg.AssetMap[address.ToHexString()]],
				SupplyApy:        supplyApy,
				BorrowApy:        borrowApy,
				BorrowLiquidity:  totalBorrowAmount.Uint64(),
				InsuranceApy:     insuranceApy,
				InsuranceAmount:  totalInsuranceAmount.Uint64(),
				CollateralFactor: marketMeta.CollateralFactorMantissa.Uint64(),
				IfCollateral:     isAssetIn,
			}
			userFlashPoolOverview.AllMarket = append(userFlashPoolOverview.AllMarket, userMarket)
		}
	}
	accountLiquidity, err := this.getAccountLiquidity(account)
	if err != nil {
		return nil, fmt.Errorf("UserFlashPoolOverviewForStore, this.getAccountLiquidity error: %s", err)
	}
	userFlashPoolOverview.BorrowLimit = accountLiquidity.Liquidity.ToBigInt().Uint64()
	if userFlashPoolOverview.SupplyBalance+userFlashPoolOverview.BorrowBalance+userFlashPoolOverview.InsuranceBalance != 0 {
		userFlashPoolOverview.NetApy = new(big.Int).Div(netApy, new(big.Int).SetUint64(userFlashPoolOverview.SupplyBalance+
			userFlashPoolOverview.BorrowBalance+userFlashPoolOverview.InsuranceBalance)).Int64()
	}
	return userFlashPoolOverview, nil
}
