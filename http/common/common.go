package common

const (
	CIRCULATINGSUPPLY = "/api/v1/wing/circulating-supply"
	TOTALSUPPLY       = "/api/v1/wing/total-supply"
)

const (
	ACTION_CIRCULATINGSUPPLY = "circulating-supply"
	ACTION_TOTALSUPPLY       = "total-supply"
)

type Wing struct {
	Total       float64
	Circulating float64
}
