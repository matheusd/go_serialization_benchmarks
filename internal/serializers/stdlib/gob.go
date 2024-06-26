package stdlib

import (
	"bytes"
	"encoding/gob"

	"github.com/alecthomas/go_serialization_benchmarks/goserbench"
)

type GobSerializer struct{}

func (g *GobSerializer) Marshal(o interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(o)
	return buf.Bytes(), err
}

func (g *GobSerializer) Unmarshal(d []byte, o interface{}) error {
	return gob.NewDecoder(bytes.NewReader(d)).Decode(o)
}

func NewGobSerializer() goserbench.Serializer {
	// registration required before first use
	gob.Register(goserbench.SmallStruct{})
	return &GobSerializer{}
}
