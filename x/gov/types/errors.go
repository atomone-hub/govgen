package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/gov module sentinel errors
var (
	ErrUnknownProposal         = sdkerrors.Register(ModuleName, 12, "unknown proposal")
	ErrInactiveProposal        = sdkerrors.Register(ModuleName, 13, "inactive proposal")
	ErrAlreadyActiveProposal   = sdkerrors.Register(ModuleName, 14, "proposal already active")
	ErrInvalidProposalContent  = sdkerrors.Register(ModuleName, 15, "invalid proposal content")
	ErrInvalidProposalType     = sdkerrors.Register(ModuleName, 16, "invalid proposal type")
	ErrInvalidVote             = sdkerrors.Register(ModuleName, 17, "invalid vote option")
	ErrInvalidGenesis          = sdkerrors.Register(ModuleName, 18, "invalid genesis state")
	ErrNoProposalHandlerExists = sdkerrors.Register(ModuleName, 19, "no handler exists for proposal type")
)
