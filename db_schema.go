package zerogen

import (
	"fmt"

	"gorm.io/gorm"
)

func GetTableSchema(db *gorm.DB, tableName string) error {

	var columns []gorm.ColumnType
	columns, err := db.Migrator().ColumnTypes(tableName)
	if err != nil {
		return err
	}

	for _, column := range columns {
		name := column.Name()
		nullable, _ := column.Nullable()
		length, _ := column.Length()
		defaultValue, _ := column.DefaultValue()
		dataType, _ := column.ColumnType()
		fmt.Printf("Column: %s, Type: %s, Nullable: %v, Length: %v, Default: %v\n", name, dataType, nullable, length, defaultValue)
	}

	return nil
}
