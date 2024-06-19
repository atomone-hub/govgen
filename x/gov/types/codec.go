package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// RegisterLegacyAminoCodec registers all the necessary types and interfaces for the
// governance module.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*Content)(nil), nil)
	cdc.RegisterConcrete(&MsgSubmitProposal{}, "govgen/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(&MsgDeposit{}, "govgen/MsgDeposit", nil)
	cdc.RegisterConcrete(&MsgVote{}, "govgen/MsgVote", nil)
	cdc.RegisterConcrete(&MsgVoteWeighted{}, "govgen/MsgVoteWeighted", nil)
	cdc.RegisterConcrete(&TextProposal{}, "govgen/TextProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSubmitProposal{},
		&MsgVote{},
		&MsgVoteWeighted{},
		&MsgDeposit{},
	)
	registry.RegisterInterface(
		"govgen.gov.v1beta1.Content",
		(*Content)(nil),
		&TextProposal{},
	)

	// Register proposal types (this is actually done in related modules, but
	// since we are using an other gov module, we need to do it manually).
	registry.RegisterImplementations(
		(*Content)(nil),
		&paramsproposal.ParameterChangeProposal{},
	)
	registry.RegisterImplementations(
		(*Content)(nil),
		&upgradetypes.SoftwareUpgradeProposal{},
	)
	registry.RegisterImplementations(
		(*Content)(nil),
		&upgradetypes.CancelSoftwareUpgradeProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// RegisterProposalTypeCodec registers an external proposal content type defined
// in another module for the internal ModuleCdc. This allows the MsgSubmitProposal
// to be correctly Amino encoded and decoded.
//
// NOTE: This should only be used for applications that are still using a concrete
// Amino codec for serialization.
func RegisterProposalTypeCodec(o interface{}, name string) {
	amino.RegisterConcrete(o, name, nil)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/gov module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/gov and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)

	// Register proposal types (this is actually done in related modules, but
	// since we are using an other gov module, we need to do it manually).
	RegisterProposalType(paramsproposal.ProposalTypeChange)
	RegisterProposalTypeCodec(&paramsproposal.ParameterChangeProposal{}, "cosmos-sdk/ParameterChangeProposal")
	RegisterProposalType(upgradetypes.ProposalTypeSoftwareUpgrade)
	RegisterProposalTypeCodec(&upgradetypes.SoftwareUpgradeProposal{}, "cosmos-sdk/SoftwareUpgradeProposal")
	RegisterProposalType(upgradetypes.ProposalTypeCancelSoftwareUpgrade)
	RegisterProposalTypeCodec(&upgradetypes.CancelSoftwareUpgradeProposal{}, "cosmos-sdk/CancelSoftwareUpgradeProposal")
}
