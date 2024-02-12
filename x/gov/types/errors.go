package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/gov module sentinel errors
var (
	ErrUnknownProposal         = sdkerrors.Register(ModuleName, 20, "unknown proposal")
	ErrInactiveProposal        = sdkerrors.Register(ModuleName, 30, "inactive proposal")
	ErrAlreadyActiveProposal   = sdkerrors.Register(ModuleName, 40, "proposal already active")
	ErrInvalidProposalContent  = sdkerrors.Register(ModuleName, 50, "invalid proposal content")
	ErrInvalidProposalType     = sdkerrors.Register(ModuleName, 60, "invalid proposal type")
	ErrInvalidVote             = sdkerrors.Register(ModuleName, 70, "invalid vote option")
	ErrInvalidGenesis          = sdkerrors.Register(ModuleName, 80, "invalid genesis state")
	ErrNoProposalHandlerExists = sdkerrors.Register(ModuleName, 90, "no handler exists for proposal type")
	ErrValidatorCannotVote     = sdkerrors.Register(ModuleName, 100, "voting for validators is disabled")
)
