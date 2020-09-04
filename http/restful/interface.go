package restful

type Web interface {
	MarketDistribution(map[string]interface{}) map[string]interface{}
	PoolDistribution(map[string]interface{}) map[string]interface{}
	GovBanner(map[string]interface{}) map[string]interface{}
	PoolBanner(map[string]interface{}) map[string]interface{}
}
