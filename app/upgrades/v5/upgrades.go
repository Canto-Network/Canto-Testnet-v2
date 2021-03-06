package v5

import (
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	ibctransferkeeper "github.com/cosmos/ibc-go/v3/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"

	feemarketv010types "github.com/Canto-Network/ethermint-v2/x/feemarket/migrations/v010/types"
	feemarketv011 "github.com/Canto-Network/ethermint-v2/x/feemarket/migrations/v011"
	feemarkettypes "github.com/Canto-Network/ethermint-v2/x/feemarket/types"

	"github.com/Canto-Network/Canto-Testnet-v2/v1/types"
)

// TestnetDenomMetadata defines the metadata for the tcanto denom on testnet
var TestnetDenomMetadata = banktypes.Metadata{
	Description: "The native EVM, governance and staking token of the canto testnet",
	DenomUnits: []*banktypes.DenomUnit{
		{
			Denom:    "atcanto",
			Exponent: 0,
			Aliases:  []string{"attotcanto"},
		},
		{
			Denom:    "tcanto",
			Exponent: 18,
		},
	},
	Base:    "atcanto",
	Display: "tcanto",
	Name:    "Testnet canto",
	Symbol:  "tcanto",
}

// CreateUpgradeHandler creates an SDK upgrade handler for v5
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	bk bankkeeper.Keeper,
	sk stakingkeeper.Keeper,
	pk paramskeeper.Keeper,
	tk ibctransferkeeper.Keeper,
	xk slashingkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)

		// modify fee market parameter defaults through global
		feemarkettypes.DefaultMinGasPrice = MainnetMinGasPrices
		feemarkettypes.DefaultMinGasMultiplier = MainnetMinGasMultiplier

		// Refs:
		// - https://docs.cosmos.network/master/building-modules/upgrade.html#registering-migrations
		// - https://docs.cosmos.network/master/migrations/chain-upgrade-guide-044.html#chain-upgrade

		// define the denom metadata for the testnet
		if types.IsTestnet(ctx.ChainID()) {
			logger.Debug("setting testnet client denomination metadata...")
			bk.SetDenomMetaData(ctx, TestnetDenomMetadata)
		}

		logger.Debug("updating Tendermint consensus params...")
		UpdateConsensusParams(ctx, sk, pk)

		logger.Debug("updating IBC transfer denom traces...")
		UpdateIBCDenomTraces(ctx, tk)

		//claims module removed

		// define from versions of the modules that have a new consensus version

		// migrate fee market module, other modules are left as-is to
		// avoid running InitGenesis.
		vm[feemarkettypes.ModuleName] = 2

		// Leave modules are as-is to avoid running InitGenesis.
		logger.Debug("running migration for fee market module (EIP-1559)...")
		return mm.RunMigrations(ctx, configurator, vm)
	}
}

// MigrateGenesis migrates exported state from v4 to v5 genesis state.
// It performs a no-op if the migration errors.
func MigrateGenesis(appState genutiltypes.AppMap, clientCtx client.Context) genutiltypes.AppMap {
	// Migrate x/feemarket.
	if appState[feemarkettypes.ModuleName] == nil {
		return appState
	}

	// unmarshal relative source genesis application state
	var oldFeeMarketState feemarketv010types.GenesisState
	if err := clientCtx.Codec.UnmarshalJSON(appState[feemarkettypes.ModuleName], &oldFeeMarketState); err != nil {
		return appState
	}

	// delete deprecated x/feemarket genesis state
	delete(appState, feemarkettypes.ModuleName)

	// Migrate relative source genesis application state and marshal it into
	// the respective key.
	newFeeMarketState := feemarketv011.MigrateJSON(oldFeeMarketState)

	feeMarketBz, err := clientCtx.Codec.MarshalJSON(&newFeeMarketState)
	if err != nil {
		return appState
	}

	appState[feemarkettypes.ModuleName] = feeMarketBz

	return appState
}

// UpdateConsensusParams updates the Tendermint Consensus Evidence params (MaxAgeDuration and
// MaxAgeNumBlocks) to match the unbonding period and use the expected avg block time based on the
// node configuration.
func UpdateConsensusParams(ctx sdk.Context, sk stakingkeeper.Keeper, pk paramskeeper.Keeper) {
	subspace, found := pk.GetSubspace(baseapp.Paramspace)
	if !found {
		return
	}

	var evidenceParams tmproto.EvidenceParams
	subspace.GetIfExists(ctx, baseapp.ParamStoreKeyEvidenceParams, &evidenceParams)

	// safety check: no-op if the evidence params is empty (shouldn't happen)
	if evidenceParams.Equal(tmproto.EvidenceParams{}) {
		return
	}

	stakingParams := sk.GetParams(ctx)
	evidenceParams.MaxAgeDuration = stakingParams.UnbondingTime

	maxAgeNumBlocks := sdk.NewInt(int64(evidenceParams.MaxAgeDuration)).QuoRaw(int64(AvgBlockTime))
	evidenceParams.MaxAgeNumBlocks = maxAgeNumBlocks.Int64()
	subspace.Set(ctx, baseapp.ParamStoreKeyEvidenceParams, evidenceParams)
}

// UpdateIBCDenomTraces iterates over current traces to check if any of them are incorrectly formed
// and corrects the trace information.
// See https://github.com/cosmos/ibc-go/blob/main/docs/migrations/support-denoms-with-slashes.md for context.
func UpdateIBCDenomTraces(ctx sdk.Context, transferKeeper ibctransferkeeper.Keeper) {
	// list of traces that must replace the old traces in store
	var newTraces []ibctransfertypes.DenomTrace
	transferKeeper.IterateDenomTraces(ctx, func(dt ibctransfertypes.DenomTrace) bool {
		// check if the new way of splitting FullDenom
		// into Trace and BaseDenom passes validation and
		// is the same as the current DenomTrace.
		// If it isn't then store the new DenomTrace in the list of new traces.
		newTrace := ibctransfertypes.ParseDenomTrace(dt.GetFullDenomPath())
		if err := newTrace.Validate(); err == nil && !equalTraces(newTrace, dt) {
			newTraces = append(newTraces, newTrace)
		}

		return false
	})

	// replace the outdated traces with the new trace information
	for _, nt := range newTraces {
		transferKeeper.SetDenomTrace(ctx, nt)
	}
}

func equalTraces(dtA, dtB ibctransfertypes.DenomTrace) bool {
	return dtA.BaseDenom == dtB.BaseDenom && dtA.Path == dtB.Path
}
