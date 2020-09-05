package oracle

import (
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
)

type OracleManager struct {
	contractAddress common.Address
	sdk             *sdk.OntologySdk
}

func NewOracleManager(contractAddress common.Address, sdk *sdk.OntologySdk) *OracleManager {
	manager := &OracleManager{
		contractAddress: contractAddress,
		sdk:             sdk,
	}

	return manager
}

func (this *OracleManager) AssetPrice(asset string) (uint64, error) {
	return this.assetPrice(asset)
}
