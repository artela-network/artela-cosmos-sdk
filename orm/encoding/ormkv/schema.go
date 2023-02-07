package ormkv

import (
	"bytes"

	"google.golang.org/protobuf/proto"

	ormv1alpha1 "cosmossdk.io/api/cosmos/orm/v1alpha1"
	"github.com/cosmos/cosmos-sdk/orm/types/ormerrors"
)

type SchemaCodec struct {
	prefix []byte
}

var schemaPrefix = []byte{0, 0}

func NewSchemaCodec(prefix []byte) *SchemaCodec {
	var p []byte
	p = append(p, prefix...)
	p = append(p, schemaPrefix...)
	return &SchemaCodec{prefix: p}
}

var _ EntryCodec = &SchemaCodec{}

func (s SchemaCodec) DecodeEntry(k, v []byte) (Entry, error) {
	if !bytes.Equal(k, s.prefix) {
		return nil, ormerrors.UnexpectedDecodePrefix
	}

	schema := &ormv1alpha1.ModuleSchemaRecord{}
	err := proto.Unmarshal(v, schema)
	if err != nil {
		return nil, err
	}

	return &SchemaEntry{Schema: schema}, nil
}

func (s SchemaCodec) EncodeEntry(entry Entry) (k, v []byte, err error) {
	schemaEntry, ok := entry.(*SchemaEntry)
	if !ok {
		return nil, nil, ormerrors.BadDecodeEntry
	}

	bz, err := proto.Marshal(schemaEntry.Schema)
	if err != nil {
		return nil, nil, err
	}

	return s.prefix, bz, nil
}

func (s SchemaCodec) Prefix() []byte {
	return s.prefix
}
