package service

import (
	"github.com/siovanus/wingServer/log"
)

func (this *Service) CirculatingSupply(param map[string]interface{}) interface{} {
	var r float64
	wing, err := this.govMgr.Wing()
	if err != nil {
		log.Errorf("CirculatingSupply error: %s", err)
	} else {
		r = wing.Circulating
		log.Infof("CirculatingSupply success")
	}

	return r
}

func (this *Service) TotalSupply(param map[string]interface{}) interface{} {
	var r float64
	wing, err := this.govMgr.Wing()
	if err != nil {
		log.Errorf("TotalSupply error: %s", err)
	} else {
		r = wing.Total
		log.Infof("TotalSupply success")
	}

	return r
}
