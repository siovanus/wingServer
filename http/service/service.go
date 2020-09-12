package service

import (
	"github.com/ontio/ontology/common"
	"os"
	"time"

	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/siovanus/wingServer/config"
	"github.com/siovanus/wingServer/log"
	"github.com/siovanus/wingServer/store"
)

type Service struct {
	sdk                  *sdk.OntologySdk
	cfg                  *config.Config
	govMgr               GovernanceManager
	fpMgr                FlashPoolManager
	store                *store.Client
	trackHeight          uint32
	listeningAddressList []common.Address
}

var ASSET = []string{"ONT", "BTC", "ETH", "DAI"}

func NewService(sdk *sdk.OntologySdk, govMgr GovernanceManager, fpMgr FlashPoolManager, store *store.Client, cfg *config.Config) *Service {
	return &Service{sdk: sdk, cfg: cfg, govMgr: govMgr, fpMgr: fpMgr, store: store}
}

func (this *Service) AddListeningAddressList() {
	oracleAddress, err := common.AddressFromHexString(this.cfg.OracleAddress)
	if err != nil {
		log.Errorf("AddListeningAddressList, common.AddressFromHexString error: %s", err)
		os.Exit(1)
	}
	allMarkets, err := this.fpMgr.GetAllMarkets()
	if err != nil {
		log.Errorf("AddListeningAddressList, this.fpMgr.GetAllMarkets error: %s", err)
		os.Exit(1)
	}
	this.listeningAddressList = append(allMarkets, oracleAddress)
}

func (this *Service) Close() {
	err := this.store.Close()
	if err != nil {
		log.Error(err)
	}

	log.Info("All connections closed. Bye!")
}

func (this *Service) Snapshot() {
	flashPoolDetail, err := this.fpMgr.FlashPoolDetailForStore()
	if err != nil {
		log.Errorf("Snapshot, this.fpMgr.FlashPoolDetailForStore error: %s", err)
	}
	err = this.store.SaveFlashPoolDetail(flashPoolDetail)
	if err != nil {
		log.Errorf("Snapshot, this.store.SaveFlashPoolDetail error: %s", err)
	}

	// FlashPoolMarketForStore
	err = this.fpMgr.FlashPoolMarketStore()
	if err != nil {
		log.Errorf("Snapshot, this.fpMgr.FlashPoolMarketStore error: %s", err)
	}
	for {
		now := time.Now()
		// 计算下一个零点
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		log.Infof("snapshot start: %v", time.Now())
		// 以下为定时执行的操作
		// FlashPoolDetailForStore
		flashPoolDetail, err := this.fpMgr.FlashPoolDetailForStore()
		if err != nil {
			log.Errorf("Snapshot, this.fpMgr.FlashPoolDetailForStore error: %s", err)
		}
		err = this.store.SaveFlashPoolDetail(flashPoolDetail)
		if err != nil {
			log.Errorf("Snapshot, this.store.SaveFlashPoolDetail error: %s", err)
		}

		// FlashPoolMarketForStore
		err = this.fpMgr.FlashPoolMarketStore()
		if err != nil {
			log.Errorf("Snapshot, this.fpMgr.FlashPoolMarketStore error: %s", err)
		}
	}
}

func (this *Service) TrackEvent() {
	trackHeight, err := this.store.LoadTrackHeight()
	if err != nil {
		log.Infof("TrackEvent, this.store.LoadTrackHeight error: %s", err)
		currentHeight, err := this.sdk.GetCurrentBlockHeight()
		if err != nil {
			log.Errorf("TrackEvent, this.sdk.GetCurrentBlockHeight error:", err)
			os.Exit(1)
		}
		this.trackHeight = currentHeight
	}
	this.trackHeight = trackHeight
	for {
		currentHeight, err := this.sdk.GetCurrentBlockHeight()
		if err != nil {
			log.Errorf("TrackEvent, this.sideSdk.GetCurrentBlockHeight error:", err)
		}
		for i := this.trackHeight + 1; i <= currentHeight; i++ {
			ifOracle, account, err := this.trackEvent(i)
			if err != nil {
				log.Errorf("TrackEvent, this.TrackOracle error:", err)
				break
			}

			if ifOracle {
				err = this.PriceFeed()
				if err != nil {
					log.Errorf("TrackEvent, this.PriceFeed error:", err)
					break
				}
			}

			if account != "" {
				err = this.StoreFlashPoolOverview(account)
				if err != nil {
					log.Errorf("TrackEvent, this.StoreFlashPoolOverview error:", err)
					break
				}
				err = this.StoreFlashPoolAllMarket()
				if err != nil {
					log.Errorf("TrackEvent, this.StoreFlashPoolAllMarket error:", err)
					break
				}
			}

			this.trackHeight++
			err = this.store.SaveTrackHeight(this.trackHeight)
			if err != nil {
				log.Errorf("TrackEvent, this.store.SaveTrackHeight error:", err)
				break
			}
		}
		time.Sleep(time.Duration(this.cfg.TrackEventInterval) * time.Second)
	}
}
