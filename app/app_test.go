package govgen_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	govgen "github.com/atomone-hub/govgen/v1/app"
	govgenhelpers "github.com/atomone-hub/govgen/v1/app/helpers"
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
