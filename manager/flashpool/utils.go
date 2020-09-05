package flashpool

import (
	"fmt"
	"github.com/ontio/ontology/common"
)

type FlashPoolBanner struct {
	Today uint64
	Share uint64
	Total uint64
}

type FlashPoolDetail struct {
	TotalSupply       uint64
	TotalSupplyRate   uint64
	SupplyMarketRank  []*MarketFund
	SupplyVolumeDaily uint64
	Supplier          uint64

	TotalBorrow       uint64
	TotalBorrowRate   uint64
	BorrowMarketRank  []*MarketFund
	BorrowVolumeDaily uint64
	Borrower          uint64

	TotalInsurance       uint64
	TotalInsuranceRate   uint64
	InsuranceMarketRank  []*MarketFund
	InsuranceVolumeDaily uint64
	Guarantor            uint64
}

type MarketFund struct {
	Icon string
	Name string
	Fund uint64
}
type FlashPoolAllMarket struct {
	FlashPoolAllMarket []*Market
}

type Market struct {
	Icon               string
	Name               string
	TotalSupply        uint64
	TotalSupplyRate    uint64
	SupplyApy          uint64
	SupplyApyRate      uint64
	TotalBorrow        uint64
	TotalBorrowRate    uint64
	BorrowApy          uint64
	BorrowApyRate      uint64
	TotalInsurance     uint64
	TotalInsuranceRate uint64
	InsuranceApy       uint64
	InsuranceApyRate   uint64
}

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

func (this *FlashPoolManager) getSupplyAmountByAccount(contractAddress, account common.Address) (uint64, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"balanceOf", []interface{}{account})
	if err != nil {
		return 0, fmt.Errorf("getSupplyAmountByAccount, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return 0, fmt.Errorf("getSupplyAmountByAccount, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	amount, eof := source.NextI128()
	if eof {
		return 0, fmt.Errorf("getSupplyAmountByAccount, source.NextI128 error")
	}
	return amount.ToBigInt().Uint64(), nil
}

func (this *FlashPoolManager) getBorrowAmountByAccount(contractAddress, account common.Address) (uint64, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"borrowBalanceStored", []interface{}{account})
	if err != nil {
		return 0, fmt.Errorf("getBorrowAmountByAccount, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return 0, fmt.Errorf("getBorrowAmountByAccount, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	amount, eof := source.NextI128()
	if eof {
		return 0, fmt.Errorf("getBorrowAmountByAccount, source.NextI128 error")
	}
	return amount.ToBigInt().Uint64(), nil
}

func (this *FlashPoolManager) getInsuranceAmountByAccount(contractAddress, account common.Address) (uint64, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"insuranceBalanceStored", []interface{}{account})
	if err != nil {
		return 0, fmt.Errorf("getInsuranceAmountByAccount, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return 0, fmt.Errorf("getInsuranceAmountByAccount, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	amount, eof := source.NextI128()
	if eof {
		return 0, fmt.Errorf("getInsuranceAmountByAccount, source.NextI128 error")
	}
	return amount.ToBigInt().Uint64(), nil
}

func (this *FlashPoolManager) getSupplyAmount(contractAddress common.Address) (uint64, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"totalSupply", []interface{}{})
	if err != nil {
		return 0, fmt.Errorf("getSupplyAmount, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return 0, fmt.Errorf("getSupplyAmount, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	amount, eof := source.NextI128()
	if eof {
		return 0, fmt.Errorf("getSupplyAmount, source.NextI128 error")
	}
	return amount.ToBigInt().Uint64(), nil
}

func (this *FlashPoolManager) getBorrowAmount(contractAddress common.Address) (uint64, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"totalBorrows", []interface{}{})
	if err != nil {
		return 0, fmt.Errorf("getBorrowAmount, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return 0, fmt.Errorf("getBorrowAmount, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	amount, eof := source.NextI128()
	if eof {
		return 0, fmt.Errorf("getBorrowAmount, source.NextI128 error")
	}
	return amount.ToBigInt().Uint64(), nil
}

func (this *FlashPoolManager) getInsuranceAmount(contractAddress common.Address) (uint64, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"totalInsurance", []interface{}{})
	if err != nil {
		return 0, fmt.Errorf("getInsuranceAmount, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return 0, fmt.Errorf("getInsuranceAmount, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	amount, eof := source.NextI128()
	if eof {
		return 0, fmt.Errorf("getInsuranceAmount, source.NextI128 error")
	}
	return amount.ToBigInt().Uint64(), nil
}

func (this *FlashPoolManager) getTotalDistribution(assetAddress common.Address) (uint64, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(this.contractAddress,
		"wingDistributedNum", []interface{}{assetAddress})
	if err != nil {
		return 0, fmt.Errorf("getTotalDistribution, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return 0, fmt.Errorf("getTotalDistribution, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	amount, eof := source.NextI128()
	if eof {
		return 0, fmt.Errorf("getTotalDistribution, source.NextI128 error")
	}
	return amount.ToBigInt().Uint64(), nil
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
	return ratePerBlock.ToBigInt().Uint64() * BlockPerYear, nil
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
	return ratePerBlock.ToBigInt().Uint64() * BlockPerYear, nil
}

func (this *FlashPoolManager) getInsuranceApy(contractAddress common.Address) (uint64, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress,
		"insuranceRatePerBlock", []interface{}{})
	if err != nil {
		return 0, fmt.Errorf("getInsuranceApy, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToByteArray()
	if err != nil {
		return 0, fmt.Errorf("getInsuranceApy, preExecResult.Result.ToByteArray error: %s", err)
	}
	source := common.NewZeroCopySource(r)
	ratePerBlock, eof := source.NextI128()
	if eof {
		return 0, fmt.Errorf("getInsuranceApy, source.NextI128 error")
	}
	return ratePerBlock.ToBigInt().Uint64() * BlockPerYear, nil
}

func (this *FlashPoolManager) flashPoolBanner() (*FlashPoolBanner, error) {
	return &FlashPoolBanner{
		Today: 8676,
		Share: 7644,
		Total: 3452636,
	}, nil
}

func (this *FlashPoolManager) flashPoolDetail() (*FlashPoolDetail, error) {
	return &FlashPoolDetail{
		TotalSupply:     86544,
		TotalSupplyRate: 8754,
		SupplyMarketRank: []*MarketFund{{Icon: "http://106.75.209.209/icon/eth_icon.svg", Name: "ETH", Fund: 2344},
			{Icon: "http://106.75.209.209/icon/asset_dai_icon.svg", Name: "DAI", Fund: 1234},
			{Icon: "http://106.75.209.209/icon/eth_icon.svg", Name: "BTC", Fund: 1233}},
		SupplyVolumeDaily: 24526,
		Supplier:          125,

		TotalBorrow:     2524,
		TotalBorrowRate: 4252,
		BorrowMarketRank: []*MarketFund{{Icon: "http://106.75.209.209/icon/eth_icon.svg", Name: "ETH", Fund: 535},
			{Icon: "http://106.75.209.209/icon/asset_dai_icon.svg", Name: "DAI", Fund: 234}},
		BorrowVolumeDaily: 3115,
		Borrower:          36,

		TotalInsurance:     6754,
		TotalInsuranceRate: 9632,
		InsuranceMarketRank: []*MarketFund{{Icon: "http://106.75.209.209/icon/eth_icon.svg", Name: "ETH", Fund: 2526},
			{Icon: "http://106.75.209.209/icon/asset_dai_icon.svg", Name: "DAI", Fund: 2458}},
		InsuranceVolumeDaily: 3277,
		Guarantor:            234,
	}, nil
}

func (this *FlashPoolManager) flashPoolAllMarket() (*FlashPoolAllMarket, error) {
	return &FlashPoolAllMarket{
		FlashPoolAllMarket: []*Market{
			{
				Icon:               "http://106.75.209.209/icon/eth_icon.svg",
				Name:               "ETH",
				TotalSupply:        2526,
				TotalSupplyRate:    2,
				SupplyApy:          4468,
				SupplyApyRate:      3,
				TotalBorrow:        25267,
				TotalBorrowRate:    23,
				BorrowApy:          563,
				BorrowApyRate:      6,
				TotalInsurance:     8265,
				TotalInsuranceRate: 6,
				InsuranceApy:       256,
				InsuranceApyRate:   9,
			},
			{
				Icon:               "http://106.75.209.209/icon/asset_dai_icon.svg",
				Name:               "DAI",
				TotalSupply:        2526,
				TotalSupplyRate:    1,
				SupplyApy:          3526,
				SupplyApyRate:      2,
				TotalBorrow:        2415,
				TotalBorrowRate:    3,
				BorrowApy:          241,
				BorrowApyRate:      4,
				TotalInsurance:     3473,
				TotalInsuranceRate: 5,
				InsuranceApy:       2541,
				InsuranceApyRate:   6,
			}},
	}, nil
}
