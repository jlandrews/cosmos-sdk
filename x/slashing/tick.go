package slashing

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
)

func NewBeginBlocker(sk Keeper) sdk.BeginBlocker {
	return func(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
		heightBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(heightBytes, uint64(req.Header.Height))
		tags := sdk.NewTags("height", heightBytes)
		for _, evidence := range req.ByzantineValidators {
			var pk crypto.PubKey
			sk.cdc.MustUnmarshalBinary(evidence.PubKey, &pk)
			sk.handleDoubleSign(ctx, evidence.Height, evidence.Time, pk)
		}
		absent := make(map[string]bool)
		for _, pubkey := range req.AbsentValidators {
			var pk crypto.PubKey
			sk.cdc.MustUnmarshalBinary(pubkey, &pk)
			absent[string(pk.Bytes())] = true
		}
		sk.stakeKeeper.IterateValidatorsBonded(ctx, func(_ int64, validator sdk.Validator) (stop bool) {
			pubkey := validator.GetPubKey()
			sk.handleValidatorSignature(ctx, pubkey, !absent[string(pubkey.Bytes())])
			return false
		})
		// TODO some tags
		return abci.ResponseBeginBlock{
			Tags: tags.ToKVPairs(),
		}
	}
}
