package types

const (
	EventTypeCreateRewardsPlan = "create_rewards_plan"

	EventTypeRewards            = "rewards"
	EventTypeCommission         = "commission"
	EventTypeWithdrawRewards    = "withdraw_rewards"
	EventTypeWithdrawCommission = "withdraw_commission"

	AttributeKeyPoolID     = "pool_id"
	AttributeKeyOperatorID = "operator_id"
	AttributeKeyServiceID  = "service_id"
	AttributeKeyDelegator  = "delegator"

	AttributeKeyAmountPerPool = "amount_per_pool"
	AttributeKeyPool          = "pool" // NOTE: it's different from the restaking pool
)
