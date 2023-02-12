package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/spf13/cast"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/aura-nw/aura/x/alliance/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	accountkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/keeper"
	ibctransferkeeper "github.com/cosmos/ibc-go/v3/modules/apps/transfer/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"
)

type (
	Keeper struct {
		cdc                 codec.BinaryCodec
		storeKey            sdk.StoreKey
		paramSpace          paramtypes.Subspace
		scopedKeeper        capabilitykeeper.ScopedKeeper
		IbcKeeper           ibckeeper.Keeper
		IcaControllerKeeper icacontrollerkeeper.Keeper
		transferKeeper      ibctransferkeeper.Keeper
		accountKeeper       accountkeeper.AccountKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	paramSpace paramtypes.Subspace,
	scopedKeeper capabilitykeeper.ScopedKeeper,
	icaControllerKeeper icacontrollerkeeper.Keeper,
	transferKeeper ibctransferkeeper.Keeper,
	ibcKeeper ibckeeper.Keeper,
	accountKeeper accountkeeper.AccountKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		paramSpace:          paramSpace,
		scopedKeeper:        scopedKeeper,
		IcaControllerKeeper: icaControllerKeeper,
		IbcKeeper:           ibcKeeper,
		transferKeeper:      transferKeeper,
		accountKeeper:       accountKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// get TTL for ICQ message
func (k Keeper) GetTtl(ctx sdk.Context) (uint64, error) {
	currentTime := ctx.BlockTime()

	// add 5 more mins to current time
	return cast.ToUint64E(currentTime.Add(time.Minute * 5).UnixNano())
}

// record for coins on osmosis to juno
func (k Keeper) SetDenomTrack(ctx sdk.Context, denomHost, denomController string) {
	storeOsmo := prefix.NewStore(ctx.KVStore(k.storeKey), types.StoreDenomHostTrack)
	storeOsmo.Set([]byte(denomHost), []byte(denomController))
	storeJuno := prefix.NewStore(ctx.KVStore(k.storeKey), types.StoreDenomControllerTrack)
	storeJuno.Set([]byte(denomController), []byte(denomHost))
}

func (k Keeper) HasHostDenomTrack(ctx sdk.Context, denomHost string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StoreDenomHostTrack)
	return store.Has([]byte(denomHost))
}
