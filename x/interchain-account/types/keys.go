package types

import "fmt"

const (
	// ModuleName defines the IBC transfer name
	ModuleName = "interchainaccount"

	PortID = "interchainaccount"

	StoreKey  = ModuleName
	RouterKey = ModuleName

	// Key to store portID in our store
	PortKey = "portID"

	QuerierRoute = ModuleName
)

func GetIdentifier(portID, channelID string) string {
	return fmt.Sprintf("%s/%s/", portID, channelID)
}
