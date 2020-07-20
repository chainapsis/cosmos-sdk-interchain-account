package types

type GenesisState struct {
	PortID string `json:"portid" yaml:"portid"`
}

func DefaultGenesis() GenesisState {
	return GenesisState{
		PortID: PortID,
	}
}
