package zerogen

import (
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

func GetTypeMappings() (map[string]TypeMapping, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	var typeMappingData []byte
	var typeMappingFile string
	zeroGenDir := filepath.Join(homeDir, ".zero-gen")
	if _, err = os.Stat(zeroGenDir); os.IsNotExist(err) {
		goto fail
	}

	typeMappingFile = filepath.Join(zeroGenDir, "type_mapping.yaml")
	if _, err = os.Stat(typeMappingFile); os.IsNotExist(err) {
		goto fail
	}
	if typeMappingData, err = os.ReadFile(typeMappingFile); err != nil {
		return nil, fmt.Errorf("failed to read type mapping file: %w", err)
	}

fail:
	if typeMappingData == nil {
		typeMappingData = typeMappingYAML
	}
	var typeMapping map[string]TypeMapping
	if err := yaml.Unmarshal(typeMappingData, &typeMapping); err != nil {
		return nil, fmt.Errorf("failed to unmarshal type mappings: %w", err)
	}
	return typeMapping, nil
}
