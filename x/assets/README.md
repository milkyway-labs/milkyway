# `x/assets`

## Abstract

`x/assets` is a Cosmos SDK module that allows the tracking of different assets and their respective on-chain and
off-chain denominations.

## Contents

* [Concepts](#concepts)
   * [Asset](#asset)
* [State](#state)
   * [Assets](#assets)
   * [TickerIndex](#tickerindex)
* [Events](#events)
   * [Handlers](#handlers)
      * [MsgRegisterAsset](#msgregisterasset)
      * [MsgDeregisterAsset](#msgderegisterasset)

## Concepts

### Asset

An asset is a representation of an off-chain token on-chain. It is composed of an off-chain ticker and an on-chain
denomination. Assets are non-unique, which means that there can be multiple `Asset` objects for the same off-chain
ticker. This is to support tokens that have the same off-chain ticker but different on-chain denominations (e.g. in the
case of assets bridged using different IBC channels).

## State

### Assets

The assets module stores the assets in state with the prefix of `0x11`.
The list of stored assets can be updated through the on-chain governance or using the authority address.

* Asset: `0x11 | Denom -> ProtocolBuffer(asset)`

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/assets/v1/models.proto#L9-L21
```

### TickerIndex

The assets module stores the mapping of off-chain tickers to on-chain denominations in state with the prefix of `0x12`.
This is used to quickly look up the on-chain denomination of an off-chain ticker.

* TickerIndex: `0x12 | Ticker | Denom -> []byte{}`

## Events

### Handlers

#### MsgRegisterAsset

|       Type       | Attribute Key | Attribute Value |
|:----------------:|:-------------:|:---------------:|
| `register_asset` |    `denom`    |    `{denom}`    |
| `register_asset` |   `ticker`    |   `{ticker}`    |
| `register_asset` |  `exponent`   |  `{exponent}`   |

#### MsgDeregisterAsset

|        Type        | Attribute Key | Attribute Value |
|:------------------:|:-------------:|:---------------:|
| `deregister_asset` |    `denom`    |    `{denom}`    |