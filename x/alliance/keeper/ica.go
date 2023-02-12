package keeper

import (
	"fmt"
	"strings"

	"github.com/aura-nw/aura/x/alliance/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
)

// SubmitTxs submits an ICA transaction containing multiple messages
// Will submit tx from fee account
func (k Keeper) SubmitTxs(
	ctx sdk.Context,
	connectionId string,
	msgs []sdk.Msg,
	timeoutTimestamp uint64,
) (uint64, error) {
	chainId, err := k.GetChainID(ctx, connectionId)
	if err != nil {
		return 0, err
	}
	owner := types.GetICAAccountOwner(chainId)
	portID, err := icatypes.NewControllerPortID(owner)
	if err != nil {
		return 0, err
	}

	k.Logger(ctx).Info(LogWithHostZone(chainId, "  Submitting ICA Tx on %s, %s with TTL: %d", portID, connectionId, timeoutTimestamp))
	for _, msg := range msgs {
		k.Logger(ctx).Info(LogWithHostZone(chainId, "    Msg: %+v", msg))
	}

	channelID, found := k.IcaControllerKeeper.GetActiveChannelID(ctx, connectionId, portID)
	if !found {
		return 0, sdkerrors.Wrapf(icatypes.ErrActiveChannelNotFound, "failed to retrieve active channel for port %s", portID)
	}

	chanCap, found := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(portID, channelID))
	if !found {
		return 0, sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	data, err := icatypes.SerializeCosmosTx(k.cdc, msgs)
	if err != nil {
		return 0, err
	}

	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
	}

	sequence, err := k.IcaControllerKeeper.SendTx(ctx, chanCap, connectionId, portID, packetData, timeoutTimestamp)
	if err != nil {
		return 0, err
	}

	return sequence, nil
}

// Alliance will invoke a ICA{IbcSend} to Orai
func (k Keeper) MsgICATransfer(ctx sdk.Context, coin sdk.Coin, oraiChannel string, receiver string) ([]sdk.Msg, error) {
	sourcePort := ibctransfertypes.PortID
	hostAddress := k.GetICAAddress(ctx)
	// DEBUG: currently set to NonNativeFeeCollectorName for debugging
	icaTimeoutNanos, err := k.GetTtl(ctx)
	if err != nil {
		return nil, err
	}

	msgs := []sdk.Msg{
		&ibctransfertypes.MsgTransfer{
			SourcePort:       sourcePort,
			SourceChannel:    oraiChannel,
			Token:            coin,
			Sender:           hostAddress,
			Receiver:         receiver,
			TimeoutTimestamp: icaTimeoutNanos,
		},
	}

	return msgs, nil
}

func LogWithHostZone(chainId string, s string, a ...any) string {
	msg := fmt.Sprintf(s, a...)
	return fmt.Sprintf("|   %-13s |  %s", strings.ToUpper(chainId), msg)
}

func (k Keeper) Alliance(ctx sdk.Context, req *types.AllianceRequest) error {
	// message creation
	msgs, err := k.MsgICATransfer(ctx, req.HostAmount, types.ORAI_AURA_CHANNEL_ID, req.Sender)
	if err != nil {
		return err
	}

	icaTimeoutNanos, err := k.GetTtl(ctx)
	if err != nil {
		return err
	}

	_, err = k.SubmitTxs(ctx, types.ORAI_AURA_CONNECTION_ID, msgs, icaTimeoutNanos)
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to submit txs")
	}

	return nil
}
