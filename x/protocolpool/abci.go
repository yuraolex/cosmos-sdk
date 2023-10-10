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

func EndBlocker(ctx sdk.Context, k *keeper.Keeper, bankKeeper types.BankKeeper) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	logger := ctx.Logger().With("module", "x/"+types.ModuleName)

	// Read the current balance of all coins in the community pool
	communityPoolBalance, err := k.GetCommunityPool(ctx)
	if err != nil {
		return err
	}

	// Get the epoch balances from the store
	previousEpochBalance, err := k.GetEpochBalance(ctx.WithBlockHeight(ctx.BlockHeight() - 1))
	if err != nil {
		panic(fmt.Sprintf("failed to get epoch balance: %v", err))
	}

	// Calculate the balance dividend
	balanceDividend, isPositive := calculateBalanceDividend(communityPoolBalance, previousEpochBalance)

	// If balanceDividend is positive, distribute funds based on the budget
	if isPositive {
		// Get the current budget from the store
		currentBudget, err := k.GetBudget(ctx)
		if err != nil {
			return err
		}

		// distribute funds based on the budget
		distributeFunds(ctx, k, balanceDividend, currentBudget)
	}

	// Update the epoch balance for the next epoch (epoch n)
	k.UpdateEpochBalance(ctx, ctx.BlockHeight(), communityPoolBalance)
}

func calculateBalanceDividend(currentBalance, previousEpochBalance sdk.Coins) (sdk.Coins, bool) {
	// Calculate the balance dividend for each denom
	balanceDividend := sdk.Coins{}
	for _, coin := range currentBalance {
		previousAmount := previousEpochBalance.AmountOf(coin.Denom)
		dividend := coin.Amount.Sub(previousAmount)
		if dividend.IsPositive() {
			balanceDividend = balanceDividend.Add(sdk.NewCoin(coin.Denom, dividend))
			return balanceDividend, true
		}
	}
	return balanceDividend, false
}

func distributeFunds(ctx sdk.Context, k *keeper.Keeper, balDividend sdk.Coins, budget map[string]math.LegacyDec) error {
	totalPercentage := math.LegacyNewDec(0)

	// Calculate the total percentage in the budget
	for _, percentage := range budget {
		totalPercentage = totalPercentage.Add(percentage)
	}

	// Check if the total percentage is exactly 1
	if !totalPercentage.Equal(math.LegacyOneDec()) {
		return fmt.Errorf("total percentage in the budget must be exactly 1")
	}

	// Calculate the amount to distribute to each address based on the budget
	for address, percentage := range budget {
		// Calculate the amount based on the percentage of the balanceDividend
		// amount := balDividend.AmountOf(balDividend[0].Denom).ToDec().Mul(percentage).TruncateInt()
		amount := balDividend.AmountOf(balDividend[0].Denom).ToLegacyDec().Mul(percentage).TruncateInt()

		// Create a coin with the calculated amount
		coin := sdk.NewCoin(balDividend[0].Denom, amount)

		// Send the funds to the specified address
		recipient := k.authKeeper.AddressCodec().StringToBytes(address)
		err := k.DistributeFromFeePool(ctx, sdk.NewCoins(coin), recipient)
		if err != nil {
			return err
		}
	}

	return nil
}
