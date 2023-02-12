package keeper

import (
	"context"

	"github.com/aura-nw/aura/x/alliance/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (server msgServer) Alliance(ctx context.Context, req *types.AllianceRequest) (*types.AllianceResponse, error) {

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := server.Keeper.Alliance(sdkCtx, req)
	if err != nil {
		return nil, err
	}

	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventAlliance,
			sdk.NewAttribute(sdk.AttributeKeySender, req.Sender),
			sdk.NewAttribute(sdk.AttributeKeyAmount, req.HostAmount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	})

	return &types.AllianceResponse{}, nil
}
