package protocolpool

import (
	"fmt"
	"time"

	"cosmossdk.io/x/protocolpool/keeper"
	"cosmossdk.io/x/protocolpool/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, k *keeper.Keeper, bankKeeper types.BankKeeper) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	logger := ctx.Logger().With("module", "x/"+types.ModuleName)

	// Read the current balance of all coins in the community pool
	communityPoolBalance, err := k.GetCommunityPool(ctx)
	if err != nil {
		return err
	}

	// Get the epoch balances from the store
	epochBalance, err := k.GetEpochBalance(ctx, ctx.BlockHeight()-1)
	if err != nil {
		panic(fmt.Sprintf("failed to get epoch balance: %v", err))
	}

	// Calculate the balance dividend
	balanceDividend := calculateBalanceDividend(communityPoolBalance, epochBalance)

	// Get the current budget from the store
	currentBudget, err := k.GetBudget(ctx)
	if err != nil {
		return err
	}

	// Transfer funds based on the budget
	transferFunds(ctx, bankKeeper, balanceDividend, currentBudget)
}

func calculateBalanceDividend(currentBalance sdk.Coins, epochBalance types.EpochBalance) sdk.Coins {
	// Perform the described calculation and return the balance dividend
	// Subtract each epoch_balance from epoch n-1 from the corresponding values in epoch n
	// If the balance dividend is positive (coins were added), return the result
	// Otherwise, return an empty Coins
	// Example:
	// return sdk.NewCoins(...)
}

func transferFunds(ctx sdk.Context, bankKeeper types.BankKeeper, balanceDividend sdk.Coins, budget types.Budget) {
	// Transfer funds based on the budget
	// Multiply the balance dividend by each value in the budget map
	// Transfer the resulting amount from the community pool to the corresponding address
	// Example:
	// for _, item := range budget.Items {
	//     amount := balanceDividend.Amount.Mul(item.Weight)
	//     err := bankKeeper.SendCoinsFromPoolToAccount(ctx, types.CommunityPoolName, item.Address, sdk.NewCoins(amount...))
	//     if err != nil {
	//         panic(fmt.Sprintf("failed to transfer funds: %v", err))
	//     }
	// }
}
