package zerogen

import (
	_ "embed"
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type ZeroGen struct {
	Dsn            string `clop:"long" usage:"database dsn"`
	ModelOutPutDir string `clop:"long" usage:"model output directory" default:"."`
	Table          string `clop:"long" usage:"table name"`
}

func GenerateGormModel(db *gorm.DB,
	tableName string,
	columnSchemas []ColumnSchema,
	typeMappings map[string]TypeMapping) (string, error) {
	// Generate the struct
	structCode, err := GenerateStruct(tableName, columnSchemas, typeMappings)
	if err != nil {
		return "", fmt.Errorf("failed to generate struct: %w", err)
	}

	return structCode, nil
}

func (z *ZeroGen) Run() error {

	var db *gorm.DB
	var err error

	if strings.Contains(z.Dsn, ":3306") {
		db, err = gorm.Open(mysql.Open(z.Dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{SingularTable: true},
			Logger:         logger.Default.LogMode(logger.Info),
		})
	} else {
		db, err = gorm.Open(postgres.Open(z.Dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{SingularTable: true},
			Logger:         logger.Default.LogMode(logger.Info),
		})
	}
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	columns, err := GetTableSchema(db, z.Table)
	if err != nil {
		return fmt.Errorf("failed to get table schema: %w", err)
	}

	// Get type mappings
	typeMappings, err := GetTypeMappings()
	if err != nil {
		return fmt.Errorf("failed to get type mappings: %w", err)
	}

	res, err := GenerateGormModel(db, z.Table, columns, typeMappings)
	if err != nil {
		return fmt.Errorf("failed to generate gorm model: %w", err)
	}

	fmt.Println(res)

	res, err = GenerateApiService(z.Table, columns, typeMappings, "api", "v1", z.Table, z.Table)
	if err != nil {
		return fmt.Errorf("failed to generate go zero api: %w", err)
	}
	fmt.Println(res)
	return nil
}
