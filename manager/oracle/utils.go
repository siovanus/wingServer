package oracle

import "fmt"

func (this *OracleManager) assetPrice(asset string) (uint64, error) {
	preExecResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(this.contractAddress,
		"getUnderlyingPrice", []interface{}{asset})
	if err != nil {
		return 0, fmt.Errorf("assetPrice, this.sdk.WasmVM.PreExecInvokeWasmVMContract error: %s", err)
	}
	r, err := preExecResult.Result.ToInteger()
	if err != nil {
		return 0, fmt.Errorf("assetPrice, preExecResult.Result.ToInteger error: %s", err)
	}
	return r.Uint64(), nil
}
