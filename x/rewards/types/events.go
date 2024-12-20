package types

const (
	EventTypeCreateRewardsPlan    = "create_rewards_plan"
	EventTypeEditRewardsPlan      = "edit_rewards_plan"
	EventTypeSetWithdrawAddress   = "set_withdraw_address"
	EventTypeRewards              = "rewards"
	EventTypeCommission           = "commission"
	EventTypeWithdrawRewards      = "withdraw_rewards"
	EventTypeWithdrawCommission   = "withdraw_commission"
	EventTypeTerminateRewardsPlan = "terminate_rewards_plan"

	AttributeKeyRewardsPlanID      = "rewards_plan_id"
	AttributeKeyWithdrawAddress    = "withdraw_address"
	AttributeKeyDelegationType     = "delegation_type"
	AttributeKeyDelegationTargetID = "delegation_target_id"
	AttributeKeyRemainingRewards   = "remaining_rewards"

	// AttributeKeyAmountPerPool represents the amount of rewards per pool (per denom).
	// See https://github.com/initia-labs/initia/blob/v0.2.10/x/distribution/types/events.go#L3-L6
	// for the reference of these attributes.
	AttributeKeyAmountPerPool = "amount_per_pool"

	// AttributeKeyPool represents the rewards pool's name (denom). Note that
	// the meaning of the name "pool" is different from the restaking pool.
	AttributeKeyPool = "pool"
)
