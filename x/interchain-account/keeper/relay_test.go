package keeper_test

import (
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/interchain-account/keeper"
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/interchain-account/types"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
)

func (suite *KeeperTestSuite) TestCreateIAAccount() {
	capName := ibctypes.ChannelCapabilityPath(testPort1, testChannel1)

	// reset
	suite.SetupTest()

	// Add counterparty info.
	suite.chainA.App.InterchainAccountKeeper.AddCounterpartyInfo(testClientIDB, keeper.CounterpartyInfo{
		CounterpartyTxCdc: suite.chainB.App.Codec(),
	})

	// create channel capability from ibc scoped keeper and claim with ia scoped keeper
	cap, err := suite.chainA.App.ScopedIBCKeeper.NewCapability(suite.chainA.GetContext(), capName)
	suite.Require().Nil(err, "could not create capability")
	err = suite.chainA.App.ScopedInterchainAccountKeeper.ClaimCapability(suite.chainA.GetContext(), cap, capName)
	suite.Require().Nil(err, "interchainaccount module could not claim capability")

	// create client, and open conn/channel
	err = suite.chainA.CreateClient(suite.chainB)
	suite.Require().Nil(err, "could not create client")
	suite.chainA.createConnection(testConnection, testConnection, testClientIDB, testClientIDA, ibctypes.OPEN)
	suite.chainA.createChannel(testPort1, testChannel1, testPort2, testChannel2, ibctypes.OPEN, ibctypes.ORDERED, testConnection)

	initialSeq := uint64(1)
	suite.chainA.App.IBCKeeper.ChannelKeeper.SetNextSequenceSend(suite.chainA.GetContext(), testPort1, testChannel1, initialSeq)

	err = suite.chainA.App.InterchainAccountKeeper.CreateInterchainAccount(suite.chainA.GetContext(), testPort1, testChannel1, testSalt)
	suite.Require().Nil(err, "could not request creating ia account")

	packet := suite.chainA.App.IBCKeeper.ChannelKeeper.GetPacketCommitment(suite.chainA.GetContext(), testPort1, testChannel1, initialSeq)
	suite.Require().Greater(len(packet), 0, "packet is empty")

	// TODO: verify a packet with proof
	err = suite.chainB.App.InterchainAccountKeeper.RegisterIBCAccount(suite.chainB.GetContext(), testPort1, testChannel1, testSalt)
	suite.Require().Nil(err)

	acc := suite.chainB.App.AccountKeeper.GetAccount(suite.chainB.GetContext(), suite.chainB.App.InterchainAccountKeeper.GenerateAddress(types.GetIdentifier(testPort1, testChannel1), testSalt))
	suite.Require().NotNil(acc, "ia account is not registered on counterparty chain")
}
