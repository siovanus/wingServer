package ifpool

import (
	"fmt"
	"math/big"

	"github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/utils"
	gov "github.com/wing-groups/wing-contract-tools/contracts/governance"
)

func (this *IFPoolManager) assetStoredPrice(asset string) (*big.Int, error) {
	price, err := this.store.LoadPrice(asset)
	if err != nil {
		return nil, fmt.Errorf("AssetStoredPrice, this.store.LoadPrice error: %s", err)
	}
	return utils.ToIntByPrecise(price.Price, this.cfg.TokenDecimal["oracle"]), nil
}

func (this *IFPoolManager) getPoolWeightInfo() (*gov.PoolWeightInfos, error) {
	contractAddress, err := common.AddressFromHexString(this.cfg.GovernanceAddress)
	if err != nil {
		return nil, fmt.Errorf("getPoolWeightInfo, ocommon.AddressFromHexString error: %s", err)
	}
	res, err := this.Sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress, "get_pool_weight_info",
		[]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("getPoolWeightInfo: %s", err)
	}
	result, err := res.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("getPoolWeightInfo: %s", err)
	}
	pools := new(gov.PoolWeightInfos)
	source := common.NewZeroCopySource(result)
	source.NextBool()
	err = pools.Deserialization(source)
	if err != nil {
		return nil, fmt.Errorf("getPoolWeightInfo: %s", err)
	}
	return pools, nil
}
