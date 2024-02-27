package govgen_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	db "github.com/tendermint/tm-db"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	govgen "github.com/atomone-hub/govgen/app"
	govgenhelpers "github.com/atomone-hub/govgen/app/helpers"
	govtypes "github.com/atomone-hub/govgen/x/gov/types"
)

type EmptyAppOptions struct{}

func (ao EmptyAppOptions) Get(_ string) interface{} {
	return nil
}

func TestGovGenApp_BlockedModuleAccountAddrs(t *testing.T) {
	app := govgen.NewGovGenApp(
		log.NewNopLogger(),
		db.NewMemDB(),
		nil,
		true,
		map[int64]bool{},
		govgen.DefaultNodeHome,
		0,
		govgen.MakeTestEncodingConfig(),
		EmptyAppOptions{},
	)
	moduleAccountAddresses := app.ModuleAccountAddrs()
	blockedAddrs := app.BlockedModuleAccountAddrs(moduleAccountAddresses)

	require.NotContains(t, blockedAddrs, authtypes.NewModuleAddress(govtypes.ModuleName).String())
}

func TestGovGenApp_Export(t *testing.T) {
	app := govgenhelpers.Setup(t)
	_, err := app.ExportAppStateAndValidators(true, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

func TestGovGenApp_InitialStakingDistribution(t *testing.T) {
	// generate 20 validators
	valset, _ := tmtypes.RandValidatorSet(20, 1)
	var (
		genesisAccounts []authtypes.GenesisAccount
		balances        []banktypes.Balance
	)
	// generate 100 accounts
	for i := 0; i < 100; i++ {
		senderPrivKey := govgenhelpers.NewPV()
		senderPubKey := senderPrivKey.PrivKey.PubKey()
		acc := authtypes.NewBaseAccount(senderPubKey.Address().Bytes(), senderPubKey, 0, 0)
		balance := banktypes.Balance{
			Address: acc.GetAddress().String(),
			Coins: sdk.NewCoins(
				// sdk.NewCoin("ugovgen", sdk.NewInt(1_000_000*tmrand.Int63n(1_000_000))),
				sdk.NewCoin("ugovgen", sdk.NewInt(1_000_000*1_000_000)),
				sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000_000_000)),
			),
		}
		balances = append(balances, balance)
		genesisAccounts = append(genesisAccounts, acc)
	}
	app := govgenhelpers.SetupWithGenesisValSet(t, valset, genesisAccounts, balances...)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{
		ChainID: fmt.Sprintf("test-chain-%s", tmrand.Str(4)),
		Height:  1,
	})

	// Checking fairness delegation distribution
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	validators := app.StakingKeeper.GetAllValidators(ctx)
	var shareReference int64
	for _, val := range validators {
		delegations := app.StakingKeeper.GetValidatorDelegations(ctx, val.GetOperator())
		if shareReference == 0 {
			// initialize the reference share, all other shares should match
			// approximately to assert godd fairness distribution.
			shareReference = delegations[0].Shares.TruncateInt64()
		}
		for _, del := range delegations {
			assert.InDelta(t, shareReference, del.Shares.TruncateInt64(), 1, "unfair share distribution")
		}
	}
	fmt.Println(len(delegations))
}
