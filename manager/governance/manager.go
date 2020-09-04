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

func (this *GovernanceManager) GovBannerOverview() (*GovBanner, error) {
	return this.govBannerOverview()
}

func (this *GovernanceManager) GovBanner() (*PoolBanner, error) {
	return this.govBanner()
}
