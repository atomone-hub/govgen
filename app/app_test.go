package govgen_test

import (
	"fmt"
	"testing"

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
	// generate 30 validators and 100 genesis accounts
	var (
		valset, _       = tmtypes.RandValidatorSet(20, 1)
		genesisAccounts []authtypes.GenesisAccount
		balances        []banktypes.Balance
	)
	for i := 0; i < 100; i++ {
		senderPrivKey := govgenhelpers.NewPV()
		senderPubKey := senderPrivKey.PrivKey.PubKey()
		acc := authtypes.NewBaseAccount(senderPubKey.Address().Bytes(), senderPubKey, 0, 0)
		balance := banktypes.Balance{
			Address: acc.GetAddress().String(),
			Coins: sdk.NewCoins(
				sdk.NewCoin("ugovgen", sdk.NewInt(100_000_000_000_000)),
				sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000_000_000)),
			),
		}
		// if i == 0 {
		// balance.Coins = balance.Coins.Add(
		// sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(1_000_000*len(valset.Validators))),
		// sdk.NewInt64Coin(sdk.DefaultBondDenom, 1_000_000),
		// )
		// }
		balances = append(balances, balance)
		genesisAccounts = append(genesisAccounts, acc)
	}
	app := govgenhelpers.SetupWithGenesisValSet(t, valset, genesisAccounts, balances...)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{
		ChainID: fmt.Sprintf("test-chain-%s", tmrand.Str(4)),
		Height:  1,
	})

	// encodingConfig := gaiaapp.MakeTestEncodingConfig()
	// encodingConfig.Amino.RegisterConcrete(&testdata.TestMsg{}, "testdata.TestMsg", nil)
	// testdata.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	fmt.Println(delegations)
}
