package governance

type GovBanner struct {
	Remain20 uint64
	Remain80 uint64
}

type PoolBanner struct {
	Daily       uint64
	Distributed uint64
}

func (this *GovernanceManager) govBannerOverview() (*GovBanner, error) {
	return &GovBanner{
		Remain20: 1450000,
		Remain80: 7650000,
	}, nil
}

func (this *GovernanceManager) govBanner() (*PoolBanner, error) {
	return &PoolBanner{
		Daily:       141513,
		Distributed: 1221312,
	}, nil
}
