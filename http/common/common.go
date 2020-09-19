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

	ASSETPRICE     = "/api/v1/assetprice"
	ASSETPRICELIST = "/api/v1/assetpricelist"
	CLAIMWING      = "/api/v1/claimwing"
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

	ACTION_ASSETPRICE     = "assetprice"
	ACTION_ASSETPRICELIST = "assetpricelist"
	ACTION_CLAIMWING      = "claimwing"
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
	Price string
}

type AssetPriceListRequest struct {
	Id        string
	AssetList []string
}

type AssetPriceListResponse struct {
	Id        string
	PriceList []string
}

type GovBannerOverview struct {
	Remain20 string
	Remain80 string
}

type GovBanner struct {
	Daily       string
	Distributed string
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
	PerDay          string
	SupplyAmount    string
	BorrowAmount    string
	InsuranceAmount string
	Total           string
}

type FlashPoolBanner struct {
	Today string
	Share string
	Total string
}

type FlashPoolDetail struct {
	TotalSupply       string
	TotalSupplyRate   string
	SupplyMarketRank  []*MarketFund
	SupplyVolumeDaily string

	TotalBorrow       string
	TotalBorrowRate   string
	BorrowMarketRank  []*MarketFund
	BorrowVolumeDaily string

	TotalInsurance       string
	TotalInsuranceRate   string
	InsuranceMarketRank  []*MarketFund
	InsuranceVolumeDaily string
}

type MarketFund struct {
	Icon string
	Name string
	Fund string
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
	BorrowLimit string
	NetApy      string
	WingAccrued string

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
	Icon                  string
	Name                  string
	SupplyBalance         string
	Apy                   string
	CollateralFactor      string
	SupplyDistribution    string
	BorrowDistribution    string
	InsuranceDistribution string
	IfCollateral          bool
}

func (this *Supply) Serialization(sink *common.ZeroCopySink) {
	sink.WriteString(this.Icon)
	sink.WriteString(this.Name)
	sink.WriteString(this.SupplyBalance)
	sink.WriteString(this.Apy)
	sink.WriteString(this.CollateralFactor)
	sink.WriteString(this.SupplyDistribution)
	sink.WriteString(this.BorrowDistribution)
	sink.WriteString(this.InsuranceDistribution)
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
	supplyBalance, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("supplyBalance deserialization error eof")
	}
	apy, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("apy deserialization error eof")
	}
	collateralFactor, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("collateralFactor deserialization error eof")
	}
	supplyDistribution, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("supplyDistribution deserialization error eof")
	}
	borrowDistribution, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("borrowDistribution deserialization error eof")
	}
	insuranceDistribution, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("insuranceDistribution deserialization error eof")
	}
	ifCollateral, irregular, eof := source.NextBool()
	if irregular || eof {
		return fmt.Errorf("ifCollateral deserialization error eof")
	}
	this.Icon = icon
	this.Name = name
	this.SupplyBalance = supplyBalance
	this.Apy = apy
	this.CollateralFactor = collateralFactor
	this.SupplyDistribution = supplyDistribution
	this.BorrowDistribution = borrowDistribution
	this.InsuranceDistribution = insuranceDistribution
	this.IfCollateral = ifCollateral
	return nil
}

type Borrow struct {
	Icon                  string
	Name                  string
	BorrowBalance         string
	Apy                   string
	Limit                 string
	CollateralFactor      string
	SupplyDistribution    string
	BorrowDistribution    string
	InsuranceDistribution string
}

func (this *Borrow) Serialization(sink *common.ZeroCopySink) {
	sink.WriteString(this.Icon)
	sink.WriteString(this.Name)
	sink.WriteString(this.BorrowBalance)
	sink.WriteString(this.Apy)
	sink.WriteString(this.Limit)
	sink.WriteString(this.CollateralFactor)
	sink.WriteString(this.SupplyDistribution)
	sink.WriteString(this.BorrowDistribution)
	sink.WriteString(this.InsuranceDistribution)
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
	borrowBalance, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("borrowBalance deserialization error eof")
	}
	apy, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("apy deserialization error eof")
	}
	limit, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("limit deserialization error eof")
	}
	collateralFactor, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("collateralFactor deserialization error eof")
	}
	supplyDistribution, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("supplyDistribution deserialization error eof")
	}
	borrowDistribution, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("borrowDistribution deserialization error eof")
	}
	insuranceDistribution, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("insuranceDistribution deserialization error eof")
	}
	this.Icon = icon
	this.Name = name
	this.BorrowBalance = borrowBalance
	this.Apy = apy
	this.Limit = limit
	this.CollateralFactor = collateralFactor
	this.SupplyDistribution = supplyDistribution
	this.BorrowDistribution = borrowDistribution
	this.InsuranceDistribution = insuranceDistribution
	return nil
}

type Insurance struct {
	Icon                  string
	Name                  string
	InsuranceBalance      string
	Apy                   string
	CollateralFactor      string
	SupplyDistribution    string
	BorrowDistribution    string
	InsuranceDistribution string
}

func (this *Insurance) Serialization(sink *common.ZeroCopySink) {
	sink.WriteString(this.Icon)
	sink.WriteString(this.Name)
	sink.WriteString(this.InsuranceBalance)
	sink.WriteString(this.Apy)
	sink.WriteString(this.CollateralFactor)
	sink.WriteString(this.SupplyDistribution)
	sink.WriteString(this.BorrowDistribution)
	sink.WriteString(this.InsuranceDistribution)
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
	insuranceBalance, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("insuranceBalance deserialization error eof")
	}
	apy, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("apy deserialization error eof")
	}
	collateralFactor, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("collateralFactor deserialization error eof")
	}
	supplyDistribution, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("supplyDistribution deserialization error eof")
	}
	borrowDistribution, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("borrowDistribution deserialization error eof")
	}
	insuranceDistribution, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("insuranceDistribution deserialization error eof")
	}
	this.Icon = icon
	this.Name = name
	this.InsuranceBalance = insuranceBalance
	this.Apy = apy
	this.CollateralFactor = collateralFactor
	this.SupplyDistribution = supplyDistribution
	this.BorrowDistribution = borrowDistribution
	this.InsuranceDistribution = insuranceDistribution
	return nil
}

type UserMarket struct {
	Icon                  string
	Name                  string
	SupplyApy             string
	BorrowApy             string
	BorrowLiquidity       string
	InsuranceApy          string
	InsuranceAmount       string
	CollateralFactor      string
	SupplyDistribution    string
	BorrowDistribution    string
	InsuranceDistribution string
	IfCollateral          bool
}

func (this *UserMarket) Serialization(sink *common.ZeroCopySink) {
	sink.WriteString(this.Icon)
	sink.WriteString(this.Name)
	sink.WriteString(this.SupplyApy)
	sink.WriteString(this.BorrowApy)
	sink.WriteString(this.BorrowLiquidity)
	sink.WriteString(this.InsuranceApy)
	sink.WriteString(this.InsuranceAmount)
	sink.WriteString(this.CollateralFactor)
	sink.WriteString(this.SupplyDistribution)
	sink.WriteString(this.BorrowDistribution)
	sink.WriteString(this.InsuranceDistribution)
	sink.WriteBool(this.IfCollateral)
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
	supplyApy, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("icon deserialization error eof")
	}
	borrowApy, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("icon deserialization error eof")
	}
	borrowLiquidity, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("icon deserialization error eof")
	}
	insuranceApy, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("icon deserialization error eof")
	}
	insuranceAmount, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("icon deserialization error eof")
	}
	collateralFactor, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("icon deserialization error eof")
	}
	supplyDistribution, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("icon deserialization error eof")
	}
	borrowDistribution, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("icon deserialization error eof")
	}
	insuranceDistribution, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("icon deserialization error eof")
	}
	ifCollateral, irregular, eof := source.NextBool()
	if irregular || eof {
		return fmt.Errorf("ifCollateral deserialization error eof")
	}
	this.Icon = icon
	this.Name = name
	this.SupplyApy = supplyApy
	this.BorrowApy = borrowApy
	this.BorrowLiquidity = borrowLiquidity
	this.InsuranceApy = insuranceApy
	this.InsuranceAmount = insuranceAmount
	this.CollateralFactor = collateralFactor
	this.SupplyDistribution = supplyDistribution
	this.BorrowDistribution = borrowDistribution
	this.InsuranceDistribution = insuranceDistribution
	this.IfCollateral = ifCollateral
	return nil
}

type FlashPoolAllMarket struct {
	FlashPoolAllMarket []*Market
}

type Market struct {
	Icon                  string
	Name                  string `gorm:"primary_key"`
	TotalSupply           string
	TotalSupplyRate       string
	SupplyApy             string
	SupplyApyRate         string
	TotalBorrow           string
	TotalBorrowRate       string
	BorrowApy             string
	BorrowApyRate         string
	TotalInsurance        string
	TotalInsuranceRate    string
	InsuranceApy          string
	InsuranceApyRate      string
	CollateralFactor      string
	SupplyDistribution    string
	BorrowDistribution    string
	InsuranceDistribution string
}

type ClaimWingRequest struct {
	Id      string
	Address string
}

type ClaimWingResponse struct {
	Id      string
	Address string
	Amount  string
}
