package types

import (
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// type ProtocolPoolConfig struct {
// 	AddressPercentage map[string]sdk.Dec `protobuf:"bytes,1,opt,name=address_percentage,json=addressPercentage,proto3" json:"address_percentage,omitempty"`
// }

// type BudgetItem struct {
// 	Address        sdk.AccAddress       `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
// 	Weight         sdk.Dec              `protobuf:"bytes,2,opt,name=weight,proto3" json:"weight"`
// 	ProtocolConfig *ProtocolPoolConfig `protobuf:"bytes,3,opt,name=protocol_pool_config,json=protocolPoolConfig,proto3" json:"protocol_pool_config,omitempty"`
// }

// type Budget struct {
// 	Items []*BudgetItem `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
// }

// NewMsgSetBudget is a constructor function for MsgSetBudget
func NewMsgSetBudget(proposer sdk.AccAddress, budget *Budget) *MsgSetBudget {
	return &MsgSetBudget{
		Proposer: proposer.String(),
		Budget:   budget,
	}
}

// ValidateBasic implements MsgSetBudget validation
func (msg *MsgSetBudget) ValidateBasic() error {
	// Implement validation logic as needed
	// Example: Check if proposer is not empty, and the budget is not nil
	if msg.Proposer == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "proposer cannot be empty")
	}
	if msg.Budget == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "budget cannot be nil")
	}

	// Add more validation as needed
	

	return nil
}
