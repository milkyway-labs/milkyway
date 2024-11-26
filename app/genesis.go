package milkyway

import (
	"encoding/json"
	"slices"

	"github.com/skip-mev/connect/v2/cmd/constants/marketmaps"
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
	coreMarkets := marketmaps.CoreMarketMap
	markets := coreMarkets.Markets

	// Sort keys so we can deterministically iterate over map items.
	keys := make([]string, 0, len(markets))
	for name := range markets {
		keys = append(keys, name)
	}
	slices.Sort(keys)

	// Initialize all markets
	var id uint64
	currencyPairGenesis := make([]oracletypes.CurrencyPairGenesis, len(markets))
	for _, market := range markets {
		currencyPairGenesis[id] = oracletypes.CurrencyPairGenesis{
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
