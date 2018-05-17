package slashing

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/stake"
	crypto "github.com/tendermint/go-crypto"
)

const (
	// MaxEvidenceAge - Max age for evidence - 21 days (3 weeks)
	// TODO Should this be a governance parameter or just modifiable with SoftwareUpgradeProposals?
	// MaxEvidenceAge = 60 * 60 * 24 * 7 * 3
	// TODO Temporarily set to 2 minutes for testnets.
	MaxEvidenceAge = 60 * 2

	// SignedBlocksWindow - sliding window for liveness slashing
	SignedBlocksWindow
)

var (
	// SlashFractionDoubleSign - currently 5% - should be a governance parameter
	SlashFractionDoubleSign = sdk.NewRat(1).Quo(sdk.NewRat(20))
)

// Keeper of the slashing store
type Keeper struct {
	storeKey    sdk.StoreKey
	cdc         *wire.Codec
	stakeKeeper stake.Keeper

	// codespace
	codespace sdk.CodespaceType
}

// NewKeeper creates a slashing keeper
func NewKeeper(cdc *wire.Codec, key sdk.StoreKey, sk stake.Keeper, codespace sdk.CodespaceType) Keeper {
	keeper := Keeper{
		storeKey:    key,
		cdc:         cdc,
		stakeKeeper: sk,
		codespace:   codespace,
	}
	return keeper
}

// handle a validator signing two blocks at the same height
func (k Keeper) handleDoubleSign(ctx sdk.Context, height int64, timestamp int64, pubkey crypto.PubKey) {
	logger := ctx.Logger().With("module", "x/slashing")
	age := ctx.BlockHeader().Time - timestamp
	if age > MaxEvidenceAge {
		logger.Info(fmt.Sprintf("Ignored double sign from %v at height %d, age of %d past max age of %d", pubkey.Address(), height, age, MaxEvidenceAge))
		return
	}
	logger.Info(fmt.Sprintf("Confirmed double sign from %v at height %d, age of %d less than max age of %d", pubkey.Address(), height, age, MaxEvidenceAge))
	validator := k.stakeKeeper.Validator(ctx, pubkey.Address())
	validator.Slash(ctx, SlashFractionDoubleSign)
	logger.Info(fmt.Sprintf("Slashed validator %s by fraction %v", validator.GetAddress(), SlashFractionDoubleSign))
}

// handle an absent validator
func (k Keeper) handleAbsentValidator(ctx sdk.Context, pubkey crypto.PubKey) {
	logger := ctx.Logger().With("module", "x/slashing")
	height := ctx.BlockHeight()
	logger.Info(fmt.Sprintf("Absent validator %v at height %d", pubkey.Address(), height))
	// store := ctx.KVStore(k.storeKey)
}

func (k Keeper) GetValidatorSigningInfo(ctx sdk.Context, address sdk.Address) ValidatorSigningInfo {
	// TODO
	panic("todo")
}

func (k Keeper) setValidatorSigningInfo(ctx sdk.Context, address sdk.Address, info ValidatorSigningInfo) {
	// TODO
	panic("todo")
}

type ValidatorSigningInfo struct {
	StartHeight int64 `json:"start_height"`
	// TODO signed blocks array, separate keys + keep counter + optimizations, update spec
}
