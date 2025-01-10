<!--
All notable changes to this project will be documented in this file.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
-->

## Version 8.0.0

### Rewards
- (a460b9f8f8beea70783c5d05559a98467943d1d5) Updated rewards plan to only support a single denom
- (04b65474c5564a712bfb7f9b80e618b535c1ab70) Added base gas fee consumption for rewards plan creation
- (ede14aa85f2779115c22aa2de86154ca1b49041d) Replaced `GetAllBalances` inside `InitGenesis`
- (62fc7f1cd047ca016ba0f9582555c932953c3820) Avoid transferring rewards when skipping the allocation
- (1addfbe145d74572f8d6e8124132251d2cc0d94f) Fixed bugs in the calculation of pool-service total delegator shares

### Restaking
- (a27d0d8f586d363921cd004b5e4003b1e9fb2ef5) Improved gas usage
- (f5815d3560599e7136c53c03ff3547f2727664b2) Added check to make sure an operator is allowed to join a service while executing `MsgJoinService`
- (96b0d54f677ce6994f0c7277f372fe29b6759290) Made sure operators that are removed from an allowlist also leave the service
- (f522eb852157e06a978e92720944d5ce86c00640) Added scaling gas costs to delegations and undelegations
- (623fa32c85277201ddcda52acd73c23db7820cb8) Updated the meaning of empty securing pools to "No pools" rather than "All pools"
- (cac1c3db7d31183a33a625ac72fd4b05829c6dbf) Optimized delegations by target id queries
- (449f0e6fc18f2b40d436ffdfe4c91ac537e54a31) Optimized the `getEligibleOperators` query
- (50af532e80818756b917558d5c24ab3c31c99729) Improve the `UserPreferences` structure and its features

### Other
- (d63822aa8a98263a9bcc6a6984d0483f13664913) Remove `SharesFromTokensTruncated` in favor of `SharesFromTokens`
- (438f1f76fbbf166df31c664257d6692ea174599d) Remove unnecessary error overrides
- (c9bc987b0aaca466a66f998b7e9c02ed324bccd3) Fixed `ParseTrustedServiceEntry`
- (e0ad42a89bd2fa200419db3afc12de285162a6a6) Added the support for store migrations inside hard fork handlers

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