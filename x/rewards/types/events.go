package types

const (
	EventTypeCreateRewardsPlan = "create_rewards_plan"

	EventTypeSetWithdrawAddress = "set_withdraw_address"
	EventTypeRewards            = "rewards"
	EventTypeCommission         = "commission"
	EventTypeWithdrawRewards    = "withdraw_rewards"
	EventTypeWithdrawCommission = "withdraw_commission"

	AttributeKeyWithdrawAddress    = "withdraw_address"
	AttributeKeyPoolID             = "pool_id"
	AttributeKeyOperatorID         = "operator_id"
	AttributeKeyDelegationType     = "delegation_type"
	AttributeKeyDelegationTargetID = "delegation_target_id"
	AttributeKeyServiceID          = "service_id"
	AttributeKeyDelegator          = "delegator"

	AttributeKeyAmountPerPool = "amount_per_pool"
	AttributeKeyPool          = "pool" // NOTE: it's different from the restaking pool
)
