package service

import (
	"fmt"
	"github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/store"
)

func (this *Service) TrackOracle(height uint32) (bool, error) {
	events, err := this.sdk.GetSmartContractEventByBlock(height)
	if err != nil {
		return false, fmt.Errorf("TrackOracle, this.sdk.GetSmartContractEventByBlock error:%s", err)
	}
	for _, event := range events {
		for _, notify := range event.Notify {
			states, ok := notify.States.([]interface{})
			if !ok {
				continue
			}
			if notify.ContractAddress != this.cfg.OracleAddress {
				continue
			}
			name, _ := states[0].(string)
			if name == "PutUnderlyingPrice" {
				return true, nil
			}
		}
	}
	return false, nil
}

func (this *Service) TrackFlash(height uint32) (string, error) {
	events, err := this.sdk.GetSmartContractEventByBlock(height)
	if err != nil {
		return "", fmt.Errorf("TrackOracle, this.sdk.GetSmartContractEventByBlock error:%s", err)
	}
	for _, event := range events {
		for _, notify := range event.Notify {
			states, ok := notify.States.([]interface{})
			if !ok {
				continue
			}
			if notify.ContractAddress != this.cfg.OracleAddress {
				continue
			}
			account := states[1].(string)
			_, err = common.AddressFromBase58(account)
			if err != nil {
				continue
			} else {
				return account, nil
			}
		}
	}
	return "", nil
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
	userFlashPoolOverview, err := this.fpMgr.UserFlashPoolOverview(account)
	if err != nil {
		return fmt.Errorf("StoreFlashPoolOverview, this.fpMgr.UserFlashPoolOverview error: %s", err)
	}
	err = this.store.SaveUserFlashPoolOverview(account, userFlashPoolOverview)
	if err != nil {
		return fmt.Errorf("StoreFlashPoolOverview, this.store.SaveUserFlashPoolOverview error: %s", err)
	}
	return nil
}
