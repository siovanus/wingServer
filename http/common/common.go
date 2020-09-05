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
