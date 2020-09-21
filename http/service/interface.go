package service

import "github.com/siovanus/wingServer/http/common"

type GovernanceManager interface {
	Wing() (*common.Wing, error)
}
