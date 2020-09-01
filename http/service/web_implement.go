package service

import (
	"github.com/siovanus/wingServer/http/common"
	"github.com/siovanus/wingServer/http/restful"
	"github.com/siovanus/wingServer/log"
	"github.com/siovanus/wingServer/utils"
)

func (this *Service) QueryData(params map[string]interface{}) map[string]interface{} {
	req := &common.QueryDataReq{}
	resp := &common.Response{}
	err := utils.ParseParams(req, params)
	if err != nil {
		resp.Error = restful.INVALID_PARAMS
		resp.Desc = err.Error()
		log.Errorf("QueryData: decode params failed, err: %s", err)
	} else {
		replaceMap := make(map[string]string)
		for _, v := range req.AddressMap {
			replaceMap[v.OldAddress] = v.NewAddress
		}
		data1, voteData1, data2, voteData2, err := this.manager.QueryData(req.StartEpoch, req.EndEpoch, req.Sum, replaceMap)
		if err != nil {
			resp.Error = restful.INTERNAL_ERROR
			resp.Desc = err.Error()
			log.Errorf("QueryData: id %s, StartEpoch %d, EndEpoch:%d, Sum: %d, err: %s", req.Id, req.StartEpoch,
				req.EndEpoch, req.Sum, err)
		} else {
			resp.Error = restful.SUCCESS
			resp.Result = &common.QueryDataResp{
				Id:        req.Id,
				Data1:     data1,
				VoteData1: voteData1,
				Data2:     data2,
				VoteData2: voteData2,
			}
			log.Infof("QueryData success, id %s, StartEpoch %d, EndEpoch:%d, Sum: %d", req.Id, req.StartEpoch,
				req.EndEpoch, req.Sum)
		}
	}
	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("QueryData: failed, err: %s", err)
	} else {
		log.Debug("QueryData: resp success")
	}
	return m
}
