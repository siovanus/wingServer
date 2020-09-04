package common

const (
	MARKETDISTRIBUTION = "/api/v1/marketdistribution"
	POOLDISTRIBUTION   = "/api/v1/pooldistribution"
	GOVBANNER          = "/api/v1/govbanner"
	POOLBANNER         = "/api/v1/poolbanner"
)

const (
	ACTION_MARKETDISTRIBUTION = "marketdistribution"
	ACTION_POOLDISTRIBUTION   = "pooldistribution"
	ACTION_GOVBANNER          = "govbanner"
	ACTION_POOLBANNER         = "poolbanner"
)

type Response struct {
	Action string      `json:"action"`
	Desc   string      `json:"desc"`
	Error  uint32      `json:"error"`
	Result interface{} `json:"result"`
}
