package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	ethermint "github.com/Canto-Network/Canto-Testnet-v2/ethermint-v2/types"
)

func CalculateEpochMintProvision(
	params Params,
	period uint64,
	epochsPerPeriod int64,
	bondedRatio sdk.Dec,
) sdk.Dec {

	minInflation := params.ExponentialCalculation.MinInflation //minInflation
	maxInflation := params.ExponentialCalculation.MaxInflation
	bondedTarget := params.ExponentialCalculation.BondingTarget
	adjustSpeed := params.ExponentialCalculation.AdjustSpeed

	//how to get the current inflation?
	bondDiff := bondedTarget.Sub(bondedRatio)
	periodProvision := bondDiff.Mul(adjustSpeed).Mul()

	//return calculated inflation in terms of periods per epoch
	if minInflation.GT(periodProvision) {
		return minInflation.Quo(sdk.NewDec(epochsPerPeriod))
	}

	if periodProvision.GT(maxInflation) {
		return maxInflation.Quo(sdk.NewDec(epochsPerPeriod))
	}

	// return this value as a mantissa
	epochProvision := periodProvision.Quo(sdk.NewDec(epochsPerPeriod))

	return epochProvision.Mul(ethermint.PowerReduction.ToDec())
}
