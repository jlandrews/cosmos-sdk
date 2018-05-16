package slashing

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
)

func NewBeginBlocker(sk Keeper) sdk.BeginBlocker {
	return func(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
		for _, evidence := range req.ByzantineValidators {
			sk.handleDoubleSign(ctx, evidence.PubKey)
		}
		for _, pubkey := range req.AbsentValidators {
			sk.handleAbsentValidator(ctx, pubkey)
		}
		return abci.ResponseBeginBlock{}
	}
}
