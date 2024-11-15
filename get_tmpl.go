package zerogen

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type TypeMapping struct {
	Default string `yaml:"default"`
	Gorm    string `yaml:"gorm"`
	Gozero  string `yaml:"gozero"`
}

//go:embed cnf/type_mapping.yaml
var typeMappingYAML []byte

//go:embed tmpl/gorm_model.tmpl
var defaultTemplate []byte

//go:embed tmpl/go_zero_api.tmpl
var goZeroApiTemplate []byte

func getFileContent(fileName string, defaultContent []byte) ([]byte, error) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
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

func GetTypeMappings() (map[string]TypeMapping, error) {
	typeMappingData, err := getFileContent("type_mapping.yaml", typeMappingYAML)
	if err != nil {
		return nil, err
	}

	var typeMapping map[string]TypeMapping
	if err := yaml.Unmarshal(typeMappingData, &typeMapping); err != nil {
		return nil, fmt.Errorf("failed to unmarshal type mappings: %w", err)
	}
	return typeMapping, nil
}

func GetGormModelTemplate() ([]byte, error) {
	return getFileContent("gorm_model.tmpl", defaultTemplate)
}

func GetGoZeroApiTemplate() ([]byte, error) {
	return getFileContent("go_zero_api.tmpl", goZeroApiTemplate)
}
