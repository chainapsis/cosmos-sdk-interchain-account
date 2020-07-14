## Status
In progress

## Context
This document outlines the design choices in the implementation of [ICS-027](https://github.com/cosmos/ics/tree/master/spec/ics-027-interchain-accounts) for the Cosmos-SDK. Interchain accounts enables an IBC connected chain(s) to: 1. Create an account on a counterparty chain that can be controlled by the logic of the account creator blockchain and 2. Send transaction messages from these interchain accounts.

This document introduces the IBC Account module that manages the accounts from other chains.

```go
type CounterpartyInfo struct {
	// This method used to marshal transaction for counterparty chain.
	SerializeTx func(data interface{}) ([]byte, error)
}

type Keeper struct {
	...
	counterpartyInfos map[string]CounterpartyInfo
    ...
	router types.Router
}
```

The most important parts of the IBC account keeper, as shown above, are the map of counterparty info and router. Because ICS-027 spec defines that the chain can send an arbitrary tx bytes to the counterparty chain, both chains should define the way that they process the caller's requests or make the tx bytes that the callee can process.

`CounterpartyInfo` contains the interface of serializing counterparty chain's tx bytes from any data, as well as the map of couterparty infos. It also contains the key such as chain-id and keeper use it to send packet. Therefore, it is necessary to know the exact counterparty chain being executed.

`SerializeCosmosTx(codec *codec.Codec)` provides a way to serialize the tx bytes from messages if the counterparty chain is based on cosmos-sdk.

The router is used to delegates the process of handling the message to a module. When a packet which requests running tx bytes is passed, it deserializes and gets the handler from the router to be passed to the handler. The keeper will check the result of all messages, and if any message returns an error, the entire transaction is aborted, and state changed rolled back.

```proto
enum Type {
    REGISTER = 0;
    RUNTX = 1;
}

message IBCAccountPacketData {
    Type type = 1;
    bytes data = 2;
}

message IBCAccountPacketAcknowledgement {
    Type type = 1;
    string chainID = 2;
    uint32 code = 3;
    string error = 4;
}
```

The example above shows the IBC packets that are used in ICS-027. `Type` indicates what action the packet is performing. When a `REGISTER` packet type is delivered, the counterparty chain will create an account with the address using the hash of {destPort}/{destChannel}/{packet.data}, assuming a duplicate prior account doesn't exist.

If the account is created successfully, it returns an acknowledgement packet to the origin chain with type `REGISTER` and code `0`. If else, it returns the acknowledgement packet with type `REGISTER` and code and result according to occured error.

When the packet of `RUNTX` type is delivered, the counterparty chain will deserialize the tx bytes (packet's data field) in a predefined way.

In this implementation for the Cosmos-SDK, it deserializes the tx bytes into slices of messages and gets the handler from the router and executes and checks the result like described above.

If the all messages are successful, it returns the acknowledgment packet to the chain with type `RUNTX` and code `0`. If else, it returns the acknowledgement packet with type `RUNTX` and the code and error of first failed message.

```go
type IBCAccountHooks interface {
	WillAccountCreate(ctx sdk.Context, chainID string, address sdk.AccAddress)
	DidAccountCreated(ctx sdk.Context, chainID string, address sdk.AccAddress)

	WillTxRun(ctx sdk.Context, chainID string, txHash []byte, data interface{})
	DidTxRun(ctx sdk.Context, chainID string, txHash []byte, data interface{})
}
```

The example above shows the hook for helping the developer that uses IBC account keeper.

Hook let the developer know the IBC account is expected to be created to the counterparty chain and when the IBC account has been successfully created on the counterparty chain.

Before sending the packet with a `IBCAccountPacketData` which type is `REGISTER`, `WillAccountCreate` is executed with the counterparty chain's chain-id and expected address. If the acknowledgement packet with the type `REGISTER` and code `0` is delivered, `DidAccountCreated` is executed with the counterparty chain's chain-id and address.

Before sending the packet with a `IBCAccountPacketData` which type is `RUNTX`, `WillTxRun` is executed with the counterparty chain's chain-id and virtual tx hash and requested data that is not serialized. If the acknowledgement packet with the type `RUNTX` and code `0` is delivered, `DidTxRun` is executed with the counterparty chain's chain-id and virtual tx hash and requested data that is not serialized. Virtual tx hash is used only for internal logic to distinguish the requested tx and it is computed by hashing the tx bytes and sequence of packet.
