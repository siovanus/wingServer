package governance

import (
	"fmt"
	"github.com/siovanus/wingServer/utils"
	"math/big"
	"testing"
	"time"
)

func TestDistributed(t *testing.T) {
	//now := uint64(time.Date(2021, time.September, 17, 16, 5, 0, 0, time.UTC).Unix())
	//gap := now - GenesisTime
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
	daily := utils.ToStringByPrecise(new(big.Int).SetUint64(DailyDistibute[index]*DaySecond), 2)
	fmt.Println(utils.ToStringByPrecise(new(big.Int).SetUint64(distributed), 2))
	fmt.Println(daily)
}
