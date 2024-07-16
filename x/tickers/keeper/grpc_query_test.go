package keeper_test

import (
	"fmt"

	"github.com/milkyway-labs/milkyway/x/tickers/types"
)

func (s *KeeperTestSuite) TestQueryTicker() {
	_, err := s.msgServer.RegisterTicker(s.Ctx, types.NewMsgRegisterTicker(s.authority, "umilk", "MILK"))
	s.Require().NoError(err)

	for _, tc := range []struct {
		name        string
		req         *types.QueryTickerRequest
		expectedErr string
		postRun     func(resp *types.QueryTickerResponse)
	}{
		{
			"successful query",
			&types.QueryTickerRequest{Denom: "umilk"},
			"",
			func(resp *types.QueryTickerResponse) {
				s.Require().Equal("MILK", resp.Ticker)
			},
		},
		{
			"invalid denom",
			&types.QueryTickerRequest{Denom: "!@#$"},
			"rpc error: code = InvalidArgument desc = invalid denom: !@#$",
			nil,
		},
		{
			"ticker not registered",
			&types.QueryTickerRequest{Denom: "uatom"},
			"rpc error: code = NotFound desc = ticker for denom uatom not registered",
			nil,
		},
	} {
		s.Run(tc.name, func() {
			resp, err := s.queryServer.Ticker(s.Ctx, tc.req)
			if tc.expectedErr == "" {
				s.Require().NoError(err)
				tc.postRun(resp)
			} else {
				s.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}

func (s *KeeperTestSuite) TestQueryDenoms() {
	_, err := s.msgServer.RegisterTicker(s.Ctx, types.NewMsgRegisterTicker(s.authority, "umilk", "MILK"))
	s.Require().NoError(err)

	_, err = s.msgServer.RegisterTicker(s.Ctx, types.NewMsgRegisterTicker(s.authority, "umilk2", "MILK"))
	s.Require().NoError(err)

	for _, tc := range []struct {
		name        string
		req         *types.QueryDenomsRequest
		expectedErr string
		postRun     func(resp *types.QueryDenomsResponse)
	}{
		{
			"successful query",
			&types.QueryDenomsRequest{Ticker: "MILK"},
			"",
			func(resp *types.QueryDenomsResponse) {
				s.Require().Equal([]string{"umilk", "umilk2"}, resp.Denoms)
			},
		},
		{
			"denoms not found",
			&types.QueryDenomsRequest{Ticker: "ATOM"},
			"",
			func(resp *types.QueryDenomsResponse) {
				s.Require().Empty(resp.Denoms)
			},
		},
		{
			"invalid ticker",
			&types.QueryDenomsRequest{},
			"rpc error: code = InvalidArgument desc = invalid ticker: empty ticker",
			func(resp *types.QueryDenomsResponse) {
				fmt.Println(resp.Denoms)
			},
		},
	} {
		s.Run(tc.name, func() {
			resp, err := s.queryServer.Denoms(s.Ctx, tc.req)
			if tc.expectedErr == "" {
				s.Require().NoError(err)
				tc.postRun(resp)
			} else {
				s.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}
