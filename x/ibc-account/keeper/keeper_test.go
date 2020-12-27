package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	ibcacctesting "github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/testing"
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
)

type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibcacctesting.Coordinator

	chainA *ibcacctesting.TestChain
	chainB *ibcacctesting.TestChain

	queryClientA types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = ibcacctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibcacctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(ibcacctesting.GetChainID(1))

	queryHelperA := baseapp.NewQueryServerTestHelper(suite.chainA.GetContext(), suite.chainA.App.InterfaceRegistry())
	types.RegisterQueryServer(queryHelperA, suite.chainA.App.IBCAccountKeeper)
	suite.queryClientA = types.NewQueryClient(queryHelperA)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
