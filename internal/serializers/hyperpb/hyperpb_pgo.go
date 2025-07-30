package hyperpb

import (
	"time"

	"buf.build/go/hyperpb"
	"github.com/alecthomas/go_serialization_benchmarks/goserbench"
	"google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

type HyperpbPGOSerializer struct {
	ty   *hyperpb.MessageType
	ctx  *hyperpb.Shared
	opts []hyperpb.UnmarshalOption
}

func (b *HyperpbPGOSerializer) Marshal(o interface{}) ([]byte, error) {
	// Note: Hyperb does not provide a marshaller, so we use the standard
	// protobuf marshaller.
	v := o.(*goserbench.SmallStruct)
	return protobufMarshal(v)
}

func (b *HyperpbPGOSerializer) Unmarshal(data []byte, o interface{}) error {
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

func compilePGO(md protoreflect.MessageDescriptor, corpus [][]byte) (*hyperpb.MessageType, error) {
	// Compile the type without any profiling information.
	msgType := hyperpb.CompileMessageDescriptor(md)

	// Construct a new profile recorder.
	profile := msgType.NewProfile()

	// Parse all of the specimens in the corpus, making sure to record a
	// profile for all of them.
	s := new(hyperpb.Shared)
	for _, specimen := range corpus {
		if err := s.NewMessage(msgType).Unmarshal(
			specimen,
			hyperpb.WithRecordProfile(profile, 1.0),
		); err != nil {
			return nil, err
		}
		s.Free()
	}

	// Recompile with the profile.
	return msgType.Recompile(profile), nil
}

func NewHyperpbPGOSerializer() goserbench.Serializer {
	// Generate a corpus of marshalled sample messages for the compiler to
	// figure out how to speed things up.
	var smallStructCorpus [][]byte
	samples := goserbench.GenerateSmallStruct(1000)
	for _, s := range samples {
		buf, err := protobufMarshal(s)
		if err != nil {
			panic(err)
		}
		smallStructCorpus = append(smallStructCorpus, buf)
	}

	// Compile a type for your message. This operation is quite slow, so it
	// should be cached, like regexp.Compile.
	ty, err := compilePGO((*SmallStruct)(nil).ProtoReflect().Descriptor(), smallStructCorpus)
	if err != nil {
		panic(err)
	}

	ctx := new(hyperpb.Shared)
	o := hyperpb.WithAllowAlias(true)

	return &HyperpbPGOSerializer{
		ty:   ty,
		ctx:  ctx,
		opts: []hyperpb.UnmarshalOption{o},
	}
}
