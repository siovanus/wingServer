package governance

import (
	"github.com/siovanus/wingServer/config"
	"github.com/siovanus/wingServer/http/common"
	"github.com/siovanus/wingServer/utils"
	"math/big"
	"time"
)

const (
	Total      = 2000000
	YearSecond = 31536000
	DaySecond  = 86400
)

var GenesisTime = uint64(time.Date(2020, time.September, 12, 0, 0, 0, 0, time.UTC).Unix())
var DailyDistibute = []uint64{6, 60, 30, 18, 6, 5, 4, 3, 2, 1, 1, 1, 1, 1}
var DistributeTime = []uint64{3 * DaySecond, 5 * DaySecond, 5 * DaySecond, 5 * DaySecond,
	YearSecond - 18*DaySecond, YearSecond, YearSecond, YearSecond,
	YearSecond, YearSecond, YearSecond, YearSecond, YearSecond, 4256064}

type GovernanceManager struct {
	cfg *config.Config
}

func NewGovernanceManager(cfg *config.Config) *GovernanceManager {
	manager := &GovernanceManager{
		cfg: cfg,
	}

	return manager
}

func (this *GovernanceManager) Wing() (*common.Wing, error) {
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

	total, _ := new(big.Float).SetString(utils.ToStringByPrecise(new(big.Int).SetUint64(distributed+Total*100), 2))
	t, _ := total.Float64()
	circulating, _ := new(big.Float).SetString(utils.ToStringByPrecise(new(big.Int).SetUint64(distributed), 2))
	c, _ := circulating.Float64()

	return &common.Wing{
		Total:       t,
		Circulating: c,
	}, nil
}
