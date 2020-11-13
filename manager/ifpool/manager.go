package ifpool

import (
	"fmt"
	"github.com/siovanus/wingServer/manager/governance"
	"github.com/siovanus/wingServer/utils"
	oscore_oracle "github.com/wing-groups/wing-contract-tools/contracts/oscore-oracle"
	"math"
	"math/big"
	"os"
	"time"

	ocommon "github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/config"
	"github.com/siovanus/wingServer/http/common"
	"github.com/siovanus/wingServer/log"
	"github.com/siovanus/wingServer/store"
	if_borrow "github.com/wing-groups/wing-contract-tools/contracts/if-borrow"
	if_ctrl "github.com/wing-groups/wing-contract-tools/contracts/if-ctrl"
	"github.com/wing-groups/wing-contract-tools/contracts/iftoken"
	"github.com/wing-groups/wing-contract-tools/contracts/iitoken"
)

var GenesisTime = time.Date(2020, time.September, 12, 0, 0, 0, 0, time.UTC).Unix()

const MaxLevel = 0x03

type IFPoolManager struct {
	cfg          *config.Config
	store        *store.Client
	Comptroller  *if_ctrl.Comptroller
	FTokenMap    map[ocommon.Address]*iftoken.IFToken
	ITokenMap    map[ocommon.Address]*iitoken.IIToken
	BorrowMap    map[ocommon.Address]*if_borrow.IfBorrowPool
	OscoreOracle *oscore_oracle.Oracle
}

func NewIFPoolManager(contractAddress, oscoreOracleAddress ocommon.Address, store *store.Client,
	cfg *config.Config) *IFPoolManager {
	comptroller, _ := if_ctrl.NewComptroller(cfg.JsonRpcAddress, contractAddress.ToHexString(), nil,
		2500, 20000)
	oracle, _ := oscore_oracle.NewOracle(cfg.JsonRpcAddress, oscoreOracleAddress.ToHexString(), nil,
		2500, 20000)
	fTokenMap := make(map[ocommon.Address]*iftoken.IFToken)
	iTokenMap := make(map[ocommon.Address]*iitoken.IIToken)
	borrowPoolMap := make(map[ocommon.Address]*if_borrow.IfBorrowPool)
	allMarket, err := comptroller.AllMarkets()
	if err != nil {
		log.Errorf("NewFlashPoolManager, comptroller.AllMarkets error: %s", err)
		os.Exit(1)
	}
	for _, name := range allMarket {
		marketInfo, err := comptroller.MarketInfo(name)
		if err != nil {
			log.Errorf("NewFlashPoolManager, comptroller.MarketInfo error: %s", err)
			os.Exit(1)
		}
		iFToken, _ := iftoken.NewIFToken(cfg.JsonRpcAddress, marketInfo.SupplyPool.ToHexString(), nil,
			2500, 20000)
		iIToken, _ := iitoken.NewIIToken(cfg.JsonRpcAddress, marketInfo.InsurancePool.ToHexString(), nil,
			2500, 20000)
		borrowPool, _ := if_borrow.NewIfBorrowPool(cfg.JsonRpcAddress, marketInfo.BorrowPool.ToHexString(), nil,
			2500, 20000)
		fTokenMap[marketInfo.SupplyPool] = iFToken
		iTokenMap[marketInfo.InsurancePool] = iIToken
		borrowPoolMap[marketInfo.BorrowPool] = borrowPool
	}

	manager := &IFPoolManager{
		cfg:          cfg,
		store:        store,
		Comptroller:  comptroller,
		FTokenMap:    fTokenMap,
		ITokenMap:    iTokenMap,
		BorrowMap:    borrowPoolMap,
		OscoreOracle: oracle,
	}

	return manager
}

func (this *IFPoolManager) StoreIFInfo() error {
	ifInfo := new(store.IFInfo)
	capacity, err := this.Comptroller.MaxSupplyValue()
	if err != nil {
		return fmt.Errorf("StoreIFInfo, this.Comptroller.MaxSupplyValue error: %s", err)
	}
	ifInfo.Cap = utils.ToStringByPrecise(capacity, this.cfg.TokenDecimal["oracle"])
	total, err := this.Comptroller.TotalSupplyValue()
	if err != nil {
		return fmt.Errorf("StoreIFInfo, this.Comptroller.TotalSupplyValue error: %s", err)
	}
	ifInfo.Total = utils.ToStringByPrecise(total, this.cfg.TokenDecimal["oracle"])
	err = this.store.SaveIFInfo(ifInfo)
	if err != nil {
		return fmt.Errorf("StoreIFInfo, this.store.SaveIFInfo error: %s", err)
	}
	return nil
}

func (this *IFPoolManager) StoreIFMarketInfo() error {
	allMarket, err := this.Comptroller.AllMarkets()
	if err != nil {
		return fmt.Errorf("StoreIFMarketInfo, this.Comptroller.AllMarkets error: %s", err)
	}
	for _, name := range allMarket {
		ifMarketInfo := new(store.IFMarketInfo)
		marketInfo, err := this.Comptroller.MarketInfo(name)
		if err != nil {
			return fmt.Errorf("StoreIFMarketInfo, this.Comptroller.MarketInfo error: %s", err)
		}
		ifMarketInfo.Name = name
		totalCash, err := this.FTokenMap[marketInfo.SupplyPool].TotalCash()
		if err != nil {
			return fmt.Errorf("StoreIFMarketInfo, this.FTokenMap[marketInfo.SupplyPool].TotalCash error: %s", err)
		}
		totalDebt, err := this.FTokenMap[marketInfo.SupplyPool].TotalDebt()
		if err != nil {
			return fmt.Errorf("StoreIFMarketInfo, this.FTokenMap[marketInfo.SupplyPool].TotalDebt error: %s", err)
		}
		interestIndex, err := this.Comptroller.InterestIndex(name)
		if err != nil {
			return fmt.Errorf("StoreIFMarketInfo, this.Comptroller.InterestIndex error: %s", err)
		}
		ifMarketInfo.TotalCash = utils.ToStringByPrecise(totalCash, 0)
		ifMarketInfo.TotalDebt = utils.ToStringByPrecise(totalDebt, 0)
		ifMarketInfo.InterestIndex = utils.ToStringByPrecise(interestIndex, 0)

		oscoreInfo, err := this.BorrowMap[marketInfo.BorrowPool].GetOscoreInfoByLevel(MaxLevel)
		if err != nil {
			return fmt.Errorf("StoreIFMarketInfo, this.BorrowMap[marketInfo.BorrowPool].GetOscoreInfoByLevel error: %s", err)
		}
		ifMarketInfo.InterestRate = oscoreInfo.InterestRate
		ifMarketInfo.CollateralFactor = oscoreInfo.CollateralFactor
		totalInsurance, err := this.ITokenMap[marketInfo.InsurancePool].TotalCash()
		if err != nil {
			return fmt.Errorf("StoreIFMarketInfo, this.ITokenMap[marketInfo.InsurancePool].TotalCash error: %s", err)
		}
		ifMarketInfo.TotalInsurance = utils.ToStringByPrecise(totalInsurance, 0)

		err = this.store.SaveIFMarketInfo(ifMarketInfo)
		if err != nil {
			return fmt.Errorf("StoreIFMarketInfo, this.store.SaveIFMarketInfo error: %s", err)
		}
	}
	return nil
}

func (this *IFPoolManager) IFPoolInfo(account string) (*common.IFPoolInfo, error) {
	ifPoolInfo := &common.IFPoolInfo{
		IFAssetList: make([]*common.IFAsset, 0),
		UserIFInfo: &common.UserIFInfo{
			Composition: make([]*common.Composition, 0),
		},
	}
	iFInfo, err := this.store.LoadIFInfo()
	if err != nil {
		return nil, fmt.Errorf("IFPoolInfo, this.store.LoadIFInfo error: %s", err)
	}
	ifPoolInfo.Cap = iFInfo.Cap
	ifPoolInfo.Total = iFInfo.Total

	allMarket, err := this.Comptroller.AllMarkets()
	if err != nil {
		return nil, fmt.Errorf("IFPoolInfo, this.Comptroller.AllMarkets error: %s", err)
	}
	totalSupplyDollar := new(big.Int)
	totalSupplyWingEarned := new(big.Int)
	totalBorrowDollar := new(big.Int)
	totalBorrowWingEarned := new(big.Int)
	totalInsuranceDollar := new(big.Int)
	totalInsuranceWingEarned := new(big.Int)
	for _, name := range allMarket {
		ifMarketInfo, err := this.store.LoadIFMarketInfo(name)
		if err != nil {
			return nil, fmt.Errorf("IFPoolInfo, this.store.LoadIFMarketInfo error: %s", err)
		}
		ifAsset := new(common.IFAsset)
		ifAsset.Name = this.cfg.IFMap[name]
		ifAsset.Icon = this.cfg.IconMap[ifAsset.Name]
		price, err := this.assetStoredPrice(name)
		if err != nil {
			return nil, fmt.Errorf("IFPoolInfo, this.assetStoredPrice error: %s", err)
		}
		ifAsset.Price = utils.ToStringByPrecise(price, this.cfg.TokenDecimal["oracle"])
		totalCash := utils.ToIntByPrecise(ifMarketInfo.TotalCash, 0)
		totalDebt := utils.ToIntByPrecise(ifMarketInfo.TotalDebt, 0)
		totalInsurance := utils.ToIntByPrecise(ifMarketInfo.TotalInsurance, 0)
		totalSupply := new(big.Int).Add(totalCash, totalDebt)
		ifAsset.TotalSupply = utils.ToStringByPrecise(totalSupply, this.cfg.TokenDecimal[ifAsset.Name])
		interestIndex := utils.ToIntByPrecise(ifMarketInfo.InterestIndex, 0)
		now := time.Now().UTC().Unix()
		ifAsset.SupplyInterestPerDay = utils.ToStringByPrecise(new(big.Int).Mul(new(big.Int).Div(interestIndex,
			new(big.Int).SetInt64(now-GenesisTime)), new(big.Int).SetUint64(governance.DaySecond)), this.cfg.TokenDecimal["interest"])
		//TODO supplyWingAPy
		ifAsset.UtilizationRate = utils.ToStringByPrecise(new(big.Int).Div(new(big.Int).Mul(totalDebt,
			new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal["interest"]))))), totalSupply), this.cfg.TokenDecimal["interest"])
		ifAsset.TotalBorrowed = utils.ToStringByPrecise(totalDebt, this.cfg.TokenDecimal[ifAsset.Name])
		//TODO BorrowWingAPY
		ifAsset.Liquidity = utils.ToStringByPrecise(totalCash, this.cfg.TokenDecimal[ifAsset.Name])
		ifAsset.BorrowCap = "500"
		ifAsset.TotalInsurance = utils.ToStringByPrecise(totalInsurance, this.cfg.TokenDecimal[ifAsset.Name])
		//TODO InsuranceWingAPY
		ifPoolInfo.IFAssetList = append(ifPoolInfo.IFAssetList, ifAsset)

		//user data
		if account != "" {
			addr, err := ocommon.AddressFromBase58(account)
			if err != nil {
				return nil, fmt.Errorf("IFPoolInfo, ocommon.AddressFromBase58 error: %s", err)
			}
			marketInfo, err := this.Comptroller.MarketInfo(name)
			if err != nil {
				return nil, fmt.Errorf("IFPoolInfo, this.Comptroller.MarketInfo error: %s", err)
			}
			assetName := this.cfg.IFMap[name]
			supplyBalance, err := this.FTokenMap[marketInfo.SupplyPool].BalanceOfUnderlying(addr)
			if err != nil {
				return nil, fmt.Errorf("IFPoolInfo, this.FTokenMap[marketInfo.SupplyPool].BalanceOfUnderlying error: %s", err)
			}
			supplyDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(supplyBalance, price),
				this.cfg.TokenDecimal["oracle"]+this.cfg.TokenDecimal[assetName]), this.cfg.TokenDecimal["USDT"])
			totalSupplyDollar = new(big.Int).Add(totalSupplyDollar, supplyDollar)
			_, supplyWingEarned, err := this.Comptroller.ClaimAllWing([]ocommon.Address{addr}, []string{name}, false, true, false, true)
			if err != nil {
				return nil, fmt.Errorf("IFPoolInfo, this.Comptroller.ClaimAllWing error: %s", err)
			}
			totalSupplyWingEarned = new(big.Int).Add(totalSupplyWingEarned, supplyWingEarned)
			_, borrowWingEarned, err := this.Comptroller.ClaimAllWing([]ocommon.Address{addr}, []string{name}, true, false, false, true)
			if err != nil {
				return nil, fmt.Errorf("IFPoolInfo, this.Comptroller.ClaimAllWing error: %s", err)
			}
			totalBorrowWingEarned = new(big.Int).Add(totalBorrowWingEarned, borrowWingEarned)
			_, insuranceWingEarned, err := this.Comptroller.ClaimAllWing([]ocommon.Address{addr}, []string{name}, false, false, true, true)
			if err != nil {
				return nil, fmt.Errorf("IFPoolInfo, this.Comptroller.ClaimAllWing error: %s", err)
			}
			totalInsuranceWingEarned = new(big.Int).Add(totalInsuranceWingEarned, insuranceWingEarned)
			insuranceBalance, err := this.ITokenMap[marketInfo.InsurancePool].BalanceOfUnderlying(addr)
			if err != nil {
				return nil, fmt.Errorf("IFPoolInfo, this.ITokenMap[marketInfo.InsurancePool].BalanceOfUnderlying error: %s", err)
			}
			insuranceDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(insuranceBalance, price),
				this.cfg.TokenDecimal["oracle"]+this.cfg.TokenDecimal[assetName]), this.cfg.TokenDecimal["USDT"])
			totalInsuranceDollar = new(big.Int).Add(totalInsuranceDollar, insuranceDollar)
			accountSnapshot, err := this.BorrowMap[marketInfo.BorrowPool].AccountSnapshot(addr)
			if err != nil {
				return nil, fmt.Errorf("IFPoolInfo, this.BorrowMap[marketInfo.BorrowPool].AccountSnapshot error: %s", err)
			}
			borrowDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(new(big.Int).Add(accountSnapshot.Principal,
				accountSnapshot.Interest), price), this.cfg.TokenDecimal["oracle"]+this.cfg.TokenDecimal[assetName]), this.cfg.TokenDecimal["USDT"])
			totalBorrowDollar = new(big.Int).Add(totalBorrowDollar, borrowDollar)
			composition := &common.Composition{
				Name:                  assetName,
				Icon:                  this.cfg.IconMap[assetName],
				SupplyBalance:         utils.ToStringByPrecise(supplyBalance, this.cfg.TokenDecimal[assetName]),
				SupplyWingEarned:      utils.ToStringByPrecise(supplyWingEarned, this.cfg.TokenDecimal["WING"]),
				BorrowWingEarned:      utils.ToStringByPrecise(borrowWingEarned, this.cfg.TokenDecimal["WING"]),
				LastBorrowTimestamp:   accountSnapshot.BorrowTime,
				InsuranceBalance:      utils.ToStringByPrecise(insuranceBalance, this.cfg.TokenDecimal[assetName]),
				InsuranceWingEarned:   utils.ToStringByPrecise(insuranceWingEarned, this.cfg.TokenDecimal["WING"]),
				CollateralBalance:     utils.ToStringByPrecise(accountSnapshot.Collateral, this.cfg.TokenDecimal[assetName]),
				BorrowUnpaidPrincipal: utils.ToStringByPrecise(accountSnapshot.Principal, this.cfg.TokenDecimal[assetName]),
				BorrowInterestBalance: utils.ToStringByPrecise(accountSnapshot.Interest, this.cfg.TokenDecimal[assetName]),
			}
			ifPoolInfo.UserIFInfo.Composition = append(ifPoolInfo.UserIFInfo.Composition, composition)
		}
	}
	if account != "" {
		ifPoolInfo.UserIFInfo.TotalSupplyDollar = utils.ToStringByPrecise(totalSupplyDollar, this.cfg.TokenDecimal["USDT"])
		ifPoolInfo.UserIFInfo.SupplyWingEarned = utils.ToStringByPrecise(totalSupplyWingEarned, this.cfg.TokenDecimal["WING"])
		ifPoolInfo.UserIFInfo.TotalBorrowDollar = utils.ToStringByPrecise(totalBorrowDollar, this.cfg.TokenDecimal["USDT"])
		ifPoolInfo.UserIFInfo.BorrowWingEarned = utils.ToStringByPrecise(totalBorrowWingEarned, this.cfg.TokenDecimal["WING"])
		ifPoolInfo.UserIFInfo.TotalInsuranceDollar = utils.ToStringByPrecise(totalInsuranceDollar, this.cfg.TokenDecimal["USDT"])
		ifPoolInfo.UserIFInfo.InsuranceWingEarned = utils.ToStringByPrecise(totalInsuranceDollar, this.cfg.TokenDecimal["USDT"])
	}
	return ifPoolInfo, nil
}

func (this *IFPoolManager) IFHistory(asset, operation string, start, end, pageNo, pageSize uint64) (*common.IFHistoryResponse, error) {
	return &common.IFHistoryResponse{
		MaxPageNum: 1,
		PageItems: []*common.IFHistory{
			{
				Name:      "pUSDT",
				Icon:      "https://app.ont.io/wing/pusdt.svg",
				Operation: "Supply",
				Timestamp: 1604026092000,
				Balance:   "32532.58",
				Dollar:    "23526.464",
				Address:   "Af3Etnp5ffrXR3swrCx9f7KuvChYLgqsTZ",
			},
			{
				Name:      "pDAI",
				Icon:      "https://app.ont.io/wing/oDAI.svg",
				Operation: "Borrow",
				Timestamp: 1604026092000,
				Balance:   "6968.58",
				Dollar:    "797.464",
				Address:   "AR36E5jLdWDKW3Yg51qDFWPGKSLvfPhbqS",
			},
		},
	}, nil
}
