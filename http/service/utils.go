package service

import (
	gocommon "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/log"
	"github.com/siovanus/wingServer/store"
)

func (this *Service) trackSnapshotEvent(events []*gocommon.SmartContactEvent) (bool, []string, error) {
	accounts := []string{}
	flag := false
	for _, event := range events {
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
	}
	return flag, accounts, nil
}

func (this *Service) PriceFeed() error {
	for _, v := range this.fpMgr.AssetMap {
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

func (this *Service) StoreUserBalance(accountStr string) {
	err := this.fpMgr.UserBalanceForStore(accountStr)
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

func (this *Service) StoreIFInfo() {
	err := this.ifMgr.StoreIFInfo()
	if err != nil {
		log.Errorf("StoreIFInfo, this.fpMgr.StoreIFInfo error: %s", err)
	}
}

func (this *Service) StoreIFMarketInfo() {
	err := this.ifMgr.StoreIFMarketInfo()
	if err != nil {
		log.Errorf("StoreIFMarketInfo, this.ifMgr.StoreIFMarketInfo error: %s", err)
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

func (this *Service) trackIfOperationEvent(height uint32, events []*gocommon.SmartContactEvent) ([]*store.IfPoolHistory, error) {
	ifPoolHistories := make([]*store.IfPoolHistory, 0)
	for _, event := range events {
		for _, notify := range event.Notify {
			var ok bool
			states, ok := notify.States.([]interface{})
			if !ok {
				continue
			}
			hexContractAddress, err := common.AddressFromHexString(notify.ContractAddress)
			if err != nil {
				log.Errorf("trackIfOperationEvent, common.AddressFromHexString error: %s", err)
			}

			fToken, isFToken := this.ifMgr.FTokenMap[hexContractAddress]
			bToken, isBToken := this.ifMgr.BorrowMap[hexContractAddress]
			iToken, isIToken := this.ifMgr.ITokenMap[hexContractAddress]
			if isFToken {
				method, _ := states[0].(string)
				if method == "Mint" {
					txHash := event.TxHash
					operation := "supply"
					addr, _ := states[1].(string)
					amount, _ := states[2].(string)
					name, err := fToken.Name()
					if err != nil {
						log.Errorf("trackIfOperationEvent, fToken.Name error: %s", err)
					}
					history, err := this.constructHistory(addr, name, operation, amount, txHash, height)
					if err != nil {
						log.Errorf("trackIfOperationEvent, this.constructHistory error: %s", err)
					}
					ifPoolHistories = append(ifPoolHistories, history)
				} else if method == "Redeem" {
					txHash := event.TxHash
					operation := "withdraw"
					addr, _ := states[1].(string)
					amount, _ := states[2].(string)
					name, err := fToken.Name()
					if err != nil {
						log.Errorf("trackIfOperationEvent, fToken.Name error: %s", err)
					}
					history, err := this.constructHistory(addr, name, operation, amount, txHash, height)
					if err != nil {
						log.Errorf("trackIfOperationEvent, this.constructHistory error: %s", err)
					}
					ifPoolHistories = append(ifPoolHistories, history)
				}
			} else if isBToken {
				method, _ := states[0].(string)
				if method == "Borrow" {
					txHash := event.TxHash
					operation := "borrow"
					addr, _ := states[1].(string)
					amount, _ := states[2].(string)
					name, err := bToken.MarketName()
					if err != nil {
						log.Errorf("trackIfOperationEvent, bToken.MarketName error: %s", err)
					}
					history, err := this.constructHistory(addr, name, operation, amount, txHash, height)
					if err != nil {
						log.Errorf("trackIfOperationEvent, this.constructHistory error: %s", err)
					}
					ifPoolHistories = append(ifPoolHistories, history)
				} else if method == "RepayBorrow" {
					txHash := event.TxHash
					operation := "repay"
					addr, _ := states[1].(string)
					amount, _ := states[3].(string)
					name, err := bToken.MarketName()
					if err != nil {
						log.Errorf("trackIfOperationEvent, bToken.MarketName error: %s", err)
					}
					history, err := this.constructHistory(addr, name, operation, amount, txHash, height)
					if err != nil {
						log.Errorf("trackIfOperationEvent, this.constructHistory error: %s", err)
					}
					ifPoolHistories = append(ifPoolHistories, history)
				}
			} else if isIToken {
				method, _ := states[0].(string)
				if method == "Mint" {
					txHash := event.TxHash
					operation := "supply"
					addr, _ := states[1].(string)
					amount, _ := states[2].(string)
					name, err := iToken.Name()
					if err != nil {
						log.Errorf("trackIfOperationEvent, iToken.Name error: %s", err)
					}
					history, err := this.constructHistory(addr, name, operation, amount, txHash, height)
					if err != nil {
						log.Errorf("trackIfOperationEvent, this.constructHistory error: %s", err)
					}
					ifPoolHistories = append(ifPoolHistories, history)
				} else if method == "Redeem" {
					txHash := event.TxHash
					operation := "withdraw"
					addr, _ := states[1].(string)
					amount, _ := states[2].(string)
					name, err := iToken.Name()
					if err != nil {
						log.Errorf("trackIfOperationEvent, iToken.Name error: %s", err)
					}
					history, err := this.constructHistory(addr, name, operation, amount, txHash, height)
					if err != nil {
						log.Errorf("trackIfOperationEvent, this.constructHistory error: %s", err)
					}
					ifPoolHistories = append(ifPoolHistories, history)
				}
			}
		}
	}
	return ifPoolHistories, nil
}

func (this *Service) constructHistory(addr string, name string, operation string, amount string, txHash string, height uint32) (*store.IfPoolHistory, error) {
	byHeight, err := this.sdk.GetBlockByHeight(height)
	if err != nil {
		log.Errorf("constructHistory, this.sdk.GetBlockByHeight error: %s", err)
	}
	timestamp := byHeight.Header.Timestamp

	history := new(store.IfPoolHistory)
	history.Address = addr
	history.Token = name
	history.Operation = operation
	history.Amount = amount
	history.TxHash = txHash
	history.Timestamp = uint64(timestamp)
	return history, nil
}

func (this *Service) StoreUserIfOperation(history *store.IfPoolHistory) {
	err := this.ifMgr.StoreUserIfOperation(history)
	if err != nil {
		log.Errorf("StoreUserIfOperation, this.fpMgr.StoreUserIfOperation error: %s", err)
	}
}
