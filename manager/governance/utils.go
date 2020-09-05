package governance

import (
	"fmt"
	"github.com/ontio/ontology/common"
)

// get 20% wing distribute count
func (this *GovernanceManager) get20WingCount() (uint64, error) {
	r, err := this.sdk.GetStorage(this.wingAddress, []byte("Community"))
	if err != nil {
		return 0, fmt.Errorf("get20WingCount, this.sdk.GetStorage error: %s", err)
	}
	return common.BigIntFromNeoBytes(r).Uint64(), nil
}

// get wing total supply
func (this *GovernanceManager) getWingTotalSupply() (uint64, error) {
	r, err := this.sdk.GetStorage(this.wingAddress, []byte("TotalSupply"))
	if err != nil {
		return 0, fmt.Errorf("getWingTotalSupply, this.sdk.GetStorage error: %s", err)
	}
	return common.BigIntFromNeoBytes(r).Uint64(), nil
}
