package hyperpb

import (
	"github.com/alecthomas/go_serialization_benchmarks/goserbench"
	"google.golang.org/protobuf/proto"
)

// protobufMarshal does the standard protobuf marshalling.
func protobufMarshal(v *goserbench.SmallStruct) ([]byte, error) {
	// Note: Hyperb does not provide a marshaller, so we use the standard
	// protobuf marshaller.
	var a SmallStruct
	a.Name = v.Name
	a.BirthDay = v.BirthDay.UnixNano()
	a.Phone = v.Phone
	a.Siblings = int32(v.Siblings)
	a.Spouse = v.Spouse
	a.Money = v.Money
	return proto.Marshal(&a)

}
