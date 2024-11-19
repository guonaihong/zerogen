package zerogen

import (
	"bytes"
	"fmt"
	"go/format"
	"path/filepath"
	"strings"
	"text/template"
)

// GenerateCRUDLogic generates CRUD logic code for a given model
func (z *ZeroGen) GenerateCRUDLogic(
	homeDir string,
	tableName string,
	columnSchema []ColumnSchema,
	packageName,
	serviceImport,
	modelImport,
	typeImport,
	modelStruct,
	modelInstanceName,
	errorMessage string) (string, error) {

	// Get all templates
	createTmpl, err := GetCreateTemplate(homeDir)
	if err != nil {
		return "", fmt.Errorf("failed to get create template: %w", err)
	}

	deleteTmpl, err := GetDeleteTemplate(homeDir)
	if err != nil {
		return "", fmt.Errorf("failed to get delete template: %w", err)
	}

	getByIdTmpl, err := GetGetByIdTemplate(homeDir)
	if err != nil {
		return "", fmt.Errorf("failed to get getbyid template: %w", err)
	}

	getListTmpl, err := GetGetListTemplate(homeDir)
	if err != nil {
		return "", fmt.Errorf("failed to get getlist template: %w", err)
	}

	updateTmpl, err := GetUpdateTemplate(homeDir)
	if err != nil {
		return "", fmt.Errorf("failed to get update template: %w", err)
	}

	idFieldType := ""
	idFieldName := ""
	idColumn := " "

	for _, v := range columnSchema {
		if strings.ToLower(v.ColumnName) == "id" {
			idFieldName = v.ColumnName
			idFieldType = v.ColumnType
		}
	}
	requestType := ToCamelCase(tableName) + "Req"
	responseType := ToCamelCase(tableName) + "Resp"

	// Prepare template data
	data := struct {
		PackageName       string
		ServiceImport     string
		ModelImport       string
		TypeImport        string
		LogicName         string
		RequestType       string
		ResponseType      string
		ModelStruct       string
		ModelInstanceName string
		ErrorMessage      string
		ColumnSchema      []ColumnSchema
		IdColumn          string
		IdFieldName       string
		IdFieldType       string
		Imports           []string
	}{
		PackageName:       packageName,
		ServiceImport:     serviceImport,
		ModelImport:       modelImport,
		TypeImport:        typeImport,
		RequestType:       requestType,
		ResponseType:      responseType,
		ModelStruct:       modelStruct,
		ModelInstanceName: modelInstanceName,
		ErrorMessage:      errorMessage,
		ColumnSchema:      columnSchema,
		IdColumn:          idColumn,
		IdFieldName:       idFieldName,
		IdFieldType:       idFieldType,
	}

	if z.ImportPathPrefix != "" {
		data.Imports = append(data.Imports, z.ImportPathPrefix+"/types")
		data.Imports = append(data.Imports, z.ImportPathPrefix+"/models")
		data.Imports = append(data.Imports, z.ImportPathPrefix+"/svc")
		if z.CopyDir != "" {
			data.Imports = append(data.Imports, z.ImportPathPrefix+"/"+filepath.Base(z.CopyDir))
		}
	}
	// Execute all templates
	var buf bytes.Buffer
	var all bytes.Buffer
	logicName := ToCamelCase(tableName)
	// Create
	data.LogicName = "Create" + logicName
	tmpl, err := template.New("createTemplate").Parse(string(createTmpl))
	if err != nil {
		return "", fmt.Errorf("failed to parse create template: %w", err)
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute create template: %w", err)
	}
	// Format the generated code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("failed to format create template: %w", err)
	}
	buf.Write(formatted)
	all.Write(formatted)

	if z.CrudLogicDir != "" {
		WriteToFile(z.CrudLogicDir, "create_"+tableName+"_logic.go", buf.Bytes())
	}
	buf.Reset()

	// Delete
	data.LogicName = "Delete" + logicName
	tmpl, err = template.New("deleteTemplate").Parse(string(deleteTmpl))
	if err != nil {
		return "", fmt.Errorf("failed to parse delete template: %w", err)
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute delete template: %w", err)
	}
	if z.CrudLogicDir != "" {
		WriteToFile(z.CrudLogicDir, "delete_"+tableName+"_logic.go", buf.Bytes())
	}
	all.Write(buf.Bytes())
	buf.Reset()

	// GetById
	data.LogicName = "Get" + logicName + "ById"
	tmpl, err = template.New("getByIdTemplate").Parse(string(getByIdTmpl))
	if err != nil {
		return "", fmt.Errorf("failed to parse getbyid template: %w", err)
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute getbyid template: %w", err)
	}
	if z.CrudLogicDir != "" {
		WriteToFile(z.CrudLogicDir, "get_"+tableName+"_by_id_logic.go", buf.Bytes())
	}
	all.Write(buf.Bytes())
	buf.Reset()

	// GetList
	data.LogicName = "Get" + logicName + "List"
	tmpl, err = template.New("getListTemplate").Parse(string(getListTmpl))
	if err != nil {
		return "", fmt.Errorf("failed to parse getlist template: %w", err)
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute getlist template: %w", err)
	}
	if z.CrudLogicDir != "" {
		WriteToFile(z.CrudLogicDir, "get_"+tableName+"_list_logic.go", buf.Bytes())
	}
	all.Write(buf.Bytes())
	buf.Reset()

	// Update
	data.LogicName = "Update" + logicName
	tmpl, err = template.New("updateTemplate").Parse(string(updateTmpl))
	if err != nil {
		return "", fmt.Errorf("failed to parse update template: %w", err)
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute update template: %w", err)
	}
	if z.CrudLogicDir != "" {
		WriteToFile(z.CrudLogicDir, "update_"+tableName+"_logic.go", buf.Bytes())
	}
	all.Write(buf.Bytes())
	buf.Reset()

	return all.String(), nil
}
