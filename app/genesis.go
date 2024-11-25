package milkyway

import (
	"encoding/json"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	tmtypes "github.com/cometbft/cometbft/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/initia-labs/initia/app/genesis_markets"
	marketmaptypes "github.com/skip-mev/connect/v2/x/marketmap/types"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibctypes "github.com/cosmos/ibc-go/v8/modules/core/types"
)

// GenesisState - The genesis state of the blockchain is represented here as a map of raw json
// messages key'd by a identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState(cdc codec.Codec, mbm module.BasicManager) GenesisState {
	return GenesisState(mbm.DefaultGenesis(cdc)).
		ConfigureIBCAllowedClients(cdc).
		AddMarketData(cdc)
}

func (genState GenesisState) AddMarketData(cdc codec.JSONCodec) GenesisState {
	var oracleGenState oracletypes.GenesisState
	cdc.MustUnmarshalJSON(genState[oracletypes.ModuleName], &oracleGenState)

	var marketGenState marketmaptypes.GenesisState
	cdc.MustUnmarshalJSON(genState[marketmaptypes.ModuleName], &marketGenState)

	// Load initial markets
	markets, err := genesis_markets.ReadMarketsFromFile(genesis_markets.GenesisMarkets)
	if err != nil {
		panic(err)
	}
	marketGenState.MarketMap = genesis_markets.ToMarketMap(markets)

	// Initialize all markets
	var id uint64
	currencyPairGenesis := make([]oracletypes.CurrencyPairGenesis, len(markets))
	for i, market := range markets {
		currencyPairGenesis[i] = oracletypes.CurrencyPairGenesis{
			CurrencyPair:      market.Ticker.CurrencyPair,
			CurrencyPairPrice: nil,
			Nonce:             0,
			Id:                id,
		}
		id++
	}

	oracleGenState.CurrencyPairGenesis = currencyPairGenesis
	oracleGenState.NextId = id

	// write the updates to genState
	genState[marketmaptypes.ModuleName] = cdc.MustMarshalJSON(&marketGenState)
	genState[oracletypes.ModuleName] = cdc.MustMarshalJSON(&oracleGenState)
	return genState
}

func (genState GenesisState) ConfigureIBCAllowedClients(cdc codec.JSONCodec) GenesisState {
	var ibcGenesis ibctypes.GenesisState
	cdc.MustUnmarshalJSON(genState[ibcexported.ModuleName], &ibcGenesis)

	allowedClients := ibcGenesis.ClientGenesis.Params.AllowedClients
	for i, client := range allowedClients {
		if client == ibcexported.Localhost {
			allowedClients = append(allowedClients[:i], allowedClients[i+1:]...)
			break
		}
	}

	ibcGenesis.ClientGenesis.Params.AllowedClients = allowedClients
	genState[ibcexported.ModuleName] = cdc.MustMarshalJSON(&ibcGenesis)

	return genState
}

// NewDefaultGenesisStateWithValidator generates the default application state with a validator.
func NewDefaultGenesisStateWithValidator(cdc codec.Codec, mbm module.BasicManager) GenesisState {
	privVal := mock.NewPV()
	pubKey, _ := privVal.GetPubKey()
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	senderPrivKey.PubKey().Address()
	acc := authtypes.NewBaseAccountWithAddress(senderPrivKey.PubKey().Address().Bytes())

	//////////////////////
	var balances []banktypes.Balance
	genesisState := NewDefaultGenesisState(cdc, mbm)
	genAccs := []authtypes.GenesisAccount{acc}
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = cdc.MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.DefaultPowerReduction

	for _, val := range valSet.Validators {
		pk, _ := cryptocodec.FromCmtPubKeyInterface(val.PubKey)
		pkAny, _ := codectypes.NewAnyWithValue(pk)
		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   sdkmath.LegacyNewDec(1),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec()),
			MinSelfDelegation: sdkmath.ZeroInt(),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress().String(), sdk.ValAddress(val.Address).String(), sdkmath.LegacyNewDec(1)))
	}
	// set validators and delegations
	stakingGenesis := stakingtypes.NewGenesisState(stakingtypes.DefaultParams(), validators, delegations)
	genesisState[stakingtypes.ModuleName] = cdc.MustMarshalJSON(stakingGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}

	for range delegations {
		// add delegated tokens to total supply
		totalSupply = totalSupply.Add(sdk.NewCoin(sdk.DefaultBondDenom, bondAmt))
	}

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, bondAmt)},
	})

	// update total supply
	bankGenesis := banktypes.NewGenesisState(
		banktypes.DefaultGenesisState().Params,
		balances,
		totalSupply,
		[]banktypes.Metadata{},
		[]banktypes.SendEnabled{},
	)
	genesisState[banktypes.ModuleName] = cdc.MustMarshalJSON(bankGenesis)

	return genesisState
}
