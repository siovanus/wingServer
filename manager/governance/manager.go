package governance

import (
	"fmt"
	sdk "github.com/ontio/ontology-go-sdk"
	ocommon "github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/config"
	"github.com/siovanus/wingServer/http/common"
	"github.com/siovanus/wingServer/utils"
	"math/big"
	"time"
)

const (
	Total      = 10000000000000000
	Total80    = 8000000
	YearSecond = 31536000
	DaySecond  = 86400
)

var GenesisTime = uint64(time.Date(2020, time.September, 12, 0, 0, 0, 0, time.UTC).Unix())
var DailyDistibute = []uint64{6, 60, 12, 12, 12, 6, 5, 4, 3, 2, 1, 1, 1, 1, 1}
var DistributeTime = []uint64{3 * DaySecond, 3*DaySecond - 8*3600 + 5*60, 5*DaySecond - (3*DaySecond - 8*3600 + 5*60),
	5 * DaySecond, 5 * DaySecond, YearSecond - 18*DaySecond, YearSecond, YearSecond,
	YearSecond, YearSecond, YearSecond, YearSecond, YearSecond, YearSecond, 4256000}

type GovernanceManager struct {
	cfg             *config.Config
	contractAddress ocommon.Address
	wingAddress     string
	sdk             *sdk.OntologySdk
}

func NewGovernanceManager(contractAddress ocommon.Address, wingAddress string, sdk *sdk.OntologySdk, cfg *config.Config) *GovernanceManager {
	manager := &GovernanceManager{
		cfg:             cfg,
		contractAddress: contractAddress,
		wingAddress:     wingAddress,
		sdk:             sdk,
	}

	return manager
}

func (this *GovernanceManager) GovBannerOverview() (*common.GovBannerOverview, error) {
	gap := uint64(time.Now().Unix()) - GenesisTime
	length := len(DailyDistibute)
	epoch := []uint64{0}
	for i := 1; i < length+1; i++ {
		epoch = append(epoch, epoch[i-1]+DistributeTime[i-1])
	}
	if gap > epoch[length] {
		gap = epoch[length]
	}
	index := 0
	for i := 0; i < len(epoch); i++ {
		if gap >= epoch[i] {
			index = i
		}
	}
	var distributed uint64 = 0
	for j := 0; j < index; j++ {
		distributed += DailyDistibute[j] * DistributeTime[j]
	}
	distributed += (gap - epoch[index]) * DailyDistibute[index]

	balance, err := this.getBalanceOf("AUKZ3KL1FRRhgcijH6DBdBtswUdtmqL8Wo")
	if err != nil {
		return nil, fmt.Errorf("GovBannerOverview, this.getBalanceOf error: %s", err)
	}
	//totalSupply, err := this.getWingTotalSupply()
	//if err != nil {
	//	return nil, fmt.Errorf("GovBannerOverview, this.getWingTotalSupply error: %s", err)
	//}
	//return &common.GovBannerOverview{
	//	Remain20: utils.ToStringByPrecise(new(big.Int).SetUint64(balance), this.cfg.TokenDecimal["WING"]),
	//	Remain80: utils.ToStringByPrecise(new(big.Int).Sub(new(big.Int).SetUint64(Total), totalSupply),
	//		this.cfg.TokenDecimal["WING"]),
	//}, nil
	remain80 := Total80*100 - distributed
	return &common.GovBannerOverview{
		Remain20: utils.ToStringByPrecise(new(big.Int).SetUint64(balance), this.cfg.TokenDecimal["WING"]),
		Remain80: utils.ToStringByPrecise(new(big.Int).SetUint64(remain80), 2),
	}, nil
}

func (this *GovernanceManager) GovBanner() (*common.GovBanner, error) {
	gap := uint64(time.Now().Unix()) - GenesisTime
	length := len(DailyDistibute)
	epoch := []uint64{0}
	for i := 1; i < length+1; i++ {
		epoch = append(epoch, epoch[i-1]+DistributeTime[i-1])
	}
	if gap > epoch[length] {
		gap = epoch[length]
	}
	index := 0
	for i := 0; i < len(epoch); i++ {
		if gap >= epoch[i] {
			index = i
		}
	}
	var distributed uint64 = 0
	for j := 0; j < index; j++ {
		distributed += DailyDistibute[j] * DistributeTime[j]
	}
	distributed += (gap - epoch[index]) * DailyDistibute[index]

	return &common.GovBanner{
		Daily:       utils.ToStringByPrecise(new(big.Int).SetUint64(DailyDistibute[index]*DaySecond), 2),
		Distributed: utils.ToStringByPrecise(new(big.Int).SetUint64(distributed), 2),
	}, nil
}
