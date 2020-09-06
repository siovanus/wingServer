package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	DEFAULT_LOG_LEVEL        = 2
	DEFAULT_CONFIG_FILE_NAME = "./config.json"
)

//Config object used by ontology-instance
type Config struct {
	JsonRpcAddress    string `json:"json_rpc_address"`
	Port              uint64 `json:"port"`
	GovernanceAddress string `json:"governance_address"`
	WingAddress       string `json:"wing_address"`
	FlashPoolAddress  string `json:"flash_pool_address"`
	OracleAddress     string `json:"oracle_address"`
	DatabaseURL       string `json:"database_url"`
}

func NewConfig(fileName string) (*Config, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal Config:%s error:%s", data, err)
	}
	return cfg, nil
}
