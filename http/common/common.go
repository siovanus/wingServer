package common

import (
	"github.com/ontio/ontology/common"
	"math/big"
)

const (
	FLASHPOOLMARKETDISTRIBUTION = "/api/v1/flashpoolmarketdistribution"
	IFPOOLMARKETDISTRIBUTION    = "/api/v1/ifpoolmarketdistribution"
	POOLDISTRIBUTION            = "/api/v1/pooldistribution"
	GOVBANNEROVERVIEW           = "/api/v1/govbanneroverview"
	GOVBANNER                   = "/api/v1/govbanner"
	RESERVES                    = "/api/v1/reserves"
	IFRESERVES                  = "/api/v1/ifreserves"
	FLASHPOOLBANNER             = "/api/v1/flashpoolbanner"
	IFPOOLBANNER                = "/api/v1/ifpoolbanner"

	FLASHPOOLDETAIL       = "/api/v1/flashpooldetail"
	FLASHPOOLALLMARKET    = "/api/v1/flashpoolallmarket"
	USERFLASHPOOLOVERVIEW = "/api/v1/userflashpooloverview"

	IFPOOLDETAIL = "/api/v1/ifpooldetail"

	BORROWADDRESSLIST = "/api/v1/borrowaddresslist"

	ASSETPRICE      = "/api/v1/assetprice"
	ASSETPRICELIST  = "/api/v1/assetpricelist"
	CLAIMWING       = "/api/v1/claimwing"
	LIQUIDATIONLIST = "/api/v1/liquidationlist"
	WINGAPYS        = "/api/v1/wingapys"

	IFPOOLINFO    = "/api/v1/if/ifpoolinfo"
	IFHOSTORY     = "/api/v1/if/ifhistory"
	IFLIQUIDATION = "/api/v1/ifliquidationlist"
)

const (
	ACTION_FLASHPOOLMARKETDISTRIBUTION = "flashpoolmarketdistribution"
	ACTION_IFPOOLMARKETDISTRIBUTION    = "ifpoolmarketdistribution"
	ACTION_POOLDISTRIBUTION            = "pooldistribution"
	ACTION_GOVBANNEROVERVIEW           = "govbanneroverview"
	ACTION_GOVBANNER                   = "govbanner"
	ACTION_RESERVES                    = "reserves"
	ACTION_IFRESERVES                  = "ifreserves"
	ACTION_FLASHPOOLBANNER             = "flashpoolbanner"
	ACTION_IFPOOLBANNER                = "ifpoolbanner"

	ACTION_FLASHPOOLDETAIL       = "flashpooldetail"
	ACTION_FLASHPOOLALLMARKET    = "flashpoolallmarket"
	ACTION_USERFLASHPOOLOVERVIEW = "userflashpooloverview"
	ACTION_IFPOOLDETAIL          = "ifpooldetail"

	ACTION_BORROWADDRESSLIST = "borrowaddresslist"

	ACTION_ASSETPRICE      = "assetprice"
	ACTION_ASSETPRICELIST  = "assetpricelist"
	ACTION_CLAIMWING       = "claimwing"
	ACTION_LIQUIDATIONLIST = "liquidationlist"
	ACTION_WINGAPYS        = "wingapys"

	ACTION_IFPOOLINFO    = "ifpoolinfo"
	ACTION_IFHISTORY     = "ifhistory"
	ACTION_IFLIQUIDATION = "ifliquidation"
)

type Response struct {
	Action string      `json:"action"`
	Desc   string      `json:"desc"`
	Error  uint32      `json:"error"`
	Result interface{} `json:"result"`
}

type AssetPriceRequest struct {
	Id    string
	Asset string
}

type AssetPriceResponse struct {
	Id    string
	Price string
}

type AssetPriceListRequest struct {
	Id        string
	AssetList []string
}

type AssetPriceListResponse struct {
	Id        string
	PriceList []string
}

type GovBannerOverview struct {
	Remain20 string
	Remain80 string
}

type GovBanner struct {
	Daily       string
	Distributed string
}

type MarketDistribution struct {
	MarketDistribution []*Distribution
}

type PoolDistribution struct {
	PoolDistribution []*Distribution
}

type Distribution struct {
	Icon            string
	Name            string
	SupplyAmount    string
	BorrowAmount    string
	InsuranceAmount string
	Total           string
}

type PoolBanner struct {
	Today string
	Share string
	Total string
}

type FlashPoolDetail struct {
	TotalSupply               string
	TotalBorrow               string
	TotalWingInsuranceBalance string
	TotalWingInsuranceDollar  string
}

type UserFlashPoolOverviewRequest struct {
	Id      string
	Address string
}

type UserFlashPoolOverviewResponse struct {
	Id                    string
	Address               string
	UserFlashPoolOverview *UserFlashPoolOverview
}

type UserFlashPoolOverview struct {
	BorrowLimit string
	NetApy      string

	CurrentSupply    []*Supply
	CurrentBorrow    []*Borrow
	CurrentInsurance []*Insurance

	AllMarket []*UserMarket
}

type Supply struct {
	Icon                  string
	Name                  string
	SupplyBalance         string
	Apy                   string
	CollateralFactor      string
	SupplyDistribution    string
	BorrowDistribution    string
	InsuranceDistribution string
	WingEarned            string
	IfCollateral          bool
}

type Borrow struct {
	Icon                  string
	Name                  string
	TotalBorrowBalance    string
	ValidBorrowBalance    string
	Apy                   string
	Limit                 string
	CollateralFactor      string
	SupplyDistribution    string
	BorrowDistribution    string
	InsuranceDistribution string
	WingEarned            string
	CollateralWing        string
}

type Insurance struct {
	Icon                  string
	Name                  string
	InsuranceBalance      string
	Apy                   string
	CollateralFactor      string
	SupplyDistribution    string
	BorrowDistribution    string
	InsuranceDistribution string
	WingEarned            string
}

type UserMarket struct {
	Icon                  string
	Name                  string
	SupplyApy             string
	BorrowApy             string
	BorrowLiquidity       string
	CollateralFactor      string
	SupplyDistribution    string
	BorrowDistribution    string
	InsuranceDistribution string
	IfCollateral          bool
}

type FlashPoolAllMarket struct {
	FlashPoolAllMarket []*Market
}

type Market struct {
	Icon                  string
	Name                  string `gorm:"primary_key"`
	TotalSupplyDollar     string
	TotalSupplyAmount     string
	SupplyApy             string
	TotalBorrowDollar     string
	TotalBorrowAmount     string
	BorrowApy             string
	TotalInsuranceDollar  string
	TotalInsuranceAmount  string
	InsuranceApy          string
	CollateralFactor      string
	SupplyDistribution    string
	BorrowDistribution    string
	InsuranceDistribution string
	ExchangeRate          string
	IExchangeRate         string
	BorrowIndex           string
}

type ClaimWingRequest struct {
	Id      string
	Address string
}

type ClaimWingResponse struct {
	Id      string
	Address string
	Amount  string
}

type Liquidation struct {
	Icon             string
	Name             string
	BorrowLimitUsed  string
	BorrowBalance    string
	BorrowDollar     string
	CollateralDollar string
	CollateralAssets []*CollateralAsset
}

type CollateralAsset struct {
	Icon    string
	Name    string
	Balance string
	Dollar  string
}

type IFPoolInfoRequest struct {
	Address string
}

type LiquidationListResponse struct {
	Id              string
	LiquidationList []*Liquidation
}

type LiquidationListRequest struct {
	Id      string
	Address string
}

type WingApys struct {
	InsuranceApy string
	WingApyList  []WingApy
}

type WingApy struct {
	AssetName    string `gorm:"primary_key"`
	SupplyApy    string
	BorrowApy    string
	InsuranceApy string
}

type Reserves struct {
	AssetReserve []*Reserve
	TotalReserve string
}

type Reserve struct {
	Name           string
	Icon           string
	ReserveFactor  string
	ReserveBalance string
	ReserveDollar  string
}

type IFPoolInfo struct {
	Total       string
	Cap         string
	IFAssetList []*IFAsset
	UserIFInfo  *UserIFInfo
}

type IFAsset struct {
	Name                 string
	Icon                 string
	Price                string
	TotalSupply          string
	SupplyInterestPerDay string
	SupplyWingAPY        string
	UtilizationRate      string
	TotalBorrowed        string
	BorrowWingAPY        string
	Liquidity            string
	BorrowCap            string
	TotalInsurance       string
	InsuranceWingAPY     string
}

type UserIFInfo struct {
	TotalSupplyDollar    string
	SupplyWingEarned     string
	TotalBorrowDollar    string
	BorrowWingEarned     string
	TotalInsuranceDollar string
	InsuranceWingEarned  string
	Composition          []*Composition
}

type Composition struct {
	Name                  string
	Icon                  string
	SupplyBalance         string
	SupplyWingEarned      string
	BorrowWingEarned      string
	LastBorrowTimestamp   uint64
	InsuranceBalance      string
	InsuranceWingEarned   string
	CollateralBalance     string
	BorrowUnpaidPrincipal string
	BorrowInterestBalance string
	BorrowExtraInterest   string
}

type IFHistoryRequest struct {
	Address        string
	Asset          string
	Operation      string
	StartTimestamp uint64
	EndTimestamp   uint64
	PageNo         uint64
	PageSize       uint64
}

type IFHistoryResponse struct {
	Count     uint64
	PageItems []*IFHistory
}

type IFHistory struct {
	Name      string
	Icon      string
	Operation string
	Timestamp uint64
	Balance   string
	Dollar    string
	Address   string
}

type IFLiquidationRequest struct {
	Start int64
	End   int64
}

type MarketUtility struct {
	UtilityMap map[common.Address]*big.Int
	Total      *big.Int
}

type IfMarketUtility struct {
	UtilityMap map[string]*big.Int
	Total      *big.Int
}

type PoolWeight struct {
	PoolStaticMap  map[common.Address]*big.Int
	PoolDynamicMap map[common.Address]*big.Int
	TotalStatic    *big.Int
	TotalDynamic   *big.Int
}

type DebtAccount struct {
	Address          string             `json:"address"`
	Debt             string             `json:"debt"`
	DebtIcon         string             `json:"debtIcon"`
	DebtAmount       string             `json:"debtAmount"`
	DebtPrice        string             `json:"debtPrice"`
	CollateralDollar string             `json:"collateralDollar"`
	CollateralAssets []*CollateralAsset `json:"collateralAssets"`
	BorrowTime       uint64             `json:"borrowTime"`
}

type IfPoolDetail struct {
	TotalSupply string
}
