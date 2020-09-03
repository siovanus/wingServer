package flashpool

type MarketDistribution struct {
	MarketDistribution []*Distribution
}

type PoolDistribution struct {
	PoolDistribution []*Distribution
}

type Distribution struct {
	Icon         string
	Name         string
	PerDay       uint64
	SupplyApy    uint64
	BorrowApy    uint64
	InsuranceApy uint64
}

func (this *FlashPoolManager) marketDistribution() (*MarketDistribution, error) {
	distribution1 := &Distribution{
		Icon:         "http://106.75.209.209/icon/eth_icon.svg",
		Name:         "oEth",
		PerDay:       234,
		SupplyApy:    6783,
		BorrowApy:    8325,
		InsuranceApy: 9517,
	}
	distribution2 := &Distribution{
		Icon:         "http://106.75.209.209/icon/asset_dai_icon.svg",
		Name:         "oDai",
		PerDay:       345,
		SupplyApy:    1574,
		BorrowApy:    4576,
		InsuranceApy: 3842,
	}
	return &MarketDistribution{MarketDistribution: []*Distribution{distribution1, distribution2}}, nil
}

func (this *FlashPoolManager) poolDistribution() (*PoolDistribution, error) {
	distribution1 := &Distribution{
		Icon:         "http://106.75.209.209/icon/flash_icon.svg",
		Name:         "Flash",
		PerDay:       231252,
		SupplyApy:    2532,
		BorrowApy:    4547,
		InsuranceApy: 1231,
	}
	distribution2 := &Distribution{
		Icon:         "http://106.75.209.209/icon/if_icon.svg",
		Name:         "IF",
		PerDay:       1431241,
		SupplyApy:    1214,
		BorrowApy:    2525,
		InsuranceApy: 7742,
	}
	return &PoolDistribution{PoolDistribution: []*Distribution{distribution1, distribution2}}, nil
}
