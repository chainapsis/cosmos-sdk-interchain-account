package keeper_test

import (
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/keeper"
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/types"
	connection "github.com/cosmos/cosmos-sdk/x/ibc/03-connection"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/24-host"
)

func (suite *KeeperTestSuite) TestCreateIAAccount() {
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
	suite.chainA.createConnection(testConnection, testConnection, testClientIDB, testClientIDA, connection.OPEN)
	suite.chainA.createChannel(testPort1, testChannel1, testPort2, testChannel2, channeltypes.OPEN, channeltypes.ORDERED, testConnection)

	initialSeq := uint64(1)
	suite.chainA.App.IBCKeeper.ChannelKeeper.SetNextSequenceSend(suite.chainA.GetContext(), testPort1, testChannel1, initialSeq)

	err = suite.chainA.App.IBCAccountKeeper.CreateInterchainAccount(suite.chainA.GetContext(), testPort1, testChannel1, testSalt)
	suite.Require().Nil(err, "could not request creating ia account")

	packet := suite.chainA.App.IBCKeeper.ChannelKeeper.GetPacketCommitment(suite.chainA.GetContext(), testPort1, testChannel1, initialSeq)
	suite.Require().Greater(len(packet), 0, "packet is empty")

	// TODO: verify a packet with proof
	err = suite.chainB.App.IBCAccountKeeper.RegisterIBCAccount(suite.chainB.GetContext(), testPort1, testChannel1, testSalt)
	suite.Require().Nil(err)

	acc := suite.chainB.App.AccountKeeper.GetAccount(suite.chainB.GetContext(), suite.chainB.App.IBCAccountKeeper.GenerateAddress(types.GetIdentifier(testPort1, testChannel1), testSalt))
	suite.Require().NotNil(acc, "ia account is not registered on counterparty chain")
}
