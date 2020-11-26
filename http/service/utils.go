package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	gocommon "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/log"
	"github.com/siovanus/wingServer/store"
	"github.com/siovanus/wingServer/utils"
	"math/big"
	"strings"
	"time"
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
		log.Errorf("StoreWingApy, this.fpMgr.WingApyForStore error: %s", err)
	}
	err = this.ifMgr.WingApyForStore()
	if err != nil {
		log.Errorf("StoreWingApy, this.ifMgr.WingApyForStore error: %s", err)
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
		var parse = true
		for _, notify := range event.Notify {
			var ok bool
			states, ok := notify.States.([]interface{})
			if !ok {
				continue
			}
			if states[0] == "Failure" {
				parse = false
				break
			}
		}
		if !parse {
			continue
		}
		//var collateraler string
		var collateralAmount string
		var collateralPName string
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
					index := strings.LastIndex(name, " ") + 1
					name = name[index:len(name)]
					pName := this.cfg.IFMap[name]
					bigAmount, b := new(big.Int).SetString(amount, 10)
					if !b {
						log.Errorf("trackIfOperationEvent, new(big.Int).SetString error")
					}
					amount = utils.ToStringByPrecise(bigAmount, this.cfg.TokenDecimal[pName])
					history, err := this.constructHistory(addr, pName, operation, amount, txHash, height, "", "", "")
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
					index := strings.LastIndex(name, " ") + 1
					name = name[index:len(name)]
					pName := this.cfg.IFMap[name]
					bigAmount, b := new(big.Int).SetString(amount, 10)
					if !b {
						log.Errorf("trackIfOperationEvent, new(big.Int).SetString error")
					}
					amount = utils.ToStringByPrecise(bigAmount, this.cfg.TokenDecimal[pName])
					history, err := this.constructHistory(addr, pName, operation, amount, txHash, height, "", "", "")
					if err != nil {
						log.Errorf("trackIfOperationEvent, this.constructHistory error: %s", err)
					}
					ifPoolHistories = append(ifPoolHistories, history)
				}
			} else if isBToken {
				method, _ := states[0].(string)
				if method == "IncreaseCollateral" {
					//collateraler, _ = states[1].(string)
					collateralAmount, _ = states[2].(string)
					collateralName, err := bToken.MarketName()
					if err != nil {
						log.Errorf("trackIfOperationEvent, bToken.MarketName error: %s", err)
					}
					collateralIndex := strings.LastIndex(collateralName, " ") + 1
					collateralName = collateralName[collateralIndex:len(collateralName)]
					collateralPName = this.cfg.IFMap[collateralName]
				} else if method == "Borrow" {
					txHash := event.TxHash
					operation := "borrow"
					addr, _ := states[1].(string)
					amount, _ := states[2].(string)
					name, err := bToken.MarketName()
					if err != nil {
						log.Errorf("trackIfOperationEvent, bToken.MarketName error: %s", err)
					}
					index := strings.LastIndex(name, " ") + 1
					name = name[index:len(name)]
					pName := this.cfg.IFMap[name]
					log.Infof("pName:%s", pName)
					log.Infof("first amount:%s", amount)
					bigAmount, b := new(big.Int).SetString(amount, 10)
					log.Infof("bigAmount:%d", bigAmount)
					if !b {
						log.Errorf("trackIfOperationEvent, new(big.Int).SetString error")
					}
					log.Infof("this.cfg.TokenDecimal[pName]:%d", this.cfg.TokenDecimal[pName])
					amount = utils.ToStringByPrecise(bigAmount, this.cfg.TokenDecimal[pName])
					log.Infof("final amount:%s", amount)
					history, err := this.constructHistory(addr, pName, operation, amount, txHash, height, "", collateralPName, collateralAmount)
					if err != nil {
						log.Errorf("trackIfOperationEvent, this.constructHistory error: %s", err)
					}
					ifPoolHistories = append(ifPoolHistories, history)
				} else if method == "RepayBorrow" {
					var remark string
					txHash := event.TxHash
					operation := "repay"
					addr, _ := states[1].(string)
					borrower, _ := states[2].(string)
					amount, _ := states[3].(string)
					repayAll, _ := states[4].(string)

					name, err := bToken.MarketName()
					if err != nil {
						log.Errorf("trackIfOperationEvent, bToken.MarketName error: %s", err)
					}
					index := strings.LastIndex(name, " ") + 1
					name = name[index:len(name)]
					pName := this.cfg.IFMap[name]
					log.Infof("pName:%s", pName)
					log.Infof("first amount:%s", amount)
					bigAmount, b := new(big.Int).SetString(amount, 10)
					log.Infof("bigAmount:%d", bigAmount)
					if !b {
						log.Errorf("trackIfOperationEvent, new(big.Int).SetString error")
					}
					log.Infof("this.cfg.TokenDecimal[pName]:%d", this.cfg.TokenDecimal[pName])
					amount = utils.ToStringByPrecise(bigAmount, this.cfg.TokenDecimal[pName])
					log.Infof("final amount:%s", amount)

					if addr == borrower {
						// 0-没还清，1-还清
						remark = repayAll
					} else {
						// 2-帮他人还款
						remark = "2"
						if repayAll == "1" {
							// 帮他人还清，记录一笔borrower还清的记录
							operation = "repayByOther"
							additionalHistory, err := this.constructHistory(borrower, pName, operation, amount, txHash, height, "1", "", "")
							if err != nil {
								log.Errorf("trackIfOperationEvent, this.constructHistory error: %s", err)
							}
							ifPoolHistories = append(ifPoolHistories, additionalHistory)
						}
					}

					history, err := this.constructHistory(addr, pName, operation, amount, txHash, height, remark, "", "")
					if err != nil {
						log.Errorf("trackIfOperationEvent, this.constructHistory error: %s", err)
					}
					ifPoolHistories = append(ifPoolHistories, history)
				}
			} else if isIToken {
				method, _ := states[0].(string)
				if method == "Mint" {
					txHash := event.TxHash
					operation := "insure"
					addr, _ := states[1].(string)
					amount, _ := states[2].(string)
					name, err := iToken.Name()
					if err != nil {
						log.Errorf("trackIfOperationEvent, iToken.Name error: %s", err)
					}
					index := strings.LastIndex(name, " ") + 1
					name = name[index:len(name)]
					pName := this.cfg.IFMap[name]
					bigAmount, b := new(big.Int).SetString(amount, 10)
					if !b {
						log.Errorf("trackIfOperationEvent, new(big.Int).SetString error")
					}
					amount = utils.ToStringByPrecise(bigAmount, this.cfg.TokenDecimal[pName])
					history, err := this.constructHistory(addr, pName, operation, amount, txHash, height, "", "", "")
					if err != nil {
						log.Errorf("trackIfOperationEvent, this.constructHistory error: %s", err)
					}
					ifPoolHistories = append(ifPoolHistories, history)
				} else if method == "Redeem" {
					txHash := event.TxHash
					operation := "surrender"
					addr, _ := states[1].(string)
					amount, _ := states[2].(string)
					name, err := iToken.Name()
					if err != nil {
						log.Errorf("trackIfOperationEvent, iToken.Name error: %s", err)
					}
					index := strings.LastIndex(name, " ") + 1
					name = name[index:len(name)]
					pName := this.cfg.IFMap[name]
					bigAmount, b := new(big.Int).SetString(amount, 10)
					if !b {
						log.Errorf("trackIfOperationEvent, new(big.Int).SetString error")
					}
					amount = utils.ToStringByPrecise(bigAmount, this.cfg.TokenDecimal[pName])
					history, err := this.constructHistory(addr, pName, operation, amount, txHash, height, "", "", "")
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

func (this *Service) constructHistory(addr string, name string, operation string, amount string, txHash string, height uint32, remark string, collateralToken string, collateralAmount string) (*store.IfPoolHistory, error) {
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
	history.Remark = remark
	history.CollateralToken = collateralToken
	history.CollateralAmount = collateralAmount
	return history, nil
}

func (this *Service) StoreUserIfOperation(history *store.IfPoolHistory) {
	err := this.ifMgr.StoreUserIfOperation(history)
	if err != nil {
		log.Errorf("StoreUserIfOperation, this.ifMgr.StoreUserIfOperation error: %s", err)
	}
}

func (this *Service) checkIfDebt() {
	now := time.Now().Unix()
	log.Infof("+++++++++++++++++++++now:%d", now)
	eightDay := 8 * this.cfg.OneDaySecond
	log.Infof("+++++++++++++++++++++eightDay:%d", eightDay)
	nineDay := 9 * this.cfg.OneDaySecond
	log.Infof("+++++++++++++++++++++nineDay:%d", nineDay)
	end := now - eightDay
	start := now - nineDay
	debtAccounts, err := this.ifMgr.CheckIfDebt(start, end)
	if err != nil {
		log.Errorf("checkIfDebt, this.ifMgr.checkIfDebt: %s", err)
	}
	reqUrl := fmt.Sprintf(this.cfg.WingBackendUrl, "if-pool/blacklist")
	log.Infof("+++++++++++++++++++++reqUrl:%s", reqUrl)
	log.Infof("+++++++++++++++++++++debtAccounts:%d", len(debtAccounts))

	reqData, err := json.Marshal(debtAccounts)
	if err != nil {
		fmt.Errorf("checkIfDebt, json.Marshal error:%s", err)
	}
	resp, err := this.httpClient.Post(reqUrl, "application/json", bytes.NewReader(reqData))
	if err != nil {
		fmt.Errorf("checkIfDebt, send http post request error:%s", err)
	}
	defer resp.Body.Close()
}
