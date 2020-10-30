package ifpool

import (
	sdk "github.com/ontio/ontology-go-sdk"
	ocommon "github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/config"
	"github.com/siovanus/wingServer/http/common"
	"github.com/siovanus/wingServer/store"
)

type IFPoolManager struct {
	cfg             *config.Config
	contractAddress ocommon.Address
	sdk             *sdk.OntologySdk
	store           *store.Client
}

func NewIFPoolManager(contractAddress ocommon.Address, sdk *sdk.OntologySdk, store *store.Client,
	cfg *config.Config) *IFPoolManager {
	manager := &IFPoolManager{
		cfg:             cfg,
		contractAddress: contractAddress,
		sdk:             sdk,
		store:           store,
	}

	return manager
}

func (this *IFPoolManager) IFPoolInfo(account string) (*common.IFPoolInfo, error) {
	IFPoolInfo := &common.IFPoolInfo{
		Total: "462000.13241",
		Cap:   "500000",
		IFAssetList: []*common.IFAssetList{
			{
				Name:                 "pUSDT",
				Icon:                 "https://app.ont.io/wing/pusdt.svg",
				Price:                "1",
				TotalSupply:          "1435.456747",
				SupplyInterestPerDay: "0.252536",
				SupplyWingAPY:        "0.2436",
				UtilizationRate:      "0.7714",
				MaximumLTV:           "0.125",
				TotalBorrowed:        "9656.25225",
				BorrowInterestPerDay: "0.00024",
				BorrowWingAPY:        "2367",
				Liquidity:            "256.4564",
				BorrowCap:            "30",
				TotalInsurance:       "1769.3536",
				InsuranceWingAPY:     "2.349",
			},
			{
				Name:                 "pSUSD",
				Icon:                 "https://app.ont.io/wing/psusd.svg",
				Price:                "1.001",
				TotalSupply:          "121435.4564747",
				SupplyInterestPerDay: "0.0253463",
				SupplyWingAPY:        "0.242536",
				UtilizationRate:      "0.786714",
				MaximumLTV:           "0.1425",
				TotalBorrowed:        "965256.25225",
				BorrowInterestPerDay: "0.000224",
				BorrowWingAPY:        "235267",
				Liquidity:            "2536.4564",
				BorrowCap:            "500",
				TotalInsurance:       "1876969.34536",
				InsuranceWingAPY:     "25.34649",
			},
			{
				Name:                 "pDAI",
				Icon:                 "https://app.ont.io/wing/oDAI.svg",
				Price:                "1.02",
				TotalSupply:          "252.3636",
				SupplyInterestPerDay: "0.0035636",
				SupplyWingAPY:        "0.24236536536",
				UtilizationRate:      "0.3536",
				MaximumLTV:           "0.3536",
				TotalBorrowed:        "35252.647",
				BorrowInterestPerDay: "0.0002724",
				BorrowWingAPY:        "79707",
				Liquidity:            "890.3242",
				BorrowCap:            "242",
				TotalInsurance:       "34535.34575536",
				InsuranceWingAPY:     "575.686",
			},
		},
	}
	if account != "" {
		IFPoolInfo.UserIFInfo = &common.UserIFInfo{
			TotalSupplyDollar:    "23526.3647",
			SupplyWingEarned:     "25.3647",
			TotalBorrowDollar:    "25364.485",
			BorrowWingEarned:     "789.36536",
			BorrowInterestPerDay: "0.837636",
			TotalInsuranceDollar: "96747.474747",
			InsuranceWingEarned:  "36796.366",
			Composition: []*common.Composition{
				{
					Name:                  "pUSDT",
					Icon:                  "https://app.ont.io/wing/pusdt.svg",
					SupplyBalance:         "42526.3636",
					SupplyWingEarned:      "242.2525",
					BorrowWingEarned:      "235.3677",
					LastBorrowTimestamp:   "1604026092000",
					InsuranceBalance:      "141536.47",
					InsuranceWingEarned:   "14.25265",
					CollateralName:        "pSUSD",
					CollateralIcon:        "https://app.ont.io/wing/psusd.svg",
					CollateralBalance:     "242.236363",
					BorrowUnpaidPrincipal: "2452.3636",
					BorrowInterestBalance: "242.242",
				},
				{
					Name:                  "pSUSD",
					Icon:                  "https://app.ont.io/wing/psusd.svg",
					SupplyBalance:         "235546.4647",
					SupplyWingEarned:      "22.225",
					BorrowWingEarned:      "25.377",
					LastBorrowTimestamp:   "1604026082000",
					InsuranceBalance:      "14136.47",
					InsuranceWingEarned:   "1.2265",
					CollateralName:        "pDAI",
					CollateralIcon:        "https://app.ont.io/wing/oDAI.svg",
					CollateralBalance:     "24.6868",
					BorrowUnpaidPrincipal: "0",
					BorrowInterestBalance: "0",
				},
				{
					Name:                  "pDAI",
					Icon:                  "https://app.ont.io/wing/oDAI.svg",
					SupplyBalance:         "235544566.464467",
					SupplyWingEarned:      "224.225",
					BorrowWingEarned:      "265.37697",
					LastBorrowTimestamp:   "1604026082000",
					InsuranceBalance:      "696.47",
					InsuranceWingEarned:   "1141.2265",
					CollateralName:        "pUSDT",
					CollateralIcon:        "https://app.ont.io/wing/pusdt.svg",
					CollateralBalance:     "2242.57578",
					BorrowUnpaidPrincipal: "0",
					BorrowInterestBalance: "0",
				},
			},
		}
	}
	return IFPoolInfo, nil
}

func (this *IFPoolManager) IFHistory(asset, operation string, start, end, pageNo, pageSize uint64) (*common.IFHistoryResponse, error) {
	return &common.IFHistoryResponse{
		MaxPageNum: 1,
		PageItems: []*common.IFHistory{
			{
				Name:      "pUSDT",
				Icon:      "https://app.ont.io/wing/pusdt.svg",
				Operation: "Supply",
				Timestamp: 1604026092000,
				Balance:   "32532.58",
				Dollar:    "23526.464",
				Address:   "Af3Etnp5ffrXR3swrCx9f7KuvChYLgqsTZ",
			},
			{
				Name:      "pDAI",
				Icon:      "https://app.ont.io/wing/oDAI.svg",
				Operation: "Borrow",
				Timestamp: 1604026092000,
				Balance:   "6968.58",
				Dollar:    "797.464",
				Address:   "AR36E5jLdWDKW3Yg51qDFWPGKSLvfPhbqS",
			},
		},
	}, nil
}
