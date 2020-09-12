package service

import (
	"fmt"
	"github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/log"
	"github.com/siovanus/wingServer/store"
)

func (this *Service) trackEvent(height uint32) (bool, []string, error) {
	accounts := []string{}
	events, err := this.sdk.GetSmartContractEventByBlock(height)
	if err != nil {
		return false, accounts, fmt.Errorf("TrackOracle, this.sdk.GetSmartContractEventByBlock error:%s", err)
	}
	flag := false
	for _, event := range events {
		for _, notify := range event.Notify {
			states, ok := notify.States.([]interface{})
			if !ok {
				continue
			}
			if !listContains(this.listeningAddressList, notify.ContractAddress) {
				continue
			}
			name, _ := states[0].(string)
			if name == "PutUnderlyingPrice" {
				flag = true
			}
			if len(states) > 1 {
				a, ok := states[1].(string)
				if !ok {
					continue
				}
				_, err = common.AddressFromBase58(a)
				if err != nil {
					continue
				} else {
					if listContains(this.cfg.SystemContract, a) {
						continue
					}
					if listContains(accounts, a) {
						continue
					}
					accounts = append(accounts, a)
				}
			}
		}
	}
	return flag, accounts, nil
}

func (this *Service) PriceFeed() error {
	for _, v := range ASSET {
		data, err := this.fpMgr.AssetPrice(v)
		if err != nil {
			log.Errorf("PriceFeed, this.fpMgr.AssetPrice error: %s", err)
			return err
		}
		price := &store.Price{
			Name:  v,
			Price: data,
		}
		err = this.store.SavePrice(price)
		if err != nil {
			log.Errorf("PriceFeed, this.store.SavePrice error: %s", err)
			return err
		}
	}
	return nil
}

func (this *Service) StoreFlashPoolOverview(account string) {
	userFlashPoolOverview, err := this.fpMgr.UserFlashPoolOverviewForStore(account)
	if err != nil {
		log.Errorf("StoreFlashPoolOverview, this.fpMgr.UserFlashPoolOverviewForStore error: %s", err)
	}
	err = this.store.SaveUserFlashPoolOverview(account, userFlashPoolOverview)
	if err != nil {
		log.Errorf("StoreFlashPoolOverview, this.store.SaveUserFlashPoolOverview error: %s", err)
	}
}

func (this *Service) StoreFlashPoolAllMarket() error {
	flashPoolAllMarket, err := this.fpMgr.FlashPoolAllMarketForStore()
	if err != nil {
		log.Errorf("StoreFlashPoolAllMarket, this.fpMgr.FlashPoolAllMarketForStore error: %s", err)
		return err
	}
	for _, v := range flashPoolAllMarket.FlashPoolAllMarket {
		err = this.store.SaveFlashMarket(v)
		if err != nil {
			log.Errorf("StoreFlashPoolOverview, this.store.SaveFlashMarket error: %s", err)
			return err
		}
	}
	return nil
}

func listContains(list []string, arg string) bool {
	for _, v := range list {
		if arg == v {
			return true
		}
	}
	return false
}