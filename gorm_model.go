package zerogen

import (
	"bytes"
	"fmt"
	"go/format"
	"text/template"
)

// GoType returns the Go type for a given database column type based on the framework type mapping,
// using "default" mapping if no specific framework mapping exists.
func GoType(columnType string,
	nullable bool,
	typeMappings map[string]TypeMapping,
	framework string,
	usePtr bool,
) (string, string, TypeMapping) {
	frameworkMappings, exists := typeMappings[columnType]
	if !exists {
		return "string", "", TypeMapping{} // 默认字符串类型
	}

	goType := frameworkMappings.Default
	if framework == "gorm" && frameworkMappings.Gorm != "" {
		goType = frameworkMappings.Gorm
	} else if framework == "gozero" && frameworkMappings.Gozero != "" {
		goType = frameworkMappings.Gozero
	} else if frameworkMappings.Gorm != "" {
		// 默认情况下，使用 Gorm 的类型
		// 比如copy
		goType = frameworkMappings.Gorm
	}

	if usePtr {
		if nullable && goType != "string" {
			goType = "*" + goType // 如果字段可为空，使用指针类型
		}
	}

	imporPath := frameworkMappings.GormImportPath
	if framework == "gozero" {
		// TODO
		imporPath = ""
	}
	return goType, imporPath, frameworkMappings
}

// GenerateStruct generates a struct based on the schema info using Go templates
func GenerateGormModel(
	modelPkgName string,
	homeDir string,
	tableName string,
	columns []ColumnSchema,
	typeMappings map[string]TypeMapping) (string, error) {
	tmplContent, err := GetGormModelTemplate(homeDir)
	if err != nil {
		return "", fmt.Errorf("failed to get gorm model template: %w", err)
	}

	tmpl, err := template.New("structTemplate").Parse(string(tmplContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Prepare data for the template
	data := struct {
		StructName  string
		TableName   string
		PackageName string
		Imports     []string
		Columns     []struct {
			FieldName  string
			FieldType  string
			ColumnName string
			Nullable   bool
		}
	}{
		StructName:  ToCamelCase(tableName),
		TableName:   tableName,
		PackageName: modelPkgName,
	}

	// Fill column data with converted names and types
	for _, col := range columns {
		goType, goImportPath, _ := GoType(col.ColumnType, col.Nullable, typeMappings, "gorm", false)
		data.Columns = append(data.Columns, struct {
			FieldName  string
			FieldType  string
			ColumnName string
			Nullable   bool
		}{
			FieldName:  ToCamelCase(col.ColumnName),
			FieldType:  goType,
			ColumnName: col.ColumnName,
			Nullable:   col.Nullable,
		})
		if goImportPath != "" {
			data.Imports = append(data.Imports, goImportPath)
		}
	}

	// Execute the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	formattedOutput, err := format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("failed to format generated code: %w", err)
	}
	return string(formattedOutput), nil

}
