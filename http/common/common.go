package common

import "github.com/siovanus/wingServer/manager/flashpool"

const (
	MARKETDISTRIBUTION = "/api/v1/marketdistribution"
	POOLDISTRIBUTION = "/api/v1/pooldistribution"
)

const (
	ACTION_MARKETDISTRIBUTION = "marketdistribution"
	ACTION_POOLDISTRIBUTION = "pooldistribution"
)

type Response struct {
	Action string      `json:"action"`
	Desc   string      `json:"desc"`
	Error  uint32      `json:"error"`
	Result interface{} `json:"result"`
}

type MarketDistributionResp struct {
	MarketDistribution *flashpool.MarketDistribution `json:"market_distribution"`
}

type PoolDistributionResp struct {
	PoolDistributionResp *flashpool.PoolDistribution `json:"pool_distribution_resp"`
}