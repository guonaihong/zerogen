package zerogen

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"
)

func GenerateCopyFuncs(homeDir string, copyDir string, tableName string, schema []ColumnSchema, typeMappings map[string]TypeMapping) (string, error) {
	structInfo := schemaToStruct(tableName, schema, typeMappings)

	// Get the last directory name from homeDir
	lastDir := filepath.Base(copyDir)
	data := struct {
		Structs     []StructInfo
		PackageName string
	}{
		Structs:     []StructInfo{structInfo},
		PackageName: lastDir,
	}

	tmpl, err := GetCopyTemplate(homeDir)
	if err != nil {
		return "", fmt.Errorf("failed to get copy template: %w", err)
	}

	tpl, err := template.New("copyTemplate").Funcs(funcMap).Parse(string(tmpl))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
