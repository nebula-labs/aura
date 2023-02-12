#!/bin/bash

KEY="test"
CHAINID="aura-1"
KEYRING="test"
MONIKER="localtestnet"
KEYALGO="secp256k1"
LOGLEVEL="info"

# retrieve all args
WILL_RECOVER=0
WILL_INSTALL=0
WILL_CONTINUE=0
# $# is to check number of arguments
if [ $# -gt 0 ];
then
    # $@ is for getting list of arguments
    for arg in "$@"; do
        case $arg in
        --recover)
            WILL_RECOVER=1
            shift
            ;;
        --install)
            WILL_INSTALL=1
            shift
            ;;
        --continue)
            WILL_CONTINUE=1
            shift
            ;;
        *)
            printf >&2 "wrong argument somewhere"; exit 1;
            ;;
        esac
    done
fi

# continue running if everything is configured
if [ $WILL_CONTINUE -eq 1 ];
then
    # Start the node (remove the --pruning=nothing flag if historical queries are not needed)
    aurad start --pruning=nothing --log_level $LOGLEVEL --minimum-gas-prices=0.0001uaura --p2p.laddr tcp://0.0.0.0:2280 --rpc.laddr tcp://0.0.0.0:2281 --grpc.address 0.0.0.0:2282 --grpc-web.address 0.0.0.0:2283
    exit 1;
fi

# validate dependencies are installed
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }
command -v toml > /dev/null 2>&1 || { echo >&2 "toml not installed. More info: https://github.com/mrijken/toml-cli"; exit 1; }

# install aurad if not exist
if [ $WILL_INSTALL -eq 0 ];
then 
    command -v aurad > /dev/null 2>&1 || { echo >&1 "installing aurad"; make install; }
else
    echo >&1 "installing aurad"
    rm -rf $HOME/.aura*
    rm scripts/mnemonic.txt
    make install
fi

aurad config keyring-backend $KEYRING
aurad config chain-id $CHAINID

# determine if user wants to recorver or create new
MNEMONIC=""
if [ $WILL_RECOVER -eq 0 ];
then
    MNEMONIC=$(aurad keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO --output json | jq -r '.mnemonic')
else
    MNEMONIC=$(aurad keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO --recover --output json | jq -r '.mnemonic')
fi

echo "MNEMONIC for $(aurad keys show $KEY -a --keyring-backend $KEYRING) = $MNEMONIC" >> scripts/mnemonic.txt

echo >&1 "\n"

# init chain
aurad init $MONIKER --chain-id $CHAINID

# Change parameter token denominations to uaura
cat $HOME/.aura/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="uaura"' > $HOME/.aura/config/tmp_genesis.json && mv $HOME/.aura/config/tmp_genesis.json $HOME/.aura/config/genesis.json
cat $HOME/.aura/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="uaura"' > $HOME/.aura/config/tmp_genesis.json && mv $HOME/.aura/config/tmp_genesis.json $HOME/.aura/config/genesis.json
cat $HOME/.aura/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="uaura"' > $HOME/.aura/config/tmp_genesis.json && mv $HOME/.aura/config/tmp_genesis.json $HOME/.aura/config/genesis.json
cat $HOME/.aura/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="uaura"' > $HOME/.aura/config/tmp_genesis.json && mv $HOME/.aura/config/tmp_genesis.json $HOME/.aura/config/genesis.json

# Set gas limit in genesis
# cat $HOME/.aura/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="10000000"' > $HOME/.aura/config/tmp_genesis.json && mv $HOME/.aura/config/tmp_genesis.json $HOME/.aura/config/genesis.json

# enable rest server and swagger
toml set --toml-path $HOME/.aura/config/app.toml api.swagger true
toml set --toml-path $HOME/.aura/config/app.toml api.enable true
toml set --toml-path $HOME/.aura/config/app.toml api.address tcp://0.0.0.0:1310
toml set --toml-path $HOME/.aura/config/client.toml node tcp://0.0.0.0:2281

# create more test key
MNEMONIC_1=$(aurad keys add test1 --keyring-backend $KEYRING --algo $KEYALGO --output json | jq -r '.mnemonic')
TO_ADDRESS=$(aurad keys show test1 -a --keyring-backend $KEYRING)
echo "MNEMONIC for $TO_ADDRESS = $MNEMONIC_1" >> scripts/mnemonic.txt

# Allocate genesis accounts (cosmos formatted addresses)
aurad add-genesis-account $KEY 1000000000000uaura --keyring-backend $KEYRING
aurad add-genesis-account test1 1000000000000uaura --keyring-backend $KEYRING

# Sign genesis transaction
aurad gentx $KEY 1000000uaura --keyring-backend $KEYRING --chain-id $CHAINID

# Collect genesis tx
aurad collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
aurad validate-genesis

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
aurad start --pruning=nothing --log_level $LOGLEVEL --minimum-gas-prices=0.0001uaura --p2p.laddr tcp://0.0.0.0:2280 --rpc.laddr tcp://0.0.0.0:2281 --grpc.address 0.0.0.0:2282 --grpc-web.address 0.0.0.0:2283