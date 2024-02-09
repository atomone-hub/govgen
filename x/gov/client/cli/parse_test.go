package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/testutil"
)

func TestParseSubmitProposalFlags(t *testing.T) {
	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "Test Proposal",
  "description": "My awesome proposal",
  "type": "Text",
  "deposit": "1000test"
}
`)

	badJSON := testutil.WriteToNewTempFile(t, "bad json")
	fs := NewCmdSubmitProposal().Flags()

	// nonexistent json
	fs.Set(FlagProposal, "fileDoesNotExist") //nolint: errcheck
	_, err := parseSubmitProposalFlags(fs)
	require.Error(t, err)

	// invalid json
	fs.Set(FlagProposal, badJSON.Name()) //nolint: errcheck
	_, err = parseSubmitProposalFlags(fs)
	require.Error(t, err)

	// ok json
	fs.Set(FlagProposal, okJSON.Name()) //nolint: errcheck
	proposal1, err := parseSubmitProposalFlags(fs)
	require.Nil(t, err, "unexpected error")
	require.Equal(t, "Test Proposal", proposal1.Title)
	require.Equal(t, "My awesome proposal", proposal1.Description)
	require.Equal(t, "Text", proposal1.Type)
	require.Equal(t, "1000test", proposal1.Deposit)

	// flags that can't be used with --proposal
	for _, incompatibleFlag := range ProposalFlags {
		fs.Set(incompatibleFlag, "some value") //nolint: errcheck
		_, err := parseSubmitProposalFlags(fs)
		require.Error(t, err)
		fs.Set(incompatibleFlag, "") //nolint: errcheck
	}

	// no --proposal, only flags
	fs.Set(FlagProposal, "")                       //nolint: errcheck
	fs.Set(FlagTitle, proposal1.Title)             //nolint: errcheck
	fs.Set(FlagDescription, proposal1.Description) //nolint: errcheck
	fs.Set(FlagProposalType, proposal1.Type)       //nolint: errcheck
	fs.Set(FlagDeposit, proposal1.Deposit)         //nolint: errcheck
	proposal2, err := parseSubmitProposalFlags(fs)

	require.Nil(t, err, "unexpected error")
	require.Equal(t, proposal1.Title, proposal2.Title)
	require.Equal(t, proposal1.Description, proposal2.Description)
	require.Equal(t, proposal1.Type, proposal2.Type)
	require.Equal(t, proposal1.Deposit, proposal2.Deposit)

	err = okJSON.Close()
	require.Nil(t, err, "unexpected error")
	err = badJSON.Close()
	require.Nil(t, err, "unexpected error")
}
