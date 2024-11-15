package zerogen

import (
	"fmt"

	"gorm.io/gorm"
)

func GetTableSchema(db *gorm.DB, tableName string) ([]ColumnSchema, error) {
	// Get the table schema
	columns, err := db.Migrator().ColumnTypes(tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get column types: %w", err)
	}

	// Convert columns to ColumnSchema
	var columnSchemas []ColumnSchema
	for _, column := range columns {
		name := column.Name()
		nullable, _ := column.Nullable()
		length, _ := column.Length()
		dataType, _ := column.ColumnType()
		defaultValue, ok := column.DefaultValue()
		var defaultValueStr *string
		if ok {
			defaultValueStr = &defaultValue
		}

		columnSchemas = append(columnSchemas, ColumnSchema{
			ColumnName:   name,
			ColumnType:   dataType,
			Nullable:     nullable,
			Length:       length,
			DefaultValue: defaultValueStr,
		})
	}

	return columnSchemas, nil
}
