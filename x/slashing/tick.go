package slashing

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
)

func NewBeginBlocker(sk Keeper) sdk.BeginBlocker {
	return func(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
		for _, evidence := range req.ByzantineValidators {
			var pk crypto.PubKey
			sk.cdc.MustUnmarshalBinary(evidence.PubKey, &pk)
			sk.handleDoubleSign(ctx, evidence.Height, evidence.Time, pk)
		}
		for _, pubkey := range req.AbsentValidators {
			var pk crypto.PubKey
			sk.cdc.MustUnmarshalBinary(pubkey, &pk)
			sk.handleAbsentValidator(ctx, pk)
		}
		return abci.ResponseBeginBlock{}
	}
}
