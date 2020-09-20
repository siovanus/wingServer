package service

import (
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
	listeningAddressList []string
	assetList            []string
}

func NewService(sdk *sdk.OntologySdk, govMgr GovernanceManager, fpMgr FlashPoolManager, store *store.Client, cfg *config.Config) *Service {
	return &Service{sdk: sdk, cfg: cfg, govMgr: govMgr, fpMgr: fpMgr, store: store}
}

func (this *Service) AddListeningAddressList() {
	allMarkets, err := this.fpMgr.GetAllMarkets()
	if err != nil {
		log.Errorf("AddListeningAddressList, this.fpMgr.GetAllMarkets error: %s", err)
		os.Exit(1)
	}
	for _, v := range allMarkets {
		this.assetList = append(this.assetList, this.cfg.OracleMap[v.ToHexString()])
		this.listeningAddressList = append(this.listeningAddressList, v.ToHexString())
		addr, err := this.fpMgr.GetInsuranceAddress(v)
		if err != nil {
			log.Errorf("AddListeningAddressList, this.fpMgr.GetInsuranceAddress error: %s", err)
			os.Exit(1)
		}
		this.listeningAddressList = append(this.listeningAddressList, addr.ToHexString())
	}
	this.listeningAddressList = append(this.listeningAddressList, this.cfg.WingAddress)
	this.listeningAddressList = append(this.listeningAddressList, this.cfg.GovernanceAddress)
	this.listeningAddressList = append(this.listeningAddressList, this.cfg.OracleAddress)
	this.listeningAddressList = append(this.listeningAddressList, this.cfg.FlashPoolAddress)
}

func (this *Service) Close() {
	err := this.store.Close()
	if err != nil {
		log.Error(err)
	}

	log.Info("All connections closed. Bye!")
}

func (this *Service) SnapshotDaily() {
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
	//init
	err := this.PriceFeed()
	if err != nil {
		log.Errorf("TrackEvent, this.PriceFeed error:", err)
		os.Exit(1)
	}
	err = this.StoreFlashPoolAllMarket()
	if err != nil {
		log.Errorf("TrackEvent, this.StoreFlashPoolAllMarket error:", err)
		os.Exit(1)
	}

	trackHeight, err := this.store.LoadTrackHeight()
	if err != nil {
		log.Infof("TrackEvent, this.store.LoadTrackHeight error: %s", err)
		currentHeight, err := this.sdk.GetCurrentBlockHeight()
		if err != nil {
			log.Errorf("TrackEvent, this.sdk.GetCurrentBlockHeight error:", err)
			os.Exit(1)
		}
		this.trackHeight = currentHeight
	} else {
		this.trackHeight = trackHeight
	}
	for {
		currentHeight, err := this.sdk.GetCurrentBlockHeight()
		if err != nil {
			log.Errorf("TrackEvent, this.sideSdk.GetCurrentBlockHeight error:", err)
		}
		for i := this.trackHeight + 1; i <= currentHeight; i++ {
			log.Infof("TrackEvent, parse block: %d", i)
			ifOracle, accounts, err := this.trackSnapshotEvent(i)
			if err != nil {
				log.Errorf("TrackEvent, this.TrackOracle error:", err)
				break
			}

			if ifOracle {
				log.Infof("TrackEvent, this.PriceFeed")
				go this.PriceFeed()
			}

			if len(accounts) != 0 {
				for _, v := range accounts {
					log.Infof("TrackEvent, account: %s", v)
					go this.StoreUserBalance(v)
				}
			}

			this.trackHeight++
			err = this.store.SaveTrackHeight(this.trackHeight)
			if err != nil {
				log.Errorf("TrackEvent, this.store.SaveTrackHeight error:", err)
				break
			}
		}
		time.Sleep(time.Second * time.Duration(this.cfg.ScanInterval))
	}
}

func (this *Service) SnapshotMinute() {
	for {
		err := this.StoreFlashPoolAllMarket()
		if err != nil {
			log.Errorf("SnapshotMinute, this.StoreFlashPoolAllMarket error:", err)
		}
		time.Sleep(time.Second * time.Duration(this.cfg.SnapshotInterval))
	}
}
