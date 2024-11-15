package zerogen

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
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

func getTemplateContent(templateFile string) ([]byte, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	var templateData []byte
	var templateFilePath string
	zeroGenDir := filepath.Join(homeDir, ".zero-gen")
	if _, err = os.Stat(zeroGenDir); os.IsNotExist(err) {
		goto fail
	}

	templateFilePath = filepath.Join(zeroGenDir, templateFile)
	if _, err = os.Stat(templateFilePath); os.IsNotExist(err) {
		goto fail
	}
	if templateData, err = os.ReadFile(templateFilePath); err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

fail:
	if templateData == nil {
		templateData = []byte(defaultTemplate) // Assuming defaultTemplate is a variable containing the default template content
	}
	return templateData, nil
}

func (z *ZeroGen) GenerateGormModel(db *gorm.DB, tableName string, templateFile string) (string, error) {
	// Get the table schema
	var columns []gorm.ColumnType
	columns, err := db.Migrator().ColumnTypes(tableName)
	if err != nil {
		return "", fmt.Errorf("failed to get column types: %w", err)
	}

	// Convert columns to ColumnSchema
	var columnSchemas []ColumnSchema
	for _, column := range columns {
		name := column.Name()
		nullable, _ := column.Nullable()
		length, _ := column.Length()
		defaultValue, _ := column.DefaultValue()
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

	// Get type mappings
	typeMappings, err := GetTypeMappings()
	if err != nil {
		return "", fmt.Errorf("failed to get type mappings: %w", err)
	}

	// Generate the struct
	structCode, err := GenerateStruct(tableName, columnSchemas, typeMappings, templateFile)
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

	res, err := z.GenerateGormModel(db, z.Table, "gorm_model.tmpl")
	if err != nil {
		return fmt.Errorf("failed to generate gorm model: %w", err)
	}
	fmt.Println(res)
	return nil
}
