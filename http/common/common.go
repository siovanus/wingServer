package common

import (
	"fmt"
	"github.com/ontio/ontology/common"
)

const (
	FLASHPOOLMARKETDISTRIBUTION = "/api/v1/flashpoolmarketdistribution"
	POOLDISTRIBUTION            = "/api/v1/pooldistribution"
	GOVBANNEROVERVIEW           = "/api/v1/govbanneroverview"
	GOVBANNER                   = "/api/v1/govbanner"
	FLASHPOOLBANNER             = "/api/v1/flashpoolbanner"

	FLASHPOOLDETAIL       = "/api/v1/flashpooldetail"
	FLASHPOOLALLMARKET    = "/api/v1/flashpoolallmarket"
	USERFLASHPOOLOVERVIEW = "/api/v1/userflashpooloverview"

	ASSETPRICE = "/api/v1/assetprice"
)

const (
	ACTION_FLASHPOOLMARKETDISTRIBUTION = "flashpoolmarketdistribution"
	ACTION_POOLDISTRIBUTION            = "pooldistribution"
	ACTION_GOVBANNEROVERVIEW           = "govbanneroverview"
	ACTION_GOVBANNER                   = "govbanner"
	ACTION_FLASHPOOLBANNER             = "flashpoolbanner"

	ACTION_FLASHPOOLDETAIL       = "flashpooldetail"
	ACTION_FLASHPOOLALLMARKET    = "flashpoolallmarket"
	ACTION_USERFLASHPOOLOVERVIEW = "userflashpooloverview"

	ACTION_ASSETPRICE = "assetprice"
)

type Response struct {
	Action string      `json:"action"`
	Desc   string      `json:"desc"`
	Error  uint32      `json:"error"`
	Result interface{} `json:"result"`
}

type AssetPriceRequest struct {
	Id    string
	Asset string
}

type AssetPriceResponse struct {
	Id    string
	Price uint64
}

type GovBannerOverview struct {
	Remain20 uint64
	Remain80 uint64
}

type GovBanner struct {
	Daily       uint64
	Distributed uint64
}

type FlashPoolMarketDistribution struct {
	FlashPoolMarketDistribution []*Distribution
}

type PoolDistribution struct {
	PoolDistribution []*Distribution
}

type Distribution struct {
	Icon            string
	Name            string
	PerDay          uint64
	SupplyAmount    uint64
	BorrowAmount    uint64
	InsuranceAmount uint64
	Total           uint64
}

type FlashPoolBanner struct {
	Today uint64
	Share uint64
	Total uint64
}

type FlashPoolDetail struct {
	TotalSupply       uint64
	TotalSupplyRate   int64
	SupplyMarketRank  []*MarketFund
	SupplyVolumeDaily int64

	TotalBorrow       uint64
	TotalBorrowRate   int64
	BorrowMarketRank  []*MarketFund
	BorrowVolumeDaily int64

	TotalInsurance       uint64
	TotalInsuranceRate   int64
	InsuranceMarketRank  []*MarketFund
	InsuranceVolumeDaily int64
}

type MarketFund struct {
	Icon string
	Name string
	Fund uint64
}

type UserFlashPoolOverviewRequest struct {
	Id      string
	Address string
}

type UserFlashPoolOverviewResponse struct {
	Id                    string
	Address               string
	UserFlashPoolOverview *UserFlashPoolOverview
}

type UserFlashPoolOverview struct {
	SupplyBalance    uint64
	BorrowBalance    uint64
	InsuranceBalance uint64
	BorrowLimit      uint64
	NetApy           int64

	CurrentSupply    []*Supply
	CurrentBorrow    []*Borrow
	CurrentInsurance []*Insurance

	AllMarket []*UserMarket
}

func (this *UserFlashPoolOverview) HalfSerialization(sink *common.ZeroCopySink) {
	sink.WriteUint64(uint64(len(this.CurrentSupply)))
	for _, v := range this.CurrentSupply {
		v.Serialization(sink)
	}

	sink.WriteUint64(uint64(len(this.CurrentBorrow)))
	for _, v := range this.CurrentBorrow {
		v.Serialization(sink)
	}

	sink.WriteUint64(uint64(len(this.CurrentInsurance)))
	for _, v := range this.CurrentInsurance {
		v.Serialization(sink)
	}

	sink.WriteUint64(uint64(len(this.AllMarket)))
	for _, v := range this.AllMarket {
		v.Serialization(sink)
	}
}

func (this *UserFlashPoolOverview) HalfDeserialization(source *common.ZeroCopySource) error {
	l, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("currentSupply length deserialization error eof")
	}
	currentSupply := make([]*Supply, 0)
	for i := 0; uint64(i) < l; i++ {
		t := new(Supply)
		err := t.Deserialization(source)
		if err != nil {
			return fmt.Errorf("currentSupply deserialization error eof")
		}
		currentSupply = append(currentSupply, t)
	}

	l, eof = source.NextUint64()
	if eof {
		return fmt.Errorf("currentBorrow length deserialization error eof")
	}
	currentBorrow := make([]*Borrow, 0)
	for i := 0; uint64(i) < l; i++ {
		t := new(Borrow)
		err := t.Deserialization(source)
		if err != nil {
			return fmt.Errorf("currentBorrow deserialization error eof")
		}
		currentBorrow = append(currentBorrow, t)
	}

	l, eof = source.NextUint64()
	if eof {
		return fmt.Errorf("currentInsurance length deserialization error eof")
	}
	currentInsurance := make([]*Insurance, 0)
	for i := 0; uint64(i) < l; i++ {
		t := new(Insurance)
		err := t.Deserialization(source)
		if err != nil {
			return fmt.Errorf("currentInsurance deserialization error eof")
		}
		currentInsurance = append(currentInsurance, t)
	}

	l, eof = source.NextUint64()
	if eof {
		return fmt.Errorf("allMarket length deserialization error eof")
	}
	allMarket := make([]*UserMarket, 0)
	for i := 0; uint64(i) < l; i++ {
		t := new(UserMarket)
		err := t.Deserialization(source)
		if err != nil {
			return fmt.Errorf("allMarket deserialization error eof")
		}
		allMarket = append(allMarket, t)
	}

	this.CurrentSupply = currentSupply
	this.CurrentBorrow = currentBorrow
	this.CurrentInsurance = currentInsurance
	this.AllMarket = allMarket
	return nil
}

type Supply struct {
	Icon             string
	Name             string
	SupplyDollar     uint64
	SupplyBalance    uint64
	Apy              uint64
	CollateralFactor uint64
	IfCollateral     bool
}

func (this *Supply) Serialization(sink *common.ZeroCopySink) {
	sink.WriteString(this.Icon)
	sink.WriteString(this.Name)
	sink.WriteUint64(this.SupplyDollar)
	sink.WriteUint64(this.SupplyBalance)
	sink.WriteUint64(this.Apy)
	sink.WriteUint64(this.CollateralFactor)
	sink.WriteBool(this.IfCollateral)
}

func (this *Supply) Deserialization(source *common.ZeroCopySource) error {
	icon, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("icon deserialization error eof")
	}
	name, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("name deserialization error eof")
	}
	supplyDollar, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("supplyDollar deserialization error eof")
	}
	supplyBalance, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("supplyBalance deserialization error eof")
	}
	apy, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("apy deserialization error eof")
	}
	collateralFactor, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("collateralFactor deserialization error eof")
	}
	ifCollateral, irregular, eof := source.NextBool()
	if irregular || eof {
		return fmt.Errorf("ifCollateral deserialization error eof")
	}
	this.Icon = icon
	this.Name = name
	this.SupplyDollar = supplyDollar
	this.SupplyBalance = supplyBalance
	this.Apy = apy
	this.CollateralFactor = collateralFactor
	this.IfCollateral = ifCollateral
	return nil
}

type Borrow struct {
	Icon             string
	Name             string
	BorrowDollar     uint64
	BorrowBalance    uint64
	Apy              uint64
	Limit            uint64
	CollateralFactor uint64
}

func (this *Borrow) Serialization(sink *common.ZeroCopySink) {
	sink.WriteString(this.Icon)
	sink.WriteString(this.Name)
	sink.WriteUint64(this.BorrowDollar)
	sink.WriteUint64(this.BorrowBalance)
	sink.WriteUint64(this.Apy)
	sink.WriteUint64(this.Limit)
	sink.WriteUint64(this.CollateralFactor)
}

func (this *Borrow) Deserialization(source *common.ZeroCopySource) error {
	icon, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("icon deserialization error eof")
	}
	name, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("name deserialization error eof")
	}
	borrowDollar, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("borrowDollar deserialization error eof")
	}
	borrowBalance, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("borrowBalance deserialization error eof")
	}
	apy, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("apy deserialization error eof")
	}
	limit, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("limit deserialization error eof")
	}
	collateralFactor, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("collateralFactor deserialization error eof")
	}
	this.Icon = icon
	this.Name = name
	this.BorrowDollar = borrowDollar
	this.BorrowBalance = borrowBalance
	this.Apy = apy
	this.Limit = limit
	this.CollateralFactor = collateralFactor
	return nil
}

type Insurance struct {
	Icon             string
	Name             string
	InsuranceDollar  uint64
	InsuranceBalance uint64
	Apy              uint64
	CollateralFactor uint64
}

func (this *Insurance) Serialization(sink *common.ZeroCopySink) {
	sink.WriteString(this.Icon)
	sink.WriteString(this.Name)
	sink.WriteUint64(this.InsuranceDollar)
	sink.WriteUint64(this.InsuranceBalance)
	sink.WriteUint64(this.Apy)
	sink.WriteUint64(this.CollateralFactor)
}

func (this *Insurance) Deserialization(source *common.ZeroCopySource) error {
	icon, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("icon deserialization error eof")
	}
	name, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("name deserialization error eof")
	}
	insuranceDollar, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("insuranceDollar deserialization error eof")
	}
	insuranceBalance, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("insuranceBalance deserialization error eof")
	}
	apy, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("apy deserialization error eof")
	}
	collateralFactor, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("collateralFactor deserialization error eof")
	}
	this.Icon = icon
	this.Name = name
	this.InsuranceDollar = insuranceDollar
	this.InsuranceBalance = insuranceBalance
	this.Apy = apy
	this.CollateralFactor = collateralFactor
	return nil
}

type UserMarket struct {
	Icon             string
	Name             string
	SupplyApy        uint64
	BorrowApy        uint64
	BorrowLiquidity  uint64
	InsuranceApy     uint64
	InsuranceAmount  uint64
	CollateralFactor uint64
}

func (this *UserMarket) Serialization(sink *common.ZeroCopySink) {
	sink.WriteString(this.Icon)
	sink.WriteString(this.Name)
	sink.WriteUint64(this.SupplyApy)
	sink.WriteUint64(this.BorrowApy)
	sink.WriteUint64(this.BorrowLiquidity)
	sink.WriteUint64(this.InsuranceApy)
	sink.WriteUint64(this.InsuranceAmount)
	sink.WriteUint64(this.CollateralFactor)
}

func (this *UserMarket) Deserialization(source *common.ZeroCopySource) error {
	icon, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("icon deserialization error eof")
	}
	name, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("name deserialization error eof")
	}
	supplyApy, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("supplyApy deserialization error eof")
	}
	borrowApy, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("borrowApy deserialization error eof")
	}
	borrowLiquidity, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("borrowLiquidity deserialization error eof")
	}
	insuranceApy, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("insuranceApy deserialization error eof")
	}
	insuranceAmount, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("insuranceAmount deserialization error eof")
	}
	collateralFactor, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("collateralFactor deserialization error eof")
	}
	this.Icon = icon
	this.Name = name
	this.SupplyApy = supplyApy
	this.BorrowApy = borrowApy
	this.BorrowLiquidity = borrowLiquidity
	this.InsuranceApy = insuranceApy
	this.InsuranceAmount = insuranceAmount
	this.CollateralFactor = collateralFactor
	return nil
}

type FlashPoolAllMarket struct {
	FlashPoolAllMarket []*Market
}

type Market struct {
	Icon               string
	Name               string `gorm:"primary_key"`
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
