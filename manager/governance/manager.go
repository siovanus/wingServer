package governance

import (
	"math/big"
	"time"

	ontology_go_sdk "github.com/ontio/ontology-go-sdk"
	ocommon "github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/config"
	"github.com/siovanus/wingServer/http/common"
	"github.com/siovanus/wingServer/utils"
)

const (
	Total             = 2000000000000000
	YearSecond        = 31536000
	DaySecond         = 86400
	ZeroAddress       = "AFmseVrdL9f9oyCzZefL9tG6UbvhPbdYzM"
	FoundationAddress = "AUKZ3KL1FRRhgcijH6DBdBtswUdtmqL8Wo"
)

var GenesisTime = uint64(time.Date(2020, time.September, 12, 0, 0, 0, 0, time.UTC).Unix())
var DailyDistibute = []uint64{60000000, 600000000, 300000000, 180000000, 60000000, 50000000, 40000000, 30000000,
	20000000, 10000000, 10000000, 10000000, 10000000, 10000000}
var DistributeTime = []uint64{3 * DaySecond, 5 * DaySecond, 5 * DaySecond, 5 * DaySecond,
	YearSecond - 18*DaySecond, YearSecond, YearSecond, YearSecond,
	YearSecond, YearSecond, YearSecond, YearSecond, YearSecond, 4256064}

type GovernanceManager struct {
	cfg *config.Config
	sdk *ontology_go_sdk.OntologySdk
}

func NewGovernanceManager(cfg *config.Config, sdk *ontology_go_sdk.OntologySdk) *GovernanceManager {
	manager := &GovernanceManager{
		cfg: cfg,
		sdk: sdk,
	}

	return manager
}

func (this *GovernanceManager) Wing() (*common.Wing, error) {
	wingAddress, err := ocommon.AddressFromHexString(this.cfg.WingAddress)
	if err != nil {
		return nil, err
	}
	address, err := ocommon.AddressFromBase58(ZeroAddress)
	if err != nil {
		return nil, err
	}
	result, err := this.sdk.NeoVM.PreExecInvokeNeoVMContract(wingAddress, []interface{}{"balanceOf", []interface{}{address}})
	if err != nil {
		return nil, err
	}
	burned, err := result.Result.ToInteger()
	if err != nil {
		return nil, err
	}

	fAddress, err := ocommon.AddressFromBase58(FoundationAddress)
	if err != nil {
		return nil, err
	}
	fResult, err := this.sdk.NeoVM.PreExecInvokeNeoVMContract(wingAddress, []interface{}{"balanceOf", []interface{}{fAddress}})
	if err != nil {
		return nil, err
	}
	fBalance, err := fResult.Result.ToInteger()
	if err != nil {
		return nil, err
	}

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

	circulating, _ := new(big.Float).SetString(utils.ToStringByPrecise(new(big.Int).Sub(new(big.Int).Sub(new(big.Int).SetUint64(
		distributed+Total), fBalance), burned), 9))
	c, _ := circulating.Float64()
	total, _ := new(big.Float).SetString(utils.ToStringByPrecise(new(big.Int).Add(new(big.Int).Sub(new(big.Int).Sub(new(big.Int).SetUint64(
		distributed+Total), fBalance), burned), fBalance), 9))
	t, _ := total.Float64()

	return &common.Wing{
		Total:       t,
		Circulating: c,
	}, nil
}
