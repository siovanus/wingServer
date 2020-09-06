package service

import (
	"github.com/siovanus/wingServer/log"
	"github.com/siovanus/wingServer/store"
	"time"
)

type Service struct {
	govMgr GovernanceManager
	fpMgr  FlashPoolManager
	store  *store.Client
}

func NewService(govMgr GovernanceManager, fpMgr FlashPoolManager, store *store.Client) *Service {
	return &Service{govMgr: govMgr, fpMgr: fpMgr, store: store}
}

func (this *Service) Close() {
	err := this.store.Close()
	if err != nil {
		log.Error(err)
	}

	log.Info("All connections closed. Bye!")
}

func (this *Service) Snapshot() {
	for {
		now := time.Now()
		// 计算下一个零点
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		log.Infof("snapshot start: %v", time.Now())
		//以下为定时执行的操作
		flashPoolDetail, err := this.fpMgr.FlashPoolDetailForStore()
		if err != nil {
			log.Errorf("Snapshot, this.fpMgr.FlashPoolDetailForStore error: %s", err)
		}
		err = this.store.SaveFlashPoolDetail(flashPoolDetail)
		if err != nil {
			log.Errorf("Snapshot, this.store.SaveFlashPoolDetail error: %s", err)
		}
	}
}
