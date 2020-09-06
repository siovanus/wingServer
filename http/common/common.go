package common

const (
	FLASHPOOLMARKETDISTRIBUTION = "/api/v1/flashpoolmarketdistribution"
	POOLDISTRIBUTION            = "/api/v1/pooldistribution"
	GOVBANNEROVERVIEW           = "/api/v1/govbanneroverview"
	GOVBANNER                   = "/api/v1/govbanner"
	FLASHPOOLBANNER             = "/api/v1/flashpoolbanner"

	FLASHPOOLDETAIL       = "/api/v1/flashpooldetail"
	FLASHPOOLALLMARKET    = "/api/v1/flashpoolallmarket"
	USERFLASHPOOLOVERVIEW = "/api/v1/userflashpooloverview"

	ASSETPRICE = "/api/v1/assetprice"
)

const (
	ACTION_FLASHPOOLMARKETDISTRIBUTION = "flashpoolmarketdistribution"
	ACTION_POOLDISTRIBUTION            = "pooldistribution"
	ACTION_GOVBANNEROVERVIEW           = "govbanneroverview"
	ACTION_GOVBANNER                   = "govbanner"
	ACTION_FLASHPOOLBANNER             = "flashpoolbanner"

	ACTION_FLASHPOOLDETAIL       = "flashpooldetail"
	ACTION_FLASHPOOLALLMARKET    = "flashpoolallmarket"
	ACTION_USERFLASHPOOLOVERVIEW = "userflashpooloverview"

	ACTION_ASSETPRICE = "assetprice"
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
	Price uint64
}

type GovBannerOverview struct {
	Remain20 uint64
	Remain80 uint64
}

type GovBanner struct {
	Daily       uint64
	Distributed uint64
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
	PerDay          uint64
	SupplyAmount    uint64
	BorrowAmount    uint64
	InsuranceAmount uint64
	Total           uint64
}

type FlashPoolBanner struct {
	Today uint64
	Share uint64
	Total uint64
}

type FlashPoolDetail struct {
	TotalSupply       uint64
	TotalSupplyRate   int64
	SupplyMarketRank  []*MarketFund
	SupplyVolumeDaily uint64

	TotalBorrow       uint64
	TotalBorrowRate   int64
	BorrowMarketRank  []*MarketFund
	BorrowVolumeDaily uint64

	TotalInsurance       uint64
	TotalInsuranceRate   int64
	InsuranceMarketRank  []*MarketFund
	InsuranceVolumeDaily uint64
}

type MarketFund struct {
	Icon string
	Name string
	Fund uint64
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
	SupplyBalance    uint64
	BorrowBalance    uint64
	InsuranceBalance uint64
	BorrowLimit      uint64
	NetApy           uint64

	CurrentSupply    []*Supply
	CurrentBorrow    []*Borrow
	CurrentInsurance []*Insurance

	AllMarket []*UserMarket
}

type Supply struct {
	Icon          string
	Name          string
	SupplyDollar  uint64
	SupplyBalance uint64
	Apy           uint64
	Earned        uint64
	IfCollateral  bool
}

type Borrow struct {
	Icon          string
	Name          string
	BorrowDollar  uint64
	BorrowBalance uint64
	Apy           uint64
	Accrued       uint64
	Limit         uint64
}

type Insurance struct {
	Icon             string
	Name             string
	InsuranceDollar  uint64
	InsuranceBalance uint64
	Apy              uint64
}

type UserMarket struct {
	Icon            string
	Name            string
	IfCollateral    bool
	SupplyApy       uint64
	BorrowApy       uint64
	BorrowLiquidity uint64
	InsuranceApy    uint64
	InsuranceAmount uint64
}

type FlashPoolAllMarket struct {
	FlashPoolAllMarket []*Market
}

type Market struct {
	Icon               string
	Name               string
	TotalSupply        uint64
	TotalSupplyRate    uint64
	SupplyApy          uint64
	SupplyApyRate      uint64
	TotalBorrow        uint64
	TotalBorrowRate    uint64
	BorrowApy          uint64
	BorrowApyRate      uint64
	TotalInsurance     uint64
	TotalInsuranceRate uint64
	InsuranceApy       uint64
	InsuranceApyRate   uint64
}
