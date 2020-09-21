package service

import (
	"github.com/siovanus/wingServer/config"
	"github.com/siovanus/wingServer/log"
)

type Service struct {
	cfg    *config.Config
	govMgr GovernanceManager
}

func NewService(govMgr GovernanceManager, cfg *config.Config) *Service {
	return &Service{cfg: cfg, govMgr: govMgr}
}

func (this *Service) Close() {
	log.Info("All connections closed. Bye!")
}
