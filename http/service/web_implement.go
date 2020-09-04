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
		resp.Result = marketDistribution
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
		resp.Result = poolDistribution
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

func (this *Service) GovBanner(param map[string]interface{}) map[string]interface{} {
	resp := &common.Response{}
	govBanner, err := this.govMgr.GovBanner()
	if err != nil {
		resp.Error = restful.INTERNAL_ERROR
		resp.Desc = err.Error()
		log.Errorf("GovBanner error: %s", err)
	} else {
		resp.Error = restful.SUCCESS
		resp.Result = govBanner
		log.Infof("GovBanner success")
	}

	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("GovBanner: failed, err: %s", err)
	} else {
		log.Debug("GovBanner: resp success")
	}
	return m
}

func (this *Service) PoolBanner(param map[string]interface{}) map[string]interface{} {
	resp := &common.Response{}
	poolBanner, err := this.fpMgr.PoolBanner()
	if err != nil {
		resp.Error = restful.INTERNAL_ERROR
		resp.Desc = err.Error()
		log.Errorf("PoolBanner error: %s", err)
	} else {
		resp.Error = restful.SUCCESS
		resp.Result = poolBanner
		log.Infof("PoolBanner success")
	}

	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("PoolBanner: failed, err: %s", err)
	} else {
		log.Debug("PoolBanner: resp success")
	}
	return m
}
