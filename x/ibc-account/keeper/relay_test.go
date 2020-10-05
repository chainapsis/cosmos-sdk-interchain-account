package keeper_test

import (
	ibcaccountkeeper "github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/keeper"
	ibcaccountmocktypes "github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/testing/mock/types"
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/02-client/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/24-host"
	"github.com/cosmos/cosmos-sdk/x/ibc/exported"
)

func (suite *KeeperTestSuite) TestTryRegisterIBCAccount() {
	// init connection and channel between chain A and chain B
	_, _, connA, connB := suite.coordinator.SetupClientConnections(suite.chainA, suite.chainB, exported.Tendermint)
	channelA, channelB := suite.coordinator.CreateIBCAccountChannels(suite.chainA, suite.chainB, connA, connB, channeltypes.ORDERED)

	// assume that chain A try to register IBC Account to chain B
	msg := &ibcaccountmocktypes.MsgTryRegisterIBCAccount{
		SourcePort:    channelA.PortID,
		SourceChannel: channelA.ID,
		Salt:          []byte("test"),
		TimeoutHeight: clienttypes.Height{
			EpochNumber: 0,
			EpochHeight: 100,
		},
		TimeoutTimestamp: 0,
		Sender:           suite.chainA.SenderAccount.GetAddress(),
	}
	err := suite.coordinator.SendMsg(suite.chainA, suite.chainB, channelB.ClientID, msg)
	suite.Require().NoError(err) // message committed

	packetData := types.IBCAccountPacketData{
		Type: types.Type_REGISTER,
		Data: []byte("test"),
	}
	packet := channeltypes.NewPacket(packetData.GetBytes(), 1, channelA.PortID, channelA.ID, channelB.PortID, channelB.ID, clienttypes.NewHeight(0, 100), 0)

	// get proof of packet commitment from chainA
	packetKey := host.KeyPacketCommitment(packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())
	proof, proofHeight := suite.chainA.QueryProof(packetKey)

	recvMsg := channeltypes.NewMsgRecvPacket(packet, proof, proofHeight, suite.chainB.SenderAccount.GetAddress())
	err = suite.coordinator.SendMsg(suite.chainB, suite.chainA, channelA.ClientID, recvMsg)
	suite.Require().NoError(err) // message committed

	acc := suite.chainB.App.AccountKeeper.GetAccount(suite.chainB.GetContext(), suite.chainB.App.IBCAccountKeeper.GenerateAddress(types.GetIdentifier(channelB.PortID, channelB.ID), []byte("test")))
	suite.Require().NotNil(acc)
}

func (suite *KeeperTestSuite) TestRunTx() {
	// init connection and channel between chain A and chain B
	_, _, connA, connB := suite.coordinator.SetupClientConnections(suite.chainA, suite.chainB, exported.Tendermint)
	channelA, channelB := suite.coordinator.CreateIBCAccountChannels(suite.chainA, suite.chainB, connA, connB, channeltypes.ORDERED)

	func() {
		// assume that chain A try to register IBC Account to chain B
		msg := &ibcaccountmocktypes.MsgTryRegisterIBCAccount{
			SourcePort:    channelA.PortID,
			SourceChannel: channelA.ID,
			Salt:          []byte("test"),
			TimeoutHeight: clienttypes.Height{
				EpochNumber: 0,
				EpochHeight: 100,
			},
			TimeoutTimestamp: 0,
			Sender:           suite.chainA.SenderAccount.GetAddress(),
		}
		err := suite.coordinator.SendMsg(suite.chainA, suite.chainB, channelB.ClientID, msg)
		suite.Require().NoError(err) // message committed

		packetData := types.IBCAccountPacketData{
			Type: types.Type_REGISTER,
			Data: []byte("test"),
		}
		packet := channeltypes.NewPacket(packetData.GetBytes(), 1, channelA.PortID, channelA.ID, channelB.PortID, channelB.ID, clienttypes.NewHeight(0, 100), 0)

		// get proof of packet commitment from chainA
		packetKey := host.KeyPacketCommitment(packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())
		proof, proofHeight := suite.chainA.QueryProof(packetKey)

		recvMsg := channeltypes.NewMsgRecvPacket(packet, proof, proofHeight, suite.chainB.SenderAccount.GetAddress())
		err = suite.coordinator.SendMsg(suite.chainB, suite.chainA, channelA.ClientID, recvMsg)
		suite.Require().NoError(err) // message committed

		acc := suite.chainB.App.AccountKeeper.GetAccount(suite.chainB.GetContext(), suite.chainB.App.IBCAccountKeeper.GenerateAddress(types.GetIdentifier(channelB.PortID, channelB.ID), []byte("test")))
		suite.Require().NotNil(acc)
	}()

	acc := suite.chainB.App.AccountKeeper.GetAccount(suite.chainB.GetContext(), suite.chainB.App.IBCAccountKeeper.GenerateAddress(types.GetIdentifier(channelB.PortID, channelB.ID), []byte("test")))

	// mint the token to IBC account on Chain B
	err := suite.chainB.App.BankKeeper.AddCoins(suite.chainB.GetContext(), acc.GetAddress(), sdk.Coins{sdk.Coin{
		Denom:  "test",
		Amount: sdk.NewInt(10000),
	}})
	suite.Require().NoError(err)

	toAddress, err := sdk.AccAddressFromHex("0000000000000000000000000000000000000000")
	suite.Require().NoError(err)
	// try to run tx that sends the token from IBC account to other account
	msg := &ibcaccountmocktypes.MsgTryRunTxMsgSend{
		SourcePort:    channelA.PortID,
		SourceChannel: channelA.ID,
		TimeoutHeight: clienttypes.Height{
			EpochNumber: 0,
			EpochHeight: 100,
		},
		TimeoutTimestamp: 0,
		FromAddress:      acc.GetAddress(),
		ToAddress:        toAddress,
		Amount: sdk.Coins{sdk.Coin{
			Denom:  "test",
			Amount: sdk.NewInt(5000),
		}},
		Sender: suite.chainA.SenderAccount.GetAddress(),
	}
	err = suite.coordinator.SendMsg(suite.chainA, suite.chainB, channelB.ClientID, msg)
	suite.Require().NoError(err) // message committed

	txBytes, err := ibcaccountkeeper.SerializeCosmosTx(suite.chainB.App.AppCodec(), suite.chainB.App.InterfaceRegistry())([]sdk.Msg{
		banktypes.NewMsgSend(acc.GetAddress(), toAddress, sdk.Coins{sdk.Coin{
			Denom:  "test",
			Amount: sdk.NewInt(5000),
		}}),
	})
	packetData := types.IBCAccountPacketData{
		Type: types.Type_RUNTX,
		Data: txBytes,
	}
	packet := channeltypes.NewPacket(packetData.GetBytes(), 2, channelA.PortID, channelA.ID, channelB.PortID, channelB.ID, clienttypes.NewHeight(0, 100), 0)

	// get proof of packet commitment from chainA
	packetKey := host.KeyPacketCommitment(packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())
	proof, proofHeight := suite.chainA.QueryProof(packetKey)

	recvMsg := channeltypes.NewMsgRecvPacket(packet, proof, proofHeight, suite.chainB.SenderAccount.GetAddress())
	err = suite.coordinator.SendMsg(suite.chainB, suite.chainA, channelA.ClientID, recvMsg)
	suite.Require().NoError(err) // message committed

	// check that the balance has been transfered
	bal := suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), acc.GetAddress(), "test")
	suite.Equal("5000", bal.Amount.String())
	bal = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), toAddress, "test")
	suite.Equal("5000", bal.Amount.String())
}
