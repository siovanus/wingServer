package ifpool

import (
	"fmt"
	sdk "github.com/ontio/ontology-go-sdk"
	ocommon "github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/config"
	"github.com/siovanus/wingServer/http/common"
	"github.com/siovanus/wingServer/store"
	"github.com/siovanus/wingServer/utils"
	"math"
	"math/big"
)

type IFPoolManager struct {
	cfg             *config.Config
	contractAddress ocommon.Address
	oracleAddress   ocommon.Address
	sdk             *sdk.OntologySdk
	store           *store.Client
}

func NewIFPoolManager(contractAddress, oracleAddress ocommon.Address, sdk *sdk.OntologySdk,
	store *store.Client, cfg *config.Config) *IFPoolManager {
	manager := &IFPoolManager{
		cfg:             cfg,
		contractAddress: contractAddress,
		oracleAddress:   oracleAddress,
		sdk:             sdk,
		store:           store,
	}

	return manager
}

func (this *IFPoolManager) AssetStoredPrice(asset string) (*big.Int, error) {
	if asset == "USDT" {
		return new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal["oracle"])))), nil
	}
	price, err := this.store.LoadPrice(asset)
	if err != nil {
		return nil, fmt.Errorf("AssetStoredPrice, this.store.LoadPrice error: %s", err)
	}
	return utils.ToIntByPrecise(price.Price, this.cfg.TokenDecimal["oracle"]), nil
}

func (this *IFPoolManager) IFPoolOverview() (*common.IFPoolOverview, error) {
	return &common.IFPoolOverview{
		Total: "462000.13241",
		Cap:   "500000",
		IFAssetList: []*common.IFAssetList{
			{Name: "pSUSD", Icon: "https://app.ont.io/wing/psusd.svg", SupplyBalance: "400000.13241", BorrowInterestPerDay: "0.0005", Liquidity: "2325.67"},
			{Name: "pUSDT", Icon: "https://app.ont.io/wing/pusdt.svg", SupplyBalance: "61000", BorrowInterestPerDay: "0.0005", Liquidity: "23525.456"},
			{Name: "pDAI", Icon: "https://app.ont.io/wing/oDAI.svg", SupplyBalance: "1000", BorrowInterestPerDay: "0.0005", Liquidity: "47474.67"},
		},
	}, nil
}

func (this *IFPoolManager) IFMarketDetail(market string) (*common.IFMarketDetail, error) {
	return &common.IFMarketDetail{
		Name:                 market,
		Icon:                 "https://app.ont.io/wing/pusdt.svg",
		TotalSupply:          "121435.4564747",
		SupplyWingAPY:        "0.242536",
		UtilizationRate:      "0.786714",
		MaximumLTV:           "0.1425",
		TotalBorrowed:        "965256.25225",
		BorrowInterestPerDay: "0.000224",
		BorrowWingAPY:        "235267",
		Liquidity:            "2536.4564",
		BorrowCap:            "500",
		TotalInsurance:       "1876969.34536",
		InsuranceWingAPY:     "25.34649",
	}, nil
}

func (this *IFPoolManager) UserIFInfo(account string) (*common.UserIFInfo, error) {
	return &common.UserIFInfo{
		TotalSupplyDollar:    "23526.3647",
		SupplyWingEarned:     "25.3647",
		TotalBorrowDollar:    "25364.485",
		BorrowWingEarned:     "789.36536",
		BorrowInterestPerDay: "0.837636",
		TotalInsuranceDollar: "96747.474747",
		InsuranceWingEarned:  "36796.366",
		Composition: []*common.Composition{
			{Operation: "supply", Name: "pSUSD", Icon: "https://app.ont.io/wing/psusd.svg", Balance: "2536.367", IfCanOp: true},
			{Operation: "supply", Name: "pUSDT", Icon: "https://app.ont.io/wing/pusdt.svg", Balance: "42526.3636", IfCanOp: true},
			{Operation: "borrow", Name: "pSUSD", Icon: "https://app.ont.io/wing/psusd.svg", Balance: "235546.4647", IfCanOp: false},
			{Operation: "insurance", Name: "pDAI", Icon: "https://app.ont.io/wing/oDAI.svg", Balance: "4756.252", IfCanOp: true},
		},
	}, nil
}

func (this *IFPoolManager) UserIFMarketInfo(account, market string) (*common.UserIFMarketInfo, error) {
	return &common.UserIFMarketInfo{
		SupplyBalance:       "583525.256",
		SupplyWingEarned:    "376.4373",
		BorrowBalance:       "4727.45747",
		Limit:               "0.276525",
		BorrowWingEarned:    "12.3637",
		InsuranceBalance:    "4225.3637",
		InsuranceWingEarned: "12.637",
	}, nil
}
