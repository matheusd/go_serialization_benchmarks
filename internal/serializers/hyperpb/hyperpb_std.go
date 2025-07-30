package hyperpb

import (
	"time"

	"buf.build/go/hyperpb"
	"github.com/alecthomas/go_serialization_benchmarks/goserbench"
	"google.golang.org/protobuf/proto"
)

type HyperpbStdSerializer struct {
	ty *hyperpb.MessageType
}

func (b *HyperpbStdSerializer) Marshal(o interface{}) ([]byte, error) {
	// Note: Hyperb does not provide a marshaller, so we use the standard
	// protobuf marshaller.
	v := o.(*goserbench.SmallStruct)
	return protobufMarshal(v)
}

func (b *HyperpbStdSerializer) Unmarshal(data []byte, o interface{}) error {
	// Unmarshal it, just how you normally would.
	msg := hyperpb.NewMessage(b.ty)
	if err := proto.Unmarshal(data, msg); err != nil {
		return err
	}

	// Use reflection to read some fields.
	fields := b.ty.Descriptor().Fields()

	a := o.(*goserbench.SmallStruct)
	a.Name = msg.Get(fields.ByName("name")).String()
	a.BirthDay = time.Unix(0, msg.Get(fields.ByName("birthDay")).Int())
	a.Phone = msg.Get(fields.ByName("phone")).String()
	a.Siblings = int(msg.Get(fields.ByName("siblings")).Int())
	a.Spouse = msg.Get(fields.ByName("spouse")).Bool()
	a.Money = msg.Get(fields.ByName("money")).Float()

	return nil
}

func NewHyperpbStdSerializer() goserbench.Serializer {
	// Compile a type for your message. This operation is quite slow, so it
	// should be cached, like regexp.Compile.
	ty := hyperpb.CompileMessageDescriptor((*SmallStruct)(nil).ProtoReflect().Descriptor())

	return &HyperpbStdSerializer{ty: ty}
}
