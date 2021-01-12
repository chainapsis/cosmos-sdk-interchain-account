# Bootstrap an Interchain Accounts Testnet using Starport

Interchain Accounts currently experimentally supports [Starport v0.13.1](https://github.com/tendermint/starport/releases/tag/v0.13.1). However, it should be noted that stability is not guaranteed as the module import functionality is an experimental feature.

## Testing ICA with Starport

### 1. Start two chains with Starport

Open two [Gitpod](https://gitpod.io/#https://github.com/chainapsis/cosmos-sdk-interchain-account) workspaces (Create fresh workspace → create) to start two instances of Starport chains: `foo` and `bar`.

**Workspace 1**

```bash
starport app github.com/foo/foo --sdk-version stargate

cd foo
```

**Workspace 2**

```bash
starport app github.com/bar/bar --sdk-version stargate

cd bar
```

### 2. Import the IBC Account module to the two chains

**Workspace 1/2**

```bash
starport-ica module import
starport-ica module import mock
```

### 3. Start the two chains

**Workspace 1/2**

```bash
starport serve
```

### 4. Check the relayer information

After bootstrapping each chain, each workspace terminal will display the relayer information as follows. Note that the value displayed will be different from the example.

**Workspace 1/2 Console**

```bash
✨ Relayer info: eyJDaGFpbklEIjoiYmFyIiwiTW5lbW9uaWMiOiJmcm9zdCByYXpvciBoYWxmIGxhdW5kcnkgcHJvZml0IHdpc2UgdG9uZSBibHVzaCBzdXJnZSBrZWVwIHRvZ2V0aGVyIHNsaWNlIHlvdXRoIHRydXRoIGVubGlzdCBjdXBib2FyZCBhYnNvcmIgc2VlZCBzZXJpZXMgZG91YmxlIHZpbGxhZ2UgdG9uZ3VlIGZsYXNoIGdvcmlsbGEiLCJSUENBZGRyZXNzIjoiaHR0cHM6Ly8yNjY1Ny1jNzllNDk2ZC1kZDk4LTQ4MWQtOTlmZi1jZGQ4OTA2NWQ4MWIud3MtZXUwMS5naXRwb2QuaW86NDQzIn0
```

### 5. Add Chain Foo

Head over to Workspace 2 and open a new terminal window. Input the following code to add the IBC connected chain. Make sure to use the relayer information shown in step 4 instead of the value provided in the example.

**Workspace 2**

```bash
cd ../bar
starport chain add {workspace1-relayer-info}
```

### 6. Modify the configuration

At this stage, IBC transfers connections are open. However, as Starport currently only natively supports IBC transfers we need to change the relayer configuration.

Refer to the bar-foo IBC path by checking the `config.yaml` file as follows:

**Workspace 2**

```bash
vi ~/.relayer/config/config.yaml
```

Manually add a new path by appending the following example path information underneath the existing `bar-foo` path (Note that the `client-id` and `connection-id` doesn't have to be changed from the example below, but you must change the `channel-id` name):

```yaml
global:
  timeout: 10s
  lite-cache-size: 20
chains:
- key: testkey
  chain-id: bar
  rpc-addr: https://26657-b444783c-780f-4b19-ad73-aa02cf4309df.ws-us03.gitpod.io:443
  account-prefix: cosmos
  gas-adjustment: 1.5
  trusting-period: 336h
- key: testkey
  chain-id: foo
  rpc-addr: https://26657-e8783cc5-17eb-44b5-990c-584a9705271e.ws-us03.gitpod.io:443
  account-prefix: cosmos
  gas-adjustment: 1.5
  trusting-period: 336h
paths:
  bar-foo:
    src:
      chain-id: bar
      client-id: b4f6056a-1344-4e5d-9d20-93ea896dfb6c
      connection-id: f3cab589-db3e-4e19-8108-4ea4c39a32f7
      channel-id: test
      port-id: transfer
      order: unordered
      version: ics20-1
    dst:
      chain-id: foo
      client-id: 75284422-ec2a-447d-9764-c59fdc744e1a
      connection-id: e683f7bb-a49b-4c97-83a5-d342ad17d28e
      channel-id: test
      port-id: transfer
      order: unordered
      version: ics20-1
    strategy:
      type: naive
//append the bar-foo-ica path information as shown here
//but change the channel-id
  bar-foo-ica:
    src:
      chain-id: bar
      client-id: 1599fbea-43a0-4e8b-9c04-3f7e5cd11f94
      connection-id: 1f397527-e94e-4362-90cc-bc58e5aabffd
      channel-id: {your-channel-id-name}
      port-id: ibcaccount
      order: ordered
      version: ics27-1
    dst:
      chain-id: foo
      client-id: 15da2512-e2b3-4607-8355-74c994925cc9
      connection-id: 1f4c0fc3-0479-43a7-84c4-21d602c7ed98
      channel-id: {your-channel-id-name}
      port-id: ibcaccount
      order: ordered
      version: ics27-1
    strategy:
        type: naive
```

Link the paths:

```bash
rly tx link bar-foo-ica
```

Once the linking is successful, you can use the `mock` module to run some test transactions that use interchain accounts.

### 7. Running test transactions

First, you need to register an IBCAccount on chain `foo` that chain `bar` manages. Use the destination chain's `channel-id` as shown in `config.yaml`.

**Workspace 2**

```bash
bard tx ibcaccount register ibcaccount {dst.channel-id} test --from bar --absolute-timeouts --packet-timeout-height "0-1000000" --packet-timeout-timestamp 0
```

Now you can use the relayer to send over the packet information to the chain that the IBC account will be created on.

**Workspace 2**

```bash
rly tx relay-packets bar-foo-ica
rly tx relay-acknowledgements bar-foo-ica
```

Now let's check if the IBC account has been registered on chain `foo`.

Open a new terminal window in workspace 1

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

Check the balance of the interchain account on chain `foo`.

**Workspace 1**

```bash
food q bank balances {ibc_account_address}
```

Congratulations! You have now successfully sent a local transaction from an interchain account on chain `bar` through the interchain accounts IBC message.
