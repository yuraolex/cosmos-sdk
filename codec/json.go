package codec

import (
	"bytes"
	"fmt"

	"github.com/cosmos/gogoproto/jsonpb"
	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/protobuf/encoding/protojson"
	protov2 "google.golang.org/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/codec/types"
)

var defaultJM = &jsonpb.Marshaler{OrigName: true, EmitDefaults: true, AnyResolver: nil}

// ProtoMarshalJSON provides an auxiliary function to return Proto3 JSON encoded
// bytes of a message.
func ProtoMarshalJSON(msg proto.Message, resolver jsonpb.AnyResolver) ([]byte, error) {
	switch protoMsg := msg.(type) {
	case protov2.Message:
		return protojson.Marshal(protoMsg)
	case proto.Message:
		// We use the OrigName because camel casing fields just doesn't make sense.
		// EmitDefaults is also often the more expected behavior for CLI users
		jm := defaultJM
		//jm2, err := protojson.Marshal(msg)
		if resolver != nil {
			jm = &jsonpb.Marshaler{OrigName: true, EmitDefaults: true, AnyResolver: resolver}
		}
		err := types.UnpackInterfaces(msg, types.ProtoJSONPacker{JSONPBMarshaler: jm})
		if err != nil {
			return nil, err
		}

		buf := new(bytes.Buffer)
		if err := jm.Marshal(buf, msg); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	default:
		return nil, fmt.Errorf("cannot proto marshal unsupported type: %T", msg)
	}
}
