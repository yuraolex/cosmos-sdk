package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) initializeClaim(ctx context.Context, recipient sdk.AccAddress) error {
}

func (k Keeper) withdrawFunds(ctx context.Context, recipient sdk.AccAddress) (sdk.Coins, error) {
	period, err := k.IncrementBudgetPeriod(ctx, recipient)
	if err != nil {
		return nil, err
	}

	// end current tranche and increment Budget period
	remainingTranches, err := k.decrementTranchePeriod(ctx, recipient)
	if err != nil {
		return nil, err
	}

	funds, err := k.CalculateFunds(ctx, recipient, endingPeiod)
	if err != nil {
		return nil, err
	}
}

func (k Keeper) IncrementBudgetPeriod(ctx context.Context, recipient sdk.AccAddress) error {
	budget, err := k.Budget.Get(ctx, recipient)
	if err != nil {
		return err
	}

	// fetch current distributed funds

	// 
}

// decrement the tranches count for a budget, and delete if zero tranches remain
func (k Keeper) decrementTranchePeriod(ctx context.Context, recipient sdk.AccAddress) (int64, error) {
	budget, err := k.Budget.Get(ctx, recipient)
	if err != nil {
		return 0, err
	}

	if budget.RemainingTranches == 0 {
		panic("cannot set negative tranche count")
	}
	budget.RemainingTranches--
	if budget.RemainingTranches == 0 {
		err := k.Budget.Remove(ctx, recipient)
		if err != nil {
			return 0, err
		}
	}
	err = k.Budget.Set(ctx, recipient, budget)
	if err != nil {
		return 0, err
	}
	return budget.RemainingTranches, nil
}
