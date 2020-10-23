package flashpool

import (
	"fmt"
	"github.com/siovanus/wingServer/utils"
	"math"
	"math/big"
	"time"

	sdk "github.com/ontio/ontology-go-sdk"
	ocommon "github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/config"
	"github.com/siovanus/wingServer/http/common"
	"github.com/siovanus/wingServer/manager/governance"
	"github.com/siovanus/wingServer/store"
)

const (
	BlockPerYear = 60 * 60 * 24 * 365 * 2 / 5
)

var GAP = new(big.Int).SetUint64(198684465873214)

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

func (this *FlashPoolManager) AssetPrice(asset string) (string, error) {
	price, err := this.assetPrice(asset)
	if err != nil {
		return "", fmt.Errorf("AssetPrice, this.assetPrice error: %s", err)
	}
	return utils.ToStringByPrecise(price, this.cfg.TokenDecimal["oracle"]), nil
}

func (this *FlashPoolManager) AssetStoredPrice(asset string) (*big.Int, error) {
	if asset == "USDT" {
		return new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal["oracle"])))), nil
	}
	price, err := this.store.LoadPrice(asset)
	if err != nil {
		return nil, fmt.Errorf("AssetStoredPrice, this.store.LoadPrice error: %s", err)
	}
	return utils.ToIntByPrecise(price.Price, this.cfg.TokenDecimal["oracle"]), nil
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
			return nil, fmt.Errorf("FlashPoolMarketDistribution, this.store.LoadFlashMarket error: %s", err)
		}
		supplyAmount := market.TotalSupplyDollar
		borrowAmount := market.TotalBorrowDollar
		insuranceAmount := market.TotalInsuranceDollar

		totalDistribution, err := this.getTotalDistribution(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolMarketDistribution, this.getTotalDistribution error: %s", err)
		}
		distributedDay := (uint64(time.Now().Unix()) - governance.GenesisTime) / governance.DaySecond
		distribution := &common.Distribution{
			Icon: this.cfg.IconMap[this.cfg.AssetMap[address.ToHexString()]],
			Name: this.cfg.AssetMap[address.ToHexString()],
			// totalDistribution / distributedDay
			PerDay: utils.ToStringByPrecise(new(big.Int).Div(totalDistribution,
				new(big.Int).SetUint64(distributedDay)), this.cfg.TokenDecimal["WING"]),
			SupplyAmount:    supplyAmount,
			BorrowAmount:    borrowAmount,
			InsuranceAmount: insuranceAmount,
			Total:           utils.ToStringByPrecise(totalDistribution, this.cfg.TokenDecimal["WING"]),
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
	s := new(big.Int).SetUint64(0)
	b := new(big.Int).SetUint64(0)
	i := new(big.Int).SetUint64(0)
	d := new(big.Int).SetUint64(0)
	for _, address := range allMarkets {
		market, err := this.store.LoadFlashMarket(this.cfg.AssetMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("PoolDistribution, this.store.LoadFlashMarket error: %s", err)
		}
		supplyAmount := market.TotalSupplyDollar
		borrowAmount := market.TotalBorrowDollar
		insuranceAmount := market.TotalInsuranceDollar

		totalDistribution, err := this.getTotalDistribution(address)
		if err != nil {
			return nil, fmt.Errorf("PoolDistribution, this.getTotalDistribution error: %s", err)
		}

		// supplyAmount * price
		s = new(big.Int).Add(s, utils.ToIntByPrecise(supplyAmount, this.cfg.TokenDecimal["pUSDT"]))
		// borrowAmount * price
		b = new(big.Int).Add(b, utils.ToIntByPrecise(borrowAmount, this.cfg.TokenDecimal["pUSDT"]))
		// insuranceAmount * price
		i = new(big.Int).Add(i, utils.ToIntByPrecise(insuranceAmount, this.cfg.TokenDecimal["pUSDT"]))
		d = new(big.Int).Add(d, totalDistribution)
	}
	distribution.Name = "Flash"
	distribution.Icon = this.cfg.IconMap[distribution.Name]
	distributedDay := new(big.Int).SetUint64((uint64(time.Now().Unix()) - governance.GenesisTime) / governance.DaySecond)
	distribution.SupplyAmount = utils.ToStringByPrecise(s, this.cfg.TokenDecimal["pUSDT"])
	distribution.BorrowAmount = utils.ToStringByPrecise(b, this.cfg.TokenDecimal["pUSDT"])
	distribution.InsuranceAmount = utils.ToStringByPrecise(i, this.cfg.TokenDecimal["pUSDT"])
	distribution.PerDay = utils.ToStringByPrecise(new(big.Int).Div(d, distributedDay), this.cfg.TokenDecimal["WING"])
	distribution.Total = utils.ToStringByPrecise(d, this.cfg.TokenDecimal["WING"])
	return distribution, nil
}

func (this *FlashPoolManager) FlashPoolBanner() (*common.FlashPoolBanner, error) {
	gap := uint64(time.Now().Unix()) - governance.GenesisTime
	length := len(governance.DailyDistibute)
	epoch := []uint64{0}
	for i := 1; i < length+1; i++ {
		epoch = append(epoch, epoch[i-1]+governance.DistributeTime[i-1])
	}
	if gap > epoch[length] {
		gap = epoch[length]
	}
	index := 0
	for i := 0; i < len(epoch); i++ {
		if gap >= epoch[i] {
			index = i
		}
	}

	allMarkets, err := this.GetAllMarkets()
	if err != nil {
		return nil, fmt.Errorf("FlashPoolBanner, this.GetAllMarkets error: %s", err)
	}
	total := new(big.Int).SetUint64(0)
	for _, address := range allMarkets {
		totalDistribution, err := this.getTotalDistribution(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolBanner, this.getTotalDistribution error: %s", err)
		}
		total = new(big.Int).Add(total, totalDistribution)
	}
	today := governance.DailyDistibute[index] * governance.DaySecond
	share := new(big.Int).SetUint64(0)
	if total.Uint64() != 0 {
		t := new(big.Int).Mul(new(big.Int).Mul(new(big.Int).SetUint64(today),
			new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal["WING"]))))),
			new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal["percentage"])))))
		share = new(big.Int).Div(new(big.Int).Div(t, new(big.Int).SetUint64(100)), total)
	}

	return &common.FlashPoolBanner{
		Today: utils.ToStringByPrecise(new(big.Int).SetUint64(today), 2),
		Share: utils.ToStringByPrecise(share, this.cfg.TokenDecimal["percentage"]),
		Total: utils.ToStringByPrecise(total, this.cfg.TokenDecimal["WING"]),
	}, nil
}

func (this *FlashPoolManager) FlashPoolDetail() (*common.FlashPoolDetail, error) {
	allMarkets, err := this.GetAllMarkets()
	if err != nil {
		return nil, fmt.Errorf("FlashPoolDetail, this.GetAllMarkets error: %s", err)
	}
	flashPoolDetail := new(common.FlashPoolDetail)
	s := new(big.Int).SetUint64(0)
	b := new(big.Int).SetUint64(0)
	i := new(big.Int).SetUint64(0)
	for _, address := range allMarkets {
		market, err := this.store.LoadFlashMarket(this.cfg.AssetMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("FlashPoolDetail, this.store.LoadFlashMarket error: %s", err)
		}
		supplyAmount := market.TotalSupplyDollar
		borrowAmount := market.TotalBorrowDollar
		insuranceAmount := market.TotalInsuranceDollar

		// supplyAmount * price
		// borrowAmount * price
		// insuranceAmount * price
		supplyDollar := utils.ToIntByPrecise(supplyAmount, this.cfg.TokenDecimal["pUSDT"])
		borrowDollar := utils.ToIntByPrecise(borrowAmount, this.cfg.TokenDecimal["pUSDT"])
		insuranceDollar := utils.ToIntByPrecise(insuranceAmount, this.cfg.TokenDecimal["pUSDT"])
		s = new(big.Int).Add(s, supplyDollar)
		b = new(big.Int).Add(b, borrowDollar)
		i = new(big.Int).Add(i, insuranceDollar)
	}

	flashPoolDetail.TotalSupply = utils.ToStringByPrecise(s, this.cfg.TokenDecimal["pUSDT"])
	flashPoolDetail.TotalBorrow = utils.ToStringByPrecise(b, this.cfg.TokenDecimal["pUSDT"])
	flashPoolDetail.TotalInsurance = utils.ToStringByPrecise(i, this.cfg.TokenDecimal["pUSDT"])

	return flashPoolDetail, nil
}

func (this *FlashPoolManager) FlashPoolDetailForStore() (*store.FlashPoolDetail, error) {
	allMarkets, err := this.GetAllMarkets()
	if err != nil {
		return nil, fmt.Errorf("FlashPoolDetailForStore, this.GetAllMarkets error: %s", err)
	}
	flashPoolDetail := new(store.FlashPoolDetail)
	s := new(big.Int).SetUint64(0)
	b := new(big.Int).SetUint64(0)
	i := new(big.Int).SetUint64(0)
	for _, address := range allMarkets {
		name := this.cfg.AssetMap[address.ToHexString()]
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
		priceStr, err := this.AssetPrice(this.cfg.OracleMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("FlashPoolDetailForStore, this.AssetPrice error: %s", err)
		}
		price := utils.ToIntByPrecise(priceStr, this.cfg.TokenDecimal["oracle"])
		// supplyAmount * price
		// borrowAmount * price
		// insuranceAmount * price
		supplyDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(supplyAmount, price),
			this.cfg.TokenDecimal[name]), this.cfg.TokenDecimal["pUSDT"])
		borrowDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(borrowAmount, price),
			this.cfg.TokenDecimal[name]), this.cfg.TokenDecimal["pUSDT"])
		insuranceDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(insuranceAmount, price),
			this.cfg.TokenDecimal[name]), this.cfg.TokenDecimal["pUSDT"])
		s = new(big.Int).Add(s, supplyDollar)
		b = new(big.Int).Add(b, borrowDollar)
		i = new(big.Int).Add(i, insuranceDollar)
	}
	flashPoolDetail.Timestamp = uint64(time.Now().Unix())
	flashPoolDetail.TotalSupply = utils.ToStringByPrecise(s, this.cfg.TokenDecimal["pUSDT"]+this.cfg.TokenDecimal["oracle"])
	flashPoolDetail.TotalBorrow = utils.ToStringByPrecise(b, this.cfg.TokenDecimal["pUSDT"]+this.cfg.TokenDecimal["oracle"])
	flashPoolDetail.TotalInsurance = utils.ToStringByPrecise(i, this.cfg.TokenDecimal["pUSDT"]+this.cfg.TokenDecimal["oracle"])
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
		supplyDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(supplyAmount, price),
			this.cfg.TokenDecimal[name]), this.cfg.TokenDecimal["pUSDT"])
		borrowDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(borrowAmount, price),
			this.cfg.TokenDecimal[name]), this.cfg.TokenDecimal["pUSDT"])
		insuranceDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(insuranceAmount, price),
			this.cfg.TokenDecimal[name]), this.cfg.TokenDecimal["pUSDT"])
		flashPoolMarket.Name = name
		flashPoolMarket.TotalSupply = utils.ToStringByPrecise(supplyDollar, this.cfg.TokenDecimal["pUSDT"]+this.cfg.TokenDecimal["oracle"])
		flashPoolMarket.TotalBorrow = utils.ToStringByPrecise(borrowDollar, this.cfg.TokenDecimal["pUSDT"]+this.cfg.TokenDecimal["oracle"])
		flashPoolMarket.TotalInsurance = utils.ToStringByPrecise(insuranceDollar, this.cfg.TokenDecimal["pUSDT"]+this.cfg.TokenDecimal["oracle"])
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
		name := this.cfg.AssetMap[address.ToHexString()]
		supplyAmount, err := this.getSupplyAmount(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolAllMarketForStore, this.getSupplyAmount error: %s", err)
		}
		borrowAmount, err := this.getBorrowAmount(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolAllMarketForStore, this.getBorrowAmount error: %s", err)
		}
		insuranceAmount, err := this.getInsuranceAmount(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolAllMarketForStore, this.getInsuranceAmount error: %s", err)
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
		//insuranceApy, err := this.getInsuranceApy(address)
		//if err != nil {
		//	return nil, fmt.Errorf("FlashPoolAllMarketForStore, this.getInsuranceApy error: %s", err)
		//}
		marketMeta, err := this.getMarketMeta(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolAllMarketForStore, this.getMarketMeta error: %s", err)
		}
		exchangeRate, err := this.getExchangeRate(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolAllMarketForStore, this.getExchangeRate error: %s", err)
		}
		borrowIndex, err := this.getBorrowIndex(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolAllMarketForStore, this.getExchangeRate error: %s", err)
		}

		market := new(common.Market)
		market.Name = this.cfg.AssetMap[address.ToHexString()]
		market.Icon = this.cfg.IconMap[market.Name]

		market.TotalSupplyAmount = utils.ToStringByPrecise(supplyAmount, this.cfg.TokenDecimal[name])
		market.TotalBorrowAmount = utils.ToStringByPrecise(borrowAmount, this.cfg.TokenDecimal[name])
		market.TotalInsuranceAmount = utils.ToStringByPrecise(insuranceAmount, this.cfg.TokenDecimal[name])
		// supplyAmount * price
		// borrowAmount * price
		// insuranceAmount * price
		market.TotalSupplyDollar = utils.ToStringByPrecise(new(big.Int).Mul(supplyAmount, price),
			this.cfg.TokenDecimal[name]+this.cfg.TokenDecimal["oracle"])
		market.TotalBorrowDollar = utils.ToStringByPrecise(new(big.Int).Mul(borrowAmount, price),
			this.cfg.TokenDecimal[name]+this.cfg.TokenDecimal["oracle"])
		market.TotalInsuranceDollar = utils.ToStringByPrecise(new(big.Int).Mul(insuranceAmount, price),
			this.cfg.TokenDecimal[name]+this.cfg.TokenDecimal["oracle"])
		market.CollateralFactor = utils.ToStringByPrecise(marketMeta.CollateralFactorMantissa, this.cfg.TokenDecimal["flash"])
		market.SupplyApy = utils.ToStringByPrecise(supplyApy, this.cfg.TokenDecimal["flash"])
		market.BorrowApy = utils.ToStringByPrecise(borrowApy, this.cfg.TokenDecimal["flash"])
		market.ExchangeRate = utils.ToStringByPrecise(exchangeRate, 0)
		market.BorrowIndex = utils.ToStringByPrecise(borrowIndex, 0)
		//market.InsuranceApy = utils.ToStringByPrecise(insuranceApy, this.cfg.TokenDecimal["flash"])
		flashPoolAllMarket.FlashPoolAllMarket = append(flashPoolAllMarket.FlashPoolAllMarket, market)
	}
	return flashPoolAllMarket, nil
}

func (this *FlashPoolManager) UserFlashPoolOverview(accountStr string) (*common.UserFlashPoolOverview, error) {
	if accountStr != "" {
		return this.userFlashPoolOverview(accountStr)
	} else {
		return this.allMarket()
	}
}

func (this *FlashPoolManager) allMarket() (*common.UserFlashPoolOverview, error) {
	allMarkets, err := this.GetAllMarkets()
	if err != nil {
		return nil, fmt.Errorf("AllMarket, this.GetAllMarkets error: %s", err)
	}
	userFlashPoolOverview := &common.UserFlashPoolOverview{
		CurrentSupply:    make([]*common.Supply, 0),
		CurrentBorrow:    make([]*common.Borrow, 0),
		CurrentInsurance: make([]*common.Insurance, 0),
		AllMarket:        make([]*common.UserMarket, 0),
	}

	for _, address := range allMarkets {
		assetName := this.cfg.AssetMap[address.ToHexString()]
		market, err := this.store.LoadFlashMarket(assetName)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverview, this.store.LoadAssetApy error: %s", err)
		}

		borrowApy := utils.ToIntByPrecise(market.BorrowApy, this.cfg.TokenDecimal["flash"])
		insuranceApy := utils.ToIntByPrecise(market.InsuranceApy, this.cfg.TokenDecimal["flash"])
		borrowAmount := utils.ToIntByPrecise(market.TotalBorrowAmount, this.cfg.TokenDecimal[assetName])
		supplyAmount := utils.ToIntByPrecise(market.TotalSupplyAmount, this.cfg.TokenDecimal[assetName])
		supplyApy := new(big.Int)
		if borrowAmount.Uint64() != 0 {
			ratio := new(big.Int).Div(new(big.Int).Mul(borrowAmount,
				new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal["flash"]))))), supplyAmount)
			supplyApy = new(big.Int).Div(new(big.Int).Mul(new(big.Int).Mul(borrowApy, ratio),
				new(big.Int).SetUint64(85)), new(big.Int).SetUint64(100))
		}

		userMarket := &common.UserMarket{
			Name:      this.cfg.AssetMap[address.ToHexString()],
			Icon:      this.cfg.IconMap[this.cfg.AssetMap[address.ToHexString()]],
			SupplyApy: utils.ToStringByPrecise(supplyApy, 2*this.cfg.TokenDecimal["flash"]),
			BorrowApy: utils.ToStringByPrecise(borrowApy, this.cfg.TokenDecimal["flash"]),
			BorrowLiquidity: utils.ToStringByPrecise(new(big.Int).Sub(utils.ToIntByPrecise(market.TotalSupplyAmount,
				this.cfg.TokenDecimal[assetName]), utils.ToIntByPrecise(market.TotalBorrowAmount,
				this.cfg.TokenDecimal[assetName])), this.cfg.TokenDecimal[assetName]),
			InsuranceApy:     utils.ToStringByPrecise(insuranceApy, this.cfg.TokenDecimal["flash"]),
			InsuranceAmount:  market.TotalInsuranceAmount,
			CollateralFactor: market.CollateralFactor,
			IfCollateral:     false,
		}
		userFlashPoolOverview.AllMarket = append(userFlashPoolOverview.AllMarket, userMarket)
	}
	return userFlashPoolOverview, nil
}

func (this *FlashPoolManager) userFlashPoolOverview(accountStr string) (*common.UserFlashPoolOverview, error) {
	account, err := ocommon.AddressFromBase58(accountStr)
	if err != nil {
		return nil, fmt.Errorf("UserFlashPoolOverview, ocommon.AddressFromBase58 error: %s", err)
	}
	allMarkets, err := this.GetAllMarkets()
	if err != nil {
		return nil, fmt.Errorf("UserFlashPoolOverview, this.GetAllMarkets error: %s", err)
	}
	userFlashPoolOverview := &common.UserFlashPoolOverview{
		CurrentSupply:    make([]*common.Supply, 0),
		CurrentBorrow:    make([]*common.Borrow, 0),
		CurrentInsurance: make([]*common.Insurance, 0),
		AllMarket:        make([]*common.UserMarket, 0),
	}
	accountLiquidity, err := this.getAccountLiquidity(account)
	if err != nil {
		return nil, fmt.Errorf("UserFlashPoolOverview, this.getAccountLiquidity error: %s", err)
	}
	userFlashPoolOverview.BorrowLimit = utils.ToStringByPrecise(new(big.Int).Sub(accountLiquidity.Liquidity.ToBigInt(),
		accountLiquidity.Shortfall.ToBigInt()), this.cfg.TokenDecimal["oracle"])

	netApy := new(big.Int).SetUint64(0)
	s := new(big.Int).SetUint64(0)
	b := new(big.Int).SetUint64(0)
	i := new(big.Int).SetUint64(0)
	for _, address := range allMarkets {
		assetName := this.cfg.AssetMap[address.ToHexString()]
		price, err := this.AssetStoredPrice(this.cfg.OracleMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverview, this.AssetStoredPrice error: %s", err)
		}
		userAssetBalance, err := this.store.LoadUserAssetBalance(accountStr, assetName)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverview, this.store.LoadUserBalance error: %s", err)
		}
		flashMarket, err := this.store.LoadFlashMarket(assetName)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverview, this.store.LoadFlashMarket error: %s", err)
		}

		borrowBalance := new(big.Int).Div(new(big.Int).Mul(utils.ToIntByPrecise(userAssetBalance.BorrowAmount, 0),
			utils.ToIntByPrecise(flashMarket.BorrowIndex, 0)),
			utils.ToIntByPrecise(userAssetBalance.BorrowIndex, 0))
		borrowDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(borrowBalance, price),
			this.cfg.TokenDecimal[assetName]), this.cfg.TokenDecimal["pUSDT"])
		b = new(big.Int).Add(b, borrowDollar)
	}
	for _, address := range allMarkets {
		assetName := this.cfg.AssetMap[address.ToHexString()]
		market, err := this.store.LoadFlashMarket(assetName)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverview, this.store.LoadAssetApy error: %s", err)
		}
		price, err := this.AssetStoredPrice(this.cfg.OracleMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverview, this.AssetStoredPrice error: %s", err)
		}
		userAssetBalance, err := this.store.LoadUserAssetBalance(accountStr, assetName)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverview, this.store.LoadUserBalance error: %s", err)
		}
		flashMarket, err := this.store.LoadFlashMarket(assetName)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverview, this.store.LoadFlashMarket error: %s", err)
		}

		supplyAmount := new(big.Int).Mul(utils.ToIntByPrecise(userAssetBalance.FToken, 0),
			utils.ToIntByPrecise(flashMarket.ExchangeRate, 0))
		borrowAmount := new(big.Int).Div(new(big.Int).Mul(utils.ToIntByPrecise(userAssetBalance.BorrowAmount, 0),
			utils.ToIntByPrecise(flashMarket.BorrowIndex, 0)),
			utils.ToIntByPrecise(userAssetBalance.BorrowIndex, 0))
		insuranceAmount := new(big.Int).Mul(utils.ToIntByPrecise(userAssetBalance.FToken, 0),
			utils.ToIntByPrecise(flashMarket.ExchangeRate, 0))
		// supplyAmount * price
		// borrowAmount * price
		// insuranceAmount * price
		supplyDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(supplyAmount, price),
			this.cfg.TokenDecimal[assetName]+this.cfg.TokenDecimal["flash"]), this.cfg.TokenDecimal["pUSDT"])
		borrowDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(borrowAmount, price),
			this.cfg.TokenDecimal[assetName]), this.cfg.TokenDecimal["pUSDT"])
		insuranceDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(insuranceAmount, price),
			this.cfg.TokenDecimal[assetName]+this.cfg.TokenDecimal["flash"]), this.cfg.TokenDecimal["pUSDT"])
		s = new(big.Int).Add(s, supplyDollar)
		i = new(big.Int).Add(i, insuranceDollar)
		supplyApy := utils.ToIntByPrecise(market.SupplyApy, this.cfg.TokenDecimal["flash"])
		borrowApy := utils.ToIntByPrecise(market.BorrowApy, this.cfg.TokenDecimal["flash"])
		insuranceApy := utils.ToIntByPrecise(market.InsuranceApy, this.cfg.TokenDecimal["flash"])
		sa := new(big.Int).Mul(supplyDollar, supplyApy)
		ia := new(big.Int).Mul(insuranceDollar, insuranceApy)
		ba := new(big.Int).Mul(borrowDollar, borrowApy)
		netApy = new(big.Int).Add(netApy, new(big.Int).Sub(new(big.Int).Add(sa, ia), ba))

		claimWingAtMarket, err := this.getClaimWingAtMarket(account, []interface{}{address})
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverview, this.getClaimWingAtMarket account %s asset %s error: %s",
				account.ToBase58(), address.ToHexString(), err)
		}

		if supplyAmount.Uint64() != 0 {
			supply := &common.Supply{
				Name:             this.cfg.AssetMap[address.ToHexString()],
				Icon:             this.cfg.IconMap[this.cfg.AssetMap[address.ToHexString()]],
				SupplyBalance:    utils.ToStringByPrecise(supplyAmount, this.cfg.TokenDecimal[assetName]),
				Apy:              utils.ToStringByPrecise(supplyApy, this.cfg.TokenDecimal["flash"]),
				CollateralFactor: market.CollateralFactor,
				WingEarned:       utils.ToStringByPrecise(claimWingAtMarket, this.cfg.TokenDecimal["WING"]),
				IfCollateral:     userAssetBalance.IfCollateral,
			}
			userFlashPoolOverview.CurrentSupply = append(userFlashPoolOverview.CurrentSupply, supply)
		}
		if borrowAmount.Uint64() != 0 {
			borrow := &common.Borrow{
				Name:             this.cfg.AssetMap[address.ToHexString()],
				Icon:             this.cfg.IconMap[this.cfg.AssetMap[address.ToHexString()]],
				BorrowBalance:    utils.ToStringByPrecise(borrowAmount, this.cfg.TokenDecimal[assetName]),
				Apy:              utils.ToStringByPrecise(borrowApy, this.cfg.TokenDecimal["flash"]),
				WingEarned:       utils.ToStringByPrecise(claimWingAtMarket, this.cfg.TokenDecimal["WING"]),
				CollateralFactor: market.CollateralFactor,
			}
			if accountLiquidity.Liquidity.ToBigInt().Uint64() != 0 {
				borrowLimit := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Sub(accountLiquidity.Liquidity.ToBigInt(),
					accountLiquidity.Shortfall.ToBigInt()), this.cfg.TokenDecimal["oracle"]), this.cfg.TokenDecimal["oracle"]+this.cfg.TokenDecimal["pUSDT"])
				totalLimit := new(big.Int).Add(borrowLimit, b)
				borrow.Limit = utils.ToStringByPrecise(new(big.Int).Div(new(big.Int).Mul(borrowDollar, new(big.Int).SetUint64(
					uint64(math.Pow10(int(this.cfg.TokenDecimal["percentage"]))))), totalLimit), this.cfg.TokenDecimal["percentage"])
			}
			userFlashPoolOverview.CurrentBorrow = append(userFlashPoolOverview.CurrentBorrow, borrow)
		}
		if insuranceAmount.Uint64() != 0 {
			insurance := &common.Insurance{
				Name:             this.cfg.AssetMap[address.ToHexString()],
				Icon:             this.cfg.IconMap[this.cfg.AssetMap[address.ToHexString()]],
				InsuranceBalance: utils.ToStringByPrecise(insuranceAmount, this.cfg.TokenDecimal[assetName]),
				Apy:              utils.ToStringByPrecise(insuranceApy, this.cfg.TokenDecimal["flash"]),
				WingEarned:       utils.ToStringByPrecise(claimWingAtMarket, this.cfg.TokenDecimal["WING"]),
				CollateralFactor: market.CollateralFactor,
			}
			userFlashPoolOverview.CurrentInsurance = append(userFlashPoolOverview.CurrentInsurance, insurance)
		}

		userMarket := &common.UserMarket{
			Name:      this.cfg.AssetMap[address.ToHexString()],
			Icon:      this.cfg.IconMap[this.cfg.AssetMap[address.ToHexString()]],
			SupplyApy: utils.ToStringByPrecise(supplyApy, this.cfg.TokenDecimal["flash"]),
			BorrowApy: utils.ToStringByPrecise(borrowApy, this.cfg.TokenDecimal["flash"]),
			BorrowLiquidity: utils.ToStringByPrecise(new(big.Int).Sub(utils.ToIntByPrecise(market.TotalSupplyAmount,
				this.cfg.TokenDecimal[assetName]), utils.ToIntByPrecise(market.TotalBorrowAmount,
				this.cfg.TokenDecimal[assetName])), this.cfg.TokenDecimal[assetName]),
			InsuranceApy:     utils.ToStringByPrecise(insuranceApy, this.cfg.TokenDecimal["flash"]),
			InsuranceAmount:  market.TotalInsuranceAmount,
			CollateralFactor: market.CollateralFactor,
			IfCollateral:     userAssetBalance.IfCollateral,
		}
		userFlashPoolOverview.AllMarket = append(userFlashPoolOverview.AllMarket, userMarket)
	}
	total := new(big.Int).Add(s, i)
	if total.Uint64() != 0 {
		userFlashPoolOverview.NetApy = utils.ToStringByPrecise(new(big.Int).Div(netApy, total), this.cfg.TokenDecimal["flash"])
	}
	return userFlashPoolOverview, nil
}

func (this *FlashPoolManager) UserBalanceForStore(accountStr, borrowAmount, borrowIndex string) error {
	account, err := ocommon.AddressFromBase58(accountStr)
	if err != nil {
		return fmt.Errorf("UserBalanceForStore, ocommon.AddressFromBase58 error: %s", err)
	}
	allMarkets, err := this.GetAllMarkets()
	if err != nil {
		return fmt.Errorf("UserBalanceForStore, this.GetAllMarkets error: %s", err)
	}
	assetsIn, _ := this.getAssetsIn(account)
	for _, address := range allMarkets {
		fToken, err := this.getFTokenAmount(address, account)
		if err != nil {
			return fmt.Errorf("UserBalanceForStore, this.getSupplyAmountByAccount error: %s", err)
		}
		iToken, err := this.getITokenAmount(address, account)
		if err != nil {
			return fmt.Errorf("UserBalanceForStore, this.getInsuranceAmountByAccount error: %s", err)
		}
		name := this.cfg.AssetMap[address.ToHexString()]
		isAssetIn := false
		for _, a := range assetsIn {
			if address == a {
				isAssetIn = true
				break
			}
		}
		if borrowAmount == "" && borrowIndex == "" {
			userAssetBalance, err := this.store.LoadUserAssetBalance(accountStr, address.ToHexString())
			if err != nil {
				return fmt.Errorf("UserBalanceForStore, this.store.LoadUserAssetBalance error: %s", err)
			}
			borrowAmount = userAssetBalance.BorrowAmount
			borrowIndex = userAssetBalance.BorrowIndex
		}
		userBalance := &store.UserAssetBalance{
			UserAddress:  accountStr,
			AssetAddress: address.ToHexString(),
			AssetName:    name,
			Icon:         this.cfg.IconMap[name],
			FToken:       utils.ToStringByPrecise(fToken, 0),
			BorrowAmount: borrowAmount,
			BorrowIndex:  borrowIndex,
			Itoken:       utils.ToStringByPrecise(iToken, 0),
			IfCollateral: isAssetIn,
		}
		err = this.store.SaveUserAssetBalance(userBalance)
		if err != nil {
			return fmt.Errorf("UserBalanceForStore, this.store.SaveUserAssetBalance error: %s", err)
		}
	}
	return nil
}

func (this *FlashPoolManager) ClaimWing(accountStr string) (string, error) {
	account, err := ocommon.AddressFromBase58(accountStr)
	if err != nil {
		return "", fmt.Errorf("ClaimWing, ocommon.AddressFromBase58 error: %s", err)
	}
	amount, err := this.getClaimWing(account)
	if err != nil {
		return "", fmt.Errorf("ClaimWing, this.getClaimWing error: %s", err)
	}
	return utils.ToStringByPrecise(amount, this.cfg.TokenDecimal["WING"]), nil
}

func (this *FlashPoolManager) BorrowAddressList() ([]store.UserAssetBalance, error) {
	borrowUsers, err := this.store.LoadBorrowUsers()
	if err != nil {
		return nil, fmt.Errorf("BorrowAddressList, this.store.LoadBorrowUsers error: %s", err)
	}
	return borrowUsers, nil
}

func (this *FlashPoolManager) LiquidationList(accountStr string) ([]*common.Liquidation, error) {
	account, err := ocommon.AddressFromBase58(accountStr)
	if err != nil {
		return nil, fmt.Errorf("LiquidationList, ocommon.AddressFromBase58 error: %s", err)
	}
	userBalance, err := this.store.LoadUserBalance(accountStr)
	if err != nil {
		return nil, fmt.Errorf("LiquidationList, this.store.LoadUserBalance error: %s", err)
	}
	liquidationList := make([]*common.Liquidation, 0)
	totalBorrowDollar := new(big.Int)
	totalCollateralDollar := new(big.Int)
	collateralAssets := make([]*common.CollateralAsset, 0)
	for _, v := range userBalance {
		price, err := this.AssetStoredPrice(this.cfg.OracleMap[v.AssetAddress])
		if err != nil {
			return nil, fmt.Errorf("LiquidationList, this.AssetStoredPrice error: %s", err)
		}
		flashMarket, err := this.store.LoadFlashMarket(v.AssetName)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverview, this.store.LoadFlashMarket error: %s", err)
		}
		borrowBalance := new(big.Int).Div(new(big.Int).Mul(utils.ToIntByPrecise(v.BorrowAmount, 0),
			utils.ToIntByPrecise(flashMarket.BorrowIndex, 0)),
			utils.ToIntByPrecise(v.BorrowIndex, 0))
		borrowDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(borrowBalance, price),
			this.cfg.TokenDecimal[v.AssetName]), this.cfg.TokenDecimal["pETH"])
		totalBorrowDollar = new(big.Int).Add(totalBorrowDollar, borrowDollar)
		supplyBalance := new(big.Int).Mul(utils.ToIntByPrecise(v.FToken, 0),
			utils.ToIntByPrecise(flashMarket.ExchangeRate, 0))
		if v.AssetName != "WING" && v.IfCollateral && supplyBalance.Uint64() != 0 {
			supplyDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(supplyBalance, price),
				this.cfg.TokenDecimal[v.AssetName]+this.cfg.TokenDecimal["flash"]), this.cfg.TokenDecimal["pETH"])
		if v.IfCollateral && v.SupplyBalance != "0" {
			supplyDollar := new(big.Int).Mul(utils.ToIntByPrecise(v.SupplyBalance, this.cfg.TokenDecimal["pETH"]), price)
			totalCollateralDollar = new(big.Int).Add(totalCollateralDollar, supplyDollar)
			collateralAsset := &common.CollateralAsset{
				Icon:    v.Icon,
				Name:    v.AssetName,
				Balance: utils.ToStringByPrecise(supplyBalance, this.cfg.TokenDecimal[v.AssetName]+this.cfg.TokenDecimal["flash"]),
				Dollar:  utils.ToStringByPrecise(supplyDollar, this.cfg.TokenDecimal["pETH"]+this.cfg.TokenDecimal["oracle"]),
			}
			collateralAssets = append(collateralAssets, collateralAsset)
		}

		if borrowBalance.Uint64() != 0 {
			price, err := this.AssetStoredPrice(this.cfg.OracleMap[v.AssetAddress])
			if err != nil {
				return nil, fmt.Errorf("LiquidationList, this.AssetStoredPrice error: %s", err)
			}
			liquidation := &common.Liquidation{
				Icon:             v.Icon,
				Name:             v.AssetName,
				BorrowBalance:    utils.ToStringByPrecise(borrowBalance, this.cfg.TokenDecimal[v.AssetName]),
				CollateralAssets: collateralAssets,
			}
			liquidation.BorrowDollar = utils.ToStringByPrecise(new(big.Int).Mul(borrowBalance, price),
				this.cfg.TokenDecimal[v.AssetName]+this.cfg.TokenDecimal["oracle"])
			accountLiquidity, err := this.getAccountLiquidity(account)
			if err != nil {
				return nil, fmt.Errorf("LiquidationList, this.getAccountLiquidity error: %s", err)
			}
			liquidity := new(big.Int).Sub(accountLiquidity.Liquidity.ToBigInt(), accountLiquidity.Shortfall.ToBigInt())
			limit := new(big.Int).Add(utils.ToIntByPrecise(utils.ToStringByPrecise(liquidity, this.cfg.TokenDecimal["oracle"]),
				this.cfg.TokenDecimal["pETH"]+this.cfg.TokenDecimal["oracle"]), totalBorrowDollar)
			if limit.Uint64() != 0 {
				liquidation.BorrowLimitUsed = utils.ToStringByPrecise(new(big.Int).Div(new(big.Int).Mul(totalBorrowDollar,
					new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal["percentage"]))))), limit),
					this.cfg.TokenDecimal["percentage"])
			} else {
				liquidation.BorrowLimitUsed = "10"
			}
			liquidation.CollateralDollar = utils.ToStringByPrecise(totalCollateralDollar, this.cfg.TokenDecimal["pETH"]+this.cfg.TokenDecimal["oracle"])
			liquidationList = append(liquidationList, liquidation)
		}
	}
	return liquidationList, nil
}

func (this *FlashPoolManager) WingApyForStore() error {
	allMarkets, err := this.GetAllMarkets()
	if err != nil {
		return fmt.Errorf("WingApy, this.GetAllMarkets error: %s", err)
	}
	for _, address := range allMarkets {
		wingSpeeds, err := this.getWingSpeeds(address)
		if err != nil {
			return fmt.Errorf("WingApy, this.getWingSpeeds error: %s", err)
		}
		wingSBIPortion, err := this.getWingSBIPortion(address)
		if err != nil {
			return fmt.Errorf("WingApy, this.getWingSBIPortion error: %s", err)
		}
		totalPortion := new(big.Int).Add(wingSBIPortion.InsurancePortion.ToBigInt(),
			new(big.Int).Add(wingSBIPortion.SupplyPortion.ToBigInt(), wingSBIPortion.BorrowPortion.ToBigInt()))
		price, err := this.AssetStoredPrice("WING")
		if err != nil {
			return fmt.Errorf("WingApy, this.AssetStoredPrice error: %s", err)
		}
		market, err := this.store.LoadFlashMarket(this.cfg.AssetMap[address.ToHexString()])
		if err != nil {
			return fmt.Errorf("WingApy, this.store.LoadFlashMarket error: %s", err)
		}
		totalSupplyDollar := utils.ToIntByPrecise(market.TotalSupplyDollar, this.cfg.TokenDecimal["pUSDT"])
		totalBorrowDollar := utils.ToIntByPrecise(market.TotalBorrowDollar, this.cfg.TokenDecimal["pUSDT"])
		totalInsuranceDollar := utils.ToIntByPrecise(market.TotalInsuranceDollar, this.cfg.TokenDecimal["pUSDT"])
		var supplyApy, borrowApy, insuranceApy string
		if totalSupplyDollar.Uint64() != 0 {
			supplyApy = utils.ToStringByPrecise(new(big.Int).Div(new(big.Int).Div(new(big.Int).Mul(new(big.Int).Mul(new(big.Int).Mul(new(big.Int).Mul(wingSpeeds,
				wingSBIPortion.SupplyPortion.ToBigInt()), price), new(big.Int).SetUint64(governance.YearSecond)),
				new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal["pUSDT"]))))), totalPortion),
				totalSupplyDollar), this.cfg.TokenDecimal["oracle"]+this.cfg.TokenDecimal["WING"])
		}
		if totalBorrowDollar.Uint64() != 0 {
			borrowApy = utils.ToStringByPrecise(new(big.Int).Div(new(big.Int).Div(new(big.Int).Mul(new(big.Int).Mul(new(big.Int).Mul(new(big.Int).Mul(wingSpeeds,
				wingSBIPortion.BorrowPortion.ToBigInt()), price), new(big.Int).SetUint64(governance.YearSecond)),
				new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal["pUSDT"]))))), totalPortion),
				totalBorrowDollar), this.cfg.TokenDecimal["oracle"]+this.cfg.TokenDecimal["WING"])
		}
		if totalInsuranceDollar.Uint64() != 0 {
			insuranceApy = utils.ToStringByPrecise(new(big.Int).Div(new(big.Int).Div(new(big.Int).Mul(new(big.Int).Mul(new(big.Int).Mul(new(big.Int).Mul(wingSpeeds,
				wingSBIPortion.InsurancePortion.ToBigInt()), price), new(big.Int).SetUint64(governance.YearSecond)),
				new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal["pUSDT"]))))), totalPortion),
				totalInsuranceDollar), this.cfg.TokenDecimal["oracle"]+this.cfg.TokenDecimal["WING"])
		}
		wingApy := &common.WingApy{
			AssetName:    this.cfg.AssetMap[address.ToHexString()],
			SupplyApy:    supplyApy,
			BorrowApy:    borrowApy,
			InsuranceApy: insuranceApy,
		}
		err = this.store.SaveWingApy(wingApy)
		if err != nil {
			return fmt.Errorf("WingApy, this.store.SaveWingApy error: %s", err)
		}
	}
	return nil
}

func (this *FlashPoolManager) WingApys() ([]common.WingApy, error) {
	wingApys, err := this.store.LoadWingApys()
	if err != nil {
		return nil, fmt.Errorf("WingApy, this.store.LoadWingApys error: %s", err)
	}
	return wingApys, nil
}

func (this *FlashPoolManager) Reserves() (*common.Reserves, error) {
	allMarkets, err := this.GetAllMarkets()
	if err != nil {
		return nil, fmt.Errorf("TotalReserve, this.GetAllMarkets error: %s", err)
	}
	totalReserve := new(big.Int)
	reserves := &common.Reserves{
		AssetReserve: make([]*common.Reserve, 0),
	}
	for _, address := range allMarkets {
		name := this.cfg.AssetMap[address.ToHexString()]
		price, err := this.AssetStoredPrice(this.cfg.OracleMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("TotalReserve, this.AssetStoredPrice error: %s", err)
		}
		reserveBalance, err := this.getTotalReserves(address)
		if err != nil {
			return nil, fmt.Errorf("TotalReserve, this.getTotalReserves error: %s", err)
		}
		reserveBalanceStr := utils.ToStringByPrecise(reserveBalance, this.cfg.TokenDecimal[name])
		reserveDollarStr := utils.ToStringByPrecise(new(big.Int).Mul(price, reserveBalance),
			this.cfg.TokenDecimal[name]+this.cfg.TokenDecimal["oracle"])
		assetReserve := &common.Reserve{
			Name:           name,
			Icon:           this.cfg.IconMap[name],
			ReserveFactor:  "0.15",
			ReserveBalance: reserveBalanceStr,
			ReserveDollar:  reserveDollarStr,
		}
		reserves.AssetReserve = append(reserves.AssetReserve, assetReserve)

		delta := utils.ToIntByPrecise(reserveDollarStr, this.cfg.TokenDecimal["pUSDT"])
		totalReserve = new(big.Int).Add(totalReserve, delta)
	}
	reserves.TotalReserve = utils.ToStringByPrecise(totalReserve, this.cfg.TokenDecimal["pUSDT"])
	return reserves, nil
}
