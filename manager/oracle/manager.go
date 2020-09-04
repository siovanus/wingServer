package oracle

import (
	"github.com/ontio/ontology/common"
)

type OracleManager struct {
	contractAddress common.Address
}

func NewOracleManager(contractAddress common.Address) *OracleManager {
	manager := &OracleManager{
		contractAddress,
	}

	return manager
}

func (this *OracleManager) AssetPrice(asset string) (uint64, error) {
	return this.assetPrice(asset)
}
