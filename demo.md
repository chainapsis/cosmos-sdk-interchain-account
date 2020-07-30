Build executables by ```make install```

Build relayer.
```
git clone https://github.com/chainapsis/relayer
cd relayer
make install
```

Initialize two local chains.
```
// Remember the mnemonic key.
democli keys add val
democli keys add etc

demod init --chain-id test-1 test-1 --home ~/.demo-test-1
sed -i -e 's/2665/1665/g' ~/.demo-test-1/config/config.toml
sed -i -e 's#localhost:6060#localhost:6061#g' ~/.demo-test-1/config/config.toml

demod add-genesis-account $(democli keys show val -a) 100000000000stake --home ~/.demo-test-1
demod add-genesis-account $(democli keys show etc -a) 100000000000stake --home ~/.demo-test-1
demod gentx --name val --amount 10000000000stake --home ~/.demo-test-1
demod collect-gentxs --home ~/.demo-test-1

demod init --chain-id test-2 test-2 --home ~/.demo-test-2
demod add-genesis-account $(democli keys show val -a) 100000000000stake --home ~/.demo-test-2
demod add-genesis-account $(democli keys show etc -a) 100000000000stake --home ~/.demo-test-2
demod gentx --name val --amount 10000000000stake --home ~/.demo-test-2
demod collect-gentxs --home ~/.demo-test-2
```

Run two chains on separate terminals.
```
demod start --pruning nothing --home ~/.demo-test-1
```
```
demod start --pruning nothing --home ~/.demo-test-2
```

Link ibc channel for two chains.
```
rly config init
// Config json file is on root folder of this repo.
rly chains add -f test-1.json
rly chains add -f test-2.json

rly keys restore test-1 relayer "{mnemonic of etc account which was made on initializing step}"
rly keys restore test-2 relayer "{mnemonic of etc account which was made on initializing step}"

rly lite init test-1 -f
rly lite init test-2 -f

rly pth gen test-1 interchainaccount test-2 interchainaccount ibcaccount

rly tx link ibcaccount
```

Start the relayer on separate terminal.
```
rly start ibcaccount
```

Test interchain account locally.
```
democli tx intertx register --from etc --chain-id test-1 --source-port interchainaccount --source-channel {find the src channel id on ~/.relayer/config/config.toml} --node tcp://localhost:16657

// Wait until relayer relays packet.

// Get the address of interchain account.
ibcaccount=$(democli query intertx ibcaccount {address of etc key} interchainaccount {find the src channel id on ~/.relayer/config/config.toml} --chain-id test-1 --node tcp://localhost:16657)
echo $ibcaccount

// Check the interchain account's balance on test-2 chain. It should be empty.
democli q bank balances $ibcaccount --chain-id test-2

// Send some assets to $ibcaccount.
democli tx send etc $ibcaccount 1000stake --chain-id test-2

// Test sending assets on intetchain account via ibc.
democli tx intertx send cosmos-sdk {other address} 100stake --from etc --chain-id test-1 --source-port interchainaccount --source-channel {find the src channel id on ~/.relayer/config/config.toml} --node tcp://localhost:16657

// Wait until relayer relays packet.

// And check if the balance has been changed.
democli q bank balances $ibcaccount --chain-id test-2
democli q bank balances {other address} --chain-id test-2
```
