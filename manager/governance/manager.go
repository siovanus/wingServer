package governance

import (
	"github.com/siovanus/wingServer/config"
)

type GovernanceManager struct {
	cfg *config.Config
}

func NewGovernanceManager(cfg *config.Config) *GovernanceManager {
	queryManager := &GovernanceManager{
		cfg: cfg,
	}

	return queryManager
}

func (this *GovernanceManager) QueryData(startEpoch uint64, endEpoch uint64, sum uint64, replaceMap map[string]string) {
	return this.getData()
}
