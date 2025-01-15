# MilkyWay Chain

[![Website](.img/cover.jpg)](https://milkyway.zone)

**<p align="center">The First and Largest Liquid Staking and Restaking for Modular Ecosystem</p>**

[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/milkyway-labs/milkyway)](https://github.com/milkyway-labs/milkyway/releases)
![License](https://img.shields.io/github/license/milkyway-labs/milkyway.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/milkyway-labs/milkyway/.svg)](https://pkg.go.dev/github.com/milkyway-labs/milkyway/)
[![Go Report](https://goreportcard.com/badge/github.com/milkyway-labs/milkyway)](https://goreportcard.com/report/github.com/milkyway-labs/milkyway)
[![Codecov](https://codecov.io/gh/milkyway-labs/milkyway/branch/main/graph/badge.svg)](https://codecov.io/gh/milkyway-labs/milkyway/branch/main)
[![Tests status](https://github.com/milkyway-labs/milkyway/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/milkyway-labs/milkyway/actions/workflows/test.yml?query=branch%3Amain+)
[![Lint status](https://github.com/milkyway-labs/milkyway/actions/workflows/lint.yml/badge.svg?branch=main)](https://github.com/milkyway-labs/milkyway/actions/workflows/lint.yml?query=branch%3Amain+)
[![Discord](https://img.shields.io/discord/1166634853576482876)](https://discord.com/invite/4ywmNE3tqq)

## What is MilkyWay?

MilkyWay is the first and largest liquid staking and restaking protocol in the modular ecosystem.

In December 2023, as strong believers in modular architecture, we began our journey by pioneering an enhanced liquid staking solution for the Celestia blockchain. We have evolved into a comprehensive solution that unifies liquidity and security across a broad range of blockchain services. MilkyWay aims to consolidate fragmented trust and security by allowing multiple services to tap into a single, flexible pool of staked assets.

This repository is the codebase of the MilkyWay Chain, a modular L1 coordination layer, built with the Cosmos SDK that is designed to secure an expanding constellation of networks and services through a trust-minimized, multi-asset, multi-chain restaking platform. It allows AVSs to “plug in” and to precisely configure asset allocation, stake distribution, rewards distribution, and slashing parameters. By restaking both native and liquid-staked assets, it diversifies security while providing a sovereign environment with direct control over consensus and validation rules.

This modular architecture frees developers and AVSs from the constraints of traditional parent chains, fostering experimentation with new governance models, incentive mechanisms, and specialized use cases. From reinforcing off-chain systems to introducing innovative staking-based financial instruments, MilkyWay sets the stage for broad ecosystem growth and adaptability. As its capabilities mature and its scope widens, the platform establishes a resilient foundation for an interconnected, composable future. By transcending conventional staking solutions and operating as a flexible, sovereign environment, MilkyWay nurtures a more adaptive, trust-minimized blockchain ecosystem. Ultimately, this dynamic approach broadens participation, strengthens interoperability, and ushers in a new era of decentralized infrastructure.

To learn more about MilkyWay, visit the [MilkyWay Documentation](https://docs.milkyway.zone/).

## Install MilkyWay Core

To find comprehensive instructions on hardware specifications, MilkyWay Core installation, full node operation, and joining a network, please refer to the [MilkyWay node tutorial](https://docs.milkyway.zone/modular-restaking/guides/consensus).

## Interact with MilkyWay

For users looking to interact with the MilkyWay blockchain without setting up a full node, the following official wallets provide a convenient interface.

- [Keplr](https://chromewebstore.google.com/detail/keplr/dmkamcknogkgcdfhhbddcghachkejeap)
- [Leap](https://chromewebstore.google.com/detail/keplr/dmkamcknogkgcdfhhbddcghachkejeap) 
- [Cosmostation](https://chromewebstore.google.com/detail/cosmostation-wallet/fpkhgmpbidmiogeglndfbkegfdlnajnf?hl=en)

Developers seeking direct blockchain interaction can download the [milkywayd](https://github.com/orgs/milkyway-labs/packages?repo_name=milkyway), the CLI and node daemon for MilkyWay chain. For detailed instructions on installing and using CLI, refer to the official MilkyWay [documentation](https://docs.milkyway.zone/).

## For Developers

For an introduction to building on MilkyWay, start with the [MilkyWay Docs](https://docs.milkyway.zone/), and explore the [Services/Operators Guide](https://docs.milkyway.zone/modular-restaking/guides) for an overview of the restaking infrastructure.

- [milkyway.proto](https://github.com/milkyway-labs/milkyway.proto): a set of TypeScript and JavaScript definitions compatible with [Cosmjs](https://github.com/cosmos/cosmjs) to easily integrate your web-app with MilkyWay

- [milkywayd](https://github.com/orgs/milkyway-labs/packages?repo_name=milkyway): the MilkyWay blockchain’s CLI and node daemon

- [mintscan](https://www.mintscan.io/milkyway): third party block explorer and chain analytics

## Resources

- [Official Website](https://milkyway.zone)
- [Documentation](https://docs.milkyway.zone/)
- [X (Formerly Twitter)](https://twitter.com/milky_way_zone)
- [Blog](https://medium.com/milkyway-zone)
- [Discord](https://discord.com/invite/4ywmNE3tqq)

## License

MilkyWay is licensed under the [Apache 2.0](LICENSE) license.

We also use portions of the codebase from other developers' software to implement some of our features:

- [Cosmos SDK](https://github.com/cosmos/cosmos-sdk), licensed under the [Apache 2.0](https://github.com/cosmos/cosmos-sdk?tab=Apache-2.0-1-ov-file)