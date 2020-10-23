package service

import (
	"fmt"
	"github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/log"
	"github.com/siovanus/wingServer/store"
)

func (this *Service) trackSnapshotEvent(height uint32) (bool, []string, error) {
	accounts := []string{}
	events, err := this.sdk.GetSmartContractEventByBlock(height)
	if err != nil {
		return false, accounts, fmt.Errorf("TrackOracle, this.sdk.GetSmartContractEventByBlock error:%s", err)
	}
	flag := false
	for _, event := range events {
		var borrowAccount, borrowAmount, borrowIndex string
		for _, notify := range event.Notify {
			var ok bool
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
				continue
			}

			if name == "Borrow" {
				borrowAccount, ok = states[1].(string)
				if !ok {
					log.Errorf("trackSnapshotEvent, Borrow event, parse account error")
					continue
				}
				borrowAmount, ok = states[2].(string)
				if !ok {
					log.Errorf("trackSnapshotEvent, Borrow event, parse borrowAmount error")
					continue
				}
				continue
			}
			if name == "AccrueInterest" {
				borrowIndex, ok = states[3].(string)
				if !ok {
					log.Errorf("trackSnapshotEvent, Borrow event, parse borrowIndex error")
					continue
				}
				continue
			}

			if len(states) > 1 {
				a, ok := states[1].(string)
				if ok {
					address, err := common.AddressFromBase58(a)
					if err == nil {
						if !listContains(this.listeningAddressList, address.ToHexString()) {
							if !listContains(accounts, a) {
								accounts = append(accounts, a)
							}
						}
					}
				}
			}
			if len(states) > 2 {
				a, ok := states[2].(string)
				if ok {
					address, err := common.AddressFromBase58(a)
					if err == nil {
						if !listContains(this.listeningAddressList, address.ToHexString()) {
							if !listContains(accounts, a) {
								accounts = append(accounts, a)
							}
						}
					}
				}
			}
		}
		if borrowAccount != "" {
			this.StoreUserBalance(borrowAccount, borrowAmount, borrowIndex)
		}
	}
	return flag, accounts, nil
}

func (this *Service) PriceFeed() error {
	for _, v := range this.assetList {
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

func (this *Service) StoreUserBalance(accountStr, borrowAmount, borrowIndex string) {
	err := this.fpMgr.UserBalanceForStore(accountStr, borrowAmount, borrowIndex)
	if err != nil {
		log.Errorf("StoreUserBalance, this.fpMgr.UserFlashPoolOverviewForStore error: %s", err)
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

func (this *Service) StoreWingApy() {
	err := this.fpMgr.WingApyForStore()
	if err != nil {
		log.Errorf("StoreUserBalance, this.fpMgr.WingApyForStore error: %s", err)
	}
}

func listContains(list []string, arg string) bool {
	for _, v := range list {
		if arg == v {
			return true
		}
	}
	return false
}
