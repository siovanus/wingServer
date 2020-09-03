package service

import (
	"github.com/siovanus/wingServer/http/common"
	"github.com/siovanus/wingServer/http/restful"
	"github.com/siovanus/wingServer/log"
	"github.com/siovanus/wingServer/utils"
)

func (this *Service) MarketDistribution(param map[string]interface{}) map[string]interface{} {
	resp := &common.Response{}
	marketDistribution, err := this.fpMgr.MarketDistribution()
	if err != nil {
		resp.Error = restful.INTERNAL_ERROR
		resp.Desc = err.Error()
		log.Errorf("MarketDistribution error: %s", err)
	} else {
		resp.Error = restful.SUCCESS
		resp.Result = &common.MarketDistributionResp{
			MarketDistribution: marketDistribution,
		}
		log.Infof("MarketDistribution success")
	}

	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("MarketDistribution: failed, err: %s", err)
	} else {
		log.Debug("MarketDistribution: resp success")
	}
	return m
}

func (this *Service) PoolDistribution(param map[string]interface{}) map[string]interface{} {
	resp := &common.Response{}
	poolDistribution, err := this.fpMgr.PoolDistribution()
	if err != nil {
		resp.Error = restful.INTERNAL_ERROR
		resp.Desc = err.Error()
		log.Errorf("PoolDistribution error: %s", err)
	} else {
		resp.Error = restful.SUCCESS
		resp.Result = &common.PoolDistributionResp{
			PoolDistributionResp: poolDistribution,
		}
		log.Infof("PoolDistribution success")
	}

	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("PoolDistribution: failed, err: %s", err)
	} else {
		log.Debug("PoolDistribution: resp success")
	}
	return m
}