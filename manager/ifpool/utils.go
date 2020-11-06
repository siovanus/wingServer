package ifpool

import (
	"fmt"
	"github.com/siovanus/wingServer/utils"
	"math/big"
)

func (this *IFPoolManager) assetStoredPrice(asset string) (*big.Int, error) {
	price, err := this.store.LoadPrice(asset)
	if err != nil {
		return nil, fmt.Errorf("AssetStoredPrice, this.store.LoadPrice error: %s", err)
	}
	return utils.ToIntByPrecise(price.Price, this.cfg.TokenDecimal["oracle"]), nil
}
