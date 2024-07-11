package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

var _ types.RestakingHooks = Hooks{}

type Hooks struct {
	k *Keeper
}

func (k *Keeper) Hooks() Hooks {
	return Hooks{k}
}

func (h Hooks) BeforePoolDelegationCreated(ctx sdk.Context, poolID uint32, delegator string) error {
	pool, found := h.k.poolsKeeper.GetPool(ctx, poolID)
	if !found {
		return poolstypes.ErrPoolNotFound
	}

	// Initialize pool info if it doesn't exist yet.
	exists, err := h.k.PoolCurrentRewards.Has(ctx, pool.ID)
	if err != nil {
		return err
	}
	if !exists {
		if err := h.k.InitializePool(ctx, pool); err != nil {
			return err
		}
	}

	_, err = h.k.IncrementPoolPeriod(ctx, pool)
	return err
}

func (h Hooks) BeforePoolDelegationSharesModified(ctx sdk.Context, poolID uint32, delegator string) error {
	pool, found := h.k.poolsKeeper.GetPool(ctx, poolID)
	if !found {
		return poolstypes.ErrPoolNotFound
	}

	// We don't have to initialize pool here because we can assume BeforePoolDelegationCreated
	// has already been called when delegation shares are being modified.

	del, found := h.k.restakingKeeper.GetPoolDelegation(ctx, poolID, delegator)
	if !found {
		return sdkerrors.ErrNotFound.Wrapf("pool delegation not found: %d, %s", poolID, delegator)
	}

	if _, err := h.k.withdrawPoolDelegationRewards(ctx, pool, del); err != nil {
		return err
	}

	return nil
}

func (h Hooks) AfterPoolDelegationModified(ctx sdk.Context, poolID uint32, delegator string) error {
	return h.k.initializePoolDelegation(ctx, poolID, delegator)
}

func (h Hooks) BeforeOperatorDelegationCreated(ctx sdk.Context, operatorID uint32, delegator string) error {
	operator, found := h.k.operatorsKeeper.GetOperator(ctx, operatorID)
	if !found {
		return operatorstypes.ErrOperatorNotFound
	}

	// Initialize operator info if it doesn't exist yet.
	exists, err := h.k.OperatorCurrentRewards.Has(ctx, operator.ID)
	if err != nil {
		return err
	}
	if !exists {
		if err := h.k.InitializeOperator(ctx, operator); err != nil {
			return err
		}
	}

	_, err = h.k.IncrementOperatorPeriod(ctx, operator)
	return err
}

func (h Hooks) BeforeOperatorDelegationSharesModified(ctx sdk.Context, operatorID uint32, delegator string) error {
	operator, found := h.k.operatorsKeeper.GetOperator(ctx, operatorID)
	if !found {
		return operatorstypes.ErrOperatorNotFound
	}

	// We don't have to initialize operator here because we can assume BeforeOperatorDelegationCreated
	// has already been called when delegation shares are being modified.

	del, found := h.k.restakingKeeper.GetOperatorDelegation(ctx, operatorID, delegator)
	if !found {
		return sdkerrors.ErrNotFound.Wrapf("operator delegation not found: %d, %s", operatorID, delegator)
	}

	if _, err := h.k.withdrawOperatorDelegationRewards(ctx, operator, del); err != nil {
		return err
	}

	return nil
}

func (h Hooks) AfterOperatorDelegationModified(ctx sdk.Context, operatorID uint32, delegator string) error {
	return h.k.initializeOperatorDelegation(ctx, operatorID, delegator)
}

func (h Hooks) BeforeServiceDelegationCreated(ctx sdk.Context, serviceID uint32, delegator string) error {
	service, found := h.k.servicesKeeper.GetService(ctx, serviceID)
	if !found {
		return servicestypes.ErrServiceNotFound
	}

	// Initialize service info if it doesn't exist yet.
	exists, err := h.k.ServiceCurrentRewards.Has(ctx, service.ID)
	if err != nil {
		return err
	}
	if !exists {
		if err := h.k.InitializeService(ctx, service); err != nil {
			return err
		}
	}

	_, err = h.k.IncrementServicePeriod(ctx, service)
	return err
}

func (h Hooks) BeforeServiceDelegationSharesModified(ctx sdk.Context, serviceID uint32, delegator string) error {
	service, found := h.k.servicesKeeper.GetService(ctx, serviceID)
	if !found {
		return servicestypes.ErrServiceNotFound
	}

	// We don't have to initialize service here because we can assume BeforeServiceDelegationCreated
	// has already been called when delegation shares are being modified.

	del, found := h.k.restakingKeeper.GetServiceDelegation(ctx, serviceID, delegator)
	if !found {
		return sdkerrors.ErrNotFound.Wrapf("service delegation not found: %d, %s", serviceID, delegator)
	}

	if _, err := h.k.withdrawServiceDelegationRewards(ctx, service, del); err != nil {
		return err
	}

	return nil
}

func (h Hooks) AfterServiceDelegationModified(ctx sdk.Context, serviceID uint32, delegator string) error {
	return h.k.initializeServiceDelegation(ctx, serviceID, delegator)
}
