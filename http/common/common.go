package common

const (
	FLASHPOOLMARKETDISTRIBUTION = "/api/v1/flashpoolmarketdistribution"
	POOLDISTRIBUTION            = "/api/v1/pooldistribution"
	GOVBANNEROVERVIEW           = "/api/v1/govbanneroverview"
	GOVBANNER                   = "/api/v1/govbanner"
	FLASHPOOLBANNER             = "/api/v1/flashpoolbanner"

	ASSETPRICE = "/api/v1/assetprice"
)

const (
	ACTION_FLASHPOOLMARKETDISTRIBUTION = "flashpoolmarketdistribution"
	ACTION_POOLDISTRIBUTION            = "pooldistribution"
	ACTION_GOVBANNEROVERVIEW           = "govbanneroverview"
	ACTION_GOVBANNER                   = "govbanner"
	ACTION_FLASHPOOLBANNER             = "flashpoolbanner"

	ACTION_ASSETPRICE = "assetprice"
)

type Response struct {
	Action string      `json:"action"`
	Desc   string      `json:"desc"`
	Error  uint32      `json:"error"`
	Result interface{} `json:"result"`
}

type AssetPriceRequest struct {
	Id    string `json:"id"`
	Asset string `json:"asset"`
}

type AssetPriceResponse struct {
	Id    string `json:"id"`
	Price uint64 `json:"price"`
}
