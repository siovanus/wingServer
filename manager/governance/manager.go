package governance

import (
	"github.com/ontio/ontology/common"
)

type GovernanceManager struct {
	contractAddress common.Address
}

func NewGovernanceManager(contractAddress common.Address) *GovernanceManager {
	manager := &GovernanceManager{
		contractAddress,
	}

	return manager
}

func (this *GovernanceManager) GovBanner() (*GovBanner, error) {
	return this.govBanner()
}
