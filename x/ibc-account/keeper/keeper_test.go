package keeper_test

import (
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/02-client/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/chainapsis/cosmos-sdk-interchain-account/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	connectiontypes "github.com/cosmos/cosmos-sdk/x/ibc/03-connection/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	ibctmtypes "github.com/cosmos/cosmos-sdk/x/ibc/07-tendermint/types"
	commitmenttypes "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// define constants used for testing
const (
	testClientIDA = "testclientIDA"
	testClientIDB = "testClientIDb"

	testConnection = "testconnectionatob"
	testPort1      = "interchainaccount"
	testPort2      = "testportid"
	testChannel1   = "firstchannel"
	testChannel2   = "secondchannel"

	testSalt = "test"

	trustingPeriod time.Duration = time.Hour * 24 * 7 * 2
	ubdPeriod      time.Duration = time.Hour * 24 * 7 * 3
	maxClockDrift  time.Duration = time.Second * 10
)

type KeeperTestSuite struct {
	suite.Suite

	chainA *TestChain
	chainB *TestChain
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.chainA = NewTestChain(testClientIDA)
	suite.chainB = NewTestChain(testClientIDB)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type TestChain struct {
	ClientID string
	App      *simapp.SimApp
	Header   ibctmtypes.Header
	Vals     *tmtypes.ValidatorSet
	Signers  []tmtypes.PrivValidator
}

func NewTestChain(clientID string) *TestChain {
	privVar := tmtypes.NewMockPV()
	pubKey, err := privVar.GetPubKey()
	if err != nil {
		panic(err)
	}

	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})
	signers := []tmtypes.PrivValidator{privVar}
	now := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)

	header := ibctmtypes.CreateTestHeader(clientID, clienttypes.Height{
		EpochNumber: 0,
		EpochHeight: 1,
	}, clienttypes.Height{
		EpochNumber: 0,
		EpochHeight: 1,
	}, now, valSet, valSet, signers)

	return &TestChain{
		ClientID: clientID,
		App:      simapp.Setup(false),
		Header:   *header,
		Vals:     valSet,
		Signers:  signers,
	}
}

// Creates simple context for testing purposes
func (chain *TestChain) GetContext() sdk.Context {
	return chain.App.BaseApp.NewContext(false, tmproto.Header{ChainID: chain.Header.SignedHeader.Header.ChainID, Height: chain.Header.SignedHeader.Header.Height})
}

func (chain *TestChain) createConnection(
	connID, counterpartyConnID, clientID, counterpartyClientID string,
	state connectiontypes.State,
) connectiontypes.ConnectionEnd {
	counterparty := connectiontypes.NewCounterparty(counterpartyClientID, counterpartyConnID, commitmenttypes.NewMerklePrefix(chain.App.IBCKeeper.ConnectionKeeper.GetCommitmentPrefix().Bytes()))
	connection := connectiontypes.ConnectionEnd{
		State:        state,
		ClientId:     clientID,
		Counterparty: counterparty,
		Versions:     connectiontypes.GetCompatibleEncodedVersions(),
	}
	ctx := chain.GetContext()
	chain.App.IBCKeeper.ConnectionKeeper.SetConnection(ctx, connID, connection)
	return connection
}

// createClient will create a client for clientChain on targetChain
func (chain *TestChain) CreateClient(client *TestChain) error {
	client.Header = nextHeader(client)
	// Commit and create a new block on appTarget to get a fresh CommitID
	client.App.Commit()
	commitID := client.App.LastCommitID()
	client.App.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: client.Header.SignedHeader.Header.Height, Time: client.Header.GetTime()}})

	// Set HistoricalInfo on client chain after Commit
	ctxClient := client.GetContext()
	validator := stakingtypes.NewValidator(
		sdk.ValAddress(client.Vals.Validators[0].Address), client.Vals.Validators[0].PubKey, stakingtypes.Description{},
	)
	validator.Status = sdk.Bonded
	validator.Tokens = sdk.NewInt(1000000) // get one voting power
	validators := []stakingtypes.Validator{validator}
	histInfo := stakingtypes.HistoricalInfo{
		Header: tmproto.Header{
			AppHash: commitID.Hash,
		},
		Valset: validators,
	}
	client.App.StakingKeeper.SetHistoricalInfo(ctxClient, client.Header.SignedHeader.Header.Height, histInfo)

	// Create target ctx
	ctxTarget := chain.GetContext()

	// create client
	clientState := ibctmtypes.NewClientState(client.ClientID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, clienttypes.Height{
		EpochNumber: 0,
		EpochHeight: client.Header.GetHeight().GetEpochHeight(),
	}, commitmenttypes.GetSDKSpecs(), false, false)
	_, err := chain.App.IBCKeeper.ClientKeeper.CreateClient(ctxTarget, client.ClientID, clientState, client.Header.ConsensusState())
	if err != nil {
		return err
	}
	return nil

	// _, _, err := simapp.SignCheckDeliver(
	// 	suite.T(),
	// 	suite.cdc,
	// 	suite.app.BaseApp,
	// 	ctx.BlockHeader(),
	// 	[]sdk.Msg{clienttypes.NewMsgCreateClient(clientID, clientexported.ClientTypeTendermint, consState, accountAddress)},
	// 	[]uint64{baseAccount.GetAccountNumber()},
	// 	[]uint64{baseAccount.GetSequence()},
	// 	true, true, accountPrivKey,
	// )
}

func (chain *TestChain) createChannel(
	portID, channelID, counterpartyPortID, counterpartyChannelID string,
	state channeltypes.State, order channeltypes.Order, connectionID string,
) channeltypes.Channel {
	counterparty := channeltypes.NewCounterparty(counterpartyPortID, counterpartyChannelID)
	channel := channeltypes.NewChannel(state, order, counterparty,
		[]string{connectionID}, "1.0",
	)
	ctx := chain.GetContext()
	chain.App.IBCKeeper.ChannelKeeper.SetChannel(ctx, portID, channelID, channel)
	return channel
}

func nextHeader(chain *TestChain) ibctmtypes.Header {
	return *ibctmtypes.CreateTestHeader(
		chain.Header.SignedHeader.Header.ChainID,
		clienttypes.Height{
			EpochNumber: 0,
			EpochHeight: uint64(chain.Header.SignedHeader.Header.Height+1),
		},
		clienttypes.Height{
			EpochNumber: 0,
			EpochHeight: uint64(chain.Header.SignedHeader.Header.Height),
		},
		chain.Header.GetTime().Add(time.Minute),
		chain.Vals,
		chain.Vals,
		chain.Signers,
	)
}
