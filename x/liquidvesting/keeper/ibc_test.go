package keeper_test

import (
	"encoding/json"
	"fmt"

	sdkmath "cosmossdk.io/math"

	"github.com/milkyway-labs/milkyway/v9/x/liquidvesting/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
)

func (suite *KeeperTestSuite) TestIBCTransferVestedTokens() {
	// Setup a user with vested tokens
	vestedCoin := sdk.NewInt64Coin("stake", 1000)
	user := sdk.AccAddress("user_______________")
	
	// Fund the user with vested tokens
	ctx := suite.ctx
	err := suite.bk.MintCoins(ctx, types.ModuleName, sdk.NewCoins(vestedCoin))
	suite.Require().NoError(err)
	
	// Create vested tokens for the user
	err = suite.k.MintLiquidVestedTokens(ctx, user, sdk.NewCoins(vestedCoin))
	suite.Require().NoError(err)
	
	// Verify user has liquid vested tokens
	liquidDenom := types.GetLiquidVestedTokenDenom("stake")
	balance := suite.bk.GetBalance(ctx, user, liquidDenom)
	suite.Require().Equal(vestedCoin.Amount, balance.Amount)
	
	// Simulate IBC transfer of the liquid vested tokens
	// This test is limited because we can't actually test the full IBC flow,
	// but we can verify the tokens can be sent and the user loses the balance
	dest := sdk.AccAddress("destination________")
	
	// Perform the transfer
	err = suite.bk.SendCoins(ctx, user, dest, sdk.NewCoins(sdk.NewCoin(liquidDenom, sdkmath.NewInt(500))))
	suite.Require().NoError(err)
	
	// Verify balances after transfer
	userBalance := suite.bk.GetBalance(ctx, user, liquidDenom)
	suite.Require().Equal(sdkmath.NewInt(500), userBalance.Amount)
	
	destBalance := suite.bk.GetBalance(ctx, dest, liquidDenom)
	suite.Require().Equal(sdkmath.NewInt(500), destBalance.Amount)
}

func (suite *KeeperTestSuite) TestIBCReceiveTokens() {
	// This test simulates receiving tokens via IBC
	// We'll mock the OnRecvPacket method directly since we can't simulate the full IBC flow
	
	// Setup a destination address
	dest := "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"
	amount := sdkmath.NewInt(1000)
	
	// Setup parameters to allow channel-0
	err := suite.k.SetParams(suite.ctx, types.Params{
		InsurancePercentage: sdkmath.LegacyNewDec(2),
		AllowedChannels:     []string{"channel-0"},
	})
	suite.Require().NoError(err)
	
	// Create a deposit insurance message
	memo := fmt.Sprintf(`{
		"liquidvesting": {
			"amounts": [{
				"depositor": "%s",
				"amount": "%s"
			}]
		}}`, dest, amount.String())
	
	// Build the data for the packet
	dataBz, err := json.Marshal(&transfertypes.FungibleTokenPacketData{
		Denom:    "foo",
		Amount:   amount.String(),
		Sender:   "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
		Receiver: authtypes.NewModuleAddress(types.ModuleName).String(),
		Memo:     memo,
	})
	suite.Require().NoError(err)
	
	// Build the packet
	packet := channeltypes.Packet{
		Data:               dataBz,
		DestinationChannel: "channel-0",
		DestinationPort:    "transfer",
	}
	
	// Process the packet
	ack := suite.ibcm.OnRecvPacket(suite.ctx, packet, suite.ak.GetModuleAddress("relayer"))
	
	// Verify success
	suite.Require().True(ack.Success())
	
	// Verify the insurance fund was updated
	insuranceFund, err := suite.k.GetUserInsuranceFundBalance(suite.ctx, dest)
	suite.Require().NoError(err)
	suite.Require().Equal("1000ibc/EB7094899ACFB7A6F2A67DB084DEE2E9A83DEFAA5DEF92D9A9814FFD9FF673FA", insuranceFund.String())
}

func (suite *KeeperTestSuite) TestKeeper_IBCHooks() {
	testCases := []struct {
		name               string
		store              func(ctx sdk.Context)
		destinationChannel string
		transferAmount     sdk.Coin
		sender             string
		receiver           string
		memo               string
		shouldErr          bool
		check              func(sdk.Context)
	}{
		{
			name:               "empty memo works as normal transfer",
			destinationChannel: "channel-0",
			transferAmount:     sdk.NewInt64Coin("foo", 1000),
			sender:             "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			receiver:           "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			memo:               "",
			shouldErr:          false,
		},
		{
			name: "sending to a normal account returns error",
			store: func(ctx sdk.Context) {
				err := suite.k.SetParams(ctx, types.Params{
					InsurancePercentage: sdkmath.LegacyNewDec(2),
					AllowedChannels:     []string{"channel-0"},
				})
				suite.Require().NoError(err)
			},
			destinationChannel: "channel-0",
			transferAmount:     sdk.NewInt64Coin("foo", 1000),
			sender:             "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			receiver:           "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			memo: fmt.Sprintf(`{
			"liquidvesting": {
				"amounts": [{
					"depositor": "%s",
					"amount": "1000"
				}]
			}}`, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			shouldErr: true,
			check: func(ctx sdk.Context) {
				// Make sure the user's insurance fund is not updated
				insuranceFund, err := suite.k.GetUserInsuranceFundBalance(ctx, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Assert().NoError(err)
				suite.Assert().Empty(insuranceFund)
			},
		},
		{
			name: "depositing more coins then received returns error",
			store: func(ctx sdk.Context) {
				err := suite.k.SetParams(ctx, types.Params{
					InsurancePercentage: sdkmath.LegacyNewDec(2),
					AllowedChannels:     []string{"channel-0"},
				})
				suite.Require().NoError(err)
			},
			destinationChannel: "channel-0",
			transferAmount:     sdk.NewInt64Coin("foo", 1000),
			sender:             "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			receiver:           authtypes.NewModuleAddress(types.ModuleName).String(),
			memo: fmt.Sprintf(`{
            "liquidvesting": {
                "amounts": [{
                    "depositor": "%s",
                    "amount": "400"
                },
                {
                    "depositor": "%s",
                    "amount": "601"
                }]
            }}`,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
			check: func(ctx sdk.Context) {
				// Make sure the user's insurance fund is not updated
				insuranceFund, err := suite.k.GetUserInsuranceFundBalance(ctx, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Assert().NoError(err)
				suite.Assert().Empty(insuranceFund)
			},
		},
		{
			name: "deposit less coins then received returns error",
			store: func(ctx sdk.Context) {
				err := suite.k.SetParams(ctx, types.Params{
					InsurancePercentage: sdkmath.LegacyNewDec(2),
					AllowedChannels:     []string{"channel-0"},
				})
				suite.Require().NoError(err)
			},
			destinationChannel: "channel-0",
			transferAmount:     sdk.NewInt64Coin("foo", 1000),
			sender:             "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			receiver:           authtypes.NewModuleAddress(types.ModuleName).String(),
			memo: fmt.Sprintf(`{
            "liquidvesting": {
                "amounts": [{
                    "depositor": "%s",
                    "amount": "300"
                },
                {
                    "depositor": "%s",
                    "amount": "600"
                }]
            }}`,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
			check: func(ctx sdk.Context) {
				// Make sure the user's insurance fund is not updated
				insuranceFund, err := suite.k.GetUserInsuranceFundBalance(ctx, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Assert().NoError(err)
				suite.Assert().Empty(insuranceFund)
			},
		},
		{
			name: "deposit from not allowed channel fails",
			store: func(ctx sdk.Context) {
				err := suite.k.SetParams(ctx, types.Params{
					InsurancePercentage: sdkmath.LegacyNewDec(2),
					AllowedChannels:     []string{"channel-1"},
				})
				suite.Require().NoError(err)
			},
			destinationChannel: "channel-0",
			transferAmount:     sdk.NewInt64Coin("foo", 1000),
			sender:             "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			receiver:           authtypes.NewModuleAddress(types.ModuleName).String(),
			memo: fmt.Sprintf(`{
            "liquidvesting": {
                "amounts": [{
                    "depositor": "%s",
                    "amount": "600"
                },
                {
                    "depositor": "%s",
                    "amount": "400"
                }]
            }}`,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd"),
			shouldErr: true,
		},
		{
			name: "correct deposit works properly",
			store: func(ctx sdk.Context) {
				// Set the sender as an allowed depositor
				err := suite.k.SetParams(ctx, types.Params{
					InsurancePercentage: sdkmath.LegacyNewDec(2),
					AllowedChannels:     []string{"channel-0"},
				})
				suite.Require().NoError(err)
			},
			destinationChannel: "channel-0",
			transferAmount:     sdk.NewInt64Coin("foo", 1000),
			sender:             "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			receiver:           authtypes.NewModuleAddress(types.ModuleName).String(),
			memo: fmt.Sprintf(`{
            "liquidvesting": {
                "amounts": [{
                    "depositor": "%s",
                    "amount": "600"
                },
                {
                    "depositor": "%s",
                    "amount": "400"
                }]
            }}`,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				// Make sure the first insurance fund is updated
				insuranceFund, err := suite.k.GetUserInsuranceFundBalance(ctx, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Assert().NoError(err)
				suite.Assert().Equal("600ibc/EB7094899ACFB7A6F2A67DB084DEE2E9A83DEFAA5DEF92D9A9814FFD9FF673FA", insuranceFund.String())

				// Make sure the second insurance fund is updated
				insuranceFund, err = suite.k.GetUserInsuranceFundBalance(ctx, "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd")
				suite.Assert().NoError(err)
				suite.Assert().Equal("400ibc/EB7094899ACFB7A6F2A67DB084DEE2E9A83DEFAA5DEF92D9A9814FFD9FF673FA", insuranceFund.String())
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			// Cache the context
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			// Build the data to be put inside the packet
			dataBz, err := json.Marshal(&transfertypes.FungibleTokenPacketData{
				Denom:    tc.transferAmount.Denom,
				Amount:   tc.transferAmount.Amount.String(),
				Sender:   tc.sender,
				Receiver: tc.receiver,
				Memo:     tc.memo,
			})
			suite.Assert().NoError(err)

			// Build the packet
			packet := channeltypes.Packet{
				Data:               dataBz,
				DestinationChannel: tc.destinationChannel,
				DestinationPort:    "transfer",
			}

			// Receive the packet
			ack := suite.ibcm.OnRecvPacket(ctx, packet, suite.ak.GetModuleAddress("relayer"))
			ack.Acknowledgement()

			if tc.shouldErr {
				suite.Assert().False(ack.Success())

				castedAck := ack.(channeltypes.Acknowledgement)
				errorResponse := castedAck.Response.(*channeltypes.Acknowledgement_Error)
				suite.Require().NotEmpty(errorResponse.Error)
			} else {
				suite.Assert().True(ack.Success())
			}

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}
