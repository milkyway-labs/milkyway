package keeper_test

import (
	"github.com/milkyway-labs/milkyway/x/tickers/types"
)

func (s *KeeperTestSuite) TestQueryAssets() {
	_, err := s.msgServer.RegisterAsset(s.Ctx, &types.MsgRegisterAsset{
		Authority: s.authority,
		Asset:     types.NewAsset("umilk", "MILK", 6),
	})
	s.Require().NoError(err)

	_, err = s.msgServer.RegisterAsset(s.Ctx, &types.MsgRegisterAsset{
		Authority: s.authority,
		Asset:     types.NewAsset("umilk2", "MILK", 6),
	})
	s.Require().NoError(err)

	_, err = s.msgServer.RegisterAsset(s.Ctx, &types.MsgRegisterAsset{
		Authority: s.authority,
		Asset:     types.NewAsset("uatom", "ATOM", 6),
	})
	s.Require().NoError(err)

	testCases := []struct {
		name        string
		req         *types.QueryAssetsRequest
		expectedErr string
		postRun     func(resp *types.QueryAssetsResponse)
	}{
		{
			"successful query",
			&types.QueryAssetsRequest{},
			"",
			func(resp *types.QueryAssetsResponse) {
				s.Require().Equal([]types.Asset{
					types.NewAsset("uatom", "ATOM", 6),
					types.NewAsset("umilk", "MILK", 6),
					types.NewAsset("umilk2", "MILK", 6),
				}, resp.Assets)
			},
		},
		{
			"successful query with ticker",
			&types.QueryAssetsRequest{Ticker: "MILK"},
			"",
			func(resp *types.QueryAssetsResponse) {
				s.Require().Equal([]types.Asset{
					types.NewAsset("umilk", "MILK", 6),
					types.NewAsset("umilk2", "MILK", 6),
				}, resp.Assets)
			},
		},
		{
			"successful query with ticker #2",
			&types.QueryAssetsRequest{Ticker: "ATOM"},
			"",
			func(resp *types.QueryAssetsResponse) {
				s.Require().Equal([]types.Asset{types.NewAsset("uatom", "ATOM", 6)}, resp.Assets)
			},
		},
		{
			"invalid ticker",
			&types.QueryAssetsRequest{Ticker: "!@#$"},
			"rpc error: code = InvalidArgument desc = bad ticker format: !@#$",
			nil,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.queryServer.Assets(s.Ctx, tc.req)
			if tc.expectedErr == "" {
				s.Require().NoError(err)
				tc.postRun(resp)
			} else {
				s.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}

func (s *KeeperTestSuite) TestQueryAsset() {
	_, err := s.msgServer.RegisterAsset(s.Ctx, &types.MsgRegisterAsset{
		Authority: s.authority,
		Asset:     types.NewAsset("umilk", "MILK", 6),
	})
	s.Require().NoError(err)

	testCases := []struct {
		name        string
		req         *types.QueryAssetRequest
		expectedErr string
		postRun     func(resp *types.QueryAssetResponse)
	}{
		{
			"successful query",
			&types.QueryAssetRequest{Denom: "umilk"},
			"",
			func(resp *types.QueryAssetResponse) {
				s.Require().Equal(types.NewAsset("umilk", "MILK", 6), resp.Asset)
			},
		},
		{
			"invalid denom",
			&types.QueryAssetRequest{Denom: "!@#$"},
			"rpc error: code = InvalidArgument desc = invalid denom: !@#$",
			nil,
		},
		{
			"ticker not registered",
			&types.QueryAssetRequest{Denom: "uatom"},
			"rpc error: code = NotFound desc = asset for denom uatom not registered",
			nil,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.queryServer.Asset(s.Ctx, tc.req)
			if tc.expectedErr == "" {
				s.Require().NoError(err)
				tc.postRun(resp)
			} else {
				s.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}
