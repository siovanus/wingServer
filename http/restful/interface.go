package restful

type Web interface {
	QueryData(map[string]interface{}) map[string]interface{}
}
