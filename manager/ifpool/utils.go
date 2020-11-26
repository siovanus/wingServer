package ifpool

import (
	"fmt"
	hcommon "github.com/siovanus/wingServer/http/common"
	"github.com/siovanus/wingServer/log"
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

func (this *IFPoolManager) getDynamicPercent() (*big.Int, error) {
	method := "get_dynamic_percent"
	params := []interface{}{}
	contractAddress, err := common.AddressFromHexString(this.cfg.GovernanceAddress)
	if err != nil {
		fmt.Errorf("getUtilities, common.AddressFromHexString: %s", err)
	}
	res, err := this.Sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress, method, params)
	if err != nil {
		fmt.Errorf("getUtilities, PreExecInvokeWasmVMContract: %s", err)
	}
	bs, err := res.Result.ToByteArray()
	if err != nil {
		fmt.Errorf("getUtilities, ToByteArray: %s", err)
	}
	source := common.NewZeroCopySource(bs)
	number, eof := source.NextI128()
	if eof {
		fmt.Errorf("getUtilities, source.NextByte: %v", err)
	}
	return number.ToBigInt(), nil
}

func (this *IFPoolManager) getUtilities() (*hcommon.IfMarketUtility, error) {
	method := "marketUtilities"
	params := []interface{}{}
	contractAddress, err := common.AddressFromHexString(this.cfg.IFPoolAddress)
	if err != nil {
		fmt.Errorf("getUtilities, common.AddressFromHexString: %s", err)
	}
	res, err := this.Sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress, method, params)
	if err != nil {
		fmt.Errorf("getUtilities, PreExecInvokeWasmVMContract: %s", err)
	}
	bs, err := res.Result.ToByteArray()
	if err != nil {
		fmt.Errorf("getUtilities, ToByteArray: %s", err)
	}
	source := common.NewZeroCopySource(bs)
	number, eof := source.NextByte()
	if eof {
		fmt.Errorf("getUtilities, source.NextByte: %v", err)
	}
	size := int(number)
	total := new(big.Int)
	utilityMap := make(map[string]*big.Int)
	for i := 0; i < size; i++ {
		name, _, _, err := source.NextString()
		if err {
			fmt.Errorf("getUtilities, source.NextAddress: %s", err)
		}
		data, err := source.NextBytes(32)
		if err {
			fmt.Errorf("getUtilities, source.NextBytes: %s", err)
		}
		utility := common.BigIntFromNeoBytes(data)
		log.Infof("____________________________market:%s utility:%d", name, utility)
		utilityMap[name] = utility
		total = new(big.Int).Add(total, utility)
	}
	return &hcommon.IfMarketUtility{
		UtilityMap: utilityMap,
		Total:      total,
	}, nil
}

func (this *IFPoolManager) getPoolWeight() (*hcommon.PoolWeight, error) {
	method := "get_pool_weight_info"
	params := []interface{}{}
	contractAddress, err := common.AddressFromHexString(this.cfg.GovernanceAddress)
	if err != nil {
		fmt.Errorf("getUtilities, common.AddressFromHexString: %s", err)
	}
	res, err := this.Sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress, method, params)
	if err != nil {
		fmt.Errorf("getUtilities, PreExecInvokeWasmVMContract: %s", err)
	}
	bs, err := res.Result.ToByteArray()
	if err != nil {
		fmt.Errorf("getUtilities, ToByteArray: %s", err)
	}
	source := common.NewZeroCopySource(bs)
	number, eof := source.NextByte()
	if eof {
		fmt.Errorf("getUtilities, source.NextByte: %v", err)
	}
	size := int(number)
	totalStatic := big.NewInt(0)
	totalDynamic := big.NewInt(0)
	poolStaticMap := make(map[common.Address]*big.Int)
	poolDynamicMap := make(map[common.Address]*big.Int)
	for i := 0; i < size; i++ {
		address, err := source.NextAddress()
		if err {
			fmt.Errorf("getUtilities, source.NextAddress: %s", err)
		}
		staticData, err := source.NextI128()
		if err {
			fmt.Errorf("getUtilities, source.NextBytes: %s", err)
		}
		staticWeight := staticData.ToBigInt()
		poolStaticMap[address] = staticWeight
		totalStatic = new(big.Int).Add(totalStatic, staticWeight)

		dynamicData, err := source.NextI128()
		if err {
			fmt.Errorf("getUtilities, source.NextBytes: %s", err)
		}
		dynamicWeight := dynamicData.ToBigInt()
		poolDynamicMap[address] = dynamicWeight
		totalDynamic = new(big.Int).Add(totalDynamic, dynamicWeight)
	}
	return &hcommon.PoolWeight{
		PoolStaticMap:  poolStaticMap,
		PoolDynamicMap: poolDynamicMap,
		TotalStatic:    totalStatic,
		TotalDynamic:   totalDynamic,
	}, nil
}
