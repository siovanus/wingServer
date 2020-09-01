package service

import "github.com/siovanus/wingServer/manager"

type QueryManager interface {
	QueryData(startEpoch uint64, endEpoch uint64, sum uint64, addressMap map[string]string) ([]*manager.Data,
		*manager.VoteData, []*manager.Data, *manager.VoteData, error)
}
