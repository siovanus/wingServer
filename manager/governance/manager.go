package governance

import (
	"fmt"
	sdk "github.com/ontio/ontology-go-sdk"
	ocommon "github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/config"
	"github.com/siovanus/wingServer/http/common"
	"time"
)

const (
	Total20         = 2000000000000000
	Total80         = 8000000000000000
	WingAmount20Per = 500000000000000
	YearSecond      = 31536000
	DaySecond       = 86400
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
	count20, err := this.get20WingCount()
	if err != nil {
		return nil, fmt.Errorf("GovBannerOverview, this.get20WingCount error: %s", err)
	}
	totalSupply, err := this.getWingTotalSupply()
	if err != nil {
		return nil, fmt.Errorf("GovBannerOverview, this.getWingTotalSupply error: %s", err)
	}
	return &common.GovBannerOverview{
		Remain20: Total20 - count20*WingAmount20Per,
		Remain80: Total80 - totalSupply - count20*WingAmount20Per,
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
		Daily:       DailyDistibute[index],
		Distributed: distributed,
	}, nil
}
