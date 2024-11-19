package zerogen

import (
	"bytes"
	"text/template"
)

type Field struct {
	FieldName        string // Go 结构体字段名（驼峰）
	FieldType        string // Go 类型
	JSONName         string // JSON 序列化字段名
	SnakeCaseName    string // Tag name in snake_case format
	CopyGoZeroToGorm string `yaml:"copyGoZeroToGorm"`
	CopyGormToGoZero string `yaml:"copyGormToGoZero"`
	Nullable         bool
}

type StructInfo struct {
	StructName   string   // Go 结构体名
	TableName    string   // 数据库表名
	Fields       []Field  // 表字段信息
	GoImportPath []string //
}

type ApiService struct {
	Prefix       string       // 路由前缀
	Group        string       // API 分组
	ServiceName  string       // 服务名称
	ModelName    string       // 数据模型名
	ResourceName string       // 路由资源名
	Structs      []StructInfo // 数据表对应的结构体
}

func GenerateApiService(
	homeDir string,
	tableName string,
	schema []ColumnSchema,
	typeMappings map[string]TypeMapping,
	prefix string,
	group string,
	serviceName string,
	resourceName string) (string, error) {
	structInfo := schemaToStruct(tableName, schema, typeMappings, "gozero")
	apiService := ApiService{
		Prefix:       prefix,
		Group:        group,
		ServiceName:  serviceName,
		ModelName:    structInfo.StructName,
		ResourceName: resourceName,
		Structs:      []StructInfo{structInfo},
	}
	tmpl, err := GetGoZeroApiTemplate(homeDir)
	if err != nil {
		return "", err
	}
	tpl, err := template.New("apiTemplate").Parse(string(tmpl))
	if err != nil {
		return "", err
	}

	// 渲染模板
	var result bytes.Buffer
	if err := tpl.Execute(&result, apiService); err != nil {
		return "", err
	}

	return result.String(), nil
}

func schemaToStruct(tableName string, schema []ColumnSchema, typeMappings map[string]TypeMapping, modelName string) StructInfo {
	structName := ToCamelCase(tableName)
	fields := []Field{}

	goImportPaths := []string{}
	for _, col := range schema {
		goType, goImportPath, typeMapping := GoType(col.ColumnType, col.Nullable, typeMappings, modelName, false)
		// if !ok {
		// 	goType = "interface{}" // 默认类型
		// }
		field := Field{
			FieldName:        ToCamelCase(col.ColumnName),
			FieldType:        goType,
			JSONName:         ToLowerCamelCase(col.ColumnName),
			SnakeCaseName:    ToSnakeCase(col.ColumnName),
			CopyGoZeroToGorm: typeMapping.CopyGoZeroToGorm,
			CopyGormToGoZero: typeMapping.CopyGormToGoZero,
			Nullable:         col.Nullable,
		}
		fields = append(fields, field)
		if modelName == "copy" {
			if typeMapping.CopyPath != "" {

				goImportPaths = append(goImportPaths, typeMapping.CopyPath)
			}
		} else if goImportPath != "" {
			goImportPaths = append(goImportPaths, goImportPath)
		}
	}

	return StructInfo{
		StructName:   structName,
		TableName:    tableName,
		Fields:       fields,
		GoImportPath: goImportPaths,
	}
}
