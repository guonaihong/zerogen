package zerogen

import (
	"bytes"
	"fmt"
	"text/template"
)

func GenerateCopyFuncs(homeDir string, tableName string, schema []ColumnSchema, typeMappings map[string]TypeMapping) (string, error) {
	structInfo := schemaToStruct(tableName, schema, typeMappings)
	data := struct {
		Structs []StructInfo
	}{
		Structs: []StructInfo{structInfo},
	}

	tmpl, err := GetCopyTemplate(homeDir)
	if err != nil {
		return "", fmt.Errorf("failed to get copy template: %w", err)
	}

	tpl, err := template.New("copyTemplate").Parse(string(tmpl))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
