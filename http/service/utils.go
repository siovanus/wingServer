package service

import (
	"fmt"
	"github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/store"
)

func (this *Service) trackEvent(height uint32) (bool, string, error) {
	events, err := this.sdk.GetSmartContractEventByBlock(height)
	if err != nil {
		return false, "", fmt.Errorf("TrackOracle, this.sdk.GetSmartContractEventByBlock error:%s", err)
	}
	flag := false
	account := ""
	for _, event := range events {
		for _, notify := range event.Notify {
			states, ok := notify.States.([]interface{})
			if !ok {
				continue
			}
			listen := false
			for _, v := range this.listeningAddressList {
				if notify.ContractAddress == v.ToHexString() {
					listen = true
				}
			}
			if !listen {
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
					account = a
				}
			}

		}
	}
	return flag, account, nil
}

func (this *Service) PriceFeed() error {
	for _, v := range ASSET {
		data, err := this.fpMgr.AssetPrice(v)
		if err != nil {
			return fmt.Errorf("PriceFeed, this.fpMgr.AssetPrice error: %s", err)
		}
		price := &store.Price{
			Name:  v,
			Price: data,
		}
		err = this.store.SavePrice(price)
		if err != nil {
			return fmt.Errorf("PriceFeed, this.store.SavePrice error: %s", err)
		}
	}
	return nil
}

func (this *Service) StoreFlashPoolOverview(account string) error {
	userFlashPoolOverview, err := this.fpMgr.UserFlashPoolOverviewForStore(account)
	if err != nil {
		return fmt.Errorf("StoreFlashPoolOverview, this.fpMgr.UserFlashPoolOverviewForStore error: %s", err)
	}
	err = this.store.SaveUserFlashPoolOverview(account, userFlashPoolOverview)
	if err != nil {
		return fmt.Errorf("StoreFlashPoolOverview, this.store.SaveUserFlashPoolOverview error: %s", err)
	}
	return nil
}

func (this *Service) StoreFlashPoolAllMarket() error {
	flashPoolAllMarket, err := this.fpMgr.FlashPoolAllMarketForStore()
	if err != nil {
		return fmt.Errorf("StoreFlashPoolAllMarket, this.fpMgr.FlashPoolAllMarketForStore error: %s", err)
	}
	for _, v := range flashPoolAllMarket.FlashPoolAllMarket {
		err = this.store.SaveFlashMarket(v)
		if err != nil {
			return fmt.Errorf("StoreFlashPoolOverview, this.store.SaveFlashMarket error: %s", err)
		}
	}
	return nil
}
