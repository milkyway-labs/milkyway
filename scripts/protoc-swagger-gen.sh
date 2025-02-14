#!/usr/bin/env bash

set -eo pipefail

# clone dependency proto files
COSMOS_URL=github.com/cosmos/cosmos-sdk
IBC_URL=github.com/cosmos/ibc-go
IBC_V=v8
COSMWASM_URL=github.com/CosmWasm/wasmd
CONNECT_URL=github.com/skip-mev/connect
CONNECT_V=v2

COSMOS_SDK_VERSION=$(cat ./go.mod | grep "$COSMOS_URL v" | sed -n -e "s/^.* //p")
IBC_VERSION=$(cat ./go.mod | grep "$IBC_URL/$IBC_V v" | sed -n -e "s/^.* //p")
COSMWASM_VERSION=$(cat ./go.mod | grep "$COSMWASM_URL v" | sed -n -e "s/^.* //p")
CONNECT_VERSION=$(cat ./go.mod | grep "$CONNECT_URL/$CONNECT_V v" | sed -n -e "s/^.* //p")

mkdir -p ./third_party
cd third_party
git clone -b $COSMOS_SDK_VERSION https://$COSMOS_URL
git clone -b $IBC_VERSION https://$IBC_URL
git clone -b $COSMWASM_VERSION https://$COSMWASM_URL
git clone -b $CONNECT_VERSION https://$CONNECT_URL
cd ..


# start generating
mkdir -p ./tmp-swagger-gen
cd proto
proto_dirs=$(find \
  ./milkyway \
  ./osmosis \
  ../third_party/cosmos-sdk/proto/cosmos \
  ../third_party/ibc-go/proto/ibc \
  ../third_party/wasmd/proto \
  ../third_party/connect/proto \
  -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  # generate swagger files (filter query files)
  query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))
  if [[ ! -z "$query_file" ]]; then
    buf generate --template buf.gen.swagger.yaml $query_file
  fi
done
cd ..

# combine swagger files
# uses nodejs package `swagger-combine`.
# all the individual swagger files need to be configured in `config.json` for merging
swagger-combine ./client/docs/config.json -o ./client/docs/swagger-ui/swagger.yaml -f yaml --continueOnConflictingPaths true --includeDefinitions true

# clean swagger files
rm -rf ./tmp-swagger-gen

# clean third party files
rm -rf ./third_party
