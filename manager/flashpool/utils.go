package flashpool

import (
	"fmt"
	"github.com/ontio/ontology/common"
	"math/big"
)

func (this *FlashPoolManager) assetPrice(asset string) (*big.Int, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(this.oracleAddress,
		"getUnderlyingPrice", []interface{}{asset})
	if err != nil {
		return nil, fmt.Errorf("assetPrice, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToInteger()
	if err != nil {
		return nil, fmt.Errorf("assetPrice, preExecResult.Result.ToInteger error: %s", err)
	}
	return r, nil
}

func (this *FlashPoolManager) GetAllMarkets() ([]common.Address, error) {
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
		"assetsIn", []interface{}{account})
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
		"balanceOfUnderlying", []interface{}{account})
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
	totalCash, err := this.getCash(contractAddress)
	if err != nil {
		return nil, fmt.Errorf("getSupplyAmount, this.getCash error: %s", err)
	}
	totalBorrows, err := this.getBorrowAmount(contractAddress)
	if err != nil {
		return nil, fmt.Errorf("getSupplyAmount, this.getBorrowAmount error: %s", err)
	}
	amount := new(big.Int).Add(totalCash, totalBorrows)

	return amount, nil
}

func (this *FlashPoolManager) getCash(contractAddress common.Address) (*big.Int, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"getCash", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("getCash, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("getCash, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	amount, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("getCash, source.NextI128 error")
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

func (this *FlashPoolManager) getTotalReserves(contractAddress common.Address) (*big.Int, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"totalReserves", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("getTotalReserves, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("getTotalReserves, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	amount, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("getTotalReserves, source.NextI128 error")
	}
	return amount.ToBigInt(), nil
}

func (this *FlashPoolManager) getInsuranceAmount(contractAddress common.Address) (*big.Int, error) {
	insuranceAddress, err := this.GetInsuranceAddress(contractAddress)
	if err != nil {
		return nil, fmt.Errorf("getInsuranceAmount, this.getInsuranceAddress error: %s", err)
	}

	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(insuranceAddress,
		"getCash", []interface{}{})
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
	if this.cfg.AssetMap[assetAddress.ToHexString()] == "pWBTC" {
		return new(big.Int).Sub(amount.ToBigInt(), GAP), nil
	} else {
		return amount.ToBigInt(), nil
	}
}

func (this *FlashPoolManager) getSupplyApy(contractAddress common.Address) (*big.Int, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"supplyRatePerBlock", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("getSupplyApy, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("getSupplyApy, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	ratePerBlock, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("getSupplyApy, source.NextI128 error")
	}
	result := new(big.Int).Mul(ratePerBlock.ToBigInt(), new(big.Int).SetUint64(BlockPerYear))
	return result, nil
}

func (this *FlashPoolManager) getBorrowRatePerBlock(contractAddress common.Address) (*big.Int, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"borrowRatePerBlock", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("getBorrowRatePerBlock, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("getBorrowRatePerBlock, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	ratePerBlock, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("getBorrowRatePerBlock, source.NextI128 error")
	}
	result := ratePerBlock.ToBigInt()
	return result, nil
}

func (this *FlashPoolManager) getBorrowApy(contractAddress common.Address) (*big.Int, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"borrowRatePerBlock", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("getBorrowApy, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("getBorrowApy, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	ratePerBlock, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("getBorrowApy, source.NextI128 error")
	}
	result := new(big.Int).Mul(ratePerBlock.ToBigInt(), new(big.Int).SetUint64(BlockPerYear))
	return result, nil
}

func (this *FlashPoolManager) GetInsuranceAddress(contractAddress common.Address) (common.Address, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"insuranceAddr", []interface{}{})
	if err != nil {
		return common.ADDRESS_EMPTY, fmt.Errorf("getInsuranceAddress, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return common.ADDRESS_EMPTY, fmt.Errorf("getInsuranceAddress, preExecResult.Result.ToByteArray error: %s", err)
	}
	insuranceAddress, err := common.AddressParseFromBytes(r)
	if err != nil {
		return common.ADDRESS_EMPTY, fmt.Errorf("getInsuranceAddress, common.AddressParseFromBytes error: %s", err)
	}
	return insuranceAddress, nil
}

func (this *FlashPoolManager) getInsuranceApy(contractAddress common.Address) (*big.Int, error) {
	insuranceAddress, err := this.GetInsuranceAddress(contractAddress)
	if err != nil {
		return nil, fmt.Errorf("getInsuranceApy, this.getInsuranceAddress error: %s", err)
	}

	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(insuranceAddress,
		"supplyRatePerBlock", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("getInsuranceApy, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("getInsuranceApy, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	ratePerBlock, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("getInsuranceApy, source.NextI128 error")
	}
	result := new(big.Int).Mul(ratePerBlock.ToBigInt(), new(big.Int).SetUint64(BlockPerYear))
	return result, nil
}

type MarketMeta struct {
	Addr          common.Address
	InsuranceAddr common.Address

	IsListed    bool
	ReceiveWing bool

	WingWeight               *big.Int
	CollateralFactorMantissa *big.Int
}

func DeserializeMarketMeta(data []byte) (*MarketMeta, error) {
	source := common.NewZeroCopySource(data)
	addr, eof := source.NextAddress()
	if eof {
		return nil, fmt.Errorf("read addr eof")
	}
	insurance, eof := source.NextAddress()
	if eof {
		return nil, fmt.Errorf("read insurance eof")
	}
	isListed, irr, eof := source.NextBool()
	if irr {
		return nil, fmt.Errorf("read isListed irr")
	}
	if eof {
		return nil, fmt.Errorf("read isListed eof")
	}
	receiveWing, irr, eof := source.NextBool()
	if irr {
		return nil, fmt.Errorf("read receiveWing irr")
	}
	if eof {
		return nil, fmt.Errorf("read receiveWing eof")
	}
	wingWeight, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("read wingWeight eof")
	}
	collateralFactor, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("read collateralFactor eof")
	}
	return &MarketMeta{
		Addr:                     addr,
		InsuranceAddr:            insurance,
		IsListed:                 isListed,
		ReceiveWing:              receiveWing,
		WingWeight:               wingWeight.ToBigInt(),
		CollateralFactorMantissa: collateralFactor.ToBigInt(),
	}, nil
}

func (this *FlashPoolManager) getMarketMeta(market common.Address) (*MarketMeta, error) {
	method := "marketMeta"
	params := []interface{}{market}
	res, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(this.contractAddress, method, params)
	if err != nil {
		return nil, fmt.Errorf("MarketMeta: %s", err)
	}
	data, err := res.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("MarketMeta: %s", err)
	}
	result, err := DeserializeMarketMeta(data)
	if err != nil {
		return nil, fmt.Errorf("MarketMeta: %s", err)
	}
	return result, nil
}

type AccountLiquidity struct {
	Error     string
	Liquidity common.I128
	Shortfall common.I128
}

func DeserializeAccountLiquidity(data []byte) (*AccountLiquidity, error) {
	source := common.NewZeroCopySource(data)
	errStr, _, ill, eof := source.NextString()
	if ill {
		return nil, fmt.Errorf("read errStr ill")
	}
	if eof {
		return nil, fmt.Errorf("read errStr eof")
	}
	liquidity, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("read liquidity eof")
	}
	shortfall, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("read shortfall eof")
	}
	return &AccountLiquidity{
		Error:     errStr,
		Liquidity: liquidity,
		Shortfall: shortfall,
	}, nil
}

// how much user can borrow
func (this *FlashPoolManager) getAccountLiquidity(account common.Address) (*AccountLiquidity, error) {
	method := "getAccountLiquidity"
	params := []interface{}{account}
	res, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(this.contractAddress, method, params)
	if err != nil {
		return nil, fmt.Errorf("GetAccountLiquidity: %s", err)
	}
	data, err := res.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("GetAccountLiquidity: %s", err)
	}
	result, err := DeserializeAccountLiquidity(data)
	if err != nil {
		return nil, fmt.Errorf("GetAccountLiquidity: %s", err)
	}
	return result, nil
}

func (this *FlashPoolManager) getWingAccrued(account common.Address) (*big.Int, error) {
	method := "wingAccrued"
	params := []interface{}{account}
	res, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(this.contractAddress, method, params)
	if err != nil {
		return nil, fmt.Errorf("getWingAccrued, this.sdk.WasmVM.PreExecInvokeWasmVMContract: %s", err)
	}
	return res.Result.ToInteger()
}

func (this *FlashPoolManager) getClaimWing(holder common.Address) (*big.Int, error) {
	method := "claimWing"
	params := []interface{}{holder}
	res, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(this.contractAddress, method, params)
	if err != nil {
		return nil, fmt.Errorf("ClaimWing, this.sdk.WasmVM.PreExecInvokeWasmVMContract: %s", err)
	}
	data, err := res.Result.ToByteArray()
	if err != nil {
		err = fmt.Errorf("ClaimWing: %s", err)
		return nil, err
	}
	source := common.NewZeroCopySource(data)
	r, eof := source.NextI128()
	if eof {
		err = fmt.Errorf("ClaimWing: read eof")
		return nil, err
	}
	return r.ToBigInt(), nil
}

func (this *FlashPoolManager) getWingSpeeds(contractAddress common.Address) (*big.Int, error) {
	method := "wingSpeeds"
	params := []interface{}{contractAddress}
	res, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(this.contractAddress, method, params)
	if err != nil {
		return nil, fmt.Errorf("getWingSpeeds, this.sdk.WasmVM.PreExecInvokeWasmVMContract: %s", err)
	}
	data, err := res.Result.ToByteArray()
	if err != nil {
		err = fmt.Errorf("getWingSpeeds: %s", err)
		return nil, err
	}
	source := common.NewZeroCopySource(data)
	r, eof := source.NextI128()
	if eof {
		err = fmt.Errorf("getWingSpeeds: read eof")
		return nil, err
	}
	return r.ToBigInt(), nil
}

type WingSBIPortion struct {
	SupplyPortion    common.I128
	BorrowPortion    common.I128
	InsurancePortion common.I128
}

func DeserializeWingSBIPortion(data []byte) (*WingSBIPortion, error) {
	source := common.NewZeroCopySource(data)
	supplyPortion, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("read supplyPortion eof")
	}
	borrowPortion, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("read borrowPortion eof")
	}
	insurancePortion, eof := source.NextI128()
	if eof {
		return nil, fmt.Errorf("read insurancePortion eof")
	}
	return &WingSBIPortion{
		SupplyPortion:    supplyPortion,
		BorrowPortion:    borrowPortion,
		InsurancePortion: insurancePortion,
	}, nil
}

func (this *FlashPoolManager) getWingSBIPortion(contractAddress common.Address) (*WingSBIPortion, error) {
	method := "wingSBIPortion"
	params := []interface{}{contractAddress}
	res, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(this.contractAddress, method, params)
	if err != nil {
		return nil, fmt.Errorf("getWingSBIPortion, this.sdk.WasmVM.PreExecInvokeWasmVMContract: %s", err)
	}
	data, err := res.Result.ToByteArray()
	if err != nil {
		err = fmt.Errorf("getWingSBIPortion: %s", err)
		return nil, err
	}
	wingSBIPortion, err := DeserializeWingSBIPortion(data)
	if err != nil {
		return nil, err
	}
	return wingSBIPortion, nil
}

func (this *FlashPoolManager) getClaimWingAtMarket(account common.Address, contractAddresses []interface{}) (*big.Int, error) {
	method := "claimWingAtMarkets"
	params := []interface{}{account, contractAddresses}
	res, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(this.contractAddress, method, params)
	if err != nil {
		return nil, fmt.Errorf("getClaimWingAtMarket, this.sdk.WasmVM.PreExecInvokeWasmVMContract: %s", err)
	}
	data, err := res.Result.ToByteArray()
	if err != nil {
		err = fmt.Errorf("getClaimWingAtMarket error: %s", err)
		return nil, err
	}
	source := common.NewZeroCopySource(data)
	r, eof := source.NextI128()
	if eof {
		err = fmt.Errorf("getClaimWingAtMarket: read eof")
		return nil, err
	}
	return r.ToBigInt(), nil
}
