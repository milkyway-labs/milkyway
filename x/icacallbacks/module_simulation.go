package icacallbacks

// XXX
//// avoid unused import issue
//var (
//	_ = sample.AccAddress
//	_ = icacallbackssimulation.FindAccount
//	_ = simappparams.StakePerAccount
//	_ = simulation.MsgEntryKind
//	_ = baseapp.Paramspace
//)
//
//const (
//// this line is used by starport scaffolding # simapp/module/const
//)
//
//// GenerateGenesisState creates a randomized GenState of the module
//func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
//	accs := make([]string, len(simState.Accounts))
//	for i, acc := range simState.Accounts {
//		accs[i] = acc.Address.String()
//	}
//	icacallbacksGenesis := types.GenesisState{
//		Params: types.DefaultParams(),
//		PortId: types.PortID,
//		// this line is used by starport scaffolding # simapp/module/genesisState
//	}
//	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&icacallbacksGenesis)
//}
//
//// ProposalContents doesn't return any content functions for governance proposals
//func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalMsg {
//	return nil
//}
//
//// RandomizedParams creates randomized  param changes for the simulator
//func (am AppModule) RandomizedParams(_ *rand.Rand) []simtypes.LegacyParamChange {
//	return []simtypes.LegacyParamChange{}
//}
//
//// RegisterStoreDecoder registers a decoder
//func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}
//
//// WeightedOperations returns the all the gov module operations with their respective weights.
//func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
//	operations := make([]simtypes.WeightedOperation, 0)
//
//	// this line is used by starport scaffolding # simapp/module/operation
//
//	return operations
//}