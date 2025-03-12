<!--
All notable changes to this project will be documented in this file.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
-->

## Version 10.0.0
### Features
#### Other
- ([\#265](https://github.com/milkyway-labs/milkyway/pull/265)) Added support for CosmWasm IBC hooks

### Dependencies
- ([\#111](https://github.com/milkyway-labs/milkyway/pull/111)) Updated `github.com/cosmos/cosmos-sdk` to `v0.50.13`

## Version 9.0.0
### Features
#### Restaking
- ([\#249](https://github.com/milkyway-labs/milkyway/pull/249)) Added a limit on the amount of unbonding entries that a target can have

### Bug Fixes
#### Liquid Vesting
- ([\#231](https://github.com/milkyway-labs/milkyway/pull/231)) Removed the strict module account balance check

#### Restaking
- ([\#230](https://github.com/milkyway-labs/milkyway/pull/230)) Filtered unbonding delegations by requested target id
- ([\#242](https://github.com/milkyway-labs/milkyway/pull/242)) Removed problematic scaling gas costs

## Version 8.1.1
### Bug fixes
#### Build
- ([793f05f](https://github.com/milkyway-labs/milkyway/commit/793f05f)) Added `CGO_ENABLED=1` to builder image
- ([37a52db](https://github.com/milkyway-labs/milkyway/commit/37a52db)) Added `linux/amd64` specification to the builder image

## Version 8.1.0-ceers
This version has been released to update the `ceers-2112` testnet to version `v8.1.0` of the software.

### Features
#### Other
- ([dc8ef23](https://github.com/milkyway-labs/milkyway/commit/dc8ef23)) Updated `v8` fork height

## Version 8.1.0
### Features
#### Other
- ([43b7910](https://github.com/milkyway-labs/milkyway/commit/43b7910)) Updated `v8` fork height

## Version 8.0.0
### Features
#### Rewards
- ([2eca460](https://github.com/milkyway-labs/milkyway/commit/2eca460)) Updated rewards plan to only support a single denom
- ([fdcc23c](https://github.com/milkyway-labs/milkyway/commit/fdcc23c)) Added base gas fee consumption for rewards plan creation
- ([a171ac0](https://github.com/milkyway-labs/milkyway/commit/a171ac0)) Avoid transferring rewards when skipping the allocation

#### Restaking
- ([987da9f](https://github.com/milkyway-labs/milkyway/commit/987da9f)) Improved gas usage
- ([1fce5da](https://github.com/milkyway-labs/milkyway/commit/1fce5da)) Added scaling gas costs to delegations and undelegations
- ([1d2e6de](https://github.com/milkyway-labs/milkyway/commit/1d2e6de)) Updated the meaning of empty securing pools to "No pools" rather than "All pools"
- ([1caccb5](https://github.com/milkyway-labs/milkyway/commit/1caccb5)) Optimized delegations by target id queries
- ([263fd65](https://github.com/milkyway-labs/milkyway/commit/263fd65)) Optimized the `getEligibleOperators` query
- ([bfcaff9](https://github.com/milkyway-labs/milkyway/commit/bfcaff9)) Improved the `UserPreferences` structure and its features
- ([64ccbd3](https://github.com/milkyway-labs/milkyway/commit/64ccbd3)) Updated the meaning of default user preferences from "Trust all services" to "Trust no service"

#### Other
- ([3f892b7](https://github.com/milkyway-labs/milkyway/commit/3f892b7)) Removed unnecessary error overrides
- ([669ce32](https://github.com/milkyway-labs/milkyway/commit/669ce32)) Added the support for store migrations inside hard fork handlers

### Bug fixes
#### Rewards
- ([d2694ea](https://github.com/milkyway-labs/milkyway/commit/d2694ea)) Replaced `GetAllBalances` inside `InitGenesis`
- ([8baa348](https://github.com/milkyway-labs/milkyway/commit/8baa348)) Fixed bugs in the calculation of pool-service total delegator shares
- ([85212d2](https://github.com/milkyway-labs/milkyway/commit/85212d2)) Fixed rewards allocations

#### Restaking
- ([2a74f58](https://github.com/milkyway-labs/milkyway/commit/2a74f58)) Added check to make sure an operator is allowed  to join a service while executing `MsgJoinService`
- ([c8f7880](https://github.com/milkyway-labs/milkyway/commit/c8f7880)) Made sure operators that are removed from an allowlist also leave the service

#### Other
- ([7090ac2](https://github.com/milkyway-labs/milkyway/commit/7090ac2)) Removed `SharesFromTokensTruncated` in favor of `SharesFromTokens`
- ([94deb84](https://github.com/milkyway-labs/milkyway/commit/94deb84)) Fixed `ParseTrustedServiceEntry`


## Version 7.0.0
### Bug fixes
#### Restaking
- ([ed8281a](https://github.com/milkyway-labs/milkyway/commit/ed8281a)) Set restaking cap to `0`

## Version 6.1.0
### Features
- ([d275ee8](https://github.com/milkyway-labs/milkyway/commit/d275ee8)) Removed the deletion of markets from upgrade
  handler

## Version 6.0-ceers
This version has been released to update the `ceers-2112` testnet to version `v6` of the software.

### Bug fixes
#### LiquidVesting
- ([adf62d4](https://github.com/milkyway-labs/milkyway/commit/adf62d4)) Properly set the `x/liquidvesting` module account

## Version 6.0.0
### Bug fixes
#### LiquidVesting
- ([\#225](https://github.com/milkyway-labs/milkyway/pull/225)) Properly initialized the module account

## Version 5.0.0
### Features
#### Other
- ([\#224](https://github.com/milkyway-labs/milkyway/pull/224)) Removed gov ante decorators to allow any proposal to be run as expedited

## Version 4.0.0
### Features
- ([\#222](https://github.com/milkyway-labs/milkyway/pull/222)) Added `v4` upgrade handler

### Dependencies
- ([\#221](https://github.com/milkyway-labs/milkyway/pull/221)) Updated `github.com/cosmos/cosmos-sdk` to `v0.50.11`

## Version 3.0.0
### Features
#### LiquidVesting
- ([\#215](https://github.com/milkyway-labs/milkyway/pull/215)) Removed `trusted_delegates` from the params

#### Other
- ([\#214](https://github.com/milkyway-labs/milkyway/pull/214)) Added v3 upgrade handler

## Version 2.0.2
### Bug fixes
#### Build
- ([a3ba245](https://github.com/milkyway-labs/milkyway/commit/a3ba245)) Fixed [buf](https://buf.build) build errors

## Version 2.0.1
### Bug fixes
#### MarketMap
- ([89f6385](https://github.com/milkyway-labs/milkyway/commit/89f6385)) Fixed the default genesis state generation

## Version 2.0.0

This is the first release of the new major version of the project.
The main change that has been made is the transition from being an L2 Optimistic Rollup to being an L1 Cosmos-SDK based
chain.

Aside from this, various bugs have been fixed and useful features has been implemented. You can see the full list of
changes [here](https://github.com/milkyway-labs/milkyway/compare/v1.6.0...v2.0.0).