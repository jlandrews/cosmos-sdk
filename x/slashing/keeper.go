package slashing

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/stake"
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

func (k Keeper) handleDoubleSign(ctx sdk.Context, pubkey []byte) {
	ctx.Logger().With("module", "slashing").Debug(fmt.Sprintf("Double sign from %s", string(pubkey)))
}

// TODO swap to pubkey, https://github.com/tendermint/abci/issues/239
func (k Keeper) handleAbsentValidator(ctx sdk.Context, pubkey []byte) {
	ctx.Logger().With("module", "slashing").Debug(fmt.Sprintf("Absent validator: %s", string(pubkey)))
}
