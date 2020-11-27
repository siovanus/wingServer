package service

import (
	"fmt"
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
		log.Errorf("FlashPoolMarketDistribution error: %s", err)
	} else {
		resp.Error = restful.SUCCESS
		resp.Result = marketDistribution
		log.Infof("FlashPoolMarketDistribution success")
	}

	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("FlashPoolMarketDistribution: failed, err: %s", err)
	} else {
		log.Debug("FlashPoolMarketDistribution: resp success")
	}
	return m
}

func (this *Service) IfPoolMarketDistribution(param map[string]interface{}) map[string]interface{} {
	resp := &common.Response{}
	marketDistribution, err := this.ifMgr.MarketDistribution()
	if err != nil {
		resp.Error = restful.INTERNAL_ERROR
		resp.Desc = err.Error()
		log.Errorf("IfPoolMarketDistribution error: %s", err)
	} else {
		resp.Error = restful.SUCCESS
		resp.Result = marketDistribution
		log.Infof("IfPoolMarketDistribution success")
	}

	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("IfPoolMarketDistribution: failed, err: %s", err)
	} else {
		log.Debug("IfPoolMarketDistribution: resp success")
	}
	return m
}

func (this *Service) PoolDistribution(param map[string]interface{}) map[string]interface{} {
	resp := &common.Response{}
	flag := true
	flashPoolDistribution, err := this.fpMgr.PoolDistribution()
	if err != nil {
		resp.Error = restful.INTERNAL_ERROR
		resp.Desc = err.Error()
		log.Errorf("PoolDistribution error: %s", err)
		flag = false
	}
	ifPoolDistribution, err := this.ifMgr.PoolDistribution()
	if err != nil {
		resp.Error = restful.INTERNAL_ERROR
		resp.Desc = err.Error()
		log.Errorf("PoolDistribution error: %s", err)
		flag = false
	}
	if flag {
		resp.Error = restful.SUCCESS
		resp.Result = &common.PoolDistribution{PoolDistribution: []*common.Distribution{flashPoolDistribution, ifPoolDistribution}}
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

func (this *Service) IfPoolDistribution(param map[string]interface{}) map[string]interface{} {
	resp := &common.Response{}
	flashPoolDistribution, err := this.ifMgr.PoolDistribution()
	if err != nil {
		resp.Error = restful.INTERNAL_ERROR
		resp.Desc = err.Error()
		log.Errorf("IfPoolDistribution error: %s", err)
	} else {
		resp.Error = restful.SUCCESS
		resp.Result = &common.PoolDistribution{PoolDistribution: []*common.Distribution{flashPoolDistribution}}
		log.Infof("IfPoolDistribution success")
	}

	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("IfPoolDistribution: failed, err: %s", err)
	} else {
		log.Debug("IfPoolDistribution: resp success")
	}
	return m
}

func (this *Service) GovBannerOverview(param map[string]interface{}) map[string]interface{} {
	resp := &common.Response{}
	govBannerOverview, err := this.govMgr.GovBannerOverview()
	if err != nil {
		resp.Error = restful.INTERNAL_ERROR
		resp.Desc = err.Error()
		log.Errorf("GovBannerOverview error: %s", err)
	} else {
		resp.Error = restful.SUCCESS
		resp.Result = govBannerOverview
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

func (this *Service) AssetPrice(param map[string]interface{}) map[string]interface{} {
	req := &common.AssetPriceRequest{}
	resp := &common.Response{}
	err := utils.ParseParams(req, param)
	if err != nil {
		resp.Error = restful.INVALID_PARAMS
		resp.Desc = err.Error()
		log.Errorf("AssetPrice: decode params failed, err: %s", err)
	} else {
		assetPrice, err := this.fpMgr.AssetPrice(req.Asset)
		if err != nil {
			resp.Error = restful.INTERNAL_ERROR
			resp.Desc = err.Error()
			log.Errorf("AssetPrice error: %s", err)
		} else {
			resp.Error = restful.SUCCESS
			resp.Result = &common.AssetPriceResponse{
				Id:    req.Id,
				Price: assetPrice,
			}
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

func (this *Service) AssetPriceList(param map[string]interface{}) map[string]interface{} {
	req := &common.AssetPriceListRequest{}
	resp := &common.Response{}
	err := utils.ParseParams(req, param)
	if err != nil {
		resp.Error = restful.INVALID_PARAMS
		resp.Desc = err.Error()
		log.Errorf("AssetPriceList: decode params failed, err: %s", err)
	} else {
		suc := true
		var asset string
		priceList := make([]string, 0)
		for _, v := range req.AssetList {
			assetPrice, err := this.fpMgr.AssetPrice(v)
			if err != nil {
				suc = false
				asset = v
			}
			priceList = append(priceList, assetPrice)
		}
		if !suc {
			resp.Error = restful.INTERNAL_ERROR
			resp.Desc = fmt.Errorf("AssetPriceList get asset price %s error", asset).Error()
			log.Errorf("AssetPriceList get asset price %s error", asset)
		} else {
			resp.Error = restful.SUCCESS
			resp.Result = &common.AssetPriceListResponse{
				Id:        req.Id,
				PriceList: priceList,
			}
			log.Infof("AssetPriceList success")
		}
	}

	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("AssetPriceList: failed, err: %s", err)
	} else {
		log.Debug("AssetPriceList: resp success")
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

func (this *Service) UserFlashPoolOverview(param map[string]interface{}) map[string]interface{} {
	req := &common.UserFlashPoolOverviewRequest{}
	resp := &common.Response{}
	err := utils.ParseParams(req, param)
	if err != nil {
		resp.Error = restful.INVALID_PARAMS
		resp.Desc = err.Error()
		log.Errorf("UserFlashPoolOverview: decode params failed, err: %s", err)
	} else {
		userFlashPoolOverview, err := this.fpMgr.UserFlashPoolOverview(req.Address)
		if err != nil {
			resp.Error = restful.INTERNAL_ERROR
			resp.Desc = err.Error()
			log.Errorf("UserFlashPoolOverview error: %s", err)
		} else {
			resp.Error = restful.SUCCESS
			resp.Result = &common.UserFlashPoolOverviewResponse{
				Id:                    req.Id,
				Address:               req.Address,
				UserFlashPoolOverview: userFlashPoolOverview,
			}
			log.Infof("UserFlashPoolOverview success")
		}
	}

	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("UserFlashPoolOverview: failed, err: %s", err)
	} else {
		log.Debug("UserFlashPoolOverview: resp success")
	}
	return m
}

func (this *Service) ClaimWing(param map[string]interface{}) map[string]interface{} {
	req := &common.ClaimWingRequest{}
	resp := &common.Response{}
	err := utils.ParseParams(req, param)
	if err != nil {
		resp.Error = restful.INVALID_PARAMS
		resp.Desc = err.Error()
		log.Errorf("ClaimWing: decode params failed, err: %s", err)
	} else {
		amount, err := this.fpMgr.ClaimWing(req.Address)
		if err != nil {
			resp.Error = restful.INTERNAL_ERROR
			resp.Desc = err.Error()
			log.Errorf("ClaimWing error: %s", err)
		} else {
			resp.Error = restful.SUCCESS
			resp.Result = &common.ClaimWingResponse{
				Id:      req.Id,
				Address: req.Address,
				Amount:  amount,
			}
			log.Infof("ClaimWing success")
		}
	}

	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("ClaimWing: failed, err: %s", err)
	} else {
		log.Debug("ClaimWing: resp success")
	}
	return m
}

func (this *Service) BorrowAddressList(param map[string]interface{}) map[string]interface{} {
	resp := &common.Response{}
	borrowAddressList, err := this.fpMgr.BorrowAddressList()
	if err != nil {
		resp.Error = restful.INTERNAL_ERROR
		resp.Desc = err.Error()
		log.Errorf("BorrowAddressList error: %s", err)
	} else {
		resp.Error = restful.SUCCESS
		resp.Result = borrowAddressList
		log.Infof("BorrowAddressList success")
	}

	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("BorrowAddressList: failed, err: %s", err)
	} else {
		log.Debug("BorrowAddressList: resp success")
	}
	return m
}

func (this *Service) LiquidationList(param map[string]interface{}) map[string]interface{} {
	req := &common.LiquidationListRequest{}
	resp := &common.Response{}
	err := utils.ParseParams(req, param)
	if err != nil {
		resp.Error = restful.INVALID_PARAMS
		resp.Desc = err.Error()
		log.Errorf("LiquidationList: decode params failed, err: %s", err)
	} else {
		liquidationList, err := this.fpMgr.LiquidationList(req.Address)
		if err != nil {
			resp.Error = restful.INTERNAL_ERROR
			resp.Desc = err.Error()
			log.Errorf("LiquidationList error: %s", err)
		} else {
			resp.Error = restful.SUCCESS
			resp.Result = &common.LiquidationListResponse{
				Id:              req.Id,
				LiquidationList: liquidationList,
			}
			log.Infof("LiquidationList success")
		}
	}

	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("LiquidationList: failed, err: %s", err)
	} else {
		log.Debug("LiquidationList: resp success")
	}
	return m
}

func (this *Service) WingApys(param map[string]interface{}) map[string]interface{} {
	resp := &common.Response{}
	wingApys, err := this.fpMgr.WingApys()
	if err != nil {
		resp.Error = restful.INTERNAL_ERROR
		resp.Desc = err.Error()
		log.Errorf("WingApys error: %s", err)
	} else {
		resp.Error = restful.SUCCESS
		resp.Result = wingApys
		log.Infof("WingApys success")
	}

	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("WingApys: failed, err: %s", err)
	} else {
		log.Debug("WingApys: resp success")
	}
	return m
}

func (this *Service) Reserves(param map[string]interface{}) map[string]interface{} {
	resp := &common.Response{}
	reserves, err := this.fpMgr.Reserves()
	if err != nil {
		resp.Error = restful.INTERNAL_ERROR
		resp.Desc = err.Error()
		log.Errorf("Reserves error: %s", err)
	} else {
		resp.Error = restful.SUCCESS
		resp.Result = reserves
		log.Infof("Reserves success")
	}

	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("Reserves: failed, err: %s", err)
	} else {
		log.Debug("Reserves: resp success")
	}
	return m
}

func (this *Service) IfReserves(param map[string]interface{}) map[string]interface{} {
	resp := &common.Response{}
	reserves, err := this.ifMgr.Reserves()
	if err != nil {
		resp.Error = restful.INTERNAL_ERROR
		resp.Desc = err.Error()
		log.Errorf("IfReserves error: %s", err)
	} else {
		resp.Error = restful.SUCCESS
		resp.Result = reserves
		log.Infof("IfReserves success")
	}

	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("IfReserves: failed, err: %s", err)
	} else {
		log.Debug("IfReserves: resp success")
	}
	return m
}

func (this *Service) IFPoolInfo(param map[string]interface{}) map[string]interface{} {
	req := &common.IFPoolInfoRequest{}
	resp := &common.Response{}
	err := utils.ParseParams(req, param)
	if err != nil {
		resp.Error = restful.INVALID_PARAMS
		resp.Desc = err.Error()
		log.Errorf("IFPoolInfo: decode params failed, err: %s", err)
	} else {
		IFPoolInfo, err := this.ifMgr.IFPoolInfo(req.Address)
		if err != nil {
			resp.Error = restful.INTERNAL_ERROR
			resp.Desc = err.Error()
			log.Errorf("IFPoolInfo error: %s", err)
		} else {
			resp.Error = restful.SUCCESS
			resp.Result = IFPoolInfo
			log.Infof("IFPoolInfo success")
		}
	}

	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("IFPoolInfo: failed, err: %s", err)
	} else {
		log.Debug("IFPoolInfo: resp success")
	}
	return m
}

func (this *Service) IFHistory(param map[string]interface{}) map[string]interface{} {
	req := &common.IFHistoryRequest{}
	resp := &common.Response{}
	err := utils.ParseParams(req, param)
	if err != nil {
		resp.Error = restful.INVALID_PARAMS
		resp.Desc = err.Error()
		log.Errorf("IFHistory: decode params failed, err: %s", err)
	} else {
		IFHistory, err := this.ifMgr.IFHistory(req.Address, req.Asset, req.Operation, req.StartTimestamp, req.EndTimestamp,
			req.PageNo, req.PageSize)
		if err != nil {
			resp.Error = restful.INTERNAL_ERROR
			resp.Desc = err.Error()
			log.Errorf("IFHistory error: %s", err)
		} else {
			resp.Error = restful.SUCCESS
			resp.Result = IFHistory
			log.Infof("IFHistory success")
		}
	}

	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("IFHistory: failed, err: %s", err)
	} else {
		log.Debug("IFHistory: resp success")
	}
	return m
}

func (this *Service) IFLiquidation(param map[string]interface{}) map[string]interface{} {
	req := &common.IFLiquidationRequest{}
	resp := &common.Response{}
	err := utils.ParseParams(req, param)
	if err != nil {
		resp.Error = restful.INVALID_PARAMS
		resp.Desc = err.Error()
		log.Errorf("IFLiquidation: decode params failed, err: %s", err)
	} else {
		debtAccounts, err := this.ifMgr.CheckIfDebt(req.Start, req.End)
		if err != nil {
			resp.Error = restful.INTERNAL_ERROR
			resp.Desc = err.Error()
			log.Errorf("IFLiquidation error: %s", err)
		} else {
			resp.Error = restful.SUCCESS
			resp.Result = debtAccounts
			log.Infof("IFLiquidation success")
		}
	}
	m, err := utils.RefactorResp(resp, resp.Error)
	if err != nil {
		log.Errorf("IFLiquidation: failed, err: %s", err)
	} else {
		log.Debug("IFLiquidation: resp success")
	}
	return m
}
