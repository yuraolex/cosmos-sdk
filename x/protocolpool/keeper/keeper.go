package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/core/store"
	"cosmossdk.io/errors"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
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
	Budget collections.Map[sdk.AccAddress, types.Budget]

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
		Budget:       collections.NewMap(sb, types.BudgetKey, "budget", sdk.AccAddressKey, codec.CollValue[types.Budget](cdc)),
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

func (k Keeper) GetEpochBalance(ctx context.Context) (sdk.Coins, error) {
	epochBal := math.NewInt(0)
	cp, err := k.GetCommunityPool(ctx)
	if err != nil {
		return nil, err
	}
	for _, coin := range cp {
		epochBal.Add(coin.Amount)
	}
	bal := sdk.NewCoins(sdk.NewCoin())
	return epochBal, nil
}

func (k Keeper) GetBudget(ctx sdk.Context) (map[string]math.LegacyDec, error) {
	var budget map[string]math.LegacyDec
	err := k.Budget.Walk(ctx, nil, func(key sdk.AccAddress, value types.Budget) (stop bool, err error) {
		for _, item := range value.Items {
			budget[item.Address] = item.Weight.Value
		}
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	return budget, nil
}

func (k Keeper) SetBudget(ctx sdk.Context, msg *types.MsgSetBudget) error {
	// Validate the budget message
	if err := k.validateBudget(ctx, msg); err != nil {
		return err
	}

	// set budget
	if err := k.Budget.Set(ctx, sdk.AccAddress(msg.Proposer), *msg.Budget); err != nil {
		return err
	}

	// Emit an event for the budget update
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetBudget,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer),
		),
	)

	return nil
}

func (k Keeper) validateBudget(ctx context.Context, msg *types.MsgSetBudget) error {
	// Perform validation checks

	// Check all account addresses exist
	for _, item := range msg.Budget.Items {
		account := k.authKeeper.GetAccount(ctx, sdk.AccAddress(item.Address))
		if account == nil {
			return fmt.Errorf("account not found: %s", item.Address)
		}
	}

	// Check the sum of all stream amounts equals exactly 1
	sum := math.LegacyNewDec(0)
	for _, item := range msg.Budget.Items {
		sum = sum.Add(item.Weight.Value)
	}
	if !sum.Equal(math.LegacyOneDec()) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "sum of budget weights must be equal to 1")
	}

	return nil
}
