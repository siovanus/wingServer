package common

import "github.com/siovanus/wingServer/manager/flashpool"

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

type UserFlashPoolOverviewRequest struct {
	Id      string
	Address string
}

type UserFlashPoolOverviewResponse struct {
	Id                    string
	Address               string
	UserFlashPoolOverview *flashpool.UserFlashPoolOverview
}
