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
	JsonRpcAddress      string            `json:"json_rpc_address"`
	Port                uint64            `json:"port"`
	GovernanceAddress   string            `json:"governance_address"`
	WingAddress         string            `json:"wing_address"`
	FlashPoolAddress    string            `json:"flash_pool_address"`
	IFPoolAddress       string            `json:"if_pool_address"`
	OracleAddress       string            `json:"oracle_address"`
	OscoreOracleAddress string            `json:"oscore_oracle_address"`
	DatabaseURL         string            `json:"database_url"`
	IconMap             map[string]string `json:"icon_map"`
	FlashAssetMap       map[string]string `json:"flash_asset_map"`
	IFMap               map[string]string `json:"if_map"`
	IFOracleMap         map[string]string `json:"if_oracle_map"`
	TrackEventInterval  uint64            `json:"track_event_interval"`
	SystemContract      []string          `json:"system_contract"`
	TokenDecimal        map[string]uint64 `json:"token_decimal"`
	ScanInterval        uint64            `json:"scan_interval"`
	SnapshotInterval    uint64            `json:"snapshot_interval"`
	OneDaySecond        int64             `json:"one_day_second"`
	MonitorIfDebt       string            `json:"monitor_if_debt"`
	WingBackendUrl      string            `json:"wing_backend_url"`
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
