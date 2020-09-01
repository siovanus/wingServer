package common

import "github.com/siovanus/wingServer/manager"

const (
	QUERYDATA = "/api/v1/querydata"
)

const (
	ACTION_QUERYDATA = "querydata"
)

type Response struct {
	Action string      `json:"action"`
	Desc   string      `json:"desc"`
	Error  uint32      `json:"error"`
	Result interface{} `json:"result"`
}

type QueryDataReq struct {
	Id         string       `json:"id"`
	StartEpoch uint64       `json:"start_epoch"`
	EndEpoch   uint64       `json:"end_epoch"`
	Sum        uint64       `json:"sum"`
	AddressMap []AddressMap `json:"address_map"`
}

type QueryDataResp struct {
	Id        string            `json:"id"`
	Data1     []*manager.Data   `json:"data_1"`
	VoteData1 *manager.VoteData `json:"vote_data_1"`
	Data2     []*manager.Data   `json:"data_2"`
	VoteData2 *manager.VoteData `json:"vote_data_2"`
}

type AddressMap struct {
	OldAddress string `json:"old_address"`
	NewAddress string `json:"new_address"`
}
