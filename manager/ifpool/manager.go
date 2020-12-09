package ifpool

import (
	"fmt"
	"math"
	"math/big"
	"os"
	"time"

	sdk "github.com/ontio/ontology-go-sdk"
	ocommon "github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/config"
	"github.com/siovanus/wingServer/http/common"
	"github.com/siovanus/wingServer/log"
	"github.com/siovanus/wingServer/manager/governance"
	"github.com/siovanus/wingServer/store"
	"github.com/siovanus/wingServer/utils"
	if_borrow "github.com/wing-groups/wing-contract-tools/contracts/if-borrow"
	if_ctrl "github.com/wing-groups/wing-contract-tools/contracts/if-ctrl"
	"github.com/wing-groups/wing-contract-tools/contracts/iftoken"
	"github.com/wing-groups/wing-contract-tools/contracts/iitoken"
	oscore_oracle "github.com/wing-groups/wing-contract-tools/contracts/oscore-oracle"
)

var GenesisTime = time.Date(2020, time.December, 9, 0, 0, 0, 0, time.UTC).Unix()

const MaxLevel uint64 = 3

type IFPoolManager struct {
	cfg               *config.Config
	store             *store.Client
	Sdk               *sdk.OntologySdk
	Comptroller       *if_ctrl.Comptroller
	FTokenMap         map[ocommon.Address]*iftoken.IFToken
	ITokenMap         map[ocommon.Address]*iitoken.IIToken
	BorrowMap         map[ocommon.Address]*if_borrow.IfBorrowPool
	OscoreOracle      *oscore_oracle.Oracle
	GovMgr            *governance.GovernanceManager
	dailyDistribution *big.Int
}

func NewIFPoolManager(sdk *sdk.OntologySdk, contractAddress, oscoreOracleAddress ocommon.Address, store *store.Client,
	cfg *config.Config, govMgr *governance.GovernanceManager) *IFPoolManager {
	comptroller, _ := if_ctrl.NewComptroller(cfg.JsonRpcAddress, contractAddress.ToHexString(), nil,
		2500, 20000)
	oracle, _ := oscore_oracle.NewOracle(cfg.JsonRpcAddress, oscoreOracleAddress.ToHexString(), nil,
		2500, 20000)
	fTokenMap := make(map[ocommon.Address]*iftoken.IFToken)
	iTokenMap := make(map[ocommon.Address]*iitoken.IIToken)
	borrowPoolMap := make(map[ocommon.Address]*if_borrow.IfBorrowPool)
	allMarket, err := comptroller.AllMarkets()
	if err != nil {
		log.Errorf("NewIFPoolManager, comptroller.AllMarkets error: %s", err)
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
		Sdk:          sdk,
		Comptroller:  comptroller,
		FTokenMap:    fTokenMap,
		ITokenMap:    iTokenMap,
		BorrowMap:    borrowPoolMap,
		OscoreOracle: oracle,
		GovMgr:       govMgr,
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
		totalInterest, err := this.BorrowMap[marketInfo.BorrowPool].TotalInterest()
		if err != nil {
			return fmt.Errorf("StoreIFMarketInfo, this.BorrowMap[marketInfo.BorrowPool].TotalInterest error: %s", err)
		}
		ifMarketInfo.TotalCash = utils.ToStringByPrecise(totalCash, 0)
		ifMarketInfo.TotalDebt = utils.ToStringByPrecise(totalDebt, 0)
		ifMarketInfo.TotalInterest = utils.ToStringByPrecise(totalInterest.ToBigInt(), 0)

		oscoreInfo, err := this.BorrowMap[marketInfo.BorrowPool].GetOscoreInfoByLevel(MaxLevel)
		if err != nil {
			return fmt.Errorf("StoreIFMarketInfo, this.BorrowMap[marketInfo.BorrowPool].GetOscoreInfoByLevel error: %s, market %s", err, name)
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

func (this *IFPoolManager) StoreUserIfOperation(history *store.IfPoolHistory) error {
	err := this.store.SaveIFHistory(history)
	if err != nil {
		return fmt.Errorf("StoreIFInfo, this.store.SaveIFInfo error: %s", err)
	}
	return nil
}

func (this *IFPoolManager) IfPoolDetail() (*common.IfPoolDetail, error) {
	allMarket, err := this.Comptroller.AllMarkets()
	if err != nil {
		return nil, fmt.Errorf("IfPoolDetail, this.Comptroller.AllMarkets error: %s", err)
	}
	totalSupplyDollar := new(big.Int)
	for _, name := range allMarket {
		ifMarketInfo, err := this.store.LoadIFMarketInfo(name)
		if err != nil {
			return nil, fmt.Errorf("IfPoolDetail, this.store.LoadIFMarketInfo error: %s", err)
		}
		price, err := this.assetStoredPrice(name)
		if err != nil {
			return nil, fmt.Errorf("IFPoolInfo, this.assetStoredPrice error: %s", err)
		}
		totalCash := utils.ToIntByPrecise(ifMarketInfo.TotalCash, 0)
		totalDebt := utils.ToIntByPrecise(ifMarketInfo.TotalDebt, 0)
		totalSupply := new(big.Int).Add(totalCash, totalDebt)
		supplyDollar := new(big.Int).Mul(totalSupply, price)
		supplyDollar = utils.ToIntByPrecise(utils.ToStringByPrecise(supplyDollar, this.cfg.TokenDecimal[this.cfg.IFMap[name]]), 0)
		totalSupplyDollar = new(big.Int).Add(totalSupplyDollar, supplyDollar)
	}
	ifPoolDetail := &common.IfPoolDetail{
		TotalSupply: utils.ToStringByPrecise(totalSupplyDollar, this.cfg.TokenDecimal["oracle"]),
	}
	return ifPoolDetail, nil
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

	addr := ocommon.ADDRESS_EMPTY
	if account != "" {
		addr, err = ocommon.AddressFromBase58(account)
		if err != nil {
			return nil, fmt.Errorf("IFPoolInfo, ocommon.AddressFromBase58 error: %s", err)
		}
	}

	for _, name := range allMarket {
		wingApy, err := this.store.LoadIfWingApy(this.cfg.IFMap[name])
		if err != nil {
			return nil, fmt.Errorf("IFPoolInfo, this.store.LoadIfWingApy error: %s", err)
		}
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
		totalInterest := utils.ToIntByPrecise(ifMarketInfo.TotalInterest, this.cfg.TokenDecimal["percentage"])
		index := new(big.Int)
		if totalSupply.Uint64() != 0 {
			index = new(big.Int).Div(totalInterest, totalSupply)
		}
		now := time.Now().UTC().Unix()
		ifAsset.SupplyInterestPerDay = utils.ToStringByPrecise(new(big.Int).Mul(new(big.Int).Div(index,
			new(big.Int).SetInt64(now-GenesisTime)), new(big.Int).SetUint64(governance.DaySecond)), this.cfg.TokenDecimal["percentage"])
		// supplyWingAPy
		if totalSupply.Uint64() != 0 {
			ifAsset.UtilizationRate = utils.ToStringByPrecise(new(big.Int).Div(new(big.Int).Mul(totalDebt,
				new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal[ifAsset.Name]))))), totalSupply), this.cfg.TokenDecimal[ifAsset.Name])
		}
		ifAsset.SupplyWingAPY = wingApy.SupplyApy
		ifAsset.TotalBorrowed = utils.ToStringByPrecise(totalDebt, this.cfg.TokenDecimal[ifAsset.Name])
		// BorrowWingAPY
		ifAsset.Liquidity = utils.ToStringByPrecise(totalCash, this.cfg.TokenDecimal[ifAsset.Name])
		ifAsset.BorrowCap = "1000"
		ifAsset.BorrowWingAPY = wingApy.BorrowApy
		ifAsset.TotalInsurance = utils.ToStringByPrecise(totalInsurance, this.cfg.TokenDecimal[ifAsset.Name])
		// InsuranceWingAPY
		ifPoolInfo.IFAssetList = append(ifPoolInfo.IFAssetList, ifAsset)
		ifAsset.InsuranceWingAPY = wingApy.InsuranceApy

		//user data
		if account != "" {
			marketInfo, err := this.Comptroller.MarketInfo(name)
			if err != nil {
				return nil, fmt.Errorf("IFPoolInfo, this.Comptroller.MarketInfo error: %s", err)
			}
			assetName := this.cfg.IFMap[name]
			supplyBalance, err := this.FTokenMap[marketInfo.SupplyPool].BalanceOfUnderlying(addr)
			markets := []string{name}
			_, interest, err := this.Comptroller.ClaimAllInterest(addr, markets, true)
			if err != nil {
				return nil, fmt.Errorf("IFPoolInfo, this.Comptroller.ClaimAllInterest error: %s", err)
			}
			supplyBalance = new(big.Int).Add(supplyBalance, interest)
			if err != nil {
				return nil, fmt.Errorf("IFPoolInfo, this.FTokenMap[marketInfo.SupplyPool].BalanceOfUnderlying error: %s", err)
			}
			supplyDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(supplyBalance, price),
				this.cfg.TokenDecimal["oracle"]+this.cfg.TokenDecimal[assetName]), this.cfg.TokenDecimal["pUSDT"])
			totalSupplyDollar = new(big.Int).Add(totalSupplyDollar, supplyDollar)
			_, supplyWingEarned, err := this.Comptroller.ClaimAllWing([]ocommon.Address{addr}, []string{name}, false, true, false, true)
			if err != nil {
				return nil, fmt.Errorf("IFPoolInfo, market:%s, supply, %s, this.Comptroller.ClaimAllWing error: %s", name, account, err)
			}
			totalSupplyWingEarned = new(big.Int).Add(totalSupplyWingEarned, supplyWingEarned)

			_, borrowWingEarned, err := this.Comptroller.ClaimAllWing([]ocommon.Address{addr}, []string{name}, true, false, false, true)
			if err != nil {
				return nil, fmt.Errorf("IFPoolInfo, market:%s, borrow, %s, this.Comptroller.ClaimAllWing error: %s", name, account, err)
			}
			totalBorrowWingEarned = new(big.Int).Add(totalBorrowWingEarned, borrowWingEarned)

			_, insuranceWingEarned, err := this.Comptroller.ClaimAllWing([]ocommon.Address{addr}, []string{name}, false, false, true, true)
			if err != nil {
				return nil, fmt.Errorf("IFPoolInfo, market:%s, insurance, %s, this.Comptroller.ClaimAllWing error: %s", name, account, err)
			}
			totalInsuranceWingEarned = new(big.Int).Add(totalInsuranceWingEarned, insuranceWingEarned)

			insuranceBalance, err := this.ITokenMap[marketInfo.InsurancePool].BalanceOfUnderlying(addr)
			if err != nil {
				return nil, fmt.Errorf("IFPoolInfo, this.ITokenMap[marketInfo.InsurancePool].BalanceOfUnderlying error: %s", err)
			}
			insuranceDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(insuranceBalance, price),
				this.cfg.TokenDecimal["oracle"]+this.cfg.TokenDecimal[assetName]), this.cfg.TokenDecimal["pUSDT"])
			totalInsuranceDollar = new(big.Int).Add(totalInsuranceDollar, insuranceDollar)
			accountSnapshot, err := this.BorrowMap[marketInfo.BorrowPool].AccountSnapshotCurrent(addr)
			if err != nil {
				return nil, fmt.Errorf("IFPoolInfo, this.BorrowMap[marketInfo.BorrowPool].AccountSnapshot error: %s", err)
			}
			borrowDollar := utils.ToIntByPrecise(utils.ToStringByPrecise(new(big.Int).Mul(new(big.Int).Add(accountSnapshot.Principal,
				accountSnapshot.Interest), price), this.cfg.TokenDecimal["oracle"]+this.cfg.TokenDecimal[assetName]), this.cfg.TokenDecimal["pUSDT"])
			totalBorrowDollar = new(big.Int).Add(totalBorrowDollar, borrowDollar)
			composition := &common.Composition{
				Name:                  assetName,
				Icon:                  this.cfg.IconMap[assetName],
				SupplyBalance:         utils.ToStringByPrecise(supplyBalance, this.cfg.TokenDecimal[assetName]),
				SupplyWingEarned:      utils.ToStringByPrecise(supplyWingEarned, this.cfg.TokenDecimal["WING"]),
				BorrowWingEarned:      utils.ToStringByPrecise(borrowWingEarned, this.cfg.TokenDecimal["WING"]),
				LastBorrowTimestamp:   accountSnapshot.BorrowDayNumber,
				InsuranceBalance:      utils.ToStringByPrecise(insuranceBalance, this.cfg.TokenDecimal[assetName]),
				InsuranceWingEarned:   utils.ToStringByPrecise(insuranceWingEarned, this.cfg.TokenDecimal["WING"]),
				CollateralBalance:     utils.ToStringByPrecise(accountSnapshot.Collateral, this.cfg.TokenDecimal[assetName]),
				BorrowUnpaidPrincipal: utils.ToStringByPrecise(accountSnapshot.Principal, this.cfg.TokenDecimal[assetName]),
				BorrowInterestBalance: utils.ToStringByPrecise(accountSnapshot.FormalInterest, this.cfg.TokenDecimal[assetName]),
				BorrowExtraInterest:   utils.ToStringByPrecise(new(big.Int).Sub(accountSnapshot.Interest, accountSnapshot.FormalInterest), this.cfg.TokenDecimal[assetName]),
			}
			ifPoolInfo.UserIFInfo.Composition = append(ifPoolInfo.UserIFInfo.Composition, composition)
		}
	}
	if account != "" {
		ifPoolInfo.UserIFInfo.TotalSupplyDollar = utils.ToStringByPrecise(totalSupplyDollar, this.cfg.TokenDecimal["pUSDT"])
		ifPoolInfo.UserIFInfo.SupplyWingEarned = utils.ToStringByPrecise(totalSupplyWingEarned, this.cfg.TokenDecimal["WING"])
		ifPoolInfo.UserIFInfo.TotalBorrowDollar = utils.ToStringByPrecise(totalBorrowDollar, this.cfg.TokenDecimal["pUSDT"])
		ifPoolInfo.UserIFInfo.BorrowWingEarned = utils.ToStringByPrecise(totalBorrowWingEarned, this.cfg.TokenDecimal["WING"])
		ifPoolInfo.UserIFInfo.TotalInsuranceDollar = utils.ToStringByPrecise(totalInsuranceDollar, this.cfg.TokenDecimal["pUSDT"])
		ifPoolInfo.UserIFInfo.InsuranceWingEarned = utils.ToStringByPrecise(totalInsuranceWingEarned, this.cfg.TokenDecimal["WING"])
	}
	return ifPoolInfo, nil
}

func (this *IFPoolManager) IFHistory(address, asset, operation string, start, end, pageNo, pageSize uint64) (*common.IFHistoryResponse, error) {
	history, err := this.store.LoadIFHistory(address, asset, operation, start, end, pageNo, pageSize)
	if err != nil {
		return nil, fmt.Errorf("IFHistory, this.store.LoadIFHistory error: %s", err)
	}
	count, err := this.store.LoadIFHistoryCount(address, asset, operation, start, end)
	if err != nil {
		return nil, fmt.Errorf("IFHistory, this.store.LoadIFHistoryCount error: %s", err)
	}
	histories := make([]*common.IFHistory, 0)
	for _, v := range history {
		price, err := this.assetStoredPrice(this.cfg.IFOracleMap[v.Token])
		if err != nil {
			log.Errorf("IFHistory, this.AssetStoredPrice error: %s", err)
		}
		amount := utils.ToIntByPrecise(v.Amount, this.cfg.TokenDecimal[v.Token])
		dollar := utils.ToStringByPrecise(new(big.Int).Mul(amount, price), this.cfg.TokenDecimal["oracle"]+this.cfg.TokenDecimal[v.Token])
		i := &common.IFHistory{
			Name:      v.Token,
			Icon:      this.cfg.IconMap[v.Token],
			Operation: v.Operation,
			Timestamp: v.Timestamp,
			Balance:   v.Amount,
			Dollar:    dollar,
			Address:   v.Address,
		}
		histories = append(histories, i)
	}
	return &common.IFHistoryResponse{
		Count:     count,
		PageItems: histories,
	}, nil
}

func (this *IFPoolManager) CheckIfDebt(start, end int64) ([]*common.DebtAccount, error) {
	history, err := this.store.LoadIFBorrowUsersInLimitDay(start, end)
	if err != nil {
		fmt.Errorf("CheckIfDebt, this.store.LoadIFBorrowUsersInLimitDay error: %s", err)
	}
	debtAccounts := make([]*common.DebtAccount, 0)
	for _, v := range history {
		collateralAssets := make([]*common.CollateralAsset, 0)
		totalCollateralDollar := new(big.Int)
		marketInfo, err := this.Comptroller.MarketInfo(this.cfg.IFOracleMap[v.Token])
		if err != nil {
			log.Errorf("CheckIfDebt, this.Comptroller.MarketInfo error: %s", err)
		}
		address := v.Address
		addr, err := ocommon.AddressFromBase58(address)
		if err != nil {
			log.Errorf("CheckIfDebt, ocommon.AddressFromBase58 error: %s", err)
		}
		accountSnapshot, err := this.BorrowMap[marketInfo.BorrowPool].AccountSnapshotCurrent(addr)
		if err != nil {
			log.Errorf("CheckIfDebt, this.BorrowMap[marketInfo.BorrowPool].AccountSnapshotCurrent error: %s", err)
		}
		principal := accountSnapshot.Principal
		interest := accountSnapshot.Interest
		debt := new(big.Int).Add(principal, interest)
		if big.NewInt(0).Cmp(debt) < 0 {
			// 逾期违约
			// 查询抵押品
			for _, marketName := range this.cfg.IFOracleMap {
				collateralMarketInfo, err := this.Comptroller.MarketInfo(marketName)
				if err != nil {
					log.Errorf("CheckIfDebt, this.Comptroller.MarketInfo error: %s", err)
				}
				collateralAccountSnapshot, err := this.BorrowMap[collateralMarketInfo.BorrowPool].AccountSnapshotCurrent(addr)
				if err != nil {
					log.Errorf("CheckIfDebt, this.BorrowMap[marketInfo.BorrowPool].AccountSnapshotCurrent error: %s", err)
				}
				collateral := collateralAccountSnapshot.Collateral
				if big.NewInt(0).Cmp(collateral) < 0 {
					collateralPrice, err := this.assetStoredPrice(marketName)
					if err != nil {
						log.Errorf("CheckIfDebt, this.AssetStoredPrice error: %s", err)
					}

					collateralDollar := utils.ToStringByPrecise(new(big.Int).Mul(collateral, collateralPrice), this.cfg.TokenDecimal["oracle"]+this.cfg.TokenDecimal[this.cfg.IFMap[marketName]])
					totalCollateralDollar = new(big.Int).Add(totalCollateralDollar, utils.ToIntByPrecise(collateralDollar, this.cfg.TokenDecimal["oracle"]))
					collateralAsset := &common.CollateralAsset{
						Icon:    this.cfg.IconMap[this.cfg.IFMap[marketName]],
						Name:    this.cfg.IFMap[marketName],
						Balance: utils.ToStringByPrecise(collateral, this.cfg.TokenDecimal[this.cfg.IFMap[marketName]]),
						Dollar:  collateralDollar,
					}
					collateralAssets = append(collateralAssets, collateralAsset)
				}
			}

			debtPrice, err := this.assetStoredPrice(this.cfg.IFOracleMap[v.Token])
			if err != nil {
				log.Errorf("CheckIfDebt, this.AssetStoredPrice error: %s", err)
			}
			debtDollar := utils.ToStringByPrecise(new(big.Int).Mul(debt, debtPrice), this.cfg.TokenDecimal["oracle"]+this.cfg.TokenDecimal[v.Token])

			i := &common.DebtAccount{
				Address:          address,
				Debt:             v.Token,
				DebtIcon:         this.cfg.IconMap[v.Token],
				DebtAmount:       utils.ToStringByPrecise(debt, this.cfg.TokenDecimal[v.Token]),
				DebtPrice:        debtDollar,
				CollateralDollar: utils.ToStringByPrecise(totalCollateralDollar, this.cfg.TokenDecimal["oracle"]),
				CollateralAssets: collateralAssets,
				BorrowTime:       accountSnapshot.BorrowDayNumber,
			}
			debtAccounts = append(debtAccounts, i)
		}
	}
	return debtAccounts, nil
}

func (this *IFPoolManager) Reserves() (*common.Reserves, error) {
	allMarket, err := this.Comptroller.AllMarkets()
	if err != nil {
		return nil, fmt.Errorf("Reserves, this.Comptroller.AllMarkets error: %s", err)
	}
	totalReserve := new(big.Int)
	reserves := &common.Reserves{
		AssetReserve: make([]*common.Reserve, 0),
	}
	for _, name := range allMarket {
		marketInfo, err := this.Comptroller.MarketInfo(name)
		if err != nil {
			return nil, fmt.Errorf("Reserves, this.Comptroller.MarketInfo error: %s", err)
		}
		assetName := this.cfg.IFMap[name]
		price, err := this.assetStoredPrice(name)
		if err != nil {
			return nil, fmt.Errorf("Reserves, this.assetStoredPrice error: %s", err)
		}
		admin, err := this.Comptroller.ReservesAddr()
		if err != nil {
			return nil, fmt.Errorf("Reserves, this.Comptroller.ReservesAddr error: %s", err)
		}
		result, err := this.Sdk.NeoVM.PreExecInvokeNeoVMContract(marketInfo.Underlying,
			[]interface{}{"balanceOf", []interface{}{admin}})
		if err != nil {
			return nil, fmt.Errorf("Reserves, this.Sdk.NeoVM.PreExecInvokeNeoVMContract error: %s", err)
		}
		reserveBalance, err := result.Result.ToInteger()
		if err != nil {
			return nil, fmt.Errorf("Reserves, result.Result.ToInteger error: %s", err)
		}
		reserveBalanceStr := utils.ToStringByPrecise(reserveBalance, this.cfg.TokenDecimal[assetName])
		reserveDollarStr := utils.ToStringByPrecise(new(big.Int).Mul(price, reserveBalance),
			this.cfg.TokenDecimal[assetName]+this.cfg.TokenDecimal["oracle"])
		reserveFactor, err := this.BorrowMap[marketInfo.BorrowPool].ReserveFactor()
		if err != nil {
			return nil, fmt.Errorf("Reserves, this.BorrowMap[marketInfo.BorrowPool].ReserveFactor error: %s", err)
		}
		assetReserve := &common.Reserve{
			Name:           assetName,
			Icon:           this.cfg.IconMap[assetName],
			ReserveFactor:  utils.ToStringByPrecise(reserveFactor, this.cfg.TokenDecimal["interest"]),
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

func (this *IFPoolManager) PoolDistribution() (*common.Distribution, error) {
	allMarket, err := this.Comptroller.AllMarkets()
	if err != nil {
		return nil, fmt.Errorf("PoolDistribution, this.Comptroller.AllMarkets error: %s", err)
	}
	distribution := new(common.Distribution)
	s := new(big.Int).SetUint64(0)
	b := new(big.Int).SetUint64(0)
	i := new(big.Int).SetUint64(0)
	d := new(big.Int).SetUint64(0)
	for _, name := range allMarket {
		assetName := this.cfg.IFMap[name]
		marketInfo, err := this.Comptroller.MarketInfo(name)
		if err != nil {
			return nil, fmt.Errorf("PoolDistribution, this.Comptroller.MarketInfo error: %s", err)
		}
		cash, err := this.FTokenMap[marketInfo.SupplyPool].TotalCash()
		if err != nil {
			return nil, fmt.Errorf("PoolDistribution, this.FTokenMap[marketInfo.SupplyPool].TotalCash error: %s", err)
		}
		borrow, err := this.FTokenMap[marketInfo.SupplyPool].TotalDebt()
		if err != nil {
			return nil, fmt.Errorf("PoolDistribution, this.FTokenMap[marketInfo.SupplyPool].TotalDebt error: %s", err)
		}
		supplyAmount := new(big.Int).Add(cash, borrow)
		borrowAmount := borrow
		insuranceAmount, err := this.ITokenMap[marketInfo.InsurancePool].TotalCash()
		if err != nil {
			return nil, fmt.Errorf("PoolDistribution, this.ITokenMap[marketInfo.InsurancePool].TotalCash error: %s", err)
		}

		supplyAmountFormal := utils.ToIntByPrecise(utils.ToStringByPrecise(supplyAmount, this.cfg.TokenDecimal[assetName]),
			this.cfg.TokenDecimal["pUSDT"])
		borrowAmountFormal := utils.ToIntByPrecise(utils.ToStringByPrecise(borrowAmount, this.cfg.TokenDecimal[assetName]),
			this.cfg.TokenDecimal["pUSDT"])
		insuranceAmountFormal := utils.ToIntByPrecise(utils.ToStringByPrecise(insuranceAmount, this.cfg.TokenDecimal[assetName]),
			this.cfg.TokenDecimal["pUSDT"])

		totalDistribution, err := this.Comptroller.WingDistributedNum(name)
		if err != nil {
			return nil, fmt.Errorf("PoolDistribution, this.Comptroller.WingDistributedNum error: %s", err)
		}
		price, err := this.assetStoredPrice(name)
		if err != nil {
			return nil, fmt.Errorf("PoolDistribution, this.assetStoredPrice error: %s", err)
		}

		supplyDollar := new(big.Int).Mul(supplyAmountFormal, price)
		borrowAmountDollar := new(big.Int).Mul(borrowAmountFormal, price)
		insuranceDollar := new(big.Int).Mul(insuranceAmountFormal, price)
		// supplyAmount * price
		s = new(big.Int).Add(s, supplyDollar)
		// borrowAmount * price
		b = new(big.Int).Add(b, borrowAmountDollar)
		// insuranceAmount * price
		i = new(big.Int).Add(i, insuranceDollar)
		d = new(big.Int).Add(d, totalDistribution)
	}
	distribution.Name = "Inclusive"
	distribution.Icon = this.cfg.IconMap[distribution.Name]
	distribution.SupplyAmount = utils.ToStringByPrecise(s, this.cfg.TokenDecimal["pUSDT"]+this.cfg.TokenDecimal["oracle"])
	distribution.BorrowAmount = utils.ToStringByPrecise(b, this.cfg.TokenDecimal["pUSDT"]+this.cfg.TokenDecimal["oracle"])
	distribution.InsuranceAmount = utils.ToStringByPrecise(i, this.cfg.TokenDecimal["pUSDT"]+this.cfg.TokenDecimal["oracle"])
	distribution.Total = utils.ToStringByPrecise(d, this.cfg.TokenDecimal["WING"])
	return distribution, nil
}

func (this *IFPoolManager) IfPoolBanner() (*common.PoolBanner, error) {
	allMarkets, err := this.Comptroller.AllMarkets()
	if err != nil {
		return nil, fmt.Errorf("FlashPoolBanner, this.GetAllMarkets error: %s", err)
	}
	total := new(big.Int).SetUint64(0)
	for _, name := range allMarkets {
		totalDistribution, err := this.Comptroller.WingDistributedNum(name)
		if err != nil {
			return nil, fmt.Errorf("IfPoolBanner, this.Comptroller.WingDistributedNum error: %s", err)
		}
		total = new(big.Int).Add(total, totalDistribution)
	}

	share := new(big.Int).SetUint64(0)
	if this.dailyDistribution.Uint64() != 0 {
		share = new(big.Int).Div(this.dailyDistribution, total)
	}

	return &common.PoolBanner{
		Today: utils.ToStringByPrecise(this.dailyDistribution, this.cfg.TokenDecimal["WING"]+this.cfg.TokenDecimal["percentage"]+2),
		Share: utils.ToStringByPrecise(share, this.cfg.TokenDecimal["percentage"]+2),
		Total: utils.ToStringByPrecise(total, this.cfg.TokenDecimal["WING"]),
	}, nil
}

func (this *IFPoolManager) MarketDistribution() (*common.MarketDistribution, error) {
	allMarket, err := this.Comptroller.AllMarkets()
	if err != nil {
		return nil, fmt.Errorf("MarketDistribution, this.Comptroller.AllMarkets error: %s", err)
	}
	ifPoolMarketDistribution := make([]*common.Distribution, 0)
	for _, name := range allMarket {
		assetName := this.cfg.IFMap[name]
		marketInfo, err := this.Comptroller.MarketInfo(name)
		if err != nil {
			return nil, fmt.Errorf("MarketDistribution, this.Comptroller.MarketInfo error: %s", err)
		}
		cash, err := this.FTokenMap[marketInfo.SupplyPool].TotalCash()
		if err != nil {
			return nil, fmt.Errorf("MarketDistribution, this.FTokenMap[marketInfo.SupplyPool].TotalCash error: %s", err)
		}
		borrow, err := this.FTokenMap[marketInfo.SupplyPool].TotalDebt()
		if err != nil {
			return nil, fmt.Errorf("MarketDistribution, this.FTokenMap[marketInfo.SupplyPool].TotalDebt error: %s", err)
		}
		supplyAmount := new(big.Int).Add(cash, borrow)
		borrowAmount := borrow
		insuranceAmount, err := this.ITokenMap[marketInfo.InsurancePool].TotalCash()
		if err != nil {
			return nil, fmt.Errorf("MarketDistribution, this.ITokenMap[marketInfo.InsurancePool].TotalCash error: %s", err)
		}

		supplyAmountFormal := utils.ToIntByPrecise(utils.ToStringByPrecise(supplyAmount, this.cfg.TokenDecimal[assetName]),
			this.cfg.TokenDecimal["pUSDT"])
		borrowAmountFormal := utils.ToIntByPrecise(utils.ToStringByPrecise(borrowAmount, this.cfg.TokenDecimal[assetName]),
			this.cfg.TokenDecimal["pUSDT"])
		insuranceAmountFormal := utils.ToIntByPrecise(utils.ToStringByPrecise(insuranceAmount, this.cfg.TokenDecimal[assetName]),
			this.cfg.TokenDecimal["pUSDT"])

		totalDistribution, err := this.Comptroller.WingDistributedNum(name)
		if err != nil {
			return nil, fmt.Errorf("MarketDistribution, this.Comptroller.WingDistributedNum error: %s", err)
		}
		price, err := this.assetStoredPrice(name)
		if err != nil {
			return nil, fmt.Errorf("MarketDistribution, this.assetStoredPrice error: %s", err)
		}

		supplyDollar := new(big.Int).Mul(supplyAmountFormal, price)
		borrowAmountDollar := new(big.Int).Mul(borrowAmountFormal, price)
		insuranceDollar := new(big.Int).Mul(insuranceAmountFormal, price)

		distribution := &common.Distribution{
			Icon:            this.cfg.IconMap[this.cfg.FlashAssetMap[name]],
			Name:            this.cfg.FlashAssetMap[name],
			SupplyAmount:    utils.ToStringByPrecise(supplyDollar, this.cfg.TokenDecimal["pUSDT"]+this.cfg.TokenDecimal["oracle"]),
			BorrowAmount:    utils.ToStringByPrecise(borrowAmountDollar, this.cfg.TokenDecimal["pUSDT"]+this.cfg.TokenDecimal["oracle"]),
			InsuranceAmount: utils.ToStringByPrecise(insuranceDollar, this.cfg.TokenDecimal["pUSDT"]+this.cfg.TokenDecimal["oracle"]),
			Total:           utils.ToStringByPrecise(totalDistribution, this.cfg.TokenDecimal["WING"]),
		}
		ifPoolMarketDistribution = append(ifPoolMarketDistribution, distribution)
	}
	return &common.MarketDistribution{MarketDistribution: ifPoolMarketDistribution}, nil
}

func (this *IFPoolManager) WingApyForStore() error {
	dynamicPercent, err := this.getDynamicPercent()
	if err != nil {
		return fmt.Errorf("IFPoolManager WingApy, this.getDynamicPercent error: %s", err)
	}
	log.Infof("dynamicPercent:%d", dynamicPercent)
	staticPercent := new(big.Int).Sub(new(big.Int).SetUint64(100), dynamicPercent)
	log.Infof("staticPercent:%d", staticPercent)

	poolWeight, err := this.getPoolWeight()
	if err != nil {
		return fmt.Errorf("IFPoolManager WingApy, this.getPoolWeight error: %s", err)
	}
	poolStaticMap := poolWeight.PoolStaticMap
	ifStaticWeight, ok := poolStaticMap[this.Comptroller.GetAddr()]
	if !ok {
		return fmt.Errorf("if StaticWeight:%d", ifStaticWeight)
	}
	log.Infof("if StaticWeight:%d", ifStaticWeight)
	totalStaticWeight := poolWeight.TotalStatic
	log.Infof("if totalStaticWeight:%d", totalStaticWeight)
	ifStaticPercent := new(big.Int).SetUint64(0)
	if totalStaticWeight.Cmp(big.NewInt(0)) != 0 {
		log.Infof("_________________________________totalStaticWeight !=0")
		ifStaticPercent = new(big.Int).Div(new(big.Int).Mul(ifStaticWeight, new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal["apyPercentage"]))))), totalStaticWeight)
	}
	log.Infof("if StaticPercent:%d", ifStaticPercent)

	poolDynamicMap := poolWeight.PoolDynamicMap
	ifDynamicWeight := poolDynamicMap[this.Comptroller.GetAddr()]
	log.Infof("if DynamicWeight:%d", ifDynamicWeight)
	totalDynamicWeight := poolWeight.TotalDynamic
	log.Infof("if totalDynamicWeight:%d", totalDynamicWeight)
	ifDynamicPercent := new(big.Int).SetUint64(0)
	if totalDynamicWeight.Cmp(big.NewInt(0)) != 0 {
		log.Infof("_________________________________totalDynamicWeight !=0")
		ifDynamicPercent = new(big.Int).Div(new(big.Int).Mul(ifDynamicWeight, new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal["apyPercentage"]))))), totalDynamicWeight)
	}
	log.Infof("if DynamicPercent:%d", ifDynamicPercent)

	utilities, err := this.getUtilities()
	if err != nil {
		return fmt.Errorf("IFPoolManager WingApy, this.getUtilities error: %s", err)
	}
	utilityMap := utilities.UtilityMap
	total := utilities.Total
	log.Infof("################################utilities.Total: %d", total)

	banner, err := this.GovMgr.GovBanner()
	if err != nil {
		return fmt.Errorf("IFPoolManager WingApy, this.GovMgr.GovBanner error: %s", err)
	}
	daily := banner.Daily
	dailyTotal := utils.ToIntByPrecise(daily, 9)
	log.Infof("origin dailyTotal:%d", dailyTotal)
	dailyTotal = new(big.Int).Div(new(big.Int).Mul(dailyTotal, new(big.Int).SetUint64(60)), new(big.Int).SetUint64(100))
	log.Infof("0.6 times dailyTotal:%d", dailyTotal)
	dailyTotal = new(big.Int).Add(new(big.Int).Mul(staticPercent, new(big.Int).Mul(dailyTotal, ifStaticPercent)), new(big.Int).Mul(dynamicPercent, new(big.Int).Mul(dailyTotal, ifDynamicPercent)))
	log.Infof("if weight dailyTotal:%d", dailyTotal)
	this.dailyDistribution = dailyTotal

	wingPrice, err := this.assetStoredPrice("WING")
	if err != nil {
		return fmt.Errorf("IFPoolManager WingApy, this.AssetStoredPrice error: %s", err)
	}

	allMarkets, err := this.Comptroller.AllMarkets()
	if err != nil {
		return fmt.Errorf("IFPoolManager WingApy, this.GetAllMarkets error: %s", err)
	}
	for _, name := range allMarkets {
		ifMarketInfo, err := this.store.LoadIFMarketInfo(name)
		if err != nil {
			fmt.Errorf("IFPoolManager WingApy, this.store.LoadIFMarketInfo error: %s", err)
		}

		totalCash := utils.ToIntByPrecise(ifMarketInfo.TotalCash, 0)
		totalDebt := utils.ToIntByPrecise(ifMarketInfo.TotalDebt, 0)
		totalInsurance := utils.ToIntByPrecise(ifMarketInfo.TotalInsurance, 0)
		totalSupply := new(big.Int).Add(totalCash, totalDebt)

		wingSBIPortion, err := this.Comptroller.WingSBIPortion(name)
		if err != nil {
			return fmt.Errorf("IFPoolManager WingApy, this.getWingSBIPortion error: %s", err)
		}

		totalPortion := new(big.Int).Add(new(big.Int).SetUint64(uint64(wingSBIPortion.InsurancePortion)),
			new(big.Int).Add(new(big.Int).SetUint64(uint64(wingSBIPortion.SupplyPortion)),
				new(big.Int).SetUint64(uint64(wingSBIPortion.BorrowPortion))))

		price, err := this.assetStoredPrice(name)
		if err != nil {
			return fmt.Errorf("IFPoolManager WingApy, this.AssetStoredPrice error: %s", err)
		}

		totalSupplyDollar := new(big.Int).Mul(totalSupply, price)
		marketInfo, err := this.Comptroller.MarketInfo(name)
		if err != nil {
			log.Errorf("WingApyForStore, this.Comptroller.MarketInfo error: %s", err)
		}
		totalValidBorrow, err := this.BorrowMap[marketInfo.BorrowPool].FormalPrincipal()
		if err != nil {
			log.Errorf("WingApyForStore, this.BorrowMap[marketInfo.BorrowPool].FormalPrincipal error: %s", err)
		}
		log.Infof("##########################name:%s,totalValidBorrow:%d", name, totalValidBorrow.ToBigInt())
		totalValidBorrowDollar := new(big.Int).Mul(totalValidBorrow.ToBigInt(), price)
		log.Infof("##########################name:%s,totalValidBorrowDollar:%d", name, totalValidBorrowDollar)
		totalInsuranceDollar := new(big.Int).Mul(totalInsurance, price)
		if totalInsuranceDollar.Uint64() == 0 {
			totalInsuranceDollar = utils.ToIntByPrecise("1", this.cfg.TokenDecimal["oracle"]+this.cfg.TokenDecimal[this.cfg.IFMap[name]])
		}
		utility, ok := utilityMap[name]
		log.Infof("##########################name:%s", name)
		log.Infof("##########################utility:%d", utility)
		log.Infof("##########################total:%d", total)
		supplyApy := "0"
		borrowApy := "0"
		insuranceApy := "0"
		if ok && totalSupplyDollar.Uint64() != 0 && utility.Cmp(big.NewInt(0)) != 0 && utility.Cmp(big.NewInt(1)) != 0 {
			log.Infof("##########################SupplyPortion:%d", wingSBIPortion.SupplyPortion)
			supplyApy = utils.ToStringByPrecise(new(big.Int).Div(new(big.Int).Div(new(big.Int).Mul(new(big.Int).Mul(new(big.Int).Mul(new(big.Int).Mul(
				new(big.Int).Div(new(big.Int).Mul(dailyTotal, utility), total),
				new(big.Int).SetUint64(uint64(wingSBIPortion.SupplyPortion))), wingPrice), new(big.Int).SetUint64(governance.YearDay)),
				new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal[this.cfg.IFMap[name]]))))), totalPortion),
				totalSupplyDollar), this.cfg.TokenDecimal["WING"]+this.cfg.TokenDecimal["apyPercentage"]+2)
		}
		if ok && totalValidBorrowDollar.Uint64() != 0 && utility.Cmp(big.NewInt(0)) != 0 && utility.Cmp(big.NewInt(1)) != 0 {
			log.Infof("##########################BorrowPortion:%d", wingSBIPortion.BorrowPortion)
			borrowApy = utils.ToStringByPrecise(new(big.Int).Div(new(big.Int).Div(new(big.Int).Mul(new(big.Int).Mul(new(big.Int).Mul(new(big.Int).Mul(
				new(big.Int).Div(new(big.Int).Mul(dailyTotal, utility), total),
				new(big.Int).SetUint64(uint64(wingSBIPortion.BorrowPortion))), wingPrice), new(big.Int).SetUint64(governance.YearDay)),
				new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal[this.cfg.IFMap[name]]))))), totalPortion),
				totalValidBorrowDollar), this.cfg.TokenDecimal["WING"]+this.cfg.TokenDecimal["apyPercentage"]+2)
		}
		if ok && totalInsuranceDollar.Uint64() != 0 && utility.Cmp(big.NewInt(0)) != 0 && utility.Cmp(big.NewInt(1)) != 0 {
			log.Infof("##########################InsurancePortion:%d", wingSBIPortion.InsurancePortion)
			insuranceApy = utils.ToStringByPrecise(new(big.Int).Div(new(big.Int).Div(new(big.Int).Mul(new(big.Int).Mul(new(big.Int).Mul(new(big.Int).Mul(
				new(big.Int).Div(new(big.Int).Mul(dailyTotal, utility), total),
				new(big.Int).SetUint64(uint64(wingSBIPortion.InsurancePortion))), wingPrice), new(big.Int).SetUint64(governance.YearDay)),
				new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal[this.cfg.IFMap[name]]))))), totalPortion),
				totalInsuranceDollar), this.cfg.TokenDecimal["WING"]+this.cfg.TokenDecimal["apyPercentage"]+2)
		}

		ifWingApy := &store.IfWingApy{
			AssetName:    this.cfg.IFMap[name],
			SupplyApy:    supplyApy,
			BorrowApy:    borrowApy,
			InsuranceApy: insuranceApy,
		}

		err = this.store.SaveIfWingApy(ifWingApy)
		if err != nil {
			return fmt.Errorf("IFPoolManager WingApy, this.store.SaveWingApy error: %s", err)
		}
	}
	return nil
}
