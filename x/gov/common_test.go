package gov_test

import (
	"bytes"
	"log"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/atomone-hub/govgen/v1/x/gov/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

var (
	valTokens                = sdk.TokensFromConsensusPower(42, sdk.DefaultPowerReduction)
	TestTextProposal         = types.NewTextProposal("Test", "description")
	TestParamsChangeProposal = paramsproposal.NewParameterChangeProposal(
		"Test", "description", []paramsproposal.ParamChange{}, // TODO test such proposal?
	)
	TestUpgradeProposal = upgradetypes.NewSoftwareUpgradeProposal(
		"Test", "description", upgradetypes.Plan{
			Name:   "plan",
			Height: 42,
		}, // TODO test such proposal?
	)
	TestCancelUpgradeProposal = upgradetypes.NewCancelSoftwareUpgradeProposal(
		"Test", "description", // TODO test such proposal?
	)
	TestDescription     = stakingtypes.NewDescription("T", "E", "S", "T", "Z")
	TestCommissionRates = stakingtypes.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
)

// SortAddresses - Sorts Addresses
func SortAddresses(addrs []sdk.AccAddress) {
	byteAddrs := make([][]byte, len(addrs))

	for i, addr := range addrs {
		byteAddrs[i] = addr.Bytes()
	}

	SortByteArrays(byteAddrs)

	for i, byteAddr := range byteAddrs {
		addrs[i] = byteAddr
	}
}

// implement `Interface` in sort package.
type sortByteArrays [][]byte

func (b sortByteArrays) Len() int {
	return len(b)
}

func (b sortByteArrays) Less(i, j int) bool {
	// bytes package already implements Comparable for []byte.
	switch bytes.Compare(b[i], b[j]) {
	case -1:
		return true
	case 0, 1:
		return false
	default:
		log.Panic("not fail-able with `bytes.Comparable` bounded [-1, 1].")
		return false
	}
}

func (b sortByteArrays) Swap(i, j int) {
	b[j], b[i] = b[i], b[j]
}

// SortByteArrays - sorts the provided byte array
func SortByteArrays(src [][]byte) [][]byte {
	sorted := sortByteArrays(src)
	sort.Sort(sorted)
	return sorted
}

const contextKeyBadProposal = "contextKeyBadProposal"

var pubkeys = []cryptotypes.PubKey{
	ed25519.GenPrivKey().PubKey(),
	ed25519.GenPrivKey().PubKey(),
	ed25519.GenPrivKey().PubKey(),
}

func createValidators(t *testing.T, stakingHandler sdk.Handler, ctx sdk.Context, addrs []sdk.ValAddress, powerAmt []int64) { //nolint: thelper
	require.True(t, len(addrs) <= len(pubkeys), "Not enough pubkeys specified at top of file.")

	for i := 0; i < len(addrs); i++ {
		valTokens := sdk.TokensFromConsensusPower(powerAmt[i], sdk.DefaultPowerReduction)
		valCreateMsg, err := stakingtypes.NewMsgCreateValidator(
			addrs[i], pubkeys[i], sdk.NewCoin(sdk.DefaultBondDenom, valTokens),
			TestDescription, TestCommissionRates, sdk.OneInt(),
		)
		require.NoError(t, err)
		handleAndCheck(t, stakingHandler, ctx, valCreateMsg)
	}
}

func handleAndCheck(t *testing.T, h sdk.Handler, ctx sdk.Context, msg sdk.Msg) { //nolint: thelper
	res, err := h(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)
}
