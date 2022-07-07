package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	"fmt"

	evmtypes "github.com/Canto-Network/ethermint-v2/x/evm/types"
)

var (
	//go:embed compiled_contracts/caller.json
	calleeJSON []byte

	// ERC20BurnableContract is the compiled ERC20Burnable contract
	CalleeContract evmtypes.CompiledContract
)

func init() {
	err := json.Unmarshal(calleeJSON, &CalleeContract)
	if err != nil {
		// panic(err)
		fmt.Println("ERROR HERE")
	}
}
