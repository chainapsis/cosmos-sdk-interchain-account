# Interchain Accounts Starport Documentation

Interchain Accounts currently experimentally supports [Starport v0.13.1](https://github.com/tendermint/starport/releases/tag/v0.13.1). However, it should be noted that stability is not guaranteed as the module import functionality is an experimental feature.

## Testing ICA with Starport

### 1. Start two chains with Starport

Open two [Gitpod](https://www.notion.so/Test-with-starport-d733ae61daeb44e8ad814af41befc1ad) workspaces to start two Starport chains `foo` and `bar`.

**Workspace 1**

```jsx
starport app github.com/foo/foo --sdk-version stargate

cd foo
```

**Workspace 2**

```jsx
starport app github.com/bar/bar --sdk-version stargate

cd bar
```

### 2. Import the IBC Account module to the two chains

**Workspace 1/2**

```jsx
starport-ica module import 
starport-ica module import mock
```

### 3. Start the two chains

**Workspace 1/2**

```jsx
starport serve
```

### 4. Check the relayer information

After bootstrapping each chain, each workspace terminal will display the relayer information as follows. Note that the value displayed will be different from the example.

**Workspace 1/2 Console**

```jsx
✨ Relayer info: eyJDaGFpbklEIjoiYmFyIiwiTW5lbW9uaWMiOiJmcm9zdCByYXpvciBoYWxmIGxhdW5kcnkgcHJvZml0IHdpc2UgdG9uZSBibHVzaCBzdXJnZSBrZWVwIHRvZ2V0aGVyIHNsaWNlIHlvdXRoIHRydXRoIGVubGlzdCBjdXBib2FyZCBhYnNvcmIgc2VlZCBzZXJpZXMgZG91YmxlIHZpbGxhZ2UgdG9uZ3VlIGZsYXNoIGdvcmlsbGEiLCJSUENBZGRyZXNzIjoiaHR0cHM6Ly8yNjY1Ny1jNzllNDk2ZC1kZDk4LTQ4MWQtOTlmZi1jZGQ4OTA2NWQ4MWIud3MtZXUwMS5naXRwb2QuaW86NDQzIn0
```

### 5. Add Chain Foo

Head over to Workspace 2 and open a new terminal window. Input the following code to add the IBC connected chain. Make sure to use the relayer information shown in step 4 instead of the value provided in the example.

```jsx
cd ../bar
starport chain add eyJDaGFpbklEIjoiYmFyIiwiTW5lbW9uaWMiOiJmcm9zdCByYXpvciBoYWxmIGxhdW5kcnkgcHJvZml0IHdpc2UgdG9uZSBibHVzaCBzdXJnZSBrZWVwIHRvZ2V0aGVyIHNsaWNlIHlvdXRoIHRydXRoIGVubGlzdCBjdXBib2FyZCBhYnNvcmIgc2VlZCBzZXJpZXMgZG91YmxlIHZpbGxhZ2UgdG9uZ3VlIGZsYXNoIGdvcmlsbGEiLCJSUENBZGRyZXNzIjoiaHR0cHM6Ly8yNjY1Ny1jNzllNDk2ZC1kZDk4LTQ4MWQtOTlmZi1jZGQ4OTA2NWQ4MWIud3MtZXUwMS5naXRwb2QuaW86NDQzIn0
```

### 6. Modify the configuration

At this stage, IBC transfers connections are open. However, as Starport currently only natively supports IBC transfers we need to change the relayer configuration.

Refer to the bar-foo IBC path by checking the `config.yaml` file as follows:

```jsx
vi ~/.relayer/config/config.yaml
```

Manually add a new path by adding the following information (Note that the `client-id` can be reused, but the `connection-id` and the `channel-id` must be different:

```jsx
bar-foo-ica:
    src:
        chain-id: bar
        client-id: 1599fbea-43a0-4e8b-9c04-3f7e5cd11f94
        connection-id: 1f397527-e94e-4362-90cc-bc58e5aabffd
        channel-id: 156b2f35-eabe-4141-8390-d3037133c0c3
        port-id: ibcaccount
        order: ordered
        version: ics27-1
    dst:
        chain-id: foo
        client-id: 15da2512-e2b3-4607-8355-74c994925cc9
        connection-id: 1f4c0fc3-0479-43a7-84c4-21d602c7ed98
        channel-id: 1568b97e-51f3-4cb6-af63-afa2404db1a5
        port-id: ibcaccount
        order: ordered
        version: ics27-1
    strategy:
        type: naive
```

Link the paths:

```jsx
rly tx link bar-foo-ica
```

Once the linking is successful, you can use the `mock` module to run some test transactions that use interchain accounts.

### 7. Running test transactions

First, you need to register an IBCAccount on chain `foo` that chain `bar` manages. Use the destination chain's `channel-id` as shown in `config.yaml`.

**Workspace 2**

```jsx
bard tx ibcaccount register ibcaccount {dst.channel-id} test --from bar --absolute-timeouts --packet-timeout-height "0-1000000" --packet-timeout-timestamp 0
```

Now you can use the relayer to send over the packet information to the chain that the IBC account will be created on.

**Workspace 2**

```jsx
rly tx relay-packets bar-foo-ica
rly tx relay-acknowledgements bar-foo-ica
```

Now let's check if the IBC account has been registered on chain `foo`.

Use the `src.channel-id` in the relayer's `config.yaml` for the next command.

**Workspace 1**

```bash
food q ibcaccount ibcaccount test ibcaccount {src.channel-id}
```

Send a small amount of tokens to the IBC account address.

**Workspace 1**

```bash
food tx bank send user2 {ibc_account_address} 100token
```

Check the balance of the IBC account.

**Workspace 1**

```bash
food q bank balances {ibc_account_address}
```

Now move to workspace 2 to use the `bar` chain to send a token from the interchain account on `foo` chain.

**Workspace 2**

```bash
bard tx ibcaccount send ibcaccount {src.channel-id} {ibc-account-address} {receiving-account-address} 50token --from bar --absolute-timeouts --packet-timeout-height "0-1000000" --packet-timeout-timestamp 0
```

Relay the packet through the relayer.

**Workspace 2**

```bash
rly tx relay-packets bar-foo-ica
rly tx relay-acknowledgements bar-foo-ica
```

Check the balance of the interchain account on chain `bar`.

```bash
food q bank balances {ibc_account_address}
```

Congratulations! You have now successfully sent a local transaction from an interchain account on chain `bar` through the interchain accounts IBC message.
