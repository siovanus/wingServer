package flashpool

import (
	"fmt"
	"github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/manager/governance"
	flash_ctrl "github.com/wing-groups/wing-contract-tools/contracts/flash-ctrl"
	"math/big"
)

func (this *FlashPoolManager) assetPrice(asset string) (*big.Int, error) {
	return this.Oracle.GetUnderlyingPrice(asset)
}

func (this *FlashPoolManager) GetAllMarkets() ([]common.Address, error) {
	return this.Comptroller.AllMarkets()
}

func (this *FlashPoolManager) getAssetsIn(account common.Address) ([]common.Address, error) {
	return this.Comptroller.AssetsIn(account)
}

func (this *FlashPoolManager) getFTokenAmount(contractAddress, account common.Address) (*big.Int, error) {
	return this.FlashTokenMap[contractAddress].BalanceOf(account)
}

func (this *FlashPoolManager) getITokenAmount(contractAddress, account common.Address) (*big.Int, error) {
	iAddress, err := this.FlashTokenMap[contractAddress].InsuranceAddr()
	if err != nil {
		return nil, fmt.Errorf("getITokenAmount, this.FlashToken.InsuranceAddr error: %s", err)
	}
	return this.FlashTokenMap[iAddress].BalanceOf(account)
}

func (this *FlashPoolManager) getBorrowAmount(contractAddress, account common.Address) (*big.Int, error) {
	return this.FlashTokenMap[contractAddress].BorrowBalanceStored(account)
}

func (this *FlashPoolManager) getTotalSupply(contractAddress common.Address) (*big.Int, error) {
	totalCash, err := this.FlashTokenMap[contractAddress].GetCash()
	if err != nil {
		return nil, fmt.Errorf("getTotalSupply, this.FlashToken.GetCash error: %s", err)
	}
	totalBorrows, err := this.FlashTokenMap[contractAddress].TotalBorrows()
	if err != nil {
		return nil, fmt.Errorf("getTotalSupply, this.FlashToken.TotalBorrows error: %s", err)
	}
	amount := new(big.Int).Add(totalCash, totalBorrows)

	return amount, nil
}

func (this *FlashPoolManager) getTotalBorrows(contractAddress common.Address) (*big.Int, error) {
	return this.FlashTokenMap[contractAddress].TotalBorrows()
}

func (this *FlashPoolManager) getTotalReserves(contractAddress common.Address) (*big.Int, error) {
	return this.FlashTokenMap[contractAddress].TotalReserves()
}

func (this *FlashPoolManager) getTotalInsurance(contractAddress common.Address) (*big.Int, error) {
	iAddress, err := this.FlashTokenMap[contractAddress].InsuranceAddr()
	if err != nil {
		return nil, fmt.Errorf("getTotalInsurance, this.FlashToken.InsuranceAddr error: %s", err)
	}
	return this.FlashTokenMap[iAddress].GetCash()
}

func (this *FlashPoolManager) getInsuranceAddress(contractAddress common.Address) (common.Address, error) {
	return this.FlashTokenMap[contractAddress].InsuranceAddr()
}

func (this *FlashPoolManager) getTotalDistribution(assetAddress common.Address) (*big.Int, error) {
	result, err := this.Comptroller.WingDistributedNum(assetAddress)
	if err != nil {
		return nil, fmt.Errorf("getTotalDistribution, this.Comptroller.WingDistributedNum error: %s", err)
	}
	if this.cfg.AssetMap[assetAddress.ToHexString()] == "pWBTC" {
		return new(big.Int).Sub(result, GAP), nil
	} else {
		return result, nil
	}
}

func (this *FlashPoolManager) getExchangeRate(contractAddress common.Address) (*big.Int, error) {
	return this.FlashTokenMap[contractAddress].ExchangeRateStored()
}

func (this *FlashPoolManager) getBorrowIndex(contractAddress common.Address) (*big.Int, error) {
	return this.FlashTokenMap[contractAddress].BorrowIndex()
}

func (this *FlashPoolManager) getReserveFactor(contractAddress common.Address) (*big.Int, error) {
	return this.FlashTokenMap[contractAddress].ReserveFactorMantissa()
}

func (this *FlashPoolManager) getSupplyApy(contractAddress common.Address) (*big.Int, error) {
	ratePerBlock, err := this.FlashTokenMap[contractAddress].SupplyRatePerBlock()
	if err != nil {
		return nil, fmt.Errorf("getSupplyApy, this.FlashToken.SupplyRatePerBlock error: %s", err)
	}

	result := new(big.Int).Mul(ratePerBlock, new(big.Int).SetUint64(governance.YearSecond))
	return result, nil
}

func (this *FlashPoolManager) getBorrowRatePerBlock(contractAddress common.Address) (*big.Int, error) {
	return this.FlashTokenMap[contractAddress].BorrowRatePerBlock()
}

func (this *FlashPoolManager) getBorrowApy(contractAddress common.Address) (*big.Int, error) {
	ratePerBlock, err := this.FlashTokenMap[contractAddress].BorrowRatePerBlock()
	if err != nil {
		return nil, fmt.Errorf("getBorrowApy, this.FlashToken.BorrowRatePerBlock error: %s", err)
	}

	result := new(big.Int).Mul(ratePerBlock, new(big.Int).SetUint64(governance.YearSecond))
	return result, nil
}

func (this *FlashPoolManager) getMarketMeta(market common.Address) (*flash_ctrl.MarketMeta, error) {
	return this.Comptroller.MarketMeta(market)
}

// how much user can borrow
func (this *FlashPoolManager) getAccountLiquidity(account common.Address) (*flash_ctrl.AccountLiquidity, error) {
	return this.Comptroller.GetAccountLiquidity(account)
}

func (this *FlashPoolManager) getWingAccrued(account common.Address) (*big.Int, error) {
	return this.Comptroller.WingAccrued(account)
}

func (this *FlashPoolManager) getClaimWing(holder common.Address) (*big.Int, error) {
	_, result, err := this.Comptroller.ClaimWing(holder, true)
	if err != nil {
		return nil, fmt.Errorf("getClaimWing, this.Comptroller.ClaimWing error: %s", err)
	}
	return result, nil
}

func (this *FlashPoolManager) getWingSpeeds(contractAddress common.Address) (*big.Int, error) {
	return this.Comptroller.WingSpeeds(contractAddress)
}

func (this *FlashPoolManager) getWingSBIPortion(contractAddress common.Address) (*flash_ctrl.WingSBI, error) {
	return this.Comptroller.WingSBIPortion(contractAddress)
}

func (this *FlashPoolManager) getClaimWingAtMarket(account common.Address, contractAddresses []common.Address) (*big.Int, error) {
	_, result, err := this.Comptroller.ClaimWingAtMarkets(account, contractAddresses, true)
	if err != nil {
		return nil, fmt.Errorf("getClaimWingAtMarket, this.Comptroller.ClaimWingAtMarkets error: %s", err)
	}
	return result, nil
}
