package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/feeabstraction module sentinel errors
var (
	ErrMarshalFailure = sdkerrors.Register(ModuleName, 4, "unable to marshal data structure")
)
