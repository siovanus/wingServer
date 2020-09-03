package restful

type Web interface {
	MarketDistribution(map[string]interface{}) map[string]interface{}
	PoolDistribution(map[string]interface{}) map[string]interface{}
}
