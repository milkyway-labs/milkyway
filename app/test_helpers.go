package milkyway

// DONTCOVER

import (
	"encoding/json"
	"testing"
	"time"

	"cosmossdk.io/log"
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/server"
	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	opchildtypes "github.com/initia-labs/OPinit/x/opchild/types"

	"github.com/milkyway-labs/milkyway/app/params"
)

// defaultConsensusParams defines the default Tendermint consensus params used in
// MilkyWayApp testing.
var defaultConsensusParams = &tmproto.ConsensusParams{
	Block: &tmproto.BlockParams{
		MaxBytes: 8000000,
		MaxGas:   1234000000,
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

func getOrCreateMemDB(db *dbm.DB) dbm.DB {
	if db != nil {
		return *db
	}
	return dbm.NewMemDB()
}

func setup(t *testing.T, db *dbm.DB, withGenesis bool) (*MilkyWayApp, GenesisState) {
	t.Helper()

	appOptions := make(simtestutil.AppOptionsMap, 0)
	appOptions[server.FlagInvCheckPeriod] = simcli.FlagPeriodValue

	encCdc := params.MakeEncodingConfig()
	app := NewMilkyWayApp(
		log.NewNopLogger(),
		getOrCreateMemDB(db),
		nil,
		false,
		map[int64]bool{},
		DefaultNodeHome,
		simtestutil.NewAppOptionsWithFlagHome(t.TempDir()),
		[]wasmkeeper.Option{},
		baseapp.SetChainID("milkyway-app"),
	)

	if withGenesis {
		return app, NewDefaultGenesisState(encCdc.Marshaler, app.ModuleBasics)
	}

	return app, GenesisState{}
}

// Setup initializes a new MilkyWayApp for testing.
// A single validator will be created and registered in opchild module.
func Setup(t *testing.T, isCheckTx bool) *MilkyWayApp {
	t.Helper()

	app, genState := setup(t, nil, true)

	// Create a validator which will be the admin of the chain as well as the
	// bridge executor.
	privVal := ed25519.GenPrivKey() // TODO: make it deterministic?
	pubKey := privVal.PubKey()
	pubKeyAny, err := codectypes.NewAnyWithValue(privVal.PubKey())
	if err != nil {
		panic(err)
	}
	validator := opchildtypes.Validator{
		Moniker:         "test-validator",
		OperatorAddress: sdk.ValAddress(privVal.PubKey().Address()).String(),
		ConsensusPubkey: pubKeyAny,
		ConsPower:       1,
	}

	// set validators and delegations
	var opchildGenesis opchildtypes.GenesisState
	app.AppCodec().MustUnmarshalJSON(genState[opchildtypes.ModuleName], &opchildGenesis)
	opchildGenesis.Params.Admin = sdk.AccAddress(pubKey.Address().Bytes()).String()
	opchildGenesis.Params.BridgeExecutors = []string{sdk.AccAddress(pubKey.Address().Bytes()).String()}
	opchildGenesis.Validators = []opchildtypes.Validator{validator}
	genState[opchildtypes.ModuleName] = app.AppCodec().MustMarshalJSON(&opchildGenesis)

	if !isCheckTx {
		genStateBytes, err := json.Marshal(genState)
		if err != nil {
			panic(err)
		}
		_, err = app.InitChain(&abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: defaultConsensusParams,
			AppStateBytes:   genStateBytes,
		})
		if err != nil {
			panic(err)
		}
	}

	return app
}

// SetupWithGenesisAccounts setup initiaapp with genesis account
func SetupWithGenesisAccounts(
	t *testing.T,
	valSet *tmtypes.ValidatorSet,
	genAccs []authtypes.GenesisAccount,
	balances ...banktypes.Balance,
) *MilkyWayApp {
	t.Helper()

	app, genesisState := setup(t, nil, true)

	if len(genAccs) == 0 {
		privAcc := secp256k1.GenPrivKey()
		genAccs = []authtypes.GenesisAccount{
			authtypes.NewBaseAccount(privAcc.PubKey().Address().Bytes(), privAcc.PubKey(), 0, 0),
		}
	}

	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	// allow empty validator
	if valSet == nil || len(valSet.Validators) == 0 {
		privVal := ed25519.GenPrivKey()
		pubKey, err := cryptocodec.ToCmtPubKeyInterface(privVal.PubKey())
		if err != nil {
			panic(err)
		}

		validator := tmtypes.NewValidator(pubKey, 1)
		valSet = tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})
	}

	validators := make([]opchildtypes.Validator, 0, len(valSet.Validators))
	for _, val := range valSet.Validators {
		pk, err := cryptocodec.FromCmtPubKeyInterface(val.PubKey)
		if err != nil {
			panic(err)
		}
		pkAny, err := codectypes.NewAnyWithValue(pk)
		if err != nil {
			panic(err)
		}

		validator := opchildtypes.Validator{
			Moniker:         "test-validator",
			OperatorAddress: sdk.ValAddress(val.Address).String(),
			ConsensusPubkey: pkAny,
			ConsPower:       1,
		}

		validators = append(validators, validator)
	}

	// set validators and delegations
	var opchildGenesis opchildtypes.GenesisState
	app.AppCodec().MustUnmarshalJSON(genesisState[opchildtypes.ModuleName], &opchildGenesis)
	opchildGenesis.Params.Admin = sdk.AccAddress(valSet.Validators[0].Address.Bytes()).String()
	opchildGenesis.Params.BridgeExecutors = []string{sdk.AccAddress(valSet.Validators[0].Address.Bytes()).String()}

	// set validators and delegations
	opchildGenesis = *opchildtypes.NewGenesisState(opchildGenesis.Params, validators, nil)
	genesisState[opchildtypes.ModuleName] = app.AppCodec().MustMarshalJSON(&opchildGenesis)

	// update total supply
	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, sdk.NewCoins(), []banktypes.Metadata{}, []banktypes.SendEnabled{})
	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		panic(err)
	}

	_, err = app.InitChain(
		&abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: defaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)
	if err != nil {
		panic(err)
	}

	_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{Height: 1})
	if err != nil {
		panic(err)
	}

	_, err = app.Commit()
	if err != nil {
		panic(err)
	}

	return app
}
