package flashpool

import (
	"fmt"
	"github.com/siovanus/wingServer/utils"
	"math"
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
		name := this.cfg.AssetMap[address.ToHexString()]
		market, err := this.store.LoadFlashMarket(this.cfg.AssetMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("FlashPoolDetail, this.store.LoadFlashMarket error: %s", err)
		}
		supplyAmount := market.TotalSupply
		borrowAmount := market.TotalBorrow
		insuranceAmount := market.TotalInsurance

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
				new(big.Int).SetUint64(distributedDay)), this.cfg.TokenDecimal[name]),
			SupplyAmount:    supplyAmount,
			BorrowAmount:    borrowAmount,
			InsuranceAmount: insuranceAmount,
			Total:           utils.ToStringByPrecise(totalDistribution, this.cfg.TokenDecimal[name]),
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
		supplyAmount := market.TotalSupply
		borrowAmount := market.TotalBorrow
		insuranceAmount := market.TotalInsurance

		totalDistribution, err := this.getTotalDistribution(address)
		if err != nil {
			return nil, fmt.Errorf("PoolDistribution, this.getTotalDistribution error: %s", err)
		}
		price, err := this.AssetStoredPrice(this.cfg.OracleMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("PoolDistribution, this.AssetStoredPrice error: %s", err)
		}
		// supplyAmount * price
		s = new(big.Int).Add(s, new(big.Int).Mul(utils.ToIntByPrecise(supplyAmount, this.cfg.TokenDecimal["oUSDT"]), price))
		// borrowAmount * price
		b = new(big.Int).Add(s, new(big.Int).Mul(utils.ToIntByPrecise(borrowAmount, this.cfg.TokenDecimal["oUSDT"]), price))
		// insuranceAmount * price
		i = new(big.Int).Add(s, new(big.Int).Mul(utils.ToIntByPrecise(insuranceAmount, this.cfg.TokenDecimal["oUSDT"]), price))
		d = new(big.Int).Add(d, totalDistribution)
	}
	distribution.Name = "Flash"
	distribution.Icon = this.cfg.IconMap[distribution.Name]
	distributedDay := new(big.Int).SetUint64((uint64(time.Now().Unix()) - governance.GenesisTime) / governance.DaySecond)
	distribution.SupplyAmount = utils.ToStringByPrecise(s, this.cfg.TokenDecimal["oUSDT"]+this.cfg.TokenDecimal["oracle"])
	distribution.BorrowAmount = utils.ToStringByPrecise(b, this.cfg.TokenDecimal["oUSDT"]+this.cfg.TokenDecimal["oracle"])
	distribution.InsuranceAmount = utils.ToStringByPrecise(i, this.cfg.TokenDecimal["oUSDT"]+this.cfg.TokenDecimal["oracle"])
	distribution.PerDay = utils.ToStringByPrecise(new(big.Int).Div(d, distributedDay), this.cfg.TokenDecimal["WING"])
	return distribution, nil
}

func (this *FlashPoolManager) FlashPoolBanner() (*common.FlashPoolBanner, error) {
	distributed := uint64(time.Now().Unix()) - governance.GenesisTime
	index := distributed/governance.YearSecond + 1

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
	today := governance.DailyDistibute[index]
	share := new(big.Int).SetUint64(0)
	if total.Uint64() != 0 {
		t := new(big.Int).Mul(new(big.Int).Mul(new(big.Int).SetUint64(today),
			new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal["WING"]))))),
			new(big.Int).SetUint64(this.cfg.TokenDecimal["percentage"]))
		share = new(big.Int).Div(t, total)
	}

	return &common.FlashPoolBanner{
		Today: new(big.Int).SetUint64(today).String(),
		Share: utils.ToStringByPrecise(share, this.cfg.TokenDecimal["percentage"]),
		Total: utils.ToStringByPrecise(total, this.cfg.TokenDecimal["WING"]),
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
	s := new(big.Int).SetUint64(0)
	b := new(big.Int).SetUint64(0)
	i := new(big.Int).SetUint64(0)
	for _, address := range allMarkets {
		market, err := this.store.LoadFlashMarket(this.cfg.AssetMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("FlashPoolDetail, this.store.LoadFlashMarket error: %s", err)
		}
		supplyAmount := market.TotalSupply
		borrowAmount := market.TotalBorrow
		insuranceAmount := market.TotalInsurance

		price, err := this.AssetStoredPrice(this.cfg.OracleMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("FlashPoolDetail, this.AssetStoredPrice error: %s", err)
		}
		// supplyAmount * price
		// borrowAmount * price
		// insuranceAmount * price
		supplyDollar := new(big.Int).Mul(utils.ToIntByPrecise(supplyAmount, this.cfg.TokenDecimal["oUSDT"]), price)
		borrowDollar := new(big.Int).Mul(utils.ToIntByPrecise(borrowAmount, this.cfg.TokenDecimal["oUSDT"]), price)
		insuranceDollar := new(big.Int).Mul(utils.ToIntByPrecise(insuranceAmount, this.cfg.TokenDecimal["oUSDT"]), price)
		s = new(big.Int).Add(s, supplyDollar)
		b = new(big.Int).Add(b, borrowDollar)
		i = new(big.Int).Add(i, insuranceDollar)

		name := this.cfg.AssetMap[address.ToHexString()]
		flashPoolDetail.SupplyMarketRank = append(flashPoolDetail.SupplyMarketRank, &common.MarketFund{
			Icon: this.cfg.IconMap[name],
			Name: name,
			Fund: utils.ToStringByPrecise(supplyDollar, this.cfg.TokenDecimal["oUSDT"]+this.cfg.TokenDecimal["oracle"]),
		})
		flashPoolDetail.BorrowMarketRank = append(flashPoolDetail.BorrowMarketRank, &common.MarketFund{
			Icon: this.cfg.IconMap[name],
			Name: name,
			Fund: utils.ToStringByPrecise(borrowDollar, this.cfg.TokenDecimal["oUSDT"]+this.cfg.TokenDecimal["oracle"]),
		})
		flashPoolDetail.InsuranceMarketRank = append(flashPoolDetail.InsuranceMarketRank, &common.MarketFund{
			Icon: this.cfg.IconMap[name],
			Name: name,
			Fund: utils.ToStringByPrecise(insuranceDollar, this.cfg.TokenDecimal["oUSDT"]+this.cfg.TokenDecimal["oracle"]),
		})
	}
	sort.SliceStable(flashPoolDetail.SupplyMarketRank, func(i, j int) bool {
		a := utils.ToIntByPrecise(flashPoolDetail.SupplyMarketRank[i].Fund, this.cfg.TokenDecimal["oUSDT"])
		b := utils.ToIntByPrecise(flashPoolDetail.SupplyMarketRank[j].Fund, this.cfg.TokenDecimal["oUSDT"])
		return a.Uint64() > b.Uint64()
	})
	sort.SliceStable(flashPoolDetail.BorrowMarketRank, func(i, j int) bool {
		a := utils.ToIntByPrecise(flashPoolDetail.BorrowMarketRank[i].Fund, this.cfg.TokenDecimal["oUSDT"])
		b := utils.ToIntByPrecise(flashPoolDetail.BorrowMarketRank[j].Fund, this.cfg.TokenDecimal["oUSDT"])
		return a.Uint64() > b.Uint64()
	})
	sort.SliceStable(flashPoolDetail.InsuranceMarketRank, func(i, j int) bool {
		a := utils.ToIntByPrecise(flashPoolDetail.InsuranceMarketRank[i].Fund, this.cfg.TokenDecimal["oUSDT"])
		b := utils.ToIntByPrecise(flashPoolDetail.InsuranceMarketRank[j].Fund, this.cfg.TokenDecimal["oUSDT"])
		return a.Uint64() > b.Uint64()
	})
	preFlashPoolDetailStore, err := this.store.LoadLatestFlashPoolDetail()
	if err != nil {
		return nil, fmt.Errorf("FlashPoolDetail, this.store.LoadLastestFlashPoolDetail error: %s", err)
	}
	flashPoolDetail.TotalSupply = utils.ToStringByPrecise(s, this.cfg.TokenDecimal["oUSDT"]+this.cfg.TokenDecimal["oracle"])
	flashPoolDetail.TotalBorrow = utils.ToStringByPrecise(b, this.cfg.TokenDecimal["oUSDT"]+this.cfg.TokenDecimal["oracle"])
	flashPoolDetail.TotalInsurance = utils.ToStringByPrecise(i, this.cfg.TokenDecimal["oUSDT"]+this.cfg.TokenDecimal["oracle"])

	flashPoolDetail.SupplyVolumeDaily = utils.ToStringByPrecise(new(big.Int).Sub(utils.ToIntByPrecise(flashPoolDetail.TotalSupply, this.cfg.TokenDecimal["oUSDT"]),
		utils.ToIntByPrecise(preFlashPoolDetailStore.TotalSupply, this.cfg.TokenDecimal["oUSDT"])), this.cfg.TokenDecimal["oUSDT"])
	flashPoolDetail.BorrowVolumeDaily = utils.ToStringByPrecise(new(big.Int).Sub(utils.ToIntByPrecise(flashPoolDetail.TotalBorrow, this.cfg.TokenDecimal["oUSDT"]),
		utils.ToIntByPrecise(preFlashPoolDetailStore.TotalBorrow, this.cfg.TokenDecimal["oUSDT"])), this.cfg.TokenDecimal["oUSDT"])
	flashPoolDetail.InsuranceVolumeDaily = utils.ToStringByPrecise(new(big.Int).Sub(utils.ToIntByPrecise(flashPoolDetail.TotalInsurance, this.cfg.TokenDecimal["oUSDT"]),
		utils.ToIntByPrecise(preFlashPoolDetailStore.TotalInsurance, this.cfg.TokenDecimal["oUSDT"])), this.cfg.TokenDecimal["oUSDT"])

	if utils.ToIntByPrecise(flashPoolDetail.TotalSupply, this.cfg.TokenDecimal["oUSDT"]).Uint64() != 0 {
		flashPoolDetail.TotalSupplyRate = utils.ToStringByPrecise(new(big.Int).Div(new(big.Int).Div(utils.ToIntByPrecise(flashPoolDetail.SupplyVolumeDaily,
			this.cfg.TokenDecimal["oUSDT"]), new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal["percentage"]))))),
			utils.ToIntByPrecise(flashPoolDetail.TotalSupply, this.cfg.TokenDecimal["oUSDT"])), this.cfg.TokenDecimal["percentage"])
	}
	if utils.ToIntByPrecise(flashPoolDetail.TotalBorrow, this.cfg.TokenDecimal["oUSDT"]).Uint64() != 0 {
		flashPoolDetail.TotalBorrowRate = utils.ToStringByPrecise(new(big.Int).Div(new(big.Int).Div(utils.ToIntByPrecise(flashPoolDetail.BorrowVolumeDaily,
			this.cfg.TokenDecimal["oUSDT"]), new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal["percentage"]))))),
			utils.ToIntByPrecise(flashPoolDetail.TotalBorrow, this.cfg.TokenDecimal["oUSDT"])), this.cfg.TokenDecimal["percentage"])
	}
	if utils.ToIntByPrecise(flashPoolDetail.TotalInsurance, this.cfg.TokenDecimal["oUSDT"]).Uint64() != 0 {
		flashPoolDetail.TotalInsuranceRate = utils.ToStringByPrecise(new(big.Int).Div(new(big.Int).Div(utils.ToIntByPrecise(flashPoolDetail.InsuranceVolumeDaily,
			this.cfg.TokenDecimal["oUSDT"]), new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal["percentage"]))))),
			utils.ToIntByPrecise(flashPoolDetail.TotalInsurance, this.cfg.TokenDecimal["oUSDT"])), this.cfg.TokenDecimal["percentage"])
	}
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
			return nil, fmt.Errorf("FlashPoolDetailForStore, this.AssetStoredPrice error: %s", err)
		}
		price := utils.ToIntByPrecise(priceStr, this.cfg.TokenDecimal["oracle"])
		// supplyAmount * price
		// borrowAmount * price
		// insuranceAmount * price
		supplyDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(supplyAmount, price),
			this.cfg.TokenDecimal[name]), this.cfg.TokenDecimal["oUSDT"])
		borrowDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(borrowAmount, price),
			this.cfg.TokenDecimal[name]), this.cfg.TokenDecimal["oUSDT"])
		insuranceDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(insuranceAmount, price),
			this.cfg.TokenDecimal[name]), this.cfg.TokenDecimal["oUSDT"])
		s = new(big.Int).Add(s, supplyDollar)
		b = new(big.Int).Add(b, borrowDollar)
		i = new(big.Int).Add(i, insuranceDollar)
	}
	flashPoolDetail.Timestamp = uint64(time.Now().Unix())
	flashPoolDetail.TotalSupply = utils.ToStringByPrecise(s, this.cfg.TokenDecimal["oUSDT"]+this.cfg.TokenDecimal["oracle"])
	flashPoolDetail.TotalBorrow = utils.ToStringByPrecise(b, this.cfg.TokenDecimal["oUSDT"]+this.cfg.TokenDecimal["oracle"])
	flashPoolDetail.TotalInsurance = utils.ToStringByPrecise(i, this.cfg.TokenDecimal["oUSDT"]+this.cfg.TokenDecimal["oracle"])
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
			this.cfg.TokenDecimal[name]), this.cfg.TokenDecimal["oUSDT"])
		borrowDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(borrowAmount, price),
			this.cfg.TokenDecimal[name]), this.cfg.TokenDecimal["oUSDT"])
		insuranceDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(insuranceAmount, price),
			this.cfg.TokenDecimal[name]), this.cfg.TokenDecimal["oUSDT"])
		flashPoolMarket.Name = name
		flashPoolMarket.TotalSupply = utils.ToStringByPrecise(supplyDollar, this.cfg.TokenDecimal["oUSDT"]+this.cfg.TokenDecimal["oracle"])
		flashPoolMarket.TotalBorrow = utils.ToStringByPrecise(borrowDollar, this.cfg.TokenDecimal["oUSDT"]+this.cfg.TokenDecimal["oracle"])
		flashPoolMarket.TotalInsurance = utils.ToStringByPrecise(insuranceDollar, this.cfg.TokenDecimal["oUSDT"]+this.cfg.TokenDecimal["oracle"])
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
		if utils.ToIntByPrecise(market.TotalSupply, this.cfg.TokenDecimal["oUSDT"]).Uint64() != 0 {
			market.TotalSupplyRate = utils.ToStringByPrecise(new(big.Int).Div(new(big.Int).Mul(new(big.Int).Sub(utils.ToIntByPrecise(market.TotalSupply,
				this.cfg.TokenDecimal["oUSDT"]), utils.ToIntByPrecise(latestFlashPoolMarket.TotalSupply,
				this.cfg.TokenDecimal["oUSDT"])), new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal["percentage"]))))),
				utils.ToIntByPrecise(market.TotalSupply, this.cfg.TokenDecimal["oUSDT"])), this.cfg.TokenDecimal["percentage"])
		}
		if utils.ToIntByPrecise(market.TotalBorrow, this.cfg.TokenDecimal["oUSDT"]).Uint64() != 0 {
			market.TotalBorrowRate = utils.ToStringByPrecise(new(big.Int).Div(new(big.Int).Mul(new(big.Int).Sub(utils.ToIntByPrecise(market.TotalBorrow,
				this.cfg.TokenDecimal["oUSDT"]), utils.ToIntByPrecise(latestFlashPoolMarket.TotalBorrow,
				this.cfg.TokenDecimal["oUSDT"])), new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal["percentage"]))))),
				utils.ToIntByPrecise(market.TotalBorrow, this.cfg.TokenDecimal["oUSDT"])), this.cfg.TokenDecimal["percentage"])
		}
		if utils.ToIntByPrecise(market.TotalInsurance, this.cfg.TokenDecimal["oUSDT"]).Uint64() != 0 {
			market.TotalInsuranceRate = utils.ToStringByPrecise(new(big.Int).Div(new(big.Int).Mul(new(big.Int).Sub(utils.ToIntByPrecise(market.TotalInsurance,
				this.cfg.TokenDecimal["oUSDT"]), utils.ToIntByPrecise(latestFlashPoolMarket.TotalInsurance,
				this.cfg.TokenDecimal["oUSDT"])), new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal["percentage"]))))),
				utils.ToIntByPrecise(market.TotalInsurance, this.cfg.TokenDecimal["oUSDT"])), this.cfg.TokenDecimal["percentage"])
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
		marketMeta, err := this.getMarketMeta(address)
		if err != nil {
			return nil, fmt.Errorf("FlashPoolAllMarketForStore, this.getMarketMeta error: %s", err)
		}

		market := new(common.Market)
		market.Name = this.cfg.AssetMap[address.ToHexString()]
		market.Icon = this.cfg.IconMap[market.Name]

		// supplyAmount * price
		// borrowAmount * price
		// insuranceAmount * price
		market.TotalSupply = utils.ToStringByPrecise(new(big.Int).Mul(supplyAmount, price),
			this.cfg.TokenDecimal[name]+this.cfg.TokenDecimal["oracle"])
		market.TotalBorrow = utils.ToStringByPrecise(new(big.Int).Mul(borrowAmount, price),
			this.cfg.TokenDecimal[name]+this.cfg.TokenDecimal["oracle"])
		market.TotalInsurance = utils.ToStringByPrecise(new(big.Int).Mul(insuranceAmount, price),
			this.cfg.TokenDecimal[name]+this.cfg.TokenDecimal["oracle"])
		market.CollateralFactor = utils.ToStringByPrecise(marketMeta.CollateralFactorMantissa, this.cfg.TokenDecimal["flash"])
		market.SupplyApy = utils.ToStringByPrecise(supplyApy, this.cfg.TokenDecimal["flash"])
		market.BorrowApy = utils.ToStringByPrecise(borrowApy, this.cfg.TokenDecimal["flash"])
		market.InsuranceApy = utils.ToStringByPrecise(insuranceApy, this.cfg.TokenDecimal["flash"])
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

			userMarket := &common.UserMarket{
				Name:                  this.cfg.AssetMap[address.ToHexString()],
				Icon:                  this.cfg.IconMap[this.cfg.AssetMap[address.ToHexString()]],
				SupplyApy:             market.SupplyApy,
				BorrowApy:             market.BorrowApy,
				BorrowLiquidity:       market.TotalBorrow,
				InsuranceApy:          market.InsuranceApy,
				InsuranceAmount:       market.TotalInsurance,
				CollateralFactor:      market.CollateralFactor,
				SupplyDistribution:    market.SupplyDistribution,
				BorrowDistribution:    market.BorrowDistribution,
				InsuranceDistribution: market.InsuranceDistribution,
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

	b := new(big.Int).SetUint64(0)
	for _, address := range allMarkets {
		name := this.cfg.AssetMap[address.ToHexString()]
		borrowAmount, err := this.getBorrowAmountByAccount(address, account)
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverviewForStore, this.getSupplyAmountByAccount error: %s", err)
		}
		price, err := this.AssetStoredPrice(this.cfg.OracleMap[address.ToHexString()])
		if err != nil {
			return nil, fmt.Errorf("UserFlashPoolOverviewForStore, this.AssetStoredPrice error: %s", err)
		}
		// borrowAmount * price
		b = new(big.Int).Add(b, new(big.Int).Mul(utils.ToIntByPrecise(
			utils.ToStringByPrecise(borrowAmount, this.cfg.TokenDecimal[name]), this.cfg.TokenDecimal["oUSDT"]), price))
	}
	userFlashPoolOverview.BorrowBalance = utils.ToStringByPrecise(b, this.cfg.TokenDecimal["oUSDT"]+this.cfg.TokenDecimal["oracle"])
	netApy := new(big.Int).SetUint64(0)

	s := new(big.Int).SetUint64(0)
	i := new(big.Int).SetUint64(0)
	w := new(big.Int).SetUint64(0)
	for _, address := range allMarkets {
		name := this.cfg.AssetMap[address.ToHexString()]
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
		supplyDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(supplyAmount, price),
			this.cfg.TokenDecimal[name]), this.cfg.TokenDecimal["oUSDT"])
		borrowDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(borrowAmount, price),
			this.cfg.TokenDecimal[name]), this.cfg.TokenDecimal["oUSDT"])
		insuranceDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(insuranceAmount, price),
			this.cfg.TokenDecimal[name]), this.cfg.TokenDecimal["oUSDT"])
		s = new(big.Int).Add(s, supplyDollar)
		i = new(big.Int).Add(i, insuranceDollar)
		w = new(big.Int).Add(w, wingAccrued)
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
		a := new(big.Int).Mul(supplyDollar, supplyApy)
		b := new(big.Int).Mul(insuranceDollar, insuranceApy)
		c := new(big.Int).Mul(borrowDollar, borrowApy)
		netApy = new(big.Int).Add(netApy, new(big.Int).Sub(new(big.Int).Add(a, b), c))

		isAssetIn := false
		for _, a := range assetsIn {
			if address == a {
				isAssetIn = true
				break
			}
		}

		if supplyAmount.Uint64() != 0 {
			supply := &common.Supply{
				Name:             this.cfg.AssetMap[address.ToHexString()],
				Icon:             this.cfg.IconMap[this.cfg.AssetMap[address.ToHexString()]],
				SupplyDollar:     utils.ToStringByPrecise(supplyDollar, this.cfg.TokenDecimal["oUSDT"]+this.cfg.TokenDecimal["oracle"]),
				SupplyBalance:    utils.ToStringByPrecise(supplyAmount, this.cfg.TokenDecimal[name]),
				Apy:              utils.ToStringByPrecise(supplyApy, this.cfg.TokenDecimal["flash"]),
				CollateralFactor: utils.ToStringByPrecise(marketMeta.CollateralFactorMantissa, this.cfg.TokenDecimal["flash"]),
				IfCollateral:     isAssetIn,
			}
			userFlashPoolOverview.CurrentSupply = append(userFlashPoolOverview.CurrentSupply, supply)
		}
		if borrowAmount.Uint64() != 0 {
			borrow := &common.Borrow{
				Name:             this.cfg.AssetMap[address.ToHexString()],
				Icon:             this.cfg.IconMap[this.cfg.AssetMap[address.ToHexString()]],
				BorrowDollar:     utils.ToStringByPrecise(borrowDollar, this.cfg.TokenDecimal["oUSDT"]+this.cfg.TokenDecimal["oracle"]),
				BorrowBalance:    utils.ToStringByPrecise(borrowAmount, this.cfg.TokenDecimal[name]),
				Apy:              utils.ToStringByPrecise(borrowApy, this.cfg.TokenDecimal["flash"]),
				CollateralFactor: utils.ToStringByPrecise(marketMeta.CollateralFactorMantissa, this.cfg.TokenDecimal["flash"]),
			}
			if b.Uint64() != 0 {
				borrow.Limit = utils.ToStringByPrecise(new(big.Int).Div(new(big.Int).Mul(borrowDollar, new(big.Int).SetUint64(
					this.cfg.TokenDecimal["percentage"])), b), this.cfg.TokenDecimal["percentage"])
			}
			userFlashPoolOverview.CurrentBorrow = append(userFlashPoolOverview.CurrentBorrow, borrow)
		}
		if insuranceAmount.Uint64() != 0 {
			insurance := &common.Insurance{
				Name:             this.cfg.AssetMap[address.ToHexString()],
				Icon:             this.cfg.IconMap[this.cfg.AssetMap[address.ToHexString()]],
				InsuranceDollar:  utils.ToStringByPrecise(insuranceDollar, this.cfg.TokenDecimal["oUSDT"]+this.cfg.TokenDecimal["oracle"]),
				InsuranceBalance: utils.ToStringByPrecise(insuranceAmount, this.cfg.TokenDecimal[name]),
				Apy:              utils.ToStringByPrecise(insuranceApy, this.cfg.TokenDecimal["flash"]),
				CollateralFactor: utils.ToStringByPrecise(marketMeta.CollateralFactorMantissa, this.cfg.TokenDecimal["flash"]),
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
				SupplyApy:        utils.ToStringByPrecise(supplyApy, this.cfg.TokenDecimal["flash"]),
				BorrowApy:        utils.ToStringByPrecise(borrowApy, this.cfg.TokenDecimal["flash"]),
				BorrowLiquidity:  utils.ToStringByPrecise(totalBorrowAmount, this.cfg.TokenDecimal[name]),
				InsuranceApy:     utils.ToStringByPrecise(insuranceApy, this.cfg.TokenDecimal["flash"]),
				InsuranceAmount:  utils.ToStringByPrecise(totalInsuranceAmount, this.cfg.TokenDecimal[name]),
				CollateralFactor: utils.ToStringByPrecise(marketMeta.CollateralFactorMantissa, this.cfg.TokenDecimal["flash"]),
				IfCollateral:     isAssetIn,
			}
			userFlashPoolOverview.AllMarket = append(userFlashPoolOverview.AllMarket, userMarket)
		}
	}
	userFlashPoolOverview.SupplyBalance = utils.ToStringByPrecise(s, this.cfg.TokenDecimal["oUSDT"]+this.cfg.TokenDecimal["oracle"])
	userFlashPoolOverview.InsuranceBalance = utils.ToStringByPrecise(i, this.cfg.TokenDecimal["oUSDT"]+this.cfg.TokenDecimal["oracle"])
	userFlashPoolOverview.WingAccrued = utils.ToStringByPrecise(w, this.cfg.TokenDecimal["WING"])
	accountLiquidity, err := this.getAccountLiquidity(account)
	if err != nil {
		return nil, fmt.Errorf("UserFlashPoolOverviewForStore, this.getAccountLiquidity error: %s", err)
	}
	userFlashPoolOverview.BorrowLimit = utils.ToStringByPrecise(accountLiquidity.Liquidity.ToBigInt(), this.cfg.TokenDecimal["flash"])
	total := new(big.Int).Add(new(big.Int).Add(s, b), i)
	if total.Uint64() != 0 {
		userFlashPoolOverview.NetApy = utils.ToStringByPrecise(new(big.Int).Div(netApy, total), this.cfg.TokenDecimal["flash"])
	}
	return userFlashPoolOverview, nil
}
