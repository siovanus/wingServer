package common

const (
	FLASHPOOLMARKETDISTRIBUTION = "/api/v1/flashpoolmarketdistribution"
	POOLDISTRIBUTION            = "/api/v1/pooldistribution"
	GOVBANNEROVERVIEW           = "/api/v1/govbanneroverview"
	GOVBANNER                   = "/api/v1/govbanner"
	RESERVES                    = "/api/v1/reserves"
	FLASHPOOLBANNER             = "/api/v1/flashpoolbanner"

	FLASHPOOLDETAIL       = "/api/v1/flashpooldetail"
	FLASHPOOLALLMARKET    = "/api/v1/flashpoolallmarket"
	USERFLASHPOOLOVERVIEW = "/api/v1/userflashpooloverview"

	BORROWADDRESSLIST = "/api/v1/borrowaddresslist"

	ASSETPRICE      = "/api/v1/assetprice"
	ASSETPRICELIST  = "/api/v1/assetpricelist"
	CLAIMWING       = "/api/v1/claimwing"
	LIQUIDATIONLIST = "/api/v1/liquidationlist"
	WINGAPYS        = "/api/v1/wingapys"

	IFPOOLINFO = "/api/v1/if/ifpoolinfo"
)

const (
	ACTION_FLASHPOOLMARKETDISTRIBUTION = "flashpoolmarketdistribution"
	ACTION_POOLDISTRIBUTION            = "pooldistribution"
	ACTION_GOVBANNEROVERVIEW           = "govbanneroverview"
	ACTION_GOVBANNER                   = "govbanner"
	ACTION_RESERVES                    = "reserves"
	ACTION_FLASHPOOLBANNER             = "flashpoolbanner"

	ACTION_FLASHPOOLDETAIL       = "flashpooldetail"
	ACTION_FLASHPOOLALLMARKET    = "flashpoolallmarket"
	ACTION_USERFLASHPOOLOVERVIEW = "userflashpooloverview"

	ACTION_BORROWADDRESSLIST = "borrowaddresslist"

	ACTION_ASSETPRICE      = "assetprice"
	ACTION_ASSETPRICELIST  = "assetpricelist"
	ACTION_CLAIMWING       = "claimwing"
	ACTION_LIQUIDATIONLIST = "liquidationlist"
	ACTION_WINGAPYS        = "wingapys"

	ACTION_IFPOOLINFO = "ifpoolinfo"
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

type FlashPoolMarketDistribution struct {
	FlashPoolMarketDistribution []*Distribution
}

type PoolDistribution struct {
	PoolDistribution []*Distribution
}

type Distribution struct {
	Icon            string
	Name            string
	PerDay          string
	SupplyAmount    string
	BorrowAmount    string
	InsuranceAmount string
	Total           string
}

type FlashPoolBanner struct {
	Today string
	Share string
	Total string
}

type FlashPoolDetail struct {
	TotalSupply    string
	TotalBorrow    string
	TotalInsurance string
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
	BorrowBalance         string
	Apy                   string
	Limit                 string
	CollateralFactor      string
	SupplyDistribution    string
	BorrowDistribution    string
	InsuranceDistribution string
	WingEarned            string
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
	InsuranceApy          string
	InsuranceAmount       string
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
	IFAssetList []*IFAssetList
	UserIFInfo  *UserIFInfo
}

type IFAssetList struct {
	Name                 string
	Icon                 string
	Price                string
	TotalSupply          string
	SupplyInterestPerDay string
	SupplyWingAPY        string
	UtilizationRate      string
	MaximumLTV           string
	TotalBorrowed        string
	BorrowInterestPerDay string
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
	BorrowInterestPerDay string
	TotalInsuranceDollar string
	InsuranceWingEarned  string
	Composition          []*Composition
}

type Composition struct {
	Name                  string
	Icon                  string
	IfCanBorrow           bool
	SupplyBalance         string
	SupplyWingEarned      string
	BorrowWingEarned      string
	InsuranceBalance      string
	InsuranceWingEarned   string
	CollateralBalance     string
	CollateralName        string
	CollateralIcon        string
	BorrowUnpaidPrincipal string
	BorrowInterestBalance string
}
