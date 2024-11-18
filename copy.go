package zerogen

import (
	"bytes"
	"fmt"
	"go/format"
	"path/filepath"
	"text/template"
)

func (z *ZeroGen) GenerateCopyFuncs(homeDir string, copyDir string, tableName string, schema []ColumnSchema, typeMappings map[string]TypeMapping) (string, error) {
	structInfo := schemaToStruct(tableName, schema, typeMappings, "copy")

	// Get the last directory name from homeDir
	lastDir := filepath.Base(copyDir)
	data := struct {
		Structs     []StructInfo
		PackageName string
		Imports     []string
	}{
		Structs:     []StructInfo{structInfo},
		PackageName: lastDir,
	}

	data.Imports = structInfo.GoImportPath
	tmpl, err := GetCopyTemplate(homeDir)
	if err != nil {
		return "", fmt.Errorf("failed to get copy template: %w", err)
	}
	if z.CopyImportPathPrefix != "" {
		data.Imports = append(data.Imports, z.CopyImportPathPrefix+"/types")
		data.Imports = append(data.Imports, z.CopyImportPathPrefix+"/models")
	}

	tpl, err := template.New("copyTemplate").Funcs(funcMap).Parse(string(tmpl))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	// Combine standard library imports and other imports
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	formattedOutput, err := format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("failed to format generated code: %w", err)
	}
	return string(formattedOutput), nil
}
