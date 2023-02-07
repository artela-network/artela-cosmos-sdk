package ormdb

import (
	"context"

	ormv1alpha1 "cosmossdk.io/api/cosmos/orm/v1alpha1"
)

func (m moduleDB) AutoMigrate(ctx context.Context) error {
	schema, err := m.readSchema(ctx)
	if err != nil {
		return err
	}

	return m.MigrateFrom(ctx, schema)
}

func (m moduleDB) MigrateFrom(ctx context.Context, oldSchema *ormv1alpha1.ModuleSchemaRecord) error {

	return m.SaveSchema(ctx)
}

func (m moduleDB) SaveSchema(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (m moduleDB) readSchema(ctx context.Context) (*ormv1alpha1.ModuleSchemaRecord, error) {
	//TODO implement me
	panic("implement me")
}

func (m moduleDB) makeSchema() *ormv1alpha1.ModuleSchemaRecord {
	return nil
}
