package service

import (
	"github.com/siovanus/wingServer/http/common"
	"github.com/siovanus/wingServer/http/restful"
	"github.com/siovanus/wingServer/log"
	"github.com/siovanus/wingServer/utils"
)

func (this *Service) FlashPoolMarketDistribution(param map[string]interface{}) map[string]interface{} {
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

func (this *Service) GovBannerOverview(param map[string]interface{}) map[string]interface{} {
	resp := &common.Response{}
	govBanner, err := this.govMgr.GovBannerOverview()
	if err != nil {
		resp.Error = restful.INTERNAL_ERROR
		resp.Desc = err.Error()
		log.Errorf("GovBannerOverview error: %s", err)
	} else {
		resp.Error = restful.SUCCESS
		resp.Result = govBanner
		log.Infof("GovBannerOverview success")
	}

	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("GovBannerOverview: failed, err: %s", err)
	} else {
		log.Debug("GovBannerOverview: resp success")
	}
	return m
}

func (this *Service) GovBanner(param map[string]interface{}) map[string]interface{} {
	resp := &common.Response{}
	poolBanner, err := this.govMgr.GovBanner()
	if err != nil {
		resp.Error = restful.INTERNAL_ERROR
		resp.Desc = err.Error()
		log.Errorf("GovBanner error: %s", err)
	} else {
		resp.Error = restful.SUCCESS
		resp.Result = poolBanner
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

func (this *Service) AssetPrice(param map[string]interface{}) map[string]interface{} {
	req := &common.AssetPriceRequest{}
	resp := &common.Response{}
	err := utils.ParseParams(req, param)
	if err != nil {
		resp.Error = restful.INVALID_PARAMS
		resp.Desc = err.Error()
		log.Errorf("AssetPrice: decode params failed, err: %s", err)
	} else {
		assetPrice, err := this.oracleMgr.AssetPrice(req.Asset)
		if err != nil {
			resp.Error = restful.INTERNAL_ERROR
			resp.Desc = err.Error()
			log.Errorf("AssetPrice error: %s", err)
		} else {
			resp.Error = restful.SUCCESS
			resp.Result = assetPrice
			log.Infof("AssetPrice success")
		}
	}

	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("AssetPrice: failed, err: %s", err)
	} else {
		log.Debug("AssetPrice: resp success")
	}
	return m
}

func (this *Service) FlashPoolBanner(param map[string]interface{}) map[string]interface{} {
	resp := &common.Response{}
	flashPoolBanner, err := this.fpMgr.FlashPoolBanner()
	if err != nil {
		resp.Error = restful.INTERNAL_ERROR
		resp.Desc = err.Error()
		log.Errorf("FlashPoolBanner error: %s", err)
	} else {
		resp.Error = restful.SUCCESS
		resp.Result = flashPoolBanner
		log.Infof("FlashPoolBanner success")
	}

	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("FlashPoolBanner: failed, err: %s", err)
	} else {
		log.Debug("FlashPoolBanner: resp success")
	}
	return m
}

func (this *Service) FlashPoolDetail(param map[string]interface{}) map[string]interface{} {
	resp := &common.Response{}
	flashPoolDetail, err := this.fpMgr.FlashPoolDetail()
	if err != nil {
		resp.Error = restful.INTERNAL_ERROR
		resp.Desc = err.Error()
		log.Errorf("FlashPoolDetail error: %s", err)
	} else {
		resp.Error = restful.SUCCESS
		resp.Result = flashPoolDetail
		log.Infof("FlashPoolDetail success")
	}

	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("FlashPoolDetail: failed, err: %s", err)
	} else {
		log.Debug("FlashPoolDetail: resp success")
	}
	return m
}

func (this *Service) FlashPoolAllMarket(param map[string]interface{}) map[string]interface{} {
	resp := &common.Response{}
	flashPoolAllMarket, err := this.fpMgr.FlashPoolAllMarket()
	if err != nil {
		resp.Error = restful.INTERNAL_ERROR
		resp.Desc = err.Error()
		log.Errorf("FlashPoolAllMarket error: %s", err)
	} else {
		resp.Error = restful.SUCCESS
		resp.Result = flashPoolAllMarket
		log.Infof("FlashPoolAllMarket success")
	}

	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("FlashPoolAllMarket: failed, err: %s", err)
	} else {
		log.Debug("FlashPoolAllMarket: resp success")
	}
	return m
}