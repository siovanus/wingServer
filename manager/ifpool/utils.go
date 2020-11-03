package ifpool

import (
	"fmt"
	"math"
	"math/big"

	"github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/utils"
)

func (this *IFPoolManager) assetStoredPrice(asset string) (*big.Int, error) {
	if asset == "USDT" {
		return new(big.Int).SetUint64(uint64(math.Pow10(int(this.cfg.TokenDecimal["oracle"])))), nil
	}
	price, err := this.store.LoadPrice(asset)
	if err != nil {
		return nil, fmt.Errorf("AssetStoredPrice, this.store.LoadPrice error: %s", err)
	}
	return utils.ToIntByPrecise(price.Price, this.cfg.TokenDecimal["oracle"]), nil
}

type Markets struct {
	SupplyPool        common.Address
	BorrowPool        common.Address
	InsurancePool     common.Address
	Underlying        common.Address
	UnderlyingDecimal uint8
	WingWeight        uint8
}

func (this *Markets) Deserialization(source *common.ZeroCopySource) error {
	supplyPool, eof := source.NextAddress()
	if eof {
		return fmt.Errorf("read supplyPool eof")
	}
	borrowPool, eof := source.NextAddress()
	if eof {
		return fmt.Errorf("read borrowPool eof")
	}
	insurancePool, eof := source.NextAddress()
	if eof {
		return fmt.Errorf("read insurancePool eof")
	}
	underlying, eof := source.NextAddress()
	if eof {
		return fmt.Errorf("read underlying eof")
	}
	underlyingDecimal, eof := source.NextUint8()
	if eof {
		return fmt.Errorf("read underlyingDecimal eof")
	}
	wingWeight, eof := source.NextUint8()
	if eof {
		return fmt.Errorf("read wingWeight eof")
	}
	this.SupplyPool = supplyPool
	this.BorrowPool = borrowPool
	this.InsurancePool = insurancePool
	this.Underlying = underlying
	this.UnderlyingDecimal = underlyingDecimal
	this.WingWeight = wingWeight
	return nil
}

func (this *IFPoolManager) GetAllMarkets() ([]*Markets, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(this.contractAddress,
		"allMarkets", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("getAllMarkets, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("getAllMarkets, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	allMarkets := make([]*Markets, 0)
	l, _, irregular, eof := source.NextVarUint()
	if irregular || eof {
		return nil, fmt.Errorf("getAllMarkets, source.NextVarUint error")
	}
	for i := 0; uint64(i) < l; i++ {
		marketName, _, irregular, eof := source.NextString()
		if irregular || eof {
			return nil, fmt.Errorf("getAllMarkets, source.NextString error")
		}
		preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(this.contractAddress,
			"marketInfo", []interface{}{marketName})
		if err != nil {
			return nil, fmt.Errorf("getAllMarkets, marketInfo, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
		}
		r, err := preExecResult.Result.ToByteArray()
		if err != nil {
			return nil, fmt.Errorf("getAllMarkets, preExecResult.Result.ToByteArray error: %s", err)
		}
		market := new(Markets)
		source := common.NewZeroCopySource(r)
		err = market.Deserialization(source)
		if err != nil {
			return nil, fmt.Errorf("getAllMarkets, market.Deserialization error: %s", err)
		}
		allMarkets = append(allMarkets, market)
	}
	return allMarkets, nil
}

func (this *IFPoolManager) getTotalSupply(contractAddress common.Address) (*big.Int, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"totalSupply", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("getTotalSupply, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("getTotalSupply, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	amount, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("getTotalSupply, source.NextI128 error")
	}
	return amount.ToBigInt(), nil
}

func (this *IFPoolManager) getExchangeRate(contractAddress common.Address) (*big.Int, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"exchangeRate", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("getExchangeRate, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("getExchangeRate, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	amount, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("getExchangeRate, source.NextI128 error")
	}
	return amount.ToBigInt(), nil
}

func (this *IFPoolManager) getTotalDebt(contractAddress common.Address) (*big.Int, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"totalDebt", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("getTotalDebt, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("getTotalDebt, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	amount, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("getTotalDebt, source.NextI128 error")
	}
	return amount.ToBigInt(), nil
}

func (this *IFPoolManager) getTotalCash(contractAddress common.Address) (*big.Int, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"totalCash", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("getTotalCash, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("getTotalCash, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	amount, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("getTotalCash, source.NextI128 error")
	}
	return amount.ToBigInt(), nil
}

func (this *IFPoolManager) getInterestIndex(contractAddress common.Address) (*big.Int, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"interestIndex", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("getInterestIndex, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("getInterestIndex, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	amount, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("getInterestIndex, source.NextI128 error")
	}
	return amount.ToBigInt(), nil
}
