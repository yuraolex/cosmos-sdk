package types

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Define the budget structure
type BudgetItem struct {
	Address sdk.AccAddress
	Weight  math.LegacyDec
}

type Budget struct {
	Items []BudgetItem
}

// Define the epoch balance structure
type EpochBalanceItem struct {
	Denom  string
	Amount math.Int
}

type EpochBalance struct {
	Items []EpochBalanceItem
}

// NewMsgSetBudget is a constructor function for MsgSetBudget
func NewMsgSetBudget(proposer sdk.AccAddress, budget Budget) *MsgSetBudget {
	return &MsgSetBudget{
		Proposer: proposer.String(),
		Budget:   budget,
	}
}

// MsgSetBudget is the message type for setting the budget through governance
type MsgSetBudget struct {
	Proposer string
	Budget   Budget
}
