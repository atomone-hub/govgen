package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/server"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/snapshots"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"

	govgen "github.com/atomone-hub/govgen/app"
	"github.com/atomone-hub/govgen/app/params"
)

// NewRootCmd creates a new root command for simd. It is called once in the
// main function.
func NewRootCmd() (*cobra.Command, params.EncodingConfig) {
	encodingConfig := govgen.MakeTestEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Codec).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithHomeDir(govgen.DefaultNodeHome).
		WithViper("")

	rootCmd := &cobra.Command{
		Use:   "govgend",
		Short: "Stargate GovGen App",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err = client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customTemplate, customGovGenConfig := initAppConfig()
			return server.InterceptConfigsPreRunHandler(cmd, customTemplate, customGovGenConfig)
		},
	}

	initRootCmd(rootCmd, encodingConfig)

	return rootCmd, encodingConfig
}

func initAppConfig() (string, interface{}) {
	srvCfg := serverconfig.DefaultConfig()
	srvCfg.StateSync.SnapshotInterval = 1000
	srvCfg.StateSync.SnapshotKeepRecent = 10

	return serverconfig.DefaultConfigTemplate, *srvCfg
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig params.EncodingConfig) {
	cfg := sdk.GetConfig()

	cfg.Seal()

	rootCmd.AddCommand(
		genutilcli.InitCmd(govgen.ModuleBasics, govgen.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, govgen.DefaultNodeHome),
		genutilcli.GenTxCmd(govgen.ModuleBasics, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, govgen.DefaultNodeHome),
		genutilcli.ValidateGenesisCmd(govgen.ModuleBasics),
		AddGenesisAccountCmd(govgen.DefaultNodeHome),
		tmcli.NewCompletionCmd(rootCmd, true),
		testnetCmd(govgen.ModuleBasics, banktypes.GenesisBalancesIterator{}),
		addDebugCommands(debug.Cmd()),
		config.Cmd(),
	)

	ac := appCreator{
		encCfg: encodingConfig,
	}
	server.AddCommands(rootCmd, govgen.DefaultNodeHome, ac.newApp, ac.appExport, addModuleInitFlags)

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		queryCommand(),
		txCommand(),
		keys.Commands(govgen.DefaultNodeHome),
	)

	rootCmd.AddCommand(server.RosettaCommand(encodingConfig.InterfaceRegistry, encodingConfig.Codec))
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetAccountCmd(),
		rpc.ValidatorCommand(),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
	)

	govgen.ModuleBasics.AddQueryCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
		getBytesToSignCommand(),
	)

	govgen.ModuleBasics.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func getBytesToSignCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bytes-to-sign [file]",
		Short: "Outputs the bytes to be signed for the input transaction",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var clientCtx client.Context

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			stdTx, err := authclient.ReadTxFromFile(clientCtx, args[0])
			if err != nil {
				return err
			}
			txFactory := tx.NewFactoryCLI(clientCtx, cmd.Flags())
			txCfg := clientCtx.TxConfig
			txBuilder, err := txCfg.WrapTxBuilder(stdTx)
			if err != nil {
				return err
			}

			from, _ := cmd.Flags().GetString(flags.FlagFrom)
			_, fromName, _, err := client.GetFromFields(clientCtx, txFactory.Keybase(), from)
			if err != nil {
				return fmt.Errorf("error getting account from keybase: %w", err)
			}
			signMode := txFactory.SignMode()
			if signMode == signing.SignMode_SIGN_MODE_UNSPECIFIED {
				// use the SignModeHandler's default mode if unspecified
				signMode = txCfg.SignModeHandler().DefaultMode()
			}
			if signMode == signing.SignMode_SIGN_MODE_DIRECT && len(txBuilder.GetTx().GetSigners()) > 1 {
				return sdkerrors.Wrap(sdkerrors.ErrNotSupported, "Signing in DIRECT mode is only supported for transactions with one signer only")
			}

			key, err := txFactory.Keybase().Key(fromName)
			if err != nil {
				return err
			}
			pubKey := key.GetPubKey()
			signerData := authsigning.SignerData{
				ChainID:       txFactory.ChainID(),
				AccountNumber: txFactory.AccountNumber(),
				Sequence:      txFactory.Sequence(),
			}

			// For SIGN_MODE_DIRECT, calling SetSignatures calls setSignerInfos on
			// TxBuilder under the hood, and SignerInfos is needed to generated the
			// sign bytes. This is the reason for setting SetSignatures here, with a
			// nil signature.
			//
			// Note: this line is not needed for SIGN_MODE_LEGACY_AMINO, but putting it
			// also doesn't affect its generated sign bytes, so for code's simplicity
			// sake, we put it here.
			sigData := signing.SingleSignatureData{
				SignMode:  signMode,
				Signature: nil,
			}
			sig := signing.SignatureV2{
				PubKey:   pubKey,
				Data:     &sigData,
				Sequence: txFactory.Sequence(),
			}
			if err := txBuilder.SetSignatures(sig); err != nil {
				return err
			}

			// Generate the bytes to be signed.
			bytesToSign, err := txCfg.SignModeHandler().GetSignBytes(signMode, signerData, txBuilder.GetTx())
			if err != nil {
				return err
			}
			cmd.Printf("Bytes to sign:\n%s\n", base64.StdEncoding.EncodeToString(bytesToSign))

			// Print the input tx with the filled `auth_info.signer_infos` field.
			json, err := clientCtx.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
			if err != nil {
				return err
			}
			cmd.Printf("Tx with filled `signer_infos`:\n%s\n", json)
			return nil
		},
	}
	cmd.Flags().String(flags.FlagChainID, "", "The network chain ID")
	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)
	cmd.MarkFlagRequired(flags.FlagChainID)
	cmd.MarkFlagRequired(flags.FlagSequence)
	cmd.MarkFlagRequired(flags.FlagAccountNumber)
	return cmd
}

type appCreator struct {
	encCfg params.EncodingConfig
}

func (ac appCreator) newApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) servertypes.Application {
	var cache sdk.MultiStorePersistentCache

	if cast.ToBool(appOpts.Get(server.FlagInterBlockCache)) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	pruningOpts, err := server.GetPruningOptionsFromFlags(appOpts)
	if err != nil {
		panic(err)
	}

	snapshotDir := filepath.Join(cast.ToString(appOpts.Get(flags.FlagHome)), "data", "snapshots")
	snapshotDB, err := sdk.NewLevelDB("metadata", snapshotDir)
	if err != nil {
		panic(err)
	}
	snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	if err != nil {
		panic(err)
	}

	return govgen.NewGovGenApp(
		logger, db, traceStore, true, skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		ac.encCfg,
		appOpts,
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(cast.ToString(appOpts.Get(server.FlagMinGasPrices))),
		baseapp.SetHaltHeight(cast.ToUint64(appOpts.Get(server.FlagHaltHeight))),
		baseapp.SetHaltTime(cast.ToUint64(appOpts.Get(server.FlagHaltTime))),
		baseapp.SetMinRetainBlocks(cast.ToUint64(appOpts.Get(server.FlagMinRetainBlocks))),
		baseapp.SetInterBlockCache(cache),
		baseapp.SetTrace(cast.ToBool(appOpts.Get(server.FlagTrace))),
		baseapp.SetIndexEvents(cast.ToStringSlice(appOpts.Get(server.FlagIndexEvents))),
		baseapp.SetSnapshotStore(snapshotStore),
		baseapp.SetSnapshotInterval(cast.ToUint64(appOpts.Get(server.FlagStateSyncSnapshotInterval))),
		baseapp.SetSnapshotKeepRecent(cast.ToUint32(appOpts.Get(server.FlagStateSyncSnapshotKeepRecent))),
		baseapp.SetIAVLCacheSize(cast.ToInt(appOpts.Get(server.FlagIAVLCacheSize))),
	)
}

func (ac appCreator) appExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
) (servertypes.ExportedApp, error) {
	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home is not set")
	}

	var loadLatest bool
	if height == -1 {
		loadLatest = true
	}

	govgenApp := govgen.NewGovGenApp(
		logger,
		db,
		traceStore,
		loadLatest,
		map[int64]bool{},
		homePath,
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		ac.encCfg,
		appOpts,
	)

	if height != -1 {
		if err := govgenApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	}

	return govgenApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs)
}
