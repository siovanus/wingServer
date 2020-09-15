package governance

import (
	"fmt"
	"testing"
	"time"
)

func TestDistributed(t *testing.T) {
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
	distributed += (gap - epoch[index]) * DailyDistibute[index+1]
	fmt.Println(distributed/100)
}
