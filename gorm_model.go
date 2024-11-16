package zerogen

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// ColumnSchema represents the schema of a database column
type ColumnSchema struct {
	ColumnName   string  // Column name, e.g., "id", "created_at"
	ColumnType   string  // Column type, e.g., "uuid", "text"
	Nullable     bool    // Column nullable, e.g., true, false
	Length       int64   // Column length, e.g., 128, 9223372036854775807
	DefaultValue *string // Default value of the column, e.g., "now()"
}

// ToCamelCase converts snake_case to CamelCase
func ToCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}

// GoType returns the Go type for a given database column type based on the framework type mapping,
// using "default" mapping if no specific framework mapping exists.
func GoType(columnType string, nullable bool, typeMappings map[string]TypeMapping, framework string) string {
	frameworkMappings, exists := typeMappings[columnType]
	if !exists {
		return "string" // 默认字符串类型
	}

	goType := frameworkMappings.Default
	if framework == "gorm" && frameworkMappings.Gorm != "" {
		goType = frameworkMappings.Gorm
	} else if framework == "gozero" && frameworkMappings.Gozero != "" {
		goType = frameworkMappings.Gozero
	}

	if nullable && goType != "string" {
		goType = "*" + goType // 如果字段可为空，使用指针类型
	}

	return goType
}

// GenerateStruct generates a struct based on the schema info using Go templates
func GenerateStruct(tableName string, columns []ColumnSchema, typeMappings map[string]TypeMapping) (string, error) {
	tmplContent, err := GetGormModelTemplate()
	if err != nil {
		return "", fmt.Errorf("failed to get gorm model template: %w", err)
	}

	tmpl, err := template.New("structTemplate").Parse(string(tmplContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Prepare data for the template
	data := struct {
		StructName string
		TableName  string
		Columns    []struct {
			FieldName  string
			FieldType  string
			ColumnName string
			Nullable   bool
		}
	}{
		StructName: ToCamelCase(tableName),
		TableName:  tableName,
	}

	// Fill column data with converted names and types
	for _, col := range columns {
		data.Columns = append(data.Columns, struct {
			FieldName  string
			FieldType  string
			ColumnName string
			Nullable   bool
		}{
			FieldName:  ToCamelCase(col.ColumnName),
			FieldType:  GoType(col.ColumnType, col.Nullable, typeMappings, "gorm"),
			ColumnName: col.ColumnName,
			Nullable:   col.Nullable,
		})
	}

	// Execute the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
