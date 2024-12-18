package abci

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/skip-mev/connect/v2/abci/strategies/aggregator"
	oracleconfig "github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/pkg/math/voteweighted"
	oracleclient "github.com/skip-mev/connect/v2/service/clients/oracle"
	servicemetrics "github.com/skip-mev/connect/v2/service/metrics"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/skip-mev/connect/v2/abci/proposals"
	compression "github.com/skip-mev/connect/v2/abci/strategies/codec"
	"github.com/skip-mev/connect/v2/abci/strategies/currencypair"
	"github.com/skip-mev/connect/v2/abci/ve"

	oraclepreblock "github.com/skip-mev/connect/v2/abci/preblock/oracle"

	"github.com/milkyway-labs/milkyway/v5/app/keepers"
)

type SetupData struct {
	ChainID       string
	Logger        log.Logger
	Keepers       keepers.AppKeepers
	ModuleManager *module.Manager
	OracleConfig  oracleconfig.AppConfig
}

type ExtensionsData struct {
	OracleClient oracleclient.OracleClient

	PrepareProposalHandler sdk.PrepareProposalHandler
	ProposalHandler        sdk.PrepareProposalHandler
	ProcessProposalHandler sdk.ProcessProposalHandler
	PreBlockHandler        sdk.PreBlocker

	ExtendVoteHandler          sdk.ExtendVoteHandler
	VerifyVoteExtensionHandler sdk.VerifyVoteExtensionHandler
}

// initializeOracle initializes the oracle client and metrics.
func initializeOracle(chainID string, cfg oracleconfig.AppConfig, logger log.Logger) (oracleclient.OracleClient, servicemetrics.Metrics, error) {
	// If app level instrumentation is enabled, then wrap the oracle service with a metrics client
	// to get metrics on the oracle service (for ABCI++). This will allow the instrumentation to track
	// latency in VerifyVoteExtension requests and more.
	oracleMetrics, err := servicemetrics.NewMetricsFromConfig(cfg, chainID)
	if err != nil {
		return nil, nil, err
	}

	// Create the oracle service.
	oracleClient, err := oracleclient.NewPriceDaemonClientFromConfig(
		cfg,
		logger.With("client", "oracle"),
		oracleMetrics,
	)
	if err != nil {
		return nil, nil, err
	}

	// Connect to the oracle service (default timeout of 5 seconds).
	go func() {
		logger.Info("attempting to start oracle client...", "address", cfg.OracleAddress)
		if err := oracleClient.Start(context.Background()); err != nil {
			logger.Error("failed to start oracle client", "err", err)
			panic(err)
		}
	}()

	return oracleClient, oracleMetrics, nil
}

func InitializeOracleABCIExtensions(data SetupData) ExtensionsData {
	// Initialize the oracle client and metrics
	oracleClient, oracleMetrics, err := initializeOracle(data.ChainID, data.OracleConfig, data.Logger)
	if err != nil {
		panic(fmt.Errorf("failed to initialize oracle client and metrics: %w", err))
	}

	// Create the proposal handler that will be used to fill proposals with
	// transactions and oracle data.
	proposalHandler := proposals.NewProposalHandler(
		data.Logger,
		baseapp.NoOpPrepareProposal(),
		baseapp.NoOpProcessProposal(),
		ve.NewDefaultValidateVoteExtensionsFn(data.Keepers.StakingKeeper),
		compression.NewCompressionVoteExtensionCodec(
			compression.NewDefaultVoteExtensionCodec(),
			compression.NewZLibCompressor(),
		),
		compression.NewCompressionExtendedCommitCodec(
			compression.NewDefaultExtendedCommitCodec(),
			compression.NewZStdCompressor(),
		),
		currencypair.NewDeltaCurrencyPairStrategy(data.Keepers.OracleKeeper),
		oracleMetrics,
	)

	// Create the aggregation function that will be used to aggregate oracle data
	// from each validator.
	aggregatorFn := voteweighted.MedianFromContext(
		data.Logger,
		data.Keepers.StakingKeeper,
		voteweighted.DefaultPowerThreshold,
	)
	veCodec := compression.NewCompressionVoteExtensionCodec(
		compression.NewDefaultVoteExtensionCodec(),
		compression.NewZLibCompressor(),
	)
	ecCodec := compression.NewCompressionExtendedCommitCodec(
		compression.NewDefaultExtendedCommitCodec(),
		compression.NewZStdCompressor(),
	)

	// Create the pre-finalize block hook that will be used to apply oracle data
	// to the state before any transactions are executed (in finalize block).
	oraclePreBlockHandler := oraclepreblock.NewOraclePreBlockHandler(
		data.Logger,
		aggregatorFn,
		data.Keepers.OracleKeeper,
		oracleMetrics,
		currencypair.NewDeltaCurrencyPairStrategy(data.Keepers.OracleKeeper), // IMPORTANT: always construct new currency pair strategy objects when functions require them as arguments.
		veCodec,
		ecCodec,
	)

	// Create the vote extensions handler that will be used to extend and verify
	// vote extensions (i.e. oracle data).
	voteExtensionsHandler := ve.NewVoteExtensionHandler(
		data.Logger,
		oracleClient,
		time.Second, // timeout
		currencypair.NewDeltaCurrencyPairStrategy(data.Keepers.OracleKeeper), // IMPORTANT: always construct new currency pair strategy objects when functions require them as arguments.
		veCodec,
		aggregator.NewOraclePriceApplier(
			aggregator.NewDefaultVoteAggregator(
				data.Logger,
				aggregatorFn,
				// we need a separate price strategy here, so that we can optimistically apply the latest prices
				// and extend our vote based on these prices
				currencypair.NewDeltaCurrencyPairStrategy(data.Keepers.OracleKeeper), // IMPORTANT: always construct new currency pair strategy objects when functions require them as arguments.
			),
			data.Keepers.OracleKeeper,
			veCodec,
			ecCodec,
			data.Logger,
		),
		oracleMetrics,
	)

	return ExtensionsData{
		OracleClient:               oracleClient,
		PrepareProposalHandler:     proposalHandler.PrepareProposalHandler(),
		ProcessProposalHandler:     proposalHandler.ProcessProposalHandler(),
		ProposalHandler:            proposalHandler.PrepareProposalHandler(),
		PreBlockHandler:            oraclePreBlockHandler.WrappedPreBlocker(data.ModuleManager),
		ExtendVoteHandler:          voteExtensionsHandler.ExtendVoteHandler(),
		VerifyVoteExtensionHandler: voteExtensionsHandler.VerifyVoteExtensionHandler(),
	}
}
