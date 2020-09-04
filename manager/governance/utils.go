package governance

type GovBanner struct {
	Remain20 uint64
	Remain80 uint64
}

func (this *GovernanceManager) govBanner() (*GovBanner, error) {
	return &GovBanner{
		Remain20: 1450000,
		Remain80: 7650000,
	}, nil
}
