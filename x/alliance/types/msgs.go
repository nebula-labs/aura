package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &AllianceRequest{}

// msg types
const (
	TypeMsgAllianceRequest = "alliance"
)

func NewMsgAllianceRequest(sender string, amount sdk.Coin) *AllianceRequest {
	return &AllianceRequest{
		Sender:     sender,
		HostAmount: amount,
	}
}

func (msg *AllianceRequest) Route() string {
	return RouterKey
}

func (msg *AllianceRequest) Type() string {
	return TypeMsgAllianceRequest
}

func (msg *AllianceRequest) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *AllianceRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *AllianceRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}
