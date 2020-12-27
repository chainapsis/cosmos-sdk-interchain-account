package keeper_test

import (
	"fmt"

	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"

	"github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
)

func (suite *KeeperTestSuite) TestQueryIBCAccount() {
	var (
		req *types.QueryIBCAccountRequest
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryIBCAccountRequest{
					Address: "",
				}
			},
			false,
		},
		{
			"invalid bech32 address",
			func() {
				req = &types.QueryIBCAccountRequest{
					Address: "cosmos1ntck6f6534u630q87jpamettes6shwgddag761",
				}
			},
			false,
		},
		{
			"unexist address",
			func() {
				req = &types.QueryIBCAccountRequest{
					Address: "cosmos1ntck6f6534u630q87jpamettes6shwgddag769",
				}
			},
			false,
		},
		{
			"success",
			func() {
				packetData := types.IBCAccountPacketData{
					Type: types.Type_REGISTER,
					Data: []byte{},
				}

				packet := channeltypes.Packet{
					Sequence:           0,
					SourcePort:         "sp",
					SourceChannel:      "sc",
					DestinationPort:    "dp",
					DestinationChannel: "dc",
					Data:               packetData.GetBytes(),
					TimeoutHeight:      clienttypes.Height{},
					TimeoutTimestamp:   0,
				}

				err := suite.chainA.App.IBCAccountKeeper.OnRecvPacket(suite.chainA.GetContext(), packet)
				if err != nil {
					panic(err)
				}

				address := suite.chainA.App.IBCAccountKeeper.GenerateAddress(types.GetIdentifier("dp", "dc"), []byte{})

				req = &types.QueryIBCAccountRequest{
					Address: sdk.AccAddress(address).String(),
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()
			ctx := sdk.WrapSDKContext(suite.chainA.GetContext())

			res, err := suite.queryClientA.IBCAccount(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
