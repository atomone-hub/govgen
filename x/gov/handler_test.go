package gov_test

import (
	"strings"
	"testing"

	"github.com/atomone-hub/govgen/v1/x/gov"
	"github.com/atomone-hub/govgen/v1/x/gov/keeper"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestInvalidMsg(t *testing.T) {
	k := keeper.Keeper{}
	h := gov.NewHandler(k)

	res, err := h(sdk.NewContext(nil, tmproto.Header{}, false, nil), testdata.NewTestMsg())
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, strings.Contains(err.Error(), "unrecognized gov message type"))
}
