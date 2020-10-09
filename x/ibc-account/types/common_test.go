package types_test

import (
	"github.com/chainapsis/cosmos-sdk-interchain-account/simapp"
)

var (
	app                   = simapp.Setup(false)
	appCodec, legacyAmino = simapp.MakeCodecs()
)
