package zerogen

import (
	"fmt"

	"gorm.io/gorm"
)

// ColumnSchema represents the schema of a database column
type ColumnSchema struct {
	ColumnName   string  // Column name, e.g., "id", "created_at"
	ColumnType   string  // Column type, e.g., "uuid", "text"
	Nullable     bool    // Column nullable, e.g., true, false
	Length       int64   // Column length, e.g., 128, 9223372036854775807
	DefaultValue *string // Default value of the column, e.g., "now()"
}

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
		columnSchema := ColumnSchema{
			ColumnName:   name,
			ColumnType:   dataType,
			Nullable:     nullable,
			Length:       length,
			DefaultValue: defaultValueStr,
		}
		columnSchemas = append(columnSchemas, columnSchema)
		fmt.Println(columnSchema)
	}

	return columnSchemas, nil
}
