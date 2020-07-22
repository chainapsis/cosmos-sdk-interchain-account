package keeper_test

import (
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/keeper"
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	connectiontypes "github.com/cosmos/cosmos-sdk/x/ibc/03-connection/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/24-host"
	"math"
)

func (suite *KeeperTestSuite) TestCreateIBCAccount() {
	suite.initChannelAtoB()

	err := suite.chainA.App.IBCAccountKeeper.TryRegisterIBCAccount(suite.chainA.GetContext(), testPort1, testChannel1, testSalt)
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
		math.MaxUint64,
		0,
	)
	suite.Require().Equal(packetCommitment, channeltypes.CommitPacket(packet))

	err = suite.chainB.App.IBCAccountKeeper.OnRecvPacket(suite.chainB.GetContext(), packet)
	suite.Require().Nil(err)

	acc := suite.chainB.App.AccountKeeper.GetAccount(suite.chainB.GetContext(), suite.chainB.App.IBCAccountKeeper.GenerateAddress(types.GetIdentifier(testPort2, testChannel2), testSalt))
	suite.Require().NotNil(acc, "ibc account is not registered on counterparty chain")
}

func (suite *KeeperTestSuite) TestRunTx() {
	suite.initChannelAtoB()

	err := suite.chainA.App.IBCAccountKeeper.TryRegisterIBCAccount(suite.chainA.GetContext(), testPort1, testChannel1, testSalt)
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
		math.MaxUint64,
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
	_, err = suite.chainB.App.BankKeeper.AddCoins(suite.chainB.GetContext(), acc.GetAddress(), mint)
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
	_, err = suite.chainA.App.IBCAccountKeeper.TryRunTx(suite.chainA.GetContext(), testPort1, testChannel1, testClientIDB, sendMsg)
	suite.Require().Nil(err)

	packetCommitment = suite.chainA.App.IBCKeeper.ChannelKeeper.GetPacketCommitment(suite.chainA.GetContext(), testPort1, testChannel1, 2)
	suite.Require().Greater(len(packetCommitment), 0, "packet commitment is empty")

	packetTxBytes, err := keeper.SerializeCosmosTx(suite.chainB.App.Codec())(sendMsg)
	suite.Require().Nil(err)
	packet = channeltypes.NewPacket(
		types.IBCAccountPacketData{Type: types.Type_RUNTX, Data: packetTxBytes}.GetBytes(),
		2,
		testPort1,
		testChannel1,
		testPort2,
		testChannel2,
		math.MaxUint64,
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

	_, err = suite.chainA.App.IBCAccountKeeper.TryRunTx(suite.chainA.GetContext(), testPort1, testChannel1, testClientIDB, sendMsg)
	suite.Require().Nil(err)

	packetCommitment = suite.chainA.App.IBCKeeper.ChannelKeeper.GetPacketCommitment(suite.chainA.GetContext(), testPort1, testChannel1, 3)
	suite.Require().Greater(len(packetCommitment), 0, "packet commitment is empty")

	packetTxBytes, err = keeper.SerializeCosmosTx(suite.chainB.App.Codec())(sendMsg)
	suite.Require().Nil(err)
	packet = channeltypes.NewPacket(
		types.IBCAccountPacketData{Type: types.Type_RUNTX, Data: packetTxBytes}.GetBytes(),
		2,
		testPort1,
		testChannel1,
		testPort2,
		testChannel2,
		math.MaxUint64,
		0,
	)
	suite.Require().Equal(packetCommitment, channeltypes.CommitPacket(packet))

	// Should fail if msg is sent from account not created by ibc account module.
	err = suite.chainB.App.IBCAccountKeeper.OnRecvPacket(suite.chainB.GetContext(), packet)
	suite.Require().NotNil(err)
}

func (suite *KeeperTestSuite) initChannelAtoB() {
	capName := host.ChannelCapabilityPath(testPort1, testChannel1)

	// Add counterparty info.
	suite.chainA.App.IBCAccountKeeper.AddCounterpartyInfo(testClientIDB, keeper.CounterpartyInfo{
		SerializeTx: keeper.SerializeCosmosTx(suite.chainB.App.Codec()),
	})

	// create channel capability from ibc scoped keeper and claim with ia scoped keeper
	cap, err := suite.chainA.App.ScopedIBCKeeper.NewCapability(suite.chainA.GetContext(), capName)
	suite.Require().Nil(err, "could not create capability")
	err = suite.chainA.App.ScopedIBCAccountKeeper.ClaimCapability(suite.chainA.GetContext(), cap, capName)
	suite.Require().Nil(err, "interchainaccount module could not claim capability")

	// create client, and open conn/channel
	err = suite.chainA.CreateClient(suite.chainB)
	suite.Require().Nil(err, "could not create client")
	suite.chainA.createConnection(testConnection, testConnection, testClientIDB, testClientIDA, connectiontypes.OPEN)
	suite.chainA.createChannel(testPort1, testChannel1, testPort2, testChannel2, channeltypes.OPEN, channeltypes.ORDERED, testConnection)

	initialSeq := uint64(1)
	suite.chainA.App.IBCKeeper.ChannelKeeper.SetNextSequenceSend(suite.chainA.GetContext(), testPort1, testChannel1, initialSeq)
}
