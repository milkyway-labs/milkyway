package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	epochstypes "github.com/milkyway-labs/milkyway/x/epochs/types"
)

const StrideEpochsPerDayEpoch = uint64(4)

func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochInfo epochstypes.EpochInfo) {
	// Update the stakeibc epoch tracker
	epochNumber, err := k.UpdateEpochTracker(ctx, epochInfo)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("Unable to update epoch tracker, err: %s", err.Error()))
		return
	}

	// Day Epoch - Process Unbondings
	if epochInfo.Identifier == epochstypes.DAY_EPOCH {
		// Initiate unbondings from any hostZone where it's appropriate
		k.InitiateAllHostZoneUnbondings(ctx, epochNumber)
		// Check previous epochs to see if unbondings finished, and sweep the tokens if so
		k.SweepAllUnbondedTokens(ctx)
		// Cleanup any records that are no longer needed
		k.CleanupEpochUnbondingRecords(ctx, epochNumber)
		// Create an empty unbonding record for this epoch
		k.CreateEpochUnbondingRecord(ctx, epochNumber)
	}

	// Stride Epoch - Process Deposits and Delegations
	if epochInfo.Identifier == epochstypes.STRIDE_EPOCH {
		// Get cadence intervals
		params := k.GetParams(ctx)
		redemptionRateInterval := params.RedemptionRateInterval
		depositInterval := params.DepositInterval
		delegationInterval := params.DelegateInterval
		reinvestInterval := params.ReinvestInterval

		// Claim accrued staking rewards at the beginning of the epoch
		k.ClaimAccruedStakingRewards(ctx)

		// Create a new deposit record for each host zone and the grab all deposit records
		k.CreateDepositRecordsForEpoch(ctx, epochNumber)
		depositRecords := k.recordsKeeper.GetAllDepositRecord(ctx)

		// TODO: move this to an external function that anyone can call, so that we don't have to call it every epoch
		k.SetWithdrawalAddress(ctx)

		// Update the redemption rate
		if epochNumber%redemptionRateInterval == 0 {
			k.UpdateRedemptionRates(ctx, depositRecords)
		}

		// Transfer deposited funds from the controller account to the delegation account on the host zone
		if epochNumber%depositInterval == 0 {
			k.TransferExistingDepositsToHostZones(ctx, epochNumber, depositRecords)
		}

		// Delegate tokens from the delegation account
		if epochNumber%delegationInterval == 0 {
			k.StakeExistingDepositsOnHostZones(ctx, epochNumber, depositRecords)
		}

		// Reinvest staking rewards
		if epochNumber%reinvestInterval == 0 { // allow a few blocks from UpdateUndelegatedBal to avoid conflicts
			k.ReinvestRewards(ctx)
		}

		// Rebalance stake according to validator weights
		// This should only be run once per day, but it should not be run on a stride epoch that
		//   overlaps the day epoch, otherwise the unbondings could cause a redelegation to fail
		// On mainnet, the stride epoch overlaps the day epoch when `epochNumber % 4 == 1`,
		//   so this will trigger the epoch before the unbonding
		if epochNumber%StrideEpochsPerDayEpoch == 0 {
			k.RebalanceAllHostZones(ctx)
		}

		// Transfers in and out of tokens for hostZones which have community pools
		k.ProcessAllCommunityPoolTokens(ctx)

		// Do transfers for all reward and swapped tokens defined by the trade routes every stride epoch
		k.TransferAllRewardTokens(ctx)
	}
	if epochInfo.Identifier == epochstypes.MINT_EPOCH {
		k.AllocateHostZoneReward(ctx)
	}
}

func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochInfo epochstypes.EpochInfo) {}

// Hooks wrapper struct for incentives keeper
type Hooks struct {
	k Keeper
}

var _ epochstypes.EpochHooks = Hooks{}

func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochInfo epochstypes.EpochInfo) {
	h.k.BeforeEpochStart(ctx, epochInfo)
}

func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochInfo epochstypes.EpochInfo) {
	h.k.AfterEpochEnd(ctx, epochInfo)
}

// SetWithdrawalAddress sets the withdrawal account address for each host zone
func (k Keeper) SetWithdrawalAddress(ctx sdk.Context) {
	k.Logger(ctx).Info("Setting Withdrawal Addresses...")

	for _, hostZone := range k.GetAllActiveHostZone(ctx) {
		err := k.SetWithdrawalAddressOnHost(ctx, hostZone)
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("Unable to set withdrawal address on %s, err: %s", hostZone.ChainId, err))
		}
	}
}

// ClaimAccruedStakingRewards allows to claim staking rewards for each host zone
func (k Keeper) ClaimAccruedStakingRewards(ctx sdk.Context) {
	k.Logger(ctx).Info("Claiming Accrued Staking Rewards...")

	for _, hostZone := range k.GetAllActiveHostZone(ctx) {
		err := k.ClaimAccruedStakingRewardsOnHost(ctx, hostZone)
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("Unable to claim accrued staking rewards on %s, err: %s", hostZone.ChainId, err))
		}
	}
}
