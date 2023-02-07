package ormtable

import (
	"context"
	"fmt"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	ormv1 "cosmossdk.io/api/cosmos/orm/v1"
	ormv1alpha1 "cosmossdk.io/api/cosmos/orm/v1alpha1"
	"github.com/cosmos/cosmos-sdk/orm/internal/fieldnames"
)

func (t *tableImpl) MigrateFrom(ctx context.Context, record *ormv1alpha1.ModuleSchemaRecord_TableRecord) (*ormv1alpha1.ModuleSchemaRecord_TableRecord, error) {
	msgName := string(t.MessageType().Descriptor().FullName())
	if msgName != record.ProtoMsgName {
		return nil, fmt.Errorf("cannot migrate from %s to %s", record.ProtoMsgName, msgName)
	}

	var oldTableDesc *ormv1.TableDescriptor
	switch desc := record.Desc.(type) {
	case *ormv1alpha1.ModuleSchemaRecord_TableRecord_Singleton:
		return nil, fmt.Errorf("cannot migrate from a singleton to a table for %s", msgName)
	case *ormv1alpha1.ModuleSchemaRecord_TableRecord_Table:
		oldTableDesc = desc.Table
	default:
		return nil, fmt.Errorf("unexpected case")
	}

	newTableDesc := t.tableDesc

	//
	// check primary key
	//
	if !keysEqual(oldTableDesc.PrimaryKey.Fields, newTableDesc.PrimaryKey.Fields) {
		return nil, fmt.Errorf("cannot change primary key of table %s from %s to %s", msgName, oldTableDesc.PrimaryKey.Fields, newTableDesc.PrimaryKey.Fields)
	}

	if !oldTableDesc.PrimaryKey.AutoIncrement && newTableDesc.PrimaryKey.AutoIncrement {
		return nil, fmt.Errorf("cannot migrate from a non-auto-increment primary key to an auto-increment primary key for %s", msgName)
	}
	// technically we can migrate from an auto-increment to a non-auto-increment table even though that is API breaking

	//
	// check indexes
	//
	oldIndexes := indexesMap(oldTableDesc.Index)
	newIndexes := indexesMap(newTableDesc.Index)
	oldIndexIds := maps.Keys(oldIndexes)
	slices.Sort(oldIndexIds)
	for _, id := range oldIndexIds {
		oldIndex := oldIndexes[id]
		newIndex, ok := newIndexes[id]
		if !ok {
			// delete removed index
			idx := t.indexesById[id]
			err := idx.DeleteBy(ctx)
			if err != nil {
				return nil, err
			}
		}

		if !keysEqual(oldIndex.Fields, newIndex.Fields) {
			return nil, fmt.Errorf("cannot change fields of index %d on table %s from %s to %s", id, msgName, oldIndex.Fields, newIndex.Fields)
		}

		if !oldIndex.Unique && newIndex.Unique {
			return nil, fmt.Errorf("cannot add unique constraint on index %d on table %s", id, msgName)
		}

		if oldIndex.Unique && !newIndex.Unique {
			return nil, fmt.Errorf("cannot remove unique constraint on index %d on table %s - unique and non-unique indexes have incompatible encodings", id, msgName)
		}
	}

	// check for newly added index
	newIndexIds := maps.Keys(newIndexes)
	slices.Sort(newIndexIds)
	for _, id := range newIndexIds {
		_, ok := oldIndexes[id]
		if ok {
			// index already existed
			continue
		}

		// build new index

		idx := t.indexesById[id].(concreteIndex)

		it, err := t.List(ctx, nil)
		if err != nil {
			return nil, err
		}

		backend, err := t.getWriteBackend(ctx)
		if err != nil {
			return nil, err
		}

		batch := newBatchIndexCommitmentWriter(backend)
		idxStore := batch.IndexStore()
		for {
			if !it.Next() {
				break
			}

			msg, err := it.GetMessage()
			if err != nil {
				return nil, err
			}

			k, v, err := idx.EncodeKVFromMessage(msg.ProtoReflect())
			if err != nil {
				return nil, err
			}

			err = idxStore.Set(k, v)
			if err != nil {
				return nil, err
			}
		}
		it.Close()

		err = batch.Write()
		if err != nil {
			return nil, err
		}
	}

	return &ormv1alpha1.ModuleSchemaRecord_TableRecord{
		Id:           t.tableId,
		ProtoMsgName: msgName,
		Desc:         &ormv1alpha1.ModuleSchemaRecord_TableRecord_Table{Table: newTableDesc},
	}, nil
}

func (t *singleton) MigrateFrom(ctx context.Context, record *ormv1alpha1.ModuleSchemaRecord_TableRecord) (*ormv1alpha1.ModuleSchemaRecord_TableRecord, error) {
	msgName := string(t.MessageType().Descriptor().FullName())
	if msgName != record.ProtoMsgName {
		return nil, fmt.Errorf("cannot migrate from %s to %s", record.ProtoMsgName, msgName)
	}

	switch record.Desc.(type) {
	case *ormv1alpha1.ModuleSchemaRecord_TableRecord_Singleton:
		return record, nil
	case *ormv1alpha1.ModuleSchemaRecord_TableRecord_Table:
		return nil, fmt.Errorf("cannot migrate from a table to a singleton for %s", msgName)
	default:
		return nil, fmt.Errorf("unexpected case")
	}
}

func keysEqual(key1, key2 string) bool {
	k1 := fieldnames.CommaSeparatedFieldNames(key1)
	k2 := fieldnames.CommaSeparatedFieldNames(key2)
	return k1.String() == k2.String()
}

func indexesMap(indexes []*ormv1.SecondaryIndexDescriptor) map[uint32]*ormv1.SecondaryIndexDescriptor {
	res := map[uint32]*ormv1.SecondaryIndexDescriptor{}
	for _, index := range indexes {
		res[index.Id] = index
	}
	return res
}
