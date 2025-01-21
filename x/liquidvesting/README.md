# `x/liquidvesting`

## Abstract

The following document specify the liquid vesting module.

This module allows the creation of synthetic tokens that represent tokens that have been locked elsewhere.

## Contents

* [Concepts](#concepts)
   * [User insurance fund](#user-insurance-fund)
   * [Locked tokens](#locked-tokens)
   * [Burn coins list](#burn-coins-list)
* [State](#state)
   * [Params](#params)
   * [User insurance funds](#user-insurance-funds)
   * [Burn coins queue](#burn-coins-queue)
* [Messages](#messages)
   * [MsgMintLockedRepresentation](#msgmintlockedrepresentation)
   * [MsgBurnLockedRepresentation](#msgburnlockedrepresentation)
   * [MsgWithdrawInsuranceFund](#msgwithdrawinsurancefund)
* [Events](#events)
   * [EndBlocker](#endblocker)
   * [Handlers](#handlers)
      * [MsgMintLockedRepresentation](#msgmintlockedrepresentation)
      * [MsgBurnLockedRepresentation](#msgburnlockedrepresentation)
      * [MsgWithdrawInsuranceFund](#msgwithdrawinsurancefund)

## Concepts

### User insurance fund

A user insurance fund is the amount of tokens that a user has accrued in the liquid vesting module as an insurance for
future slashing operations. This fund is accumulated by sending over tokens to the liquid vesting module and specifying
the user that should receive the insurance.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/liquidvesting/v1/models.proto#L12-L19
```

When trying to restake some locked tokens, the following formula is applied:

```
restakable_quantity = insurance_fund / insurance_percentage
```

Where `insurance_percentage` is the percentage of the insurance fund that is used to cover for the slashing of the
restaked and is determined by the module's params.

This means that the minimum amount of tokens that a user always has to have in their insurance fund is:

```
min_insurance_amount = restaked_tokens * insurance_percentage
```

#### Depositing into an insurance fund

In order to deposit into an insurance fund, a user must send over tokens to the liquid vesting module though an IBC
transfer message, and specify themselves as the receiver of those funds. This can be done using the following command:

```
<app-name> tx ibc-transfer transfer transfer [channel-id] milk102lq49sg6lmw2e0mw740tjldzq68v0yfgtazz5 [amount] \
  --memo '{"liquidvesting":{"amounts":[{"depositor":"[user-address]","amount":"[amount]"}]}}'
```

Note that the memo has to be parsable inside the following object structure:

```json
{
  "liquidvesting": {
    "amounts": [
      {
        "depositor": "[user-address]",
        "amount": "[amount]"
      }
    ]
  }
}
```

Where the `[amount]` has to be only the amount of tokens to be deposited, and does not have to contain the denomination.
For example, if you want to deposit `1` token to two users, you would have to use the following memo:

```json
{
  "liquidvesting": {
    "amounts": [
      {
        "depositor": "[user-1-address]",
        "amount": "1000000"
      },
      {
        "depositor": "[user-2-address]",
        "amount": "1000000"
      }
    ]
  }
}
```

Another thing to note is that only specific IBC channels are allowed when depositing funds into a user insurance fund.
The list of channels is specific inside the module's on-chain parameters and can be changed by the chain's governance.

#### Withdrawing from an insurance fund

A user can always withdraw excessive funds from their insurance fund by sending a message to the liquid vesting module.
The only condition to be able to properly withdraw funds is that after the withdrawal there have to be enough funds to
cover for slashing of tokens that have been restaked based on the insurance percentage amount set in the module params.

To give an example, let's consider the following scenario:

* Alice has an insurance fund of 50 tokens
* The insurance percentage is set to 10%
* Alice has restaked 20 tokens

Considering the insurance percentage at 10%, this means that Alice only needs to have 2 tokens in her insurance fund to
cover for the slashing of the 20 tokens that have been restaked. This means that Alice can withdraw 48 tokens from her
insurance fund without any issues. However, if she was to withdraw 49 tokens, the transaction would fail as that would
leave her with only 1 token in her insurance fund, which is not enough to cover for the slashing of the 20 tokens that
have been restaked.

### Locked tokens

Locked tokens are native tokens that have been created with the usage of the `x/tokenfactory` module. These tokens can
only be created by specific users that have been granted the "minter" role in the module.

Once minted, locked tokens can only be restaked and cannot be transferred to other users or modules.

#### Minting locked tokens

The minting of locked tokens can be performed by sending a `MsgMintLockedRepresentation` from users that have been
accredited with the "minter" role using the module's params.

#### Burning locked tokens

The burning of locked tokens can be performed by sending a `MsgBurnLockedRepresentation` from users that have been
accredited with the "burner" role using the module's params.

If a burner sends a burn message for a user who has their tokens restaked, the restaking positions will automatically be
undelegated to cover the burn. Once the undelegations complete, then the undelegated tokens will be burned as soon as
the next block is produced.

### Burn coins list

The `BurnCoinsList` contains a list of `BurnCoins` object, each one containing the information about the coins to burn
once the unbonding period of the tokens ends.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/liquidvesting/v1/models.proto#L21-L50
```

## State

### Params

The module params are stored using the `0x01` key:

* Params: `0x01 -> ProtocolBuffer(Params)`

### User insurance funds

The user insurance funds are stored using the `0x10` key:

* User insurance fund: `0x10 | UserAddress -> ProtocolBuffer(UserInsuranceFund)`

### Burn coins queue

The amount of tokens that are automatically undelegated and queued to be burned are stored using the `0x20` key:

* Burn coins queue: `0x20 -> ProtocolBuffer(BurnCoinsList)`

## Messages

### MsgMintLockedRepresentation

The `MsgMintLockedRepresentation` can be sent by authorized minters to mint new locked tokens.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/liquidvesting/v1/messages.proto#L39-L59
```

The message will fail under the following conditions:

* The user is not an authorized minter

### MsgBurnLockedRepresentation

The `MsgBurnLockedRepresentation` can be sent by authorized burners to burn locked tokens.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/liquidvesting/v1/messages.proto#L65-L84
```

The message will fail under the following conditions:

* The user is not an authorized burner
* The amount to be burned is not composed of locked tokens
* The user does not have any liquid nor restaked locked tokens

### MsgWithdrawInsuranceFund

The `MsgWithdrawInsuranceFund` can be sent by users to withdraw tokens from their insurance fund.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/liquidvesting/v1/messages.proto#L90-L105
```

The message will fail under the following conditions:

* The user does not have enough funds in their insurance fund to be withdrawn
* The remaining insurance fund amount would not cover for the restaked tokens amount

## Events

### EndBlocker

| Type                         | Attribute Key | Attribute Value              |
|------------------------------|---------------|------------------------------|
| `burn_locked_representation` | `amount`      | `{burnedTokensAmount}`       |
| `burn_locked_representation` | `user`        | `{burnedTokensOwnerAddress}` |

### Handlers

#### MsgMintLockedRepresentation

| Type                         | Attribute Key | Attribute Value     |
|------------------------------|---------------|---------------------|
| `mint_locked_representation` | `sender`      | `{minterAddress}`   |
| `mint_locked_representation` | `amount`      | `{mintedAmount}`    |
| `mint_locked_representation` | `receiver`    | `{receiverAddress}` |

#### MsgBurnLockedRepresentation

| Type                         | Attribute Key | Attribute Value       |
|------------------------------|---------------|-----------------------|
| `burn_locked_representation` | `sender`      | `{burnerAddress}`     |
| `burn_locked_representation` | `amount`      | `{burnedAmount}`      |
| `burn_locked_representation` | `user`        | `{tokenOwnerAddress}` |

#### MsgWithdrawInsuranceFund

| Type                      | Attribute Key | Attribute Value     |
|---------------------------|---------------|---------------------|
| `withdraw_insurance_fund` | `sender`      | `{userAddress}`     |
| `withdraw_insurance_fund` | `amount`      | `{withdrawnAmount}` |

## Parameters

The liquid vesting module contains the following parameters:

 ```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/liquidvesting/v1/params.proto#L9-L35
```


