package ormdb

import (
	"context"
	"fmt"

	"github.com/cosmos/gogoproto/proto"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	ormv1alpha1 "cosmossdk.io/api/cosmos/orm/v1alpha1"
	"github.com/cosmos/cosmos-sdk/orm/model/ormtable"
)

func (m moduleDB) AutoMigrate(ctx context.Context) error {
	schema, err := m.readSchema(ctx)
	if err != nil {
		return err
	}

	return m.MigrateFrom(ctx, schema)
}

func (m moduleDB) MigrateFrom(ctx context.Context, oldSchema *ormv1alpha1.ModuleSchemaRecord) error {
	if oldSchema != nil && oldSchema.Version > 0 {
		return fmt.Errorf("do not know how to migrate from a schema with version %d", oldSchema.Version)
	}

	oldFilesMap := map[uint32]*ormv1alpha1.ModuleSchemaRecord_FileRecord{}
	if oldSchema != nil {
		for _, file := range oldSchema.Files {
			oldFilesMap[file.Id] = file
		}
	}
	oldFileIds := maps.Keys(oldFilesMap)
	slices.Sort(oldFileIds)

	for _, id := range oldFileIds {
		_, ok := m.filesById[id]
		if !ok {
			// TODO remove file
			//kvStore := m.kvStoreService.OpenKVStore(ctx)
		}
	}

	newFileIds := maps.Keys(m.filesById)
	slices.Sort(newFileIds)

	newSchema := &ormv1alpha1.ModuleSchemaRecord{
		Version: 0,
		Files:   nil,
	}

	for _, id := range newFileIds {
		fileDb := m.filesById[id]
		if fileDb.storageType != ormv1alpha1.StorageType_STORAGE_TYPE_DEFAULT_UNSPECIFIED {
			// don't worry about memory and transient tables
			continue
		}

		fileSchema, err := fileDb.MigrateFrom(ctx, oldFilesMap[id])
		if err != nil {
			return err
		}

		newSchema.Files = append(newSchema.Files, fileSchema)
	}

	return m.saveSchema(ctx, newSchema)
}

func (f fileDescriptorDB) MigrateFrom(ctx context.Context, oldSchema *ormv1alpha1.ModuleSchemaRecord_FileRecord) (*ormv1alpha1.ModuleSchemaRecord_FileRecord, error) {
	if oldSchema != nil && oldSchema.ProtoFilePath != f.fileDescriptor.Path() {
		return nil, fmt.Errorf("can't migrate from schema proto file %s to %s", oldSchema.ProtoFilePath, f.fileDescriptor.Path())
	}

	oldSchemasMap := map[uint32]*ormv1alpha1.ModuleSchemaRecord_TableRecord{}
	if oldSchema != nil {
		for _, record := range oldSchema.Tables {
			oldSchemasMap[record.Id] = record
		}
	}

	oldTableIds := maps.Keys(oldSchemasMap)
	slices.Sort(oldTableIds)

	newTables := map[uint32]ormtable.Table{}
	for id, table := range f.tablesById {
		newTables[id] = table
	}
	newTableIds := maps.Keys(newTables)
	slices.Sort(newTableIds)

	newTableRecords := map[uint32]*ormv1alpha1.ModuleSchemaRecord_TableRecord{}

	for _, id := range oldTableIds {
		_, ok := newTables[id]
		if !ok {
			// TODO remove deleted table
			//kvStore := m.kvStoreService.OpenKVStore(ctx)
		}
	}

	for _, id := range newTableIds {
		oldSchema := oldSchemasMap[id]
		newTable := newTables[id]

		newTableMigrate := newTable.(interface {
			MigrateFrom(context.Context, *ormv1alpha1.ModuleSchemaRecord_TableRecord) (*ormv1alpha1.ModuleSchemaRecord_TableRecord, error)
		})

		tableRecord, err := newTableMigrate.MigrateFrom(ctx, oldSchema)
		if err != nil {
			return nil, err
		}

		newTableRecords[id] = tableRecord
	}

	return nil, nil
}

func (m moduleDB) saveSchema(ctx context.Context, schema *ormv1alpha1.ModuleSchemaRecord) error {
	bz, err := proto.Marshal(schema)
	if err != nil {
		return err
	}

	kvStore := m.kvStoreService.OpenKVStore(ctx)
	return kvStore.Set(m.schemaCodec.Prefix(), bz)
}

func (m moduleDB) readSchema(ctx context.Context) (*ormv1alpha1.ModuleSchemaRecord, error) {
	kvStore := m.kvStoreService.OpenKVStore(ctx)
	bz, err := kvStore.Get(m.schemaCodec.Prefix())
	if err != nil {
		return nil, err
	}

	schema := &ormv1alpha1.ModuleSchemaRecord{}
	err = proto.Unmarshal(bz, schema)
	if err != nil {
		return nil, err
	}

	return schema, nil
}
