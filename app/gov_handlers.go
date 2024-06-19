package govgen

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	paramsrest "github.com/cosmos/cosmos-sdk/x/params/client/rest"
	paramscutils "github.com/cosmos/cosmos-sdk/x/params/client/utils"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	upgradecli "github.com/cosmos/cosmos-sdk/x/upgrade/client/cli"
	upgraderest "github.com/cosmos/cosmos-sdk/x/upgrade/client/rest"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	govclient "github.com/atomone-hub/govgen/x/gov/client"
	"github.com/atomone-hub/govgen/x/gov/client/cli"
	govtypes "github.com/atomone-hub/govgen/x/gov/types"
)

var (
	paramsChangeProposalHandler = govclient.NewProposalHandler(
		newSubmitParamChangeProposalTxCmd,
		govclient.WrapPropposalRESTHandler(paramsrest.ProposalRESTHandler),
	)
	upgradeProposalHandler = govclient.NewProposalHandler(
		newCmdSubmitUpgradeProposal,
		govclient.WrapPropposalRESTHandler(upgraderest.ProposalRESTHandler),
	)
	cancelUpgradeProposalHandler = govclient.NewProposalHandler(
		newCmdSubmitCancelUpgradeProposal,
		govclient.WrapPropposalRESTHandler(upgraderest.ProposalRESTHandler),
	)
)

// newSubmitParamChangeProposalTxCmd returns a CLI command handler for creating
// a parameter change proposal governance transaction.
//
// NOTE: copy of x/params/client.newSubmitParamChangeProposalTxCmd() except
// that it creates a govgen.gov.MsgSubmitProposal instead of a
// cosmos.gov.MsgSubmitProposal.
func newSubmitParamChangeProposalTxCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "param-change [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a parameter change proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a parameter proposal along with an initial deposit.
The proposal details must be supplied via a JSON file. For values that contains
objects, only non-empty fields will be updated.

IMPORTANT: Currently parameter changes are evaluated but not validated, so it is
very important that any "value" change is valid (ie. correct type and within bounds)
for its respective parameter, eg. "MaxValidators" should be an integer and not a decimal.

Proper vetting of a parameter change proposal should prevent this from happening
(no deposits should occur during the governance process), but it should be noted
regardless.

Example:
$ %s tx gov submit-proposal param-change <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Staking Param Change",
  "description": "Update max validators",
  "changes": [
    {
      "subspace": "staking",
      "key": "MaxValidators",
      "value": 105
    }
  ],
  "deposit": "1000stake"
}
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			proposal, err := paramscutils.ParseParamChangeProposalJSON(clientCtx.LegacyAmino, args[0])
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			content := paramsproposal.NewParameterChangeProposal(
				proposal.Title, proposal.Description, proposal.Changes.ToParamChanges(),
			)

			deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
}

// newCmdSubmitUpgradeProposal implements a command handler for submitting a software upgrade proposal transaction.
//
// NOTE: copy of x/upgrade/client.NewCmdSubmitUpgradeProposal() except
// that it creates a govgen.gov.MsgSubmitProposal instead of a
// cosmos.gov.MsgSubmitProposal.
func newCmdSubmitUpgradeProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "software-upgrade [name] (--upgrade-height [height]) (--upgrade-info [info]) [flags]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a software upgrade proposal",
		Long: "Submit a software upgrade along with an initial deposit.\n" +
			"Please specify a unique name and height for the upgrade to take effect.\n" +
			"You may include info to reference a binary download link, in a format compatible with: https://github.com/cosmos/cosmos-sdk/tree/master/cosmovisor",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			name := args[0]
			content, err := parseArgsToContent(cmd, name)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	cmd.Flags().Int64(upgradecli.FlagUpgradeHeight, 0, "The height at which the upgrade must happen")
	cmd.Flags().String(upgradecli.FlagUpgradeInfo, "", "Optional info for the planned upgrade such as commit hash, etc.")

	return cmd
}

// newCmdSubmitCancelUpgradeProposal implements a command handler for submitting a software upgrade cancel proposal transaction.
//
// NOTE: copy of x/upgrade/client.NewCmdSubmitCancelUpgradeProposal() except
// that it creates a govgen.gov.MsgSubmitProposal instead of a
// cosmos.gov.MsgSubmitProposal.
func newCmdSubmitCancelUpgradeProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-software-upgrade [flags]",
		Args:  cobra.ExactArgs(0),
		Short: "Cancel the current software upgrade proposal",
		Long:  "Cancel a software upgrade along with an initial deposit.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			content := upgradetypes.NewCancelSoftwareUpgradeProposal(title, description)

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	cmd.MarkFlagRequired(cli.FlagTitle)       //nolint:errcheck
	cmd.MarkFlagRequired(cli.FlagDescription) //nolint:errcheck

	return cmd
}

func parseArgsToContent(cmd *cobra.Command, name string) (govtypes.Content, error) {
	title, err := cmd.Flags().GetString(cli.FlagTitle)
	if err != nil {
		return nil, err
	}

	description, err := cmd.Flags().GetString(cli.FlagDescription)
	if err != nil {
		return nil, err
	}

	height, err := cmd.Flags().GetInt64(upgradecli.FlagUpgradeHeight)
	if err != nil {
		return nil, err
	}

	info, err := cmd.Flags().GetString(upgradecli.FlagUpgradeInfo)
	if err != nil {
		return nil, err
	}

	plan := upgradetypes.Plan{Name: name, Height: height, Info: info}
	content := upgradetypes.NewSoftwareUpgradeProposal(title, description, plan)
	return content, nil
}
