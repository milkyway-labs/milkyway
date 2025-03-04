# `x/pools`

## Abstract

The following document specifies the pools module.

This module contains the logic responsible for managing the restaking pools of our protocol. The pools are used to
track the amount of tokens that each user has decided to restaking and use to provide security to multiple services.

Pools are only created by other modules, and cannot be created by outside users. The only operation that users can
perform are querying the existing pools and their details.

## Contents

* [Concepts](#concepts)
   * [Pool](#pool)
* [State](#state)
   * [Next Pool ID](#next-pool-id)
   * [Pools](#pools)
   * [Pool addresses](#pool-addresses)

## Concepts

### Pool

A pool is used to track the amount of tokens that each user has decided to restake. Each pool can only manage a single
denomination of tokens (i.e. multi-tokens pools are not allowed). At the same time, for each token denomination, a
single pool can exist at any given moment.

For easier reference throughout the codebase, each pool is uniquely identified by an ID that is automatically assigned
when the pool is created.

A pool also has a unique address that is derived from the pool ID itself. The address is used to store the tokens that
are delegated to the pool, while inside the pool object itself shares are tracked.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/pools/v1/models.proto#L10-L35
```

## State

### Next Pool ID

The pools module stores the next pool ID in state with the key `0xa1`.

* Next pool ID: `0xa1 -> uint32`

### Pools

The pools module stores the pools in state with the prefix of `0xa2`.

* Pool: `0xa2 | PoolID -> ProtocolBuffer(Pool)`

### Pool addresses

In order to know more easily and faster if a particular address represents a pool, we store the set of pool addresses
using the `0xa3` prefix.

* Pool address: `0xa3 | PoolAddress -> []byte{}`
