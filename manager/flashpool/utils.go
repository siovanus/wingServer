package flashpool

import (
	"fmt"
	"github.com/ontio/ontology/common"
	"math/big"
)

func (this *FlashPoolManager) assetPrice(asset string) (uint64, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(this.oracleAddress,
		"getUnderlyingPrice", []interface{}{asset})
	if err != nil {
		return 0, fmt.Errorf("assetPrice, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToInteger()
	if err != nil {
		return 0, fmt.Errorf("assetPrice, preExecResult.Result.ToInteger error: %s", err)
	}
	return r.Uint64(), nil
}

func (this *FlashPoolManager) getAllMarkets() ([]common.Address, error) {
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
	allMarkets := make([]common.Address, 0)
	l, _, irregular, eof := source.NextVarUint()
	if irregular || eof {
		return nil, fmt.Errorf("getAllMarkets, source.NextVarUint error")
	}
	for i := 0; uint64(i) < l; i++ {
		addr, eof := source.NextAddress()
		if eof {
			return nil, fmt.Errorf("getAllMarkets, source.NextAddress error")
		}
		allMarkets = append(allMarkets, addr)
	}
	return allMarkets, nil
}

func (this *FlashPoolManager) getAssetsIn(account common.Address) ([]common.Address, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(this.contractAddress,
		"assetsIn", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("getAssetsIn, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("getAssetsIn, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	assetsIn := make([]common.Address, 0)
	l, _, irregular, eof := source.NextVarUint()
	if irregular || eof {
		return nil, fmt.Errorf("getAssetsIn, source.NextVarUint error")
	}
	for i := 0; uint64(i) < l; i++ {
		addr, eof := source.NextAddress()
		if eof {
			return nil, fmt.Errorf("getAssetsIn, source.NextAddress error")
		}
		assetsIn = append(assetsIn, addr)
	}
	return assetsIn, nil
}

func (this *FlashPoolManager) getSupplyAmountByAccount(contractAddress, account common.Address) (*big.Int, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"balanceOf", []interface{}{account})
	if err != nil {
		return nil, fmt.Errorf("getSupplyAmountByAccount, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("getSupplyAmountByAccount, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	amount, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("getSupplyAmountByAccount, source.NextI128 error")
	}
	return amount.ToBigInt(), nil
}

func (this *FlashPoolManager) getBorrowAmountByAccount(contractAddress, account common.Address) (*big.Int, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"borrowBalanceStored", []interface{}{account})
	if err != nil {
		return nil, fmt.Errorf("getBorrowAmountByAccount, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("getBorrowAmountByAccount, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	amount, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("getBorrowAmountByAccount, source.NextI128 error")
	}
	return amount.ToBigInt(), nil
}

func (this *FlashPoolManager) getInsuranceAmountByAccount(contractAddress, account common.Address) (*big.Int, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"insuranceAddr", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("getInsuranceAmountByAccount, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("getInsuranceAmountByAccount, preExecResult.Result.ToByteArray error: %s", err)
	}
	insuranceAddress, err := common.AddressParseFromBytes(r)
	if err != nil {
		return nil, fmt.Errorf("getInsuranceAmountByAccount, common.AddressParseFromBytes error: %s", err)
	}

	preExecResult, err = this.sdk.WasmVM.PreExecInvokeWasmVMContract(insuranceAddress,
		"balanceOfUnderlying", []interface{}{account})
	if err != nil {
		return nil, fmt.Errorf("getInsuranceAmountByAccount, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err = preExecResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("getInsuranceAmountByAccount, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	amount, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("getInsuranceAmountByAccount, source.NextI128 error")
	}
	return amount.ToBigInt(), nil
}

func (this *FlashPoolManager) getSupplyAmount(contractAddress common.Address) (*big.Int, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"totalSupply", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("getSupplyAmount, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("getSupplyAmount, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	amount, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("getSupplyAmount, source.NextI128 error")
	}
	return amount.ToBigInt(), nil
}

func (this *FlashPoolManager) getBorrowAmount(contractAddress common.Address) (*big.Int, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"totalBorrows", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("getBorrowAmount, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("getBorrowAmount, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	amount, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("getBorrowAmount, source.NextI128 error")
	}
	return amount.ToBigInt(), nil
}

func (this *FlashPoolManager) getInsuranceAmount(contractAddress common.Address) (*big.Int, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"totalInsurance", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("getInsuranceAmount, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("getInsuranceAmount, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	amount, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("getInsuranceAmount, source.NextI128 error")
	}
	return amount.ToBigInt(), nil
}

func (this *FlashPoolManager) getTotalDistribution(assetAddress common.Address) (*big.Int, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(this.contractAddress,
		"wingDistributedNum", []interface{}{assetAddress})
	if err != nil {
		return nil, fmt.Errorf("getTotalDistribution, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("getTotalDistribution, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	amount, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("getTotalDistribution, source.NextI128 error")
	}
	return amount.ToBigInt(), nil
}

func (this *FlashPoolManager) getSupplyApy(contractAddress common.Address) (uint64, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"supplyRatePerBlock", []interface{}{})
	if err != nil {
		return 0, fmt.Errorf("getSupplyApy, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return 0, fmt.Errorf("getSupplyApy, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	ratePerBlock, eof := source.NextI128()
	if eof {
		return 0, fmt.Errorf("getSupplyApy, source.NextI128 error")
	}
	result := new(big.Int).Div(new(big.Int).Mul(ratePerBlock.ToBigInt(), new(big.Int).SetUint64(BlockPerYear)),
		FrontPercentageDecimal).Uint64()
	return result, nil
}

func (this *FlashPoolManager) getBorrowApy(contractAddress common.Address) (uint64, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"borrowRatePerBlock", []interface{}{})
	if err != nil {
		return 0, fmt.Errorf("getBorrowApy, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return 0, fmt.Errorf("getBorrowApy, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	ratePerBlock, eof := source.NextI128()
	if eof {
		return 0, fmt.Errorf("getBorrowApy, source.NextI128 error")
	}
	result := new(big.Int).Div(new(big.Int).Mul(ratePerBlock.ToBigInt(), new(big.Int).SetUint64(BlockPerYear)),
		FrontPercentageDecimal).Uint64()
	return result, nil
}

func (this *FlashPoolManager) getInsuranceApy(contractAddress common.Address) (uint64, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"insuranceAddr", []interface{}{})
	if err != nil {
		return 0, fmt.Errorf("getInsuranceApy, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return 0, fmt.Errorf("getInsuranceApy, preExecResult.Result.ToByteArray error: %s", err)
	}
	insuranceAddress, err := common.AddressParseFromBytes(r)
	if err != nil {
		return 0, fmt.Errorf("getInsuranceApy, common.AddressParseFromBytes error: %s", err)
	}

	preExecResult, err = this.sdk.WasmVM.PreExecInvokeWasmVMContract(insuranceAddress,
		"supplyRatePerBlock", []interface{}{})
	if err != nil {
		return 0, fmt.Errorf("getInsuranceApy, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err = preExecResult.Result.ToByteArray()
	if err != nil {
		return 0, fmt.Errorf("getInsuranceApy, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	ratePerBlock, eof := source.NextI128()
	if eof {
		return 0, fmt.Errorf("getInsuranceApy, source.NextI128 error")
	}
	result := new(big.Int).Div(new(big.Int).Mul(ratePerBlock.ToBigInt(), new(big.Int).SetUint64(BlockPerYear)),
		FrontPercentageDecimal).Uint64()
	return result, nil
}

type MarketMeta struct {
	Address          common.Address
	IsList           bool
	CollateralFactor common.I128
}

func (this *MarketMeta) Deserialization(source *common.ZeroCopySource) error {
	address, eof := source.NextAddress()
	if eof {
		return fmt.Errorf("address deserialization error eof")
	}
	isList, irregular, eof := source.NextBool()
	if irregular || eof {
		return fmt.Errorf("isList deserialization error eof")
	}
	collateralFactor, eof := source.NextI128()
	if eof {
		return fmt.Errorf("collateralFactor deserialization error eof")
	}
	this.Address = address
	this.IsList = isList
	this.CollateralFactor = collateralFactor
	return nil
}

func (this *FlashPoolManager) getMarketMeta() (*MarketMeta, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(this.contractAddress,
		"marketMeta", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("getMarketMeta, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("getMarketMeta, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	marketMeta := new(MarketMeta)
	err = marketMeta.Deserialization(source)
	if err != nil {
		return nil, fmt.Errorf("getMarketMeta, marketMeta.Deserialization error")
	}
	return marketMeta, nil
}
