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
	Dsn                  string `clop:"long" usage:"database dsn"`
	ModelDir             string `clop:"long" usage:"gorm model output directory"`
	GoZeroApiDir         string `clop:"long" usage:"go zero api output directory"`
	CopyDir              string `clop:"long" usage:"copy functions output directory"`
	CrudLogicDir         string `clop:"long" usage:"crud logic output directory"`
	Table                string `clop:"long" usage:"table name"`
	Home                 string `clop:"long" usage:"template home directory"`
	Debug                bool   `clop:"long" usage:"debug mode"`
	ModelPkgName         string `clop:"long" usage:"gorm model package name" default:"models"`
	ApiPrefix            string `clop:"long" usage:"go zero api file name prefix"`
	ServiceName          string `clop:"long" usage:"go zero api service name"`
	ApiGroup             string `clop:"long" usage:"go zero api group name"`
	CopyImportPathPrefix string `clop:"long" usage:"copy module import path prefix"`
}

func WriteToFile(dir string, fileName string, data []byte) error {
	// Check if the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Create the directory if it does not exist
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Write the data to the file
	filePath := filepath.Join(dir, fileName)
	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
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
	typeMappings, err := GetTypeMappings(z.Home)
	if err != nil {
		return fmt.Errorf("failed to get type mappings: %w", err)
	}

	res, err := GenerateGormModel(z.ModelPkgName, z.Home, z.Table, columns, typeMappings)
	if err != nil {
		return fmt.Errorf("failed to generate gorm model: %w", err)
	}

	if z.Debug {
		fmt.Println(res)
	}
	if z.ModelDir != "" {
		err = WriteToFile(z.ModelDir, "d_"+z.Table+".go", []byte(res))
		if err != nil {
			return fmt.Errorf("failed to write gorm model file: %w", err)
		}
	}

	serviceName := z.Table
	if z.ServiceName != "" {
		serviceName = z.ServiceName
	}
	groupName := z.ApiGroup
	res, err = GenerateApiService(z.Home, z.Table, columns, typeMappings, "/api/v1", groupName, serviceName, z.Table)
	if err != nil {
		return fmt.Errorf("failed to generate go zero api: %w", err)
	}
	if z.Debug {
		fmt.Println(res)
	}
	if z.GoZeroApiDir != "" {
		apiPrefix := "a_"
		if z.ApiPrefix != "" {
			apiPrefix = z.ApiPrefix
		}
		err = WriteToFile(z.GoZeroApiDir, apiPrefix+z.Table+".api", []byte(res))
		if err != nil {
			return fmt.Errorf("failed to write go zero api file: %w", err)
		}
	}

	res, err = z.GenerateCopyFuncs(z.Home, z.CopyDir, z.Table, columns, typeMappings)
	if err != nil {
		return fmt.Errorf("failed to generate copy funcs: %w", err)
	}
	if z.Debug {
		fmt.Println(res)
	}
	if z.CopyDir != "" {
		err = WriteToFile(z.CopyDir, "c_"+z.Table+".go", []byte(res))
		if err != nil {
			return fmt.Errorf("failed to write copy funcs file: %w", err)
		}
	}

	res, err = z.GenerateCRUDLogic(
		z.Home,
		z.Table,
		columns,
		"api",
		"v1",
		z.Table,
		z.Table,
		"Create"+z.Table,
		"Create"+z.Table+"Request",
		"Create"+z.Table+"Response",
		z.Table, z.Table,
		"Failed to create "+z.Table)
	if err != nil {
		return fmt.Errorf("failed to generate crud logic: %w", err)
	}
	if z.Debug {
		fmt.Println(res)
	}

	return nil
}
