package zerogen

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type TypeMapping struct {
	Default          string `yaml:"default"`
	Gorm             string `yaml:"gorm"`
	Gozero           string `yaml:"gozero"`
	GormImportPath   string `yaml:"gormImportPath"`
	CopyGoZeroToGorm string `yaml:"copyGoZeroToGorm"`
	CopyGormToGoZero string `yaml:"copyGormToGoZero"`
	CopyPath         string `yaml:"copyPath"`
}

//go:embed cnf/type_mapping.yaml
var typeMappingYAML []byte

//go:embed tmpl/gorm_model.tmpl
var defaultTemplate []byte

//go:embed tmpl/go_zero_api.tmpl
var goZeroApiTemplate []byte

//go:embed tmpl/copy.tmpl
var copyTemplate []byte

//go:embed tmpl/create.tmpl
var createTemplate []byte

//go:embed tmpl/create_hook.tmpl
var createHookTemplate []byte

//go:embed tmpl/delete.tmpl
var deleteTemplate []byte

//go:embed tmpl/getbyid.tmpl
var getByIdTemplate []byte

//go:embed tmpl/getbyid_hook.tmpl
var getByIdHookTemplate []byte

//go:embed tmpl/getlist.tmpl
var getListTemplate []byte

//go:embed tmpl/getlist_hook.tmpl
var getListHookTemplate []byte

//go:embed tmpl/update.tmpl
var updateTemplate []byte

//go:embed tmpl/update_hook.tmpl
var updateHookTemplate []byte

func getFileContent(fileName string, defaultContent []byte, homeDir string) ([]byte, error) {
	var err error
	if homeDir == "" {
		homeDir, err = os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
	}

	var fileData []byte
	var filePath string
	zeroGenDir := filepath.Join(homeDir, ".zero-gen")
	if _, err = os.Stat(zeroGenDir); os.IsNotExist(err) {
		goto useDefault
	}

	filePath = filepath.Join(zeroGenDir, fileName)
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		goto useDefault
	}
	if fileData, err = os.ReadFile(filePath); err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", fileName, err)
	}

useDefault:
	if fileData == nil {
		fileData = defaultContent
	}
	return fileData, nil
}

func GetTypeMappings(homeDir string) (map[string]TypeMapping, error) {
	typeMappingData, err := getFileContent("type_mapping.yaml", typeMappingYAML, homeDir)
	if err != nil {
		return nil, err
	}

	var typeMapping map[string]TypeMapping
	if err := yaml.Unmarshal(typeMappingData, &typeMapping); err != nil {
		return nil, fmt.Errorf("failed to unmarshal type mappings: %w", err)
	}
	return typeMapping, nil
}

func GetGormModelTemplate(homeDir string) ([]byte, error) {
	return getFileContent("gorm_model.tmpl", defaultTemplate, homeDir)
}

func GetGoZeroApiTemplate(homeDir string) ([]byte, error) {
	return getFileContent("go_zero_api.tmpl", goZeroApiTemplate, homeDir)
}

func GetCopyTemplate(homeDir string) ([]byte, error) {
	return getFileContent("copy.tmpl", copyTemplate, homeDir)
}

func GetCreateTemplate(homeDir string) ([]byte, error) {
	return getFileContent("create.tmpl", createTemplate, homeDir)
}

func GetDeleteTemplate(homeDir string) ([]byte, error) {
	return getFileContent("delete.tmpl", deleteTemplate, homeDir)
}

func GetGetByIdTemplate(homeDir string) ([]byte, error) {
	return getFileContent("getbyid.tmpl", getByIdTemplate, homeDir)
}

func GetGetListTemplate(homeDir string) ([]byte, error) {
	return getFileContent("getlist.tmpl", getListTemplate, homeDir)
}

func GetUpdateTemplate(homeDir string) ([]byte, error) {
	return getFileContent("update.tmpl", updateTemplate, homeDir)
}
