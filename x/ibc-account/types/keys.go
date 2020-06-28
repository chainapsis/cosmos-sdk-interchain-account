package types

import "fmt"

const (
	// ModuleName defines the IBC transfer name
	ModuleName = "interchainaccount"

	// Version defines the current version the IBC tranfer
	// module supports
	Version = "ics27-1"

	PortID = "interchainaccount"

	StoreKey  = ModuleName
	RouterKey = ModuleName

	// Key to store portID in our store
	PortKey = "portID"

	QuerierRoute = ModuleName
)

var (
	KeyPrefixRegisteredAccount = []byte("register")
)

func GetIdentifier(portID, channelID string) string {
	return fmt.Sprintf("%s/%s/", portID, channelID)
}
