package keeper_test

import (
	"encoding/json"
	"fmt"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
)

func (suite *KeeperTestSuite) TestKeeper_IBCHooks() {
	testCases := []struct {
		name           string
		transferAmount sdk.Coin
		sender         string
		receiver       string
		memo           string
		shouldErr      bool
		check          func(sdk.Context)
	}{
		{
			name:           "empty memo",
			transferAmount: sdk.NewInt64Coin("foo", 1000),
			sender:         "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			receiver:       "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			memo:           "",
		},
		{
			name:           "trigger by sending to a normal account",
			transferAmount: sdk.NewInt64Coin("foo", 1000),
			sender:         "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			receiver:       "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			memo: fmt.Sprintf(`{
			"liquidvesting": {
				"amounts": [{
					"depositor": "%s",
					"amount": { "amount": "1000", "denom": "foo" }
				}]
			}}`, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			shouldErr: true,
			check: func(ctx sdk.Context) {
				// Make sure the user's insurance fund is not updated
				userAddr, err := sdk.AccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Assert().NoError(err)

				insuranceFund, err := suite.k.GetUserInsuranceFundBalance(ctx, userAddr)
				suite.Assert().NoError(err)
				suite.Assert().Empty(insuranceFund)
			},
		},
		{
			name:           "transfer not received denom",
			transferAmount: sdk.NewInt64Coin("foo", 1000),
			sender:         "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			receiver:       authtypes.NewModuleAddress(types.ModuleName).String(),
			memo: fmt.Sprintf(`{"liquidvesting": {
				"amounts": [
					{
 						"depositor": "%s",
 						"amount": { "amount": "600", "denom": "bar" }
 					},
 					{
 						"depositor": "%s",
 						"amount": { "amount": "400", "denom": "bar" }
 					}
			]}}`,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
			check: func(ctx sdk.Context) {
				// Make sure the user's insurance fund is not updated
				userAddr, err := sdk.AccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Assert().NoError(err)

				insuranceFund, err := suite.k.GetUserInsuranceFundBalance(ctx, userAddr)
				suite.Assert().NoError(err)
				suite.Assert().Empty(insuranceFund)
			},
		},
		{
			name:           "multiple denoms in amount to deposit",
			transferAmount: sdk.NewInt64Coin("foo", 1000),
			sender:         "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			receiver:       authtypes.NewModuleAddress(types.ModuleName).String(),
			memo: fmt.Sprintf(`{
			"liquidvesting": {
				"amounts": [{
					"depositor": "%s",
					"amount": { "amount": "1000", "denom": "foo" }
				},
				{
					"depositor": "%s",
					"amount": { "amount": "1000", "denom": "bar" }
				}]
			}}`,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
			check: func(ctx sdk.Context) {
				// Make sure the user's insurance fund is not updated
				userAddr, err := sdk.AccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Assert().NoError(err)

				insuranceFund, err := suite.k.GetUserInsuranceFundBalance(ctx, userAddr)
				suite.Assert().NoError(err)
				suite.Assert().Empty(insuranceFund)
			},
		},
		{
			name:           "deposit more coins then received",
			transferAmount: sdk.NewInt64Coin("foo", 1000),
			sender:         "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			receiver:       authtypes.NewModuleAddress(types.ModuleName).String(),
			memo: fmt.Sprintf(`{
            "liquidvesting": {
                "amounts": [{
                    "depositor": "%s",
                    "amount": { "amount": "400", "denom": "foo" }
                },
                {
                    "depositor": "%s",
                    "amount": { "amount": "601", "denom": "foo" }
                }]
            }}`,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
			check: func(ctx sdk.Context) {
				// Make sure the user's insurance fund is not updated
				userAddr, err := sdk.AccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Assert().NoError(err)

				insuranceFund, err := suite.k.GetUserInsuranceFundBalance(ctx, userAddr)
				suite.Assert().NoError(err)
				suite.Assert().Empty(insuranceFund)
			},
		},
		{
			name:           "deposit less coins then received",
			transferAmount: sdk.NewInt64Coin("foo", 1000),
			sender:         "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			receiver:       authtypes.NewModuleAddress(types.ModuleName).String(),
			memo: fmt.Sprintf(`{
            "liquidvesting": {
                "amounts": [{
                    "depositor": "%s",
                    "amount": { "amount": "300", "denom": "foo" }
                },
                {
                    "depositor": "%s",
                    "amount": { "amount": "600", "denom": "foo" }
                }]
            }}`,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
			check: func(ctx sdk.Context) {
				// Make sure the user's insurance fund is not updated
				userAddr, err := sdk.AccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Assert().NoError(err)

				insuranceFund, err := suite.k.GetUserInsuranceFundBalance(ctx, userAddr)
				suite.Assert().NoError(err)
				suite.Assert().Empty(insuranceFund)
			},
		},
		{
			name:           "correct deposit",
			transferAmount: sdk.NewInt64Coin("foo", 1000),
			sender:         "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			receiver:       authtypes.NewModuleAddress(types.ModuleName).String(),
			memo: fmt.Sprintf(`{
            "liquidvesting": {
                "amounts": [{
                    "depositor": "%s",
                    "amount": { "amount": "600", "denom": "foo" }
                },
                {
                    "depositor": "%s",
                    "amount": { "amount": "400", "denom": "foo" }
                }]
            }}`,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			check: func(ctx sdk.Context) {
				// Make sure the first insurance fund is updated
				userAddr, err := sdk.AccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Assert().NoError(err)

				insuranceFund, err := suite.k.GetUserInsuranceFundBalance(ctx, userAddr)
				suite.Assert().NoError(err)
				suite.Assert().Equal("600ibc/EB7094899ACFB7A6F2A67DB084DEE2E9A83DEFAA5DEF92D9A9814FFD9FF673FA", insuranceFund.String())

				// Make sure the second insurance fund is updated
				userAddr, err = sdk.AccAddressFromBech32("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd")
				suite.Assert().NoError(err)

				insuranceFund, err = suite.k.GetUserInsuranceFundBalance(ctx, userAddr)
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
				DestinationChannel: "channel-0",
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
