package common

const (
	FLASHPOOLMARKETDISTRIBUTION = "/api/v1/flashpoolmarketdistribution"
	POOLDISTRIBUTION            = "/api/v1/pooldistribution"
	GOVBANNEROVERVIEW           = "/api/v1/govbanneroverview"
	GOVBANNER                   = "/api/v1/govbanner"
	FLASHPOOLBANNER             = "/api/v1/flashpoolbanner"

	FLASHPOOLALLMARKET    = "/api/v1/flashpoolallmarket"
	USERFLASHPOOLOVERVIEW = "/api/v1/userflashpooloverview"

	ASSETPRICE     = "/api/v1/assetprice"
	ASSETPRICELIST = "/api/v1/assetpricelist"
	CLAIMWING      = "/api/v1/claimwing"
)

const (
	ACTION_FLASHPOOLMARKETDISTRIBUTION = "flashpoolmarketdistribution"
	ACTION_POOLDISTRIBUTION            = "pooldistribution"
	ACTION_GOVBANNEROVERVIEW           = "govbanneroverview"
	ACTION_GOVBANNER                   = "govbanner"
	ACTION_FLASHPOOLBANNER             = "flashpoolbanner"

	ACTION_FLASHPOOLALLMARKET    = "flashpoolallmarket"
	ACTION_USERFLASHPOOLOVERVIEW = "userflashpooloverview"

	ACTION_ASSETPRICE     = "assetprice"
	ACTION_ASSETPRICELIST = "assetpricelist"
	ACTION_CLAIMWING      = "claimwing"
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

type MarketFund struct {
	Icon string
	Name string
	Fund string
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
	TotalSupply    string
	TotalBorrow    string
	TotalInsurance string

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
	TotalSupply           string
	SupplyApy             string
	TotalBorrow           string
	BorrowApy             string
	TotalInsurance        string
	InsuranceApy          string
	CollateralFactor      string
	SupplyDistribution    string
	BorrowDistribution    string
	InsuranceDistribution string
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

type BorrowAddress struct {
	AddressList []string
}

type Liquidation struct {
	Address          string
	BorrowLimitUsed  string
	BorrowBalance    string
	BorrowDollar     string
	CollateralDollar string
	CollateralAsset  []*CollateralAsset
}

type LiquidationList struct {
	Liquidation *Liquidation
}

type CollateralAsset struct {
	Icon    string
	Name    string
	Balance string
	Dollar  string
}
