package flashpool

type MarketDistribution struct {
	MarketDistribution []*Distribution `json:"market_distribution"`
}

type PoolDistribution struct {
	PoolDistribution []*Distribution `json:"pool_distribution"`
}

type FlashPoolBanner struct {
	Today uint64
	Share uint64
	Total uint64
}

type Distribution struct {
	Icon         string
	Name         string
	PerDay       uint64
	SupplyApy    uint64
	BorrowApy    uint64
	InsuranceApy uint64
	Total        uint64
}

func (this *FlashPoolManager) marketDistribution() (*MarketDistribution, error) {
	distribution1 := &Distribution{
		Icon:         "http://106.75.209.209/icon/eth_icon.svg",
		Name:         "oEth",
		PerDay:       234,
		SupplyApy:    6783,
		BorrowApy:    8325,
		InsuranceApy: 9517,
		Total:        121234,
	}
	distribution2 := &Distribution{
		Icon:         "http://106.75.209.209/icon/asset_dai_icon.svg",
		Name:         "oDai",
		PerDay:       345,
		SupplyApy:    1574,
		BorrowApy:    4576,
		InsuranceApy: 3842,
		Total:        25252,
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
		Total:        28364,
	}
	distribution2 := &Distribution{
		Icon:         "http://106.75.209.209/icon/if_icon.svg",
		Name:         "IF",
		PerDay:       1431241,
		SupplyApy:    1214,
		BorrowApy:    2525,
		InsuranceApy: 7742,
		Total:        72526,
	}
	return &PoolDistribution{PoolDistribution: []*Distribution{distribution1, distribution2}}, nil
}

func (this *FlashPoolManager) flashPoolBanner() (*FlashPoolBanner, error) {
	return &FlashPoolBanner{
		Today: 8676,
		Share: 7644,
		Total: 3452636,
	}, nil
}
