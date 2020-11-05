package ifpool

import (
	"fmt"
)

func (this *IFPoolManager) getPrice(name string) (string, error) {
	price, err := this.store.LoadPrice(name)
	if err != nil {
		return "", fmt.Errorf("getPrice, this.store.LoadPrice: %s", err)
	}
	return price.Price, nil
}
