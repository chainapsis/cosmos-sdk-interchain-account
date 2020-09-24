package keeper_test

import (
	ibcaccountmocktypes "github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/testing/mock/types"
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/types"
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

/*
func (suite *KeeperTestSuite) TestRunTx() {
	suite.initChannelAtoB()

	err := suite.chainA.App.IBCAccountKeeper.TryRegisterIBCAccount(suite.chainA.GetContext(), testPort1, testChannel1, testSalt, clienttypes.Height{
		EpochNumber: 0,
		EpochHeight: 100,
	})
	suite.Require().Nil(err, "could not request creating ia account")

	packetCommitment := suite.chainA.App.IBCKeeper.ChannelKeeper.GetPacketCommitment(suite.chainA.GetContext(), testPort1, testChannel1, 1)
	suite.Require().Greater(len(packetCommitment), 0, "packet commitment is empty")

	packet := channeltypes.NewPacket(
		types.IBCAccountPacketData{Type: types.Type_REGISTER, Data: []byte(testSalt)}.GetBytes(),
		1,
		testPort1,
		testChannel1,
		testPort2,
		testChannel2,
		clienttypes.Height{
			EpochNumber: 0,
			EpochHeight: 100,
		},
		0,
	)
	suite.Require().Equal(packetCommitment, channeltypes.CommitPacket(packet))

	err = suite.chainB.App.IBCAccountKeeper.OnRecvPacket(suite.chainB.GetContext(), packet)
	suite.Require().Nil(err)

	acc := suite.chainB.App.AccountKeeper.GetAccount(suite.chainB.GetContext(), suite.chainB.App.IBCAccountKeeper.GenerateAddress(types.GetIdentifier(testPort2, testChannel2), testSalt))
	suite.Require().NotNil(acc, "ibc account is not registered on counterparty chain")

	// Tokens to mint to ibc account.
	mint := sdk.Coins{
		sdk.Coin{
			Denom:  "test",
			Amount: sdk.NewInt(1000),
		},
	}

	// Bal should be empty.
	bal := suite.chainB.App.BankKeeper.GetAllBalances(suite.chainB.GetContext(), acc.GetAddress())
	suite.Require().Equal(sdk.Coins{}, bal)

	// Mint tokens.
	err = suite.chainB.App.BankKeeper.AddCoins(suite.chainB.GetContext(), acc.GetAddress(), mint)
	suite.Require().Nil(err)

	bal = suite.chainB.App.BankKeeper.GetAllBalances(suite.chainB.GetContext(), acc.GetAddress())
	suite.Require().Equal(mint, bal)

	testAddress := suite.chainB.App.IBCAccountKeeper.GenerateAddress(types.GetIdentifier(testPort2, "otherchannel"), testSalt)
	sendMsg := banktypes.NewMsgSend(acc.GetAddress(), testAddress, sdk.Coins{
		sdk.Coin{
			Denom:  "test",
			Amount: sdk.NewInt(500),
		},
	})
	_, err = suite.chainA.App.IBCAccountKeeper.TryRunTx(suite.chainA.GetContext(), testPort1, testChannel1, testClientIDB, sendMsg, clienttypes.Height{
		EpochNumber: 0,
		EpochHeight: 100,
	})
	suite.Require().Nil(err)

	packetCommitment = suite.chainA.App.IBCKeeper.ChannelKeeper.GetPacketCommitment(suite.chainA.GetContext(), testPort1, testChannel1, 2)
	suite.Require().Greater(len(packetCommitment), 0, "packet commitment is empty")

	packetTxBytes, err := keeper.SerializeCosmosTx(suite.chainB.App.AppCodec(), suite.chainB.App.InterfaceRegistry())(sendMsg)
	suite.Require().Nil(err)
	packet = channeltypes.NewPacket(
		types.IBCAccountPacketData{Type: types.Type_RUNTX, Data: packetTxBytes}.GetBytes(),
		2,
		testPort1,
		testChannel1,
		testPort2,
		testChannel2,
		clienttypes.Height{
			EpochNumber: 0,
			EpochHeight: 100,
		},
		0,
	)
	suite.Require().Equal(packetCommitment, channeltypes.CommitPacket(packet))

	err = suite.chainB.App.IBCAccountKeeper.OnRecvPacket(suite.chainB.GetContext(), packet)
	suite.Require().Nil(err)

	// Bal should have been transfered.
	bal = suite.chainB.App.BankKeeper.GetAllBalances(suite.chainB.GetContext(), acc.GetAddress())
	suite.Require().Equal(sdk.Coins{
		sdk.Coin{
			Denom:  "test",
			Amount: sdk.NewInt(500),
		},
	}, bal)

	bal = suite.chainB.App.BankKeeper.GetAllBalances(suite.chainB.GetContext(), testAddress)
	suite.Require().Equal(sdk.Coins{
		sdk.Coin{
			Denom:  "test",
			Amount: sdk.NewInt(500),
		},
	}, bal)

	// Test the case that msg is sent from not created by ibc account module.
	sendMsg = banktypes.NewMsgSend(testAddress, acc.GetAddress(), sdk.Coins{
		sdk.Coin{
			Denom:  "test",
			Amount: sdk.NewInt(500),
		},
	})

	_, err = suite.chainA.App.IBCAccountKeeper.TryRunTx(suite.chainA.GetContext(), testPort1, testChannel1, testClientIDB, sendMsg, clienttypes.Height{
		EpochNumber: 0,
		EpochHeight: 100,
	})
	suite.Require().Nil(err)

	packetCommitment = suite.chainA.App.IBCKeeper.ChannelKeeper.GetPacketCommitment(suite.chainA.GetContext(), testPort1, testChannel1, 3)
	suite.Require().Greater(len(packetCommitment), 0, "packet commitment is empty")

	packetTxBytes, err = keeper.SerializeCosmosTx(suite.chainB.App.AppCodec(), suite.chainB.App.InterfaceRegistry())(sendMsg)
	suite.Require().Nil(err)
	packet = channeltypes.NewPacket(
		types.IBCAccountPacketData{Type: types.Type_RUNTX, Data: packetTxBytes}.GetBytes(),
		2,
		testPort1,
		testChannel1,
		testPort2,
		testChannel2,
		clienttypes.Height{
			EpochNumber: 0,
			EpochHeight: 100,
		},
		0,
	)
	suite.Require().Equal(packetCommitment, channeltypes.CommitPacket(packet))

	// Should fail if msg is sent from account not created by ibc account module.
	err = suite.chainB.App.IBCAccountKeeper.OnRecvPacket(suite.chainB.GetContext(), packet)
	suite.Require().NotNil(err)
}
*/
