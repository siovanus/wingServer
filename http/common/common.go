package common

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
