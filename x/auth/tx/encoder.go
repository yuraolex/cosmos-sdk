package tx

import (
	"fmt"

	"github.com/cosmos/gogoproto/proto"

	protov2 "google.golang.org/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
)

// DefaultTxEncoder returns a default protobuf TxEncoder using the provided Marshaler
func DefaultTxEncoder() sdk.TxEncoder {
	return func(tx sdk.Tx) ([]byte, error) {
		txWrapper, ok := tx.(*wrapper)
		if !ok {
			return nil, fmt.Errorf("expected %T, got %T", &wrapper{}, tx)
		}

		raw := &txtypes.TxRaw{
			BodyBytes:     txWrapper.getBodyBytes(),
			AuthInfoBytes: txWrapper.getAuthInfoBytes(),
			Signatures:    txWrapper.tx.Signatures,
		}

		return proto.Marshal(raw)
	}
}

// DefaultJSONTxEncoder returns a default protobuf JSON TxEncoder using the provided Marshaler.
func DefaultJSONTxEncoder(cdc codec.ProtoCodecMarshaler) sdk.TxEncoder {
	return func(tx sdk.Tx) ([]byte, error) {
		txWrapper, ok := tx.(*wrapper)
		if ok {
			fmt.Println("txWrapper.tx.GetMsgs()", txWrapper.tx.GetMsgs())
			v2Tx := txWrapper.GetProtoTxV2()
			bz := cdc.MustMarshal(txWrapper.tx)
			err := protov2.Unmarshal(bz, v2Tx)
			if err != nil {
				return nil, err
			}

			bz, err = cdc.MarshalJSON(v2Tx)
			fmt.Println(string(bz))

			return bz, err
		}

		return nil, fmt.Errorf("expected %T, got %T", &wrapper{}, tx)
	}
}
