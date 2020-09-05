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

type UserFlashPoolOverview struct {
	SupplyBalance    uint64
	BorrowBalance    uint64
	InsuranceBalance uint64
	BorrowLimit      uint64
	NetApy           uint64

	CurrentSupply []*Supply

	AllMarket []*UserMarket
}

type Supply struct {
	Icon          string
	Name          string
	SupplyBalance uint64
	Apy           uint64
	Earned        uint64
	IfCollateral  bool
}

type UserMarket struct {
	Icon            string
	Name            string
	IfCollateral    bool
	SupplyApy       uint64
	BorrowApy       uint64
	BorrowLiquidity uint64
	InsuranceApy    uint64
	InsuranceAmount uint64
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

func (this *FlashPoolManager) getSupplyAmount(address common.Address) (uint64, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(address,
		"balanceOf", []interface{}{})
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

func (this *FlashPoolManager) getBorrowAmount(address common.Address) (uint64, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(address,
		"borrowBalanceStored", []interface{}{})
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

func (this *FlashPoolManager) getInsuranceAmount(address common.Address) (uint64, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(address,
		"insuranceBalanceStored", []interface{}{})
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

func (this *FlashPoolManager) getSupplyApy(address common.Address) (uint64, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(address,
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

func (this *FlashPoolManager) getBorrowApy(address common.Address) (uint64, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(address,
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

func (this *FlashPoolManager) getInsuranceApy(address common.Address) (uint64, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(address,
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

func (this *FlashPoolManager) userFlashPoolOverview(address string) (*UserFlashPoolOverview, error) {
	return &UserFlashPoolOverview{
		SupplyBalance:    22626,
		BorrowBalance:    23525,
		InsuranceBalance: 252,
		BorrowLimit:      2355,
		NetApy:           252,

		CurrentSupply: []*Supply{
			{
				Icon:          "http://106.75.209.209/icon/eth_icon.svg",
				Name:          "ETH",
				SupplyBalance: 2326,
				Apy:           266,
				Earned:        67,
				IfCollateral:  true,
			},
			{
				Icon:          "http://106.75.209.209/icon/eth_icon.svg",
				Name:          "ETH",
				SupplyBalance: 4627,
				Apy:           54,
				Earned:        367,
				IfCollateral:  false,
			},
		},

		AllMarket: []*UserMarket{
			{
				Icon:            "http://106.75.209.209/icon/eth_icon.svg",
				Name:            "ETH",
				IfCollateral:    true,
				SupplyApy:       4468,
				BorrowApy:       563,
				BorrowLiquidity: 255,
				InsuranceApy:    256,
				InsuranceAmount: 2526,
			},
			{
				Icon:            "http://106.75.209.209/icon/asset_dai_icon.svg",
				Name:            "DAI",
				IfCollateral:    true,
				SupplyApy:       3526,
				BorrowApy:       241,
				BorrowLiquidity: 255,
				InsuranceApy:    2541,
				InsuranceAmount: 2526,
			},
		},
	}, nil
}
