package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	govgenhelpers "github.com/atomone-hub/govgen/v1/app/helpers"
	"github.com/atomone-hub/govgen/v1/x/gov/types"
)

func TestTallyNoOneVotes(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	createValidators(t, ctx, app, []int64{5, 5, 5})

	tp := govgenhelpers.TestTextProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)

	require.False(t, passes)
	require.True(t, burnDeposits)
	require.True(t, tallyResults.Equals(types.EmptyTallyResult()))
}

func TestTallyNoQuorum(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	createValidators(t, ctx, app, []int64{2, 5, 0})
	addrs := govgenhelpers.AddTestAddrs(app, ctx, 1, sdk.NewInt(10000000))

	tp := govgenhelpers.TestTextProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	err = app.GovKeeper.AddVote(ctx, proposalID, addrs[0], types.NewNonSplitVoteOption(types.OptionYes))
	require.Nil(t, err)

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, _ := app.GovKeeper.Tally(ctx, proposal)
	require.False(t, passes)
	require.True(t, burnDeposits)
}

func TestTallyAllYes(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	_, valAddrs := createValidators(t, ctx, app, []int64{1, 1, 1})
	addrs := govgenhelpers.AddTestAddrs(app, ctx, 3, sdk.NewInt(50000000))

	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 5)
	val1, found := app.StakingKeeper.GetValidator(ctx, valAddrs[0])
	require.True(t, found)

	for _, addr := range addrs {
		_, err := app.StakingKeeper.Delegate(ctx, addr, delTokens, stakingtypes.Unbonded, val1, true)
		require.NoError(t, err)
	}

	tp := govgenhelpers.TestTextProposal

	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[0], types.NewNonSplitVoteOption(types.OptionYes)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[1], types.NewNonSplitVoteOption(types.OptionYes)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[2], types.NewNonSplitVoteOption(types.OptionYes)))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)

	require.True(t, passes)
	require.False(t, burnDeposits)
	require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
}

func TestTally51No(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	power := []int64{5, 6, 0}
	_, valAddrs := createValidators(t, ctx, app, []int64{1, 1, 1})
	addrs := govgenhelpers.AddTestAddrs(app, ctx, 3, sdk.NewInt(50000000))

	val1, found := app.StakingKeeper.GetValidator(ctx, valAddrs[0])
	require.True(t, found)

	for i := range addrs {
		delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, power[i])
		_, err := app.StakingKeeper.Delegate(ctx, addrs[i], delTokens, stakingtypes.Unbonded, val1, true)
		require.NoError(t, err)
	}

	tp := govgenhelpers.TestTextProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[0], types.NewNonSplitVoteOption(types.OptionYes)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[1], types.NewNonSplitVoteOption(types.OptionNo)))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, _ := app.GovKeeper.Tally(ctx, proposal)

	require.False(t, passes)
	require.False(t, burnDeposits)
}

func TestTally51Yes(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	power := []int64{5, 6, 0}
	_, valAddrs := createValidators(t, ctx, app, []int64{1, 1, 1})
	addrs := govgenhelpers.AddTestAddrs(app, ctx, 3, sdk.NewInt(50000000))

	val1, found := app.StakingKeeper.GetValidator(ctx, valAddrs[0])
	require.True(t, found)

	for i := range addrs {
		delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, power[i])
		_, err := app.StakingKeeper.Delegate(ctx, addrs[i], delTokens, stakingtypes.Unbonded, val1, true)
		require.NoError(t, err)
	}

	tp := govgenhelpers.TestTextProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[0], types.NewNonSplitVoteOption(types.OptionNo)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[1], types.NewNonSplitVoteOption(types.OptionYes)))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)

	require.True(t, passes)
	require.False(t, burnDeposits)
	require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
}

func TestTallyVetoed(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	power := []int64{6, 6, 7}
	_, valAddrs := createValidators(t, ctx, app, []int64{1, 1, 1})
	addrs := govgenhelpers.AddTestAddrs(app, ctx, 3, sdk.NewInt(50000000))

	val1, found := app.StakingKeeper.GetValidator(ctx, valAddrs[0])
	require.True(t, found)

	for i := range addrs {
		delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, power[i])
		_, err := app.StakingKeeper.Delegate(ctx, addrs[i], delTokens, stakingtypes.Unbonded, val1, true)
		require.NoError(t, err)
	}

	tp := govgenhelpers.TestTextProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[0], types.NewNonSplitVoteOption(types.OptionYes)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[1], types.NewNonSplitVoteOption(types.OptionYes)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[2], types.NewNonSplitVoteOption(types.OptionNoWithVeto)))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)

	require.False(t, passes)
	require.True(t, burnDeposits)
	require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
}

func TestTallyAbstainPasses(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	power := []int64{6, 6, 7}
	_, valAddrs := createValidators(t, ctx, app, []int64{1, 1, 1})
	addrs := govgenhelpers.AddTestAddrs(app, ctx, 3, sdk.NewInt(50000000))

	val1, found := app.StakingKeeper.GetValidator(ctx, valAddrs[0])
	require.True(t, found)

	for i := range addrs {
		delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, power[i])
		_, err := app.StakingKeeper.Delegate(ctx, addrs[i], delTokens, stakingtypes.Unbonded, val1, true)
		require.NoError(t, err)
	}

	tp := govgenhelpers.TestTextProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[0], types.NewNonSplitVoteOption(types.OptionAbstain)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[1], types.NewNonSplitVoteOption(types.OptionNo)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[2], types.NewNonSplitVoteOption(types.OptionYes)))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)

	require.True(t, passes)
	require.False(t, burnDeposits)
	require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
}

func TestTallyAbstainFails(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	power := []int64{6, 6, 7}
	_, valAddrs := createValidators(t, ctx, app, []int64{1, 1, 1})
	addrs := govgenhelpers.AddTestAddrs(app, ctx, 3, sdk.NewInt(50000000))

	val1, found := app.StakingKeeper.GetValidator(ctx, valAddrs[0])
	require.True(t, found)

	for i := range addrs {
		delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, power[i])
		_, err := app.StakingKeeper.Delegate(ctx, addrs[i], delTokens, stakingtypes.Unbonded, val1, true)
		require.NoError(t, err)
	}

	tp := govgenhelpers.TestTextProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[0], types.NewNonSplitVoteOption(types.OptionAbstain)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[1], types.NewNonSplitVoteOption(types.OptionYes)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[2], types.NewNonSplitVoteOption(types.OptionNo)))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)

	require.False(t, passes)
	require.False(t, burnDeposits)
	require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
}

func TestTallyNonVoter(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	power := []int64{6, 6, 7}
	_, valAddrs := createValidators(t, ctx, app, []int64{1, 1, 1})
	addrs := govgenhelpers.AddTestAddrs(app, ctx, 3, sdk.NewInt(50000000))

	val1, found := app.StakingKeeper.GetValidator(ctx, valAddrs[0])
	require.True(t, found)

	for i := range addrs {
		delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, power[i])
		_, err := app.StakingKeeper.Delegate(ctx, addrs[i], delTokens, stakingtypes.Unbonded, val1, true)
		require.NoError(t, err)
	}

	tp := govgenhelpers.TestTextProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[0], types.NewNonSplitVoteOption(types.OptionYes)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[1], types.NewNonSplitVoteOption(types.OptionNo)))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)

	require.False(t, passes)
	require.False(t, burnDeposits)
	require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
}

/*
// NOTE: these tests are disabled as voting for validators is not available and therefore
// all features around voting and delegations will not work

func TestTallyDelgatorOverride(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addrs, valAddrs := createValidators(t, ctx, app, []int64{5, 6, 7})

	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 30)
	val1, found := app.StakingKeeper.GetValidator(ctx, valAddrs[0])
	require.True(t, found)

	_, err := app.StakingKeeper.Delegate(ctx, addrs[4], delTokens, stakingtypes.Unbonded, val1, true)
	require.NoError(t, err)

	_ = staking.EndBlocker(ctx, app.StakingKeeper)

	tp := TestProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[1], types.NewNonSplitVoteOption(types.OptionYes)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[2], types.NewNonSplitVoteOption(types.OptionYes)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[3], types.NewNonSplitVoteOption(types.OptionYes)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[4], types.NewNonSplitVoteOption(types.OptionNo)))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)

	require.False(t, passes)
	require.False(t, burnDeposits)
	require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
}

func TestTallyDelgatorInherit(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addrs, vals := createValidators(t, ctx, app, []int64{5, 6, 7})

	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 30)
	val3, found := app.StakingKeeper.GetValidator(ctx, vals[2])
	require.True(t, found)

	_, err := app.StakingKeeper.Delegate(ctx, addrs[3], delTokens, stakingtypes.Unbonded, val3, true)
	require.NoError(t, err)

	_ = staking.EndBlocker(ctx, app.StakingKeeper)

	tp := TestProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[0], types.NewNonSplitVoteOption(types.OptionNo)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[1], types.NewNonSplitVoteOption(types.OptionNo)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[2], types.NewNonSplitVoteOption(types.OptionYes)))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)

	require.True(t, passes)
	require.False(t, burnDeposits)
	require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
}

func TestTallyDelgatorMultipleOverride(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addrs, vals := createValidators(t, ctx, app, []int64{5, 6, 7})

	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	val1, found := app.StakingKeeper.GetValidator(ctx, vals[0])
	require.True(t, found)
	val2, found := app.StakingKeeper.GetValidator(ctx, vals[1])
	require.True(t, found)

	_, err := app.StakingKeeper.Delegate(ctx, addrs[3], delTokens, stakingtypes.Unbonded, val1, true)
	require.NoError(t, err)
	_, err = app.StakingKeeper.Delegate(ctx, addrs[3], delTokens, stakingtypes.Unbonded, val2, true)
	require.NoError(t, err)

	_ = staking.EndBlocker(ctx, app.StakingKeeper)

	tp := TestProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[0], types.NewNonSplitVoteOption(types.OptionYes)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[1], types.NewNonSplitVoteOption(types.OptionYes)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[2], types.NewNonSplitVoteOption(types.OptionYes)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[3], types.NewNonSplitVoteOption(types.OptionNo)))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)

	require.False(t, passes)
	require.False(t, burnDeposits)
	require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
}

func TestTallyDelgatorMultipleInherit(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	createValidators(t, ctx, app, []int64{25, 6, 7})

	addrs, vals := createValidators(t, ctx, app, []int64{5, 6, 7})

	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	val2, found := app.StakingKeeper.GetValidator(ctx, vals[1])
	require.True(t, found)
	val3, found := app.StakingKeeper.GetValidator(ctx, vals[2])
	require.True(t, found)

	_, err := app.StakingKeeper.Delegate(ctx, addrs[3], delTokens, stakingtypes.Unbonded, val2, true)
	require.NoError(t, err)
	_, err = app.StakingKeeper.Delegate(ctx, addrs[3], delTokens, stakingtypes.Unbonded, val3, true)
	require.NoError(t, err)

	_ = staking.EndBlocker(ctx, app.StakingKeeper)

	tp := TestProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[0], types.NewNonSplitVoteOption(types.OptionYes)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[1], types.NewNonSplitVoteOption(types.OptionNo)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[2], types.NewNonSplitVoteOption(types.OptionNo)))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)

	require.False(t, passes)
	require.False(t, burnDeposits)
	require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
}

func TestTallyJailedValidator(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addrs, valAddrs := createValidators(t, ctx, app, []int64{25, 6, 7})

	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	val2, found := app.StakingKeeper.GetValidator(ctx, valAddrs[1])
	require.True(t, found)
	val3, found := app.StakingKeeper.GetValidator(ctx, valAddrs[2])
	require.True(t, found)

	_, err := app.StakingKeeper.Delegate(ctx, addrs[3], delTokens, stakingtypes.Unbonded, val2, true)
	require.NoError(t, err)
	_, err = app.StakingKeeper.Delegate(ctx, addrs[3], delTokens, stakingtypes.Unbonded, val3, true)
	require.NoError(t, err)

	_ = staking.EndBlocker(ctx, app.StakingKeeper)

	consAddr, err := val2.GetConsAddr()
	require.NoError(t, err)
	app.StakingKeeper.Jail(ctx, sdk.ConsAddress(consAddr.Bytes()))

	tp := TestProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[0], types.NewNonSplitVoteOption(types.OptionYes)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[1], types.NewNonSplitVoteOption(types.OptionNo)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[2], types.NewNonSplitVoteOption(types.OptionNo)))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)

	require.True(t, passes)
	require.False(t, burnDeposits)
	require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
}

func TestTallyValidatorMultipleDelegations(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addrs, valAddrs := createValidators(t, ctx, app, []int64{10, 10, 10})

	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	val2, found := app.StakingKeeper.GetValidator(ctx, valAddrs[1])
	require.True(t, found)

	_, err := app.StakingKeeper.Delegate(ctx, addrs[0], delTokens, stakingtypes.Unbonded, val2, true)
	require.NoError(t, err)

	tp := TestProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[0], types.NewNonSplitVoteOption(types.OptionYes)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[1], types.NewNonSplitVoteOption(types.OptionNo)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, addrs[2], types.NewNonSplitVoteOption(types.OptionYes)))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)

	require.True(t, passes)
	require.False(t, burnDeposits)

	expectedYes := app.StakingKeeper.TokensFromConsensusPower(ctx, 30)
	expectedAbstain := app.StakingKeeper.TokensFromConsensusPower(ctx, 0)
	expectedNo := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	expectedNoWithVeto := app.StakingKeeper.TokensFromConsensusPower(ctx, 0)
	expectedTallyResult := types.NewTallyResult(expectedYes, expectedAbstain, expectedNo, expectedNoWithVeto)

	require.True(t, tallyResults.Equals(expectedTallyResult))
}

*/
