package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	govgenhelpers "github.com/atomone-hub/govgen/app/helpers"
	"github.com/atomone-hub/govgen/x/gov/types"
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

func TestTallyOnlyValidatorsAllYes(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addrs, _ := createValidators(t, ctx, app, []int64{1, 1, 1})

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

func TestTallyOnlyValidators51No(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	valAccAddrs, _ := createValidators(t, ctx, app, []int64{5, 6, 0})

	tp := govgenhelpers.TestTextProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, valAccAddrs[0], types.NewNonSplitVoteOption(types.OptionYes)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, valAccAddrs[1], types.NewNonSplitVoteOption(types.OptionNo)))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, _ := app.GovKeeper.Tally(ctx, proposal)

	require.False(t, passes)
	require.False(t, burnDeposits)
}

func TestTallyOnlyValidators51Yes(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	valAccAddrs, _ := createValidators(t, ctx, app, []int64{5, 6, 0})

	tp := govgenhelpers.TestTextProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, valAccAddrs[0], types.NewNonSplitVoteOption(types.OptionNo)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, valAccAddrs[1], types.NewNonSplitVoteOption(types.OptionYes)))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)

	require.True(t, passes)
	require.False(t, burnDeposits)
	require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
}

func TestTallyOnlyValidatorsVetoed(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	valAccAddrs, _ := createValidators(t, ctx, app, []int64{6, 6, 7})

	tp := govgenhelpers.TestTextProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, valAccAddrs[0], types.NewNonSplitVoteOption(types.OptionYes)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, valAccAddrs[1], types.NewNonSplitVoteOption(types.OptionYes)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, valAccAddrs[2], types.NewNonSplitVoteOption(types.OptionNoWithVeto)))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)

	require.False(t, passes)
	require.True(t, burnDeposits)
	require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
}

func TestTallyOnlyValidatorsAbstainPasses(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	valAccAddrs, _ := createValidators(t, ctx, app, []int64{6, 6, 7})

	tp := govgenhelpers.TestTextProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, valAccAddrs[0], types.NewNonSplitVoteOption(types.OptionAbstain)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, valAccAddrs[1], types.NewNonSplitVoteOption(types.OptionNo)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, valAccAddrs[2], types.NewNonSplitVoteOption(types.OptionYes)))

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

	valAccAddrs, _ := createValidators(t, ctx, app, []int64{6, 6, 7})

	tp := govgenhelpers.TestTextProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, valAccAddrs[0], types.NewNonSplitVoteOption(types.OptionAbstain)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, valAccAddrs[1], types.NewNonSplitVoteOption(types.OptionYes)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, valAccAddrs[2], types.NewNonSplitVoteOption(types.OptionNo)))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)

	require.False(t, passes)
	require.False(t, burnDeposits)
	require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
}

func TestTallyOnlyValidatorsNonVoter(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	valAccAddrs, _ := createValidators(t, ctx, app, []int64{6, 6, 7})

	tp := govgenhelpers.TestTextProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, valAccAddrs[0], types.NewNonSplitVoteOption(types.OptionYes)))
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, valAccAddrs[1], types.NewNonSplitVoteOption(types.OptionNo)))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)

	require.False(t, passes)
	require.False(t, burnDeposits)
	require.False(t, tallyResults.Equals(types.EmptyTallyResult()))
}

func TestTallyDelegatorOverride(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addrs, valAddrs := createValidators(t, ctx, app, []int64{5, 6, 7})

	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 30)
	val1, found := app.StakingKeeper.GetValidator(ctx, valAddrs[0])
	require.True(t, found)

	_, err := app.StakingKeeper.Delegate(ctx, addrs[4], delTokens, stakingtypes.Unbonded, val1, true)
	require.NoError(t, err)

	_ = staking.EndBlocker(ctx, app.StakingKeeper)

	tp := govgenhelpers.TestTextProposal
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

// As validators can only vote with their own stake, delegators don't inherit votes from validators
// so the proposal fails
func TestTallyDelegatorInherit(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	valPowers := []int64{5, 6, 7}
	addrs, vals := createValidators(t, ctx, app, valPowers)

	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 30)
	val3, found := app.StakingKeeper.GetValidator(ctx, vals[2])
	require.True(t, found)

	_, err := app.StakingKeeper.Delegate(ctx, addrs[3], delTokens, stakingtypes.Unbonded, val3, true)
	require.NoError(t, err)

	_ = staking.EndBlocker(ctx, app.StakingKeeper)

	tp := govgenhelpers.TestTextProposal
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

	require.False(t, passes)
	require.False(t, burnDeposits)
	valSelfDelegations := []sdk.Int{
		app.StakingKeeper.TokensFromConsensusPower(ctx, valPowers[0]),
		app.StakingKeeper.TokensFromConsensusPower(ctx, valPowers[1]),
		app.StakingKeeper.TokensFromConsensusPower(ctx, valPowers[2]),
	}
	require.Equal(t, tallyResults.String(), types.NewTallyResult(
		valSelfDelegations[2],
		sdk.ZeroInt(),
		valSelfDelegations[0].Add(valSelfDelegations[1]),
		sdk.ZeroInt(),
	).String())
}

func TestTallyDelegatorMultipleOverride(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	valPowers := []int64{5, 6, 7}
	addrs, vals := createValidators(t, ctx, app, valPowers)

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

	tp := govgenhelpers.TestTextProposal
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
	valSelfDelegations := []sdk.Int{
		app.StakingKeeper.TokensFromConsensusPower(ctx, valPowers[0]),
		app.StakingKeeper.TokensFromConsensusPower(ctx, valPowers[1]),
		app.StakingKeeper.TokensFromConsensusPower(ctx, valPowers[2]),
	}
	require.Equal(t, tallyResults.String(), types.NewTallyResult(
		valSelfDelegations[0].Add(valSelfDelegations[1]).Add(valSelfDelegations[2]),
		sdk.ZeroInt(),
		delTokens.Add(delTokens),
		sdk.ZeroInt(),
	).String())
}

// As validators can only vote with their own stake, delegators don't inherit votes from validators
// so the proposal passes
func TestTallyDelegatorMultipleInherit(t *testing.T) {
	app := govgenhelpers.SetupNoValset(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	valPowers := []int64{5, 6, 7}
	addrs, vals := createValidators(t, ctx, app, valPowers)

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

	tp := govgenhelpers.TestTextProposal
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
	valSelfDelegations := []sdk.Int{
		app.StakingKeeper.TokensFromConsensusPower(ctx, valPowers[0]),
		app.StakingKeeper.TokensFromConsensusPower(ctx, valPowers[1]),
		app.StakingKeeper.TokensFromConsensusPower(ctx, valPowers[2]),
	}
	require.Equal(t, tallyResults.String(), types.NewTallyResult(
		valSelfDelegations[0],
		sdk.ZeroInt(),
		valSelfDelegations[1].Add(valSelfDelegations[2]),
		sdk.ZeroInt(),
	).String())
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

	tp := govgenhelpers.TestTextProposal
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

	tp := govgenhelpers.TestTextProposal
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
