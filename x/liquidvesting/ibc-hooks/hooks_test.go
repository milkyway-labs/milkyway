package ibc_hooks_test

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/stretchr/testify/require"
)

func (suite *HooksTestSuite) TestHooks_IbcTriggers() {
	user1 := suite.userAddress("user1")
	user2 := suite.userAddress("user2")

	testCases := []struct {
		name           string
		transferAmount sdk.Coin
		sender         string
		receiver       string
		memo           string
		shouldErr      bool
		errorMessage   string
		check          func(sdk.Context)
	}{
		{
			name:           "empty memo",
			transferAmount: sdk.NewInt64Coin("foo", 1000),
			sender:         user1.String(),
			receiver:       user2.String(),
			memo:           "",
			shouldErr:      false,
		},
		{
			name:           "trigger by sending to a normal account",
			transferAmount: sdk.NewInt64Coin("foo", 1000),
			sender:         user1.String(),
			receiver:       user2.String(),
			memo: fmt.Sprintf(`{
			"liquidvesting": {
				"amounts": [{
					"depositor": "%s",
					"amount": { "amount": "1000", "denom": "foo" }
				}]
			}}`, user1.String()),
			shouldErr: true,
			errorMessage: fmt.Sprintf(
				"ibc hook error: the receiver should be the module address, got: %s, expected: %s",
				user2.String(), suite.k.ModuleAddress),
		},
		{
			name:           "transfer not received denom",
			transferAmount: sdk.NewInt64Coin("foo", 1000),
			sender:         user1.String(),
			receiver:       suite.k.ModuleAddress,
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
			]}}`, user1.String(), user2.String()),
			shouldErr:    true,
			errorMessage: "ibc hook error: amount received is not equal to the amounts to deposit in the users' insurance fund",
		},
		{
			name:           "multiple denoms in amount to deposit",
			transferAmount: sdk.NewInt64Coin("foo", 1000),
			sender:         user1.String(),
			receiver:       suite.k.ModuleAddress,
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
			}}`, user1.String(), user2.String()),
			shouldErr:    true,
			errorMessage: "ibc hook error: can't deposit multiple coins",
		},
		{
			name:           "deposit more coins then received",
			transferAmount: sdk.NewInt64Coin("foo", 1000),
			sender:         user1.String(),
			receiver:       suite.k.ModuleAddress,
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
            }}`, user1.String(), user2.String()),
			shouldErr:    true,
			errorMessage: "ibc hook error: amount received is not equal to the amounts to deposit in the users' insurance fund",
		},
		{
			name:           "deposit less coins then received",
			transferAmount: sdk.NewInt64Coin("foo", 1000),
			sender:         user1.String(),
			receiver:       suite.k.ModuleAddress,
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
            }}`, user1.String(), user2.String()),
			shouldErr:    true,
			errorMessage: "ibc hook error: amount received is not equal to the amounts to deposit in the users' insurance fund",
		},
		{
			name:           "correct deposit",
			transferAmount: sdk.NewInt64Coin("foo", 1000),
			sender:         user1.String(),
			receiver:       suite.k.ModuleAddress,
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
            }}`, user1.String(), user2.String()),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				addrInsuranceFund, err := suite.k.GetUserInsuranceFundBalance(ctx, user1)
				require.NoError(suite.T(), err)
				require.Equal(suite.T(), "600foo", addrInsuranceFund.String())
				addr2InsuranceFund, err := suite.k.GetUserInsuranceFundBalance(ctx, user2)
				require.NoError(suite.T(), err)
				require.Equal(suite.T(), "400foo", addr2InsuranceFund.String())
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			data := transfertypes.FungibleTokenPacketData{
				Denom:    tc.transferAmount.Denom,
				Amount:   tc.transferAmount.Amount.String(),
				Sender:   tc.sender,
				Receiver: tc.receiver,
				Memo:     tc.memo,
			}

			dataBz, err := json.Marshal(&data)
			require.NoError(suite.T(), err)

			relayer := suite.ak.GetModuleAddress("relayer")
			ack := suite.IBCHooksMiddleware.OnRecvPacket(suite.ctx, channeltypes.Packet{
				Data: dataBz,
			}, relayer)
			ack.Acknowledgement()

			if tc.shouldErr {
				require.False(suite.T(), ack.Success())
				castedAck := ack.(channeltypes.Acknowledgement)
				errorResponse := castedAck.Response.(*channeltypes.Acknowledgement_Error)
				require.Equal(suite.T(), tc.errorMessage, errorResponse.Error)

				if tc.check != nil {
					tc.check(suite.ctx)
				}
			} else {
				require.True(suite.T(), ack.Success())
			}
		})
	}
}