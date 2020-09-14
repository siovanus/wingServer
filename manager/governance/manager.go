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
	YearSecond = 31536000
	DaySecond  = 86400
)

var GenesisTime = uint64(time.Date(2020, time.September, 8, 0, 0, 0, 0, time.UTC).Unix())
var DailyDistibute = []uint64{5184, 51840, 25920, 15552, 5184, 4320, 3456, 2592, 1728, 864, 864, 864, 864, 864}
var DistributeTime = []uint64{3 * DaySecond, 5 * DaySecond, 5 * DaySecond, 5 * DaySecond, YearSecond - 18*DaySecond, YearSecond, YearSecond,
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
	balance, err := this.getBalanceOf("AUKZ3KL1FRRhgcijH6DBdBtswUdtmqL8Wo")
	if err != nil {
		return nil, fmt.Errorf("GovBannerOverview, this.getBalanceOf error: %s", err)
	}
	totalSupply, err := this.getWingTotalSupply()
	if err != nil {
		return nil, fmt.Errorf("GovBannerOverview, this.getWingTotalSupply error: %s", err)
	}
	return &common.GovBannerOverview{
		Remain20: utils.ToStringByPrecise(new(big.Int).SetUint64(balance), this.cfg.TokenDecimal["WING"]),
		Remain80: utils.ToStringByPrecise(new(big.Int).Sub(new(big.Int).SetUint64(Total), totalSupply),
			this.cfg.TokenDecimal["WING"]),
	}, nil
}

func (this *GovernanceManager) GovBanner() (*common.GovBanner, error) {
	distributed := uint64(time.Now().Unix()) - GenesisTime
	length := len(DailyDistibute)
	epoch := []uint64{0}
	for i := 1; i < length+1; i++ {
		epoch = append(epoch, epoch[i-1]+DistributeTime[i-1])
	}
	if distributed > epoch[length] {
		distributed = epoch[length]
	}
	index := 0
	for i := 0; i < len(epoch); i++ {
		if distributed >= epoch[i] {
			index = i
		}
	}

	return &common.GovBanner{
		Daily:       new(big.Int).SetUint64(DailyDistibute[index]).String(),
		Distributed: new(big.Int).SetUint64(distributed).String(),
	}, nil
}
