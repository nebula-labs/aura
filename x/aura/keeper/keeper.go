package keeper

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/aura-nw/aura/x/aura/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/keeper"
	ibctransferkeeper "github.com/cosmos/ibc-go/v3/modules/apps/transfer/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"
)

type (
	Keeper struct {
		cdc                 codec.BinaryCodec
		storeKey            sdk.StoreKey
		memKey              sdk.StoreKey
		paramSpace          paramtypes.Subspace
		scopedKeeper        capabilitykeeper.ScopedKeeper
		IbcKeeper           ibckeeper.Keeper
		IcaControllerKeeper icacontrollerkeeper.Keeper
		transferKeeper      ibctransferkeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	paramSpace paramtypes.Subspace,
	scopedKeeper capabilitykeeper.ScopedKeeper,
	icaControllerKeeper icacontrollerkeeper.Keeper,
	transferKeeper ibctransferkeeper.Keeper,
	ibcKeeper ibckeeper.Keeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		memKey:              memKey,
		paramSpace:          paramSpace,
		scopedKeeper:        scopedKeeper,
		IcaControllerKeeper: icaControllerKeeper,
		IbcKeeper:           ibcKeeper,
		transferKeeper:      transferKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetParams returns the total set of aura parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of aura parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetMaxSupply return max supply of aura coin
func (k Keeper) GetMaxSupply(ctx sdk.Context) string {
	params := k.GetParams(ctx)
	return params.MaxSupply
}

// GetExcludeCirculatingAddr return list exclude address do not calculator for circulating
func (k Keeper) GetExcludeCirculatingAddr(ctx sdk.Context) []sdk.AccAddress {
	params := k.GetParams(ctx)
	excludeAddr := make([]sdk.AccAddress, 0, len(params.ExcludeCirculatingAddr))
	for _, addrBech32 := range params.ExcludeCirculatingAddr {
		addr, err := sdk.AccAddressFromBech32(addrBech32)
		if err != nil {
			panic(err)
		}
		excludeAddr = append(excludeAddr, addr)
	}

	return excludeAddr
}
