package keeper

import (
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"

	"github.com/aura-nw/aura/x/alliance/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	k.Logger(ctx).Info("Discovering ibc tokens from orai")

	// due to ibc denom hash, we can only accept denom directly from Osmosis
	// is there a way to get all assets on a channel
	k.transferKeeper.IterateDenomTraces(ctx, func(denomTrace ibctransfertypes.DenomTrace) bool {
		k.Logger(ctx).Info(fmt.Sprintf("Found token pair: (%s, %s) on channel %s",
			denomTrace.GetBaseDenom(), denomTrace.IBCDenom(), types.ORAI_AURA_CHANNEL_ID))

		// if an ibc denom exists, skip
		if k.HasHostDenomTrack(ctx, denomTrace.GetBaseDenom()) {
			return true
		}

		// if found out that denom belongs to orai channel_id, register denom trace
		if strings.Contains(denomTrace.GetPath(), types.ORAI_AURA_CHANNEL_ID) {
			k.Logger(ctx).Info("Registering token pair")
			k.SetDenomTrack(ctx, denomTrace.GetBaseDenom(), denomTrace.IBCDenom())
		}

		// putting ica registration behind ibc transfer ensures that connection on both parties are OPEN.
		account := types.GetICAAccountOwner(types.HOST_ZONE_CHAIN_ID)
		_, exist := k.IcaControllerKeeper.GetInterchainAccountAddress(ctx, types.ORAI_AURA_CONNECTION_ID, account)
		if !exist {
			if err := k.IcaControllerKeeper.RegisterInterchainAccount(ctx, types.ORAI_AURA_CONNECTION_ID, account); err != nil {
				k.Logger(ctx).Error(fmt.Sprintf("unable to register fee account, err: %s", err.Error()))
				return true
			}
		}

		return true
	})
}
