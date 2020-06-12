package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidOrder         = sdkerrors.Register(ModuleName, 2, "should be ordered")
	ErrUnknownPacketData    = sdkerrors.Register(ModuleName, 3, "unknown packet data")
	ErrAccountAlreadyExist  = sdkerrors.Register(ModuleName, 4, "account already exist")
	ErrUnsupportedChainType = sdkerrors.Register(ModuleName, 5, "unsupported chain type")
	ErrInvalidOutgoingData  = sdkerrors.Register(ModuleName, 6, "invalid outgoing data")
	ErrInvalidRoute         = sdkerrors.Register(ModuleName, 7, "invalid route")
)
