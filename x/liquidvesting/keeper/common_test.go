package keeper_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	bankkeeper "github.com/milkyway-labs/milkyway/x/bank/keeper"
	"github.com/milkyway-labs/milkyway/x/liquidvesting/keeper"
	"github.com/milkyway-labs/milkyway/x/liquidvesting/testutils"
	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
	operatorskeeper "github.com/milkyway-labs/milkyway/x/operators/keeper"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolskeeper "github.com/milkyway-labs/milkyway/x/pools/keeper"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	restakingkeeper "github.com/milkyway-labs/milkyway/x/restaking/keeper"
	serviceskeeper "github.com/milkyway-labs/milkyway/x/services/keeper"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

const (
	IBCDenom       = "ibc/D79E7D83AB399BFFF93433E54FAA480C191248FC556924A2A8351AE2638B3877"
	vestedIBCDenom = "vested/" + IBCDenom
)

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context

	bk *bankkeeper.Keeper
	ok *operatorskeeper.Keeper
	pk *poolskeeper.Keeper
	sk *serviceskeeper.Keeper
	rk *restakingkeeper.Keeper

	k *keeper.Keeper
}

func (suite *KeeperTestSuite) SetupTest() {
	data := testutils.NewKeeperTestData(suite.T())

	// Context and codecs
	suite.ctx = data.Context

	// Keepers
	suite.bk = &data.BankKeeper
	suite.pk = data.PoolsKeeper
	suite.ok = data.OperatorsKeeper
	suite.sk = data.ServicesKeeper
	suite.rk = data.RestakingKeeper
	suite.k = data.Keeper
}

// --------------------------------------------------------------------------------------------------------------------

// fundAccount add the given amount of coins to the account's balance
func (suite *KeeperTestSuite) fundAccount(ctx sdk.Context, address string, amount sdk.Coins) {
	// Mint the tokens in the insurance fund.
	suite.Assert().NoError(suite.bk.MintCoins(ctx, types.ModuleName, amount))

	suite.Assert().NoError(suite.bk.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, sdk.MustAccAddressFromBech32(address), amount))
}

// mintVestedRepresentation mints the vested representation of the provided amount to
// the user balance
func (suite *KeeperTestSuite) mintVestedRepresentation(address string, amount sdk.Coins) {
	accAddress, err := sdk.AccAddressFromBech32(address)
	suite.Assert().NoError(err)

	_, err = suite.k.MintVestedRepresentation(
		suite.ctx, accAddress, amount,
	)
	suite.Assert().NoError(err)
}

// fundAccountInsuranceFund add the given amount of coins to the account's insurance fund
func (suite *KeeperTestSuite) fundAccountInsuranceFund(ctx sdk.Context, address string, amount sdk.Coins) {
	// Mint the tokens in the insurance fund.
	suite.Assert().NoError(suite.bk.MintCoins(suite.ctx, types.ModuleName, amount))

	// Assign those tokens to the user insurance fund
	userAddress, err := sdk.AccAddressFromBech32(address)
	suite.Assert().NoError(err)
	suite.Assert().NoError(suite.k.AddToUserInsuranceFund(
		ctx,
		userAddress,
		amount,
	))
}

// createPool creates a test pool with the given id and denom
func (suite *KeeperTestSuite) createPool(id uint32, denom string) {
	err := suite.pk.SavePool(suite.ctx, poolstypes.Pool{
		ID:              id,
		Denom:           denom,
		Address:         poolstypes.GetPoolAddress(id).String(),
		Tokens:          sdkmath.NewInt(0),
		DelegatorShares: sdkmath.LegacyNewDec(0),
	})
	suite.Assert().NoError(err)
}

// createService creates a test service with the provided id
func (suite *KeeperTestSuite) createService(id uint32) {
	err := suite.sk.CreateService(suite.ctx, servicestypes.NewService(
		id,
		servicestypes.SERVICE_STATUS_ACTIVE,
		fmt.Sprintf("test %d", id),
		fmt.Sprintf("test service %d", id),
		"",
		"",
		fmt.Sprintf("service-%d-admin", id),
		false,
	))
	suite.Assert().NoError(err)
}

func (suite *KeeperTestSuite) createOperator(id uint32) {
	suite.Assert().NoError(suite.ok.RegisterOperator(suite.ctx, operatorstypes.NewOperator(
		id,
		operatorstypes.OPERATOR_STATUS_ACTIVE,
		fmt.Sprintf("operator-%d", id),
		"",
		"",
		fmt.Sprintf("operator-%d-admin", id))))
}
