<!--
All notable changes to this project will be documented in this file.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
-->

## Version 8.0.0

### Rewards

- ([a460b9f](https://github.com/milkyway-labs/milkyway/commit/a460b9f)) Updated rewards plan to only support a single
  denom
- ([04b6547](https://github.com/milkyway-labs/milkyway/commit/04b6547)) Added base gas fee consumption for rewards plan
  creation
- ([ede14aa](https://github.com/milkyway-labs/milkyway/commit/ede14aa)) Replaced `GetAllBalances` inside `InitGenesis`
- ([62fc7f1](https://github.com/milkyway-labs/milkyway/commit/62fc7f1)) Avoid transferring rewards when skipping the
  allocation
- ([1addfbe](https://github.com/milkyway-labs/milkyway/commit/1addfbe)) Fixed bugs in the calculation of pool-service
  total delegator shares
- ([81ceb4c](https://github.com/milkyway-labs/milkyway/commit/81ceb4c)) Fixed rewards allocations

### Restaking

- ([a27d0d8](https://github.com/milkyway-labs/milkyway/commit/a27d0d8)) Improved gas usage
- ([f5815d3](https://github.com/milkyway-labs/milkyway/commit/f5815d3)) Added check to make sure an operator is allowed
  to join a service while executing `MsgJoinService`
- ([96b0d54](https://github.com/milkyway-labs/milkyway/commit/96b0d54)) Made sure operators that are removed from an
  allowlist also leave the service
- ([f522eb8](https://github.com/milkyway-labs/milkyway/commit/f522eb8)) Added scaling gas costs to delegations and
  undelegations
- ([623fa32](https://github.com/milkyway-labs/milkyway/commit/623fa32)) Updated the meaning of empty securing pools to "
  No pools" rather than "All pools"
- ([cac1c3d](https://github.com/milkyway-labs/milkyway/commit/cac1c3d)) Optimized delegations by target id queries
- ([449f0e6](https://github.com/milkyway-labs/milkyway/commit/449f0e6)) Optimized the `getEligibleOperators` query
- ([50af532](https://github.com/milkyway-labs/milkyway/commit/50af532)) Improve the `UserPreferences` structure and its
  features
- ([76b80f8](https://github.com/milkyway-labs/milkyway/commit/76b80f8)) Updated the meaning of default user preferences
  from "Trust all services" to "Trust no service"

### Other

- ([d63822a](https://github.com/milkyway-labs/milkyway/commit/d63822a)) Remove `SharesFromTokensTruncated` in favor of
  `SharesFromTokens`
- ([438f1f7](https://github.com/milkyway-labs/milkyway/commit/438f1f7)) Remove unnecessary error overrides
- ([c9bc987](https://github.com/milkyway-labs/milkyway/commit/c9bc987)) Fixed `ParseTrustedServiceEntry`
- ([e0ad42a](https://github.com/milkyway-labs/milkyway/commit/e0ad42a)) Added the support for store migrations inside
  hard fork handlers

## Version 7.0.0

### Bug fixes
### Restaking

- ([ed8281a](https://github.com/milkyway-labs/milkyway/commit/ed8281a)) Set restaking cap to `0`

## Version 6.1.0
### Features

- ([d275ee8](https://github.com/milkyway-labs/milkyway/commit/d275ee8)) Removed the deletion of markets from upgrade
  handler

## Version 6.0-ceers
This version has been released to update the `ceers-2112` testnet to version `v6` of the software.

### Bug fixes
#### LiquidVesting

* ([adf62d4](https://github.com/milkyway-labs/milkyway/commit/adf62d4)) Properly set the `x/liquidvesting` module
  account

## Version 6.0.0
### Bug fixes
#### LiquidVesting

* ([\#225](https://github.com/milkyway-labs/milkyway/pull/225)) Properly initialized the module account

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

* ([a3ba245](https://github.com/milkyway-labs/milkyway/commit/a3ba245)) Fixed [buf](https://buf.build) build errors

## Version 2.0.1

### Bug fixes

#### x/marketmap

* ([89f6385](https://github.com/milkyway-labs/milkyway/commit/89f6385)) Fixed the default genesis state generation

## Version 2.0.0

This is the first release of the new major version of the project.
The main change that has been made is the transition from being an L2 Optimistic Rollup to being an L1 Cosmos-SDK based
chain.

Aside from this, various bugs have been fixed and useful features has been implemented. You can see the full list of
changes [here](https://github.com/milkyway-labs/milkyway/compare/v1.6.0...v2.0.0).