package restful

type Web interface {
	CirculatingSupply(map[string]interface{}) interface{}
	TotalSupply(map[string]interface{}) interface{}
}
