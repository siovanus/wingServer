package service

type Service struct {
	govMgr GovernanceManager
	fpMgr  FlashPoolManager
}

func NewService(govMgr GovernanceManager, fpMgr FlashPoolManager) *Service {
	return &Service{govMgr: govMgr, fpMgr: fpMgr}
}
