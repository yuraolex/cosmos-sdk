package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	"cosmossdk.io/x/protocolpool/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type Keeper struct {
	storeService storetypes.KVStoreService
	authKeeper   types.AccountKeeper
	bankKeeper   types.BankKeeper

	// State
	Schema collections.Schema
	Budget collections.Map[sdk.AccAddress, types.MsgSetBudgetProposal]

	authority string
}

func NewKeeper(cdc codec.BinaryCodec, storeService storetypes.KVStoreService,
	ak types.AccountKeeper, bk types.BankKeeper, authority string,
) Keeper {
	// ensure pool module account is set
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	sb := collections.NewSchemaBuilder(storeService)

	keeper := Keeper{
		storeService: storeService,
		authKeeper:   ak,
		bankKeeper:   bk,
		authority:    authority,
		Budget:       collections.NewMap(sb, types.BudgetKey, "budget", sdk.AccAddressKey, codec.CollValue[types.MsgSetBudgetProposal](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	keeper.Schema = schema

	return keeper
}

// GetAuthority returns the x/protocolpool module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With(log.ModuleKey, "x/"+types.ModuleName)
}

// FundCommunityPool allows an account to directly fund the community fund pool.
func (k Keeper) FundCommunityPool(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error {
	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, amount)
}

// DistributeFromFeePool distributes funds from the protocolpool module account to
// a receiver address.
func (k Keeper) DistributeFromFeePool(ctx context.Context, amount sdk.Coins, receiveAddr sdk.AccAddress) error {
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiveAddr, amount)
}

// GetCommunityPool get the community pool balance.
func (k Keeper) GetCommunityPool(ctx context.Context) (sdk.Coins, error) {
	moduleAccount := k.authKeeper.GetModuleAccount(ctx, types.ModuleName)
	if moduleAccount == nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", moduleAccount)
	}
	return k.bankKeeper.GetAllBalances(ctx, moduleAccount.GetAddress()), nil
}

func (k Keeper) ClaimFunds(ctx context.Context, recipient sdk.AccAddress) (sdk.Coins, error) {
	funds, err := k.withdrawFunds(ctx, recipient)
	if err != nil {
		return nil, err
	}

	// reinitialize the Claim
	err = k.initializeClaim(ctx, recipient)
	if err != nil {
		return nil, err
	}
	return funds, nil
}

func (k Keeper) HandleSetBudgetProposal(ctx context.Context, budget types.MsgSetBudgetProposal) error {
	// Validate the proposal content
	if err := k.validateSetBudgetProposal(ctx, budget); err != nil {
		return err
	}

	// Perform the logic to set the budget
	budget = types.MsgSetBudgetProposal{
		RecipientAddress:  budget.RecipientAddress,
		TotalBudget:       budget.TotalBudget,
		StartTime:         budget.StartTime,
		RemainingTranches: budget.RemainingTranches,
		Period:            budget.Period,
	}

	// Store the budget in the state
	if err := k.Budget.Set(ctx, sdk.AccAddress(budget.RecipientAddress), budget); err != nil {
		return err
	}
	// k.SetBudget(ctx, &budget)

	return nil
}

func (k Keeper) GetBudget(ctx context.Context, recipient sdk.AccAddress) (types.MsgSetBudgetProposal, error) {
	return k.Budget.Get(ctx, recipient)
}

// Validate the proposal content
func (k Keeper) validateSetBudgetProposal(ctx context.Context, proposal types.MsgSetBudgetProposal) error {
	account := k.authKeeper.GetAccount(ctx, sdk.AccAddress(proposal.RecipientAddress))
	if account == nil {
		return fmt.Errorf("account not found: %s", proposal.RecipientAddress)
	}

	if proposal.TotalBudget.IsZero() {
		return fmt.Errorf("total budget cannot be zero")
	}

	if proposal.StartTime <= 0 {
		return fmt.Errorf("start time must be a positive value")
	}

	if proposal.RemainingTranches <= 0 {
		return fmt.Errorf("remaining tranches must be a positive value")
	}

	if proposal.Period <= 0 {
		return fmt.Errorf("period must be a positive value")
	}

	return nil
}
