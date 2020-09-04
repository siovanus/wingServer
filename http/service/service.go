package service

type Service struct {
	govMgr    GovernanceManager
	fpMgr     FlashPoolManager
	oracleMgr OracleManager
}

func NewService(govMgr GovernanceManager, fpMgr FlashPoolManager, oracleMgr OracleManager) *Service {
	return &Service{govMgr: govMgr, fpMgr: fpMgr, oracleMgr: oracleMgr}
}
