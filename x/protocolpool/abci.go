package protocolpool

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	"cosmossdk.io/x/protocolpool/keeper"
	"cosmossdk.io/x/protocolpool/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Define Tranche as a simple struct
type Tranche struct {
	Period int64
	Amount math.Int
}

// Define a function to execute the budget distribution
func distributeBudget(ctx sdk.Context, keeper *keeper.Keeper, bankKeeper types.BankKeeper, budget *types.MsgSetBudgetProposal) error {
	currentTime := ctx.BlockTime().Unix()

	// Check if the start time is reached
	if currentTime < budget.StartTime {
		return fmt.Errorf("Distribution has not started yet")
	}

	// Calculate the number of periods elapsed
	periodsElapsed := (currentTime - budget.StartTime) / budget.Period

	if periodsElapsed < 1 {
		return fmt.Errorf("Distribution period has not passed yet")
	}
	// Calculate the amount to distribute
	amountToDistribute := budget.TotalBudget.Amount
	if periodsElapsed <= budget.RemainingTranches {
		// For the case with a fixed number of tranches
		amountToDistribute = budget.TotalBudget.Amount.QuoRaw(budget.RemainingTranches)
	}

	// Perform the transfer
	coins := sdk.NewCoins(sdk.NewCoin(budget.TotalBudget.Denom, amountToDistribute))
	recipient := keeper.authKeeper.AddressCodec().StringToBytes(budget.RecipientAddress)
	err := keeper.DistributeFromFeePool(ctx, coins, recipient)
	// Perform the transfer
	// err := transferFunds(sdk.AccAddress(budget.RecipientAddress), amountToDistribute, budget.TotalBudget.Denom, bankKeeper)
	if err != nil {
		return err
	}

	// Update the remaining tranches
	budget.RemainingTranches--

	// Update the budget state
	// This part would depend on the specifics of your blockchain's state management.
	// You need to store the updated budget back to the state.

	return nil
}

// func transferFunds(recipient sdk.AccAddress, amount math.Int, denom string, bankKeeper types.BankKeeper) error {
// 	coins := sdk.NewCoins(sdk.NewCoin(denom, amount))
// 	// Transfer funds to the recipient
// 	err := bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, coins)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func EndBlocker(ctx sdk.Context, k keeper.Keeper) error {
	err := k.Budget.Walk(ctx, nil, func(key sdk.AccAddress, budget types.MsgSetBudgetProposal) (stop bool, err error) {
		// Check if the budget has reached its end
		if budget.RemainingTranches <= 0 {
			// Delete the budget from the stateproposal
			// k.DeleteBudget(ctx, budget.RecipientAddress)
			err := k.Budget.Remove(ctx, sdk.AccAddress(budget.RecipientAddress))
			if err != nil {
				return true, err
			}
			k.Logger(ctx).Info(fmt.Sprintf("Budget ended for recipient: %s", budget.RecipientAddress))
		} else {
			// Transfer funds from the community pool to the recipient
			err := k.TransferFundsFromCommunityPool(ctx, budget.TotalBudget, budget.Period)
			if err != nil {
				k.Logger(ctx).Error(fmt.Sprintf("Error transferring funds proposalto recipient %s: %s", budget.RecipientAddress, err.Error()))
				return false, err // Continue iterating
			}

			// Update the remaining tranches
			budget.RemainingTranches--

			// Log the processing of the budget
			k.Logger(ctx).Info(fmt.Sprintf("Processing budget for recipient: %s", &budget.RecipientAddress))

		}
		return false, nil
	})

	// Iterate through recipients and process their budgets
	k.IterateBudgets(ctx, func(recipientAddress sdk.AccAddress, budget *types.Budget) (stop bool) {
		// Your logic to process the budget for each recipient
		// For example, you might transfer funds, update the remaining tranches, etc.

		// Check if the budget has reached its end
		if budget.RemainingTranches <= 0 {
			// Perform actions when the budget is exhausted
			// For example, transfer remaining funds back to the community pool
			// k.TransferRemainingFundsToCommunityPool(ctx, budget.TotalBudget)

			// Delete the budget from the state
			k.DeleteBudget(ctx, recipientAddress)

			// Log the end of the budget
			k.Logger(ctx).Info(fmt.Sprintf("Budget ended for recipient: %s", recipientAddress.String()))
		} else {
			// Perform actions when the budget is still active
			// For example, transfer funds, update the remaining tranches, etc.
			// ...

			// Log the processing of the budget
			k.Logger(ctx).Info(fmt.Sprintf("Processing budget for recipient: %s", recipientAddress.String()))
		}

		return false // Continue iterating
	})

	return nil
}

// EndBlocker is called at the end of every block
func _EndBlocker(ctx sdk.Context, keeper keeper.Keeper, bankKeeper types.BankKeeper) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	// Retrieve budget from the store
	budgetBytes := keeper.GetBudget(ctx) // Implement GetBudget in your Keeper
	if budgetBytes == nil {
		return
	}

	// Unmarshal the stored budget
	var budget types.Budget
	keeper.cdc.MustUnmarshalBinaryBare(budgetBytes, &budget)

	// Your logic to distribute the budget
	err := distributeBudget(ctx, keeper.bankKeeper, &budget)
	if err != nil {
		// Handle error, log, or panic depending on your use case
	}

	// Update the state, remove tranches, etc. (depends on your application logic)
	keeper.UpdateBudget(ctx, &budget) // Implement UpdateBudget in your Keeper
}
