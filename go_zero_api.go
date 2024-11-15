package zerogen

import (
	"bytes"
	"text/template"
)

type Field struct {
	FieldName string // Go 结构体字段名（驼峰）
	FieldType string // Go 类型
	JSONName  string // JSON 序列化字段名
}

type StructInfo struct {
	StructName string  // Go 结构体名
	TableName  string  // 数据库表名
	Fields     []Field // 表字段信息
}

type ApiService struct {
	Prefix       string       // 路由前缀
	Group        string       // API 分组
	ServiceName  string       // 服务名称
	ModelName    string       // 数据模型名
	ResourceName string       // 路由资源名
	Structs      []StructInfo // 数据表对应的结构体
}

func GenerateApiService(tableName string,
	schema []ColumnSchema,
	typeMappings map[string]TypeMapping,
	prefix string,
	group string,
	serviceName string,
	resourceName string) (string, error) {
	structInfo := schemaToStruct(tableName, schema, typeMappings)
	apiService := ApiService{
		Prefix:       prefix,
		Group:        group,
		ServiceName:  serviceName,
		ModelName:    structInfo.StructName,
		ResourceName: resourceName,
		Structs:      []StructInfo{structInfo},
	}
	tmpl, err := GetGoZeroApiTemplate()
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

func schemaToStruct(tableName string, schema []ColumnSchema, typeMappings map[string]TypeMapping) StructInfo {
	structName := ToCamelCase(tableName)
	fields := []Field{}

	for _, col := range schema {
		goType := GoType(col.ColumnType, col.Nullable, typeMappings)
		// if !ok {
		// 	goType = "interface{}" // 默认类型
		// }
		field := Field{
			FieldName: ToCamelCase(col.ColumnName),
			FieldType: goType,
			JSONName:  col.ColumnName,
		}
		fields = append(fields, field)
	}

	return StructInfo{
		StructName: structName,
		TableName:  tableName,
		Fields:     fields,
	}
}
