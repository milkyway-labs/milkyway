<!--
All notable changes to this project will be documented in this file.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
-->

## Version 7.0.0

### Bug fixes
### Restaking

* Set restaking cap to `0` (ed8281a9bc8c0ce5de91019d23bf788d2f4c0af2)

## Version 6.1.0
### Features

* Removed the deletion of markets from upgrade handler (d275ee8)

## Version 6.0-ceers
This version has been released to update the `ceers-2112` testnet to version `v6` of the software.

### Bug fixes
#### LiquidVesting

* Properly set the `x/liquidvesting` module
  account (https://github.com/milkyway-labs/milkyway/commit/adf62d4fd620c76f39d9fb76bb6ffada01139e93)

## Version 6.0.0
### Bug fixes
#### LiquidVesting

* Properly initialized the module account

## Version 5.0.0
### Features
#### Other

- ([\#224](https://github.com/milkyway-labs/milkyway/pull/224)) Removed gov ante decorators to allow any proposal to be
  run as expedited

## Version 4.0.0
### Features

- ([\#222](https://github.com/milkyway-labs/milkyway/pull/222)) Added `v4` upgrade handler

### Dependencies

- ([\#221](https://github.com/milkyway-labs/milkyway/pull/221)) Updated `github.com/cosmos/cosmos-sdk` to `v0.50.11`

## Version 3.0.0
### Features
#### Liquid Vesting

- ([\#215](https://github.com/milkyway-labs/milkyway/pull/215)) Removed `trusted_delegates` from the params

#### Other

- ([\#214](https://github.com/milkyway-labs/milkyway/pull/214)) Added v3 upgrade handler

## Version 2.0.2

### Bug fixes

#### Build

* Fixed [buf](https://buf.build) build errors

## Version 2.0.1

### Bug fixes

#### x/marketmap

* Fixed the default genesis state
  generation ([89f6385](https://github.com/milkyway-labs/milkyway/commit/89f638567af91e819e6ae3948823b55a24292d61))

## Version 2.0.0

This is the first release of the new major version of the project.
The main change that has been made is the transition from being an L2 Optimistic Rollup to being an L1 Cosmos-SDK based
chain.

Aside from this, various bugs have been fixed and useful features has been implemented. You can see the full list of
changes [here](https://github.com/milkyway-labs/milkyway/compare/v1.6.0...v2.0.0).