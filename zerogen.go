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

type ZeroGenCore struct {
	Dsn              string `clop:"long" usage:"database dsn" yaml:"dsn"`
	ModelDir         string `clop:"long" usage:"gorm model output directory" yaml:"modelDir"`
	GoZeroApiDir     string `clop:"long" usage:"go zero api output directory" yaml:"goZeroApiDir"`
	CopyDir          string `clop:"long" usage:"copy functions output directory" yaml:"copyDir"`
	CrudLogicDir     string `clop:"long" usage:"crud logic output directory" yaml:"crudLogicDir"`
	Table            string `clop:"long" usage:"table name" yaml:"table"`
	Home             string `clop:"long" usage:"template home directory" yaml:"home"`
	Debug            bool   `clop:"long" usage:"debug mode" yaml:"debug"`
	ModelPkgName     string `clop:"long" usage:"gorm model package name" default:"models" yaml:"modelPkgName"`
	ApiPrefix        string `clop:"long" usage:"go zero api file name prefix" yaml:"apiPrefix"`
	ServiceName      string `clop:"long" usage:"go zero api service name" yaml:"serviceName"`
	ApiGroup         string `clop:"long" usage:"go zero api group name" yaml:"apiGroup"`
	ImportPathPrefix string `clop:"long" usage:"copy module import path prefix" yaml:"importPathPrefix"`
	ApiUrlPrefix     string `clop:"long" usage:"api url prefix" default:"/api/v1" yaml:"apiUrlPrefix"`
	CreateHook       bool   `clop:"long" usage:"generate create hook" yaml:"createHook"`
	UpdateHook       bool   `clop:"long" usage:"generate update hook" yaml:"updateHook"`
	GetListHook      bool   `clop:"long" usage:"generate get list hook" yaml:"getListHook"`
	GetByIdHook      bool   `clop:"long" usage:"generate get by id hook" yaml:"getByIdHook"`
	SaveModels       *bool  `clop:"long" usage:"save models to file" yaml:"saveModels"`
	SaveApi          *bool  `clop:"long" usage:"save api to file" yaml:"saveApi"`
	SaveCopy         *bool  `clop:"long" usage:"save copy to file" yaml:"saveCopy"`
	SaveCrudLogic    *bool  `clop:"long" usage:"save crud logic to file" yaml:"saveCrudLogic"`
}
type ZeroGen struct {
	ConfigFile string `clop:"-f;long" usage:"local configuration file path"`
	ZeroGenCore
}

type copyWithProtocol struct {
	Protocol string `yaml:"protocol"`
	Cmd      string `yaml:"cmd"`
}

type after struct {
	Copy []copyWithProtocol `yaml:"copy"`
}
type zeroConfigWithAction struct {
	Table ZeroGenCore `yaml:"table"`
	After after       `yaml:"after"`
}
type ZeroConfig struct {
	Version string                 `yaml:"version"`
	Global  zeroConfigWithAction   `yaml:"global"`
	Local   []zeroConfigWithAction `yaml:"local"`
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

	var res string
	// 默认保存，如果为false，则不保存
	if z.SaveModels == nil || *z.SaveModels {
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
	}

	serviceName := z.Table
	if z.ServiceName != "" {
		serviceName = z.ServiceName
	}
	groupName := z.ApiGroup
	if z.SaveApi == nil || *z.SaveApi {
		res, err = GenerateApiService(z.Home, z.Table, columns, typeMappings, z.ApiUrlPrefix, groupName, serviceName, z.Table)
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
	}

	if z.SaveCopy == nil || *z.SaveCopy {
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
	}

	if z.SaveCrudLogic == nil || *z.SaveCrudLogic {
		res, err = z.GenerateCRUDLogic(
			columns,
			z.ApiGroup)
		if err != nil {
			return fmt.Errorf("failed to generate crud logic: %w", err)
		}
		if z.Debug {
			fmt.Println(res)
		}
	}

	return nil
}
