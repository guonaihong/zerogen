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
	columnSchema []ColumnSchema,
	packageName string) (string, error) {

	homeDir := z.Home
	tableName := z.Table
	modelInstanceName := ToLowerCamelCase(tableName)
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
			idFieldName = ToCamelCase(v.ColumnName)
			idFieldType = v.ColumnType
		}
	}
	requestType := ToCamelCase(tableName) + "Req"
	responseType := ToCamelCase(tableName) + "Resp"

	// Get type mappings
	typeMappings, err := GetTypeMappings(z.Home)
	if err != nil {
		return "", fmt.Errorf("failed to get type mappings: %w", err)
	}
	structInfo := schemaToStruct(tableName, columnSchema, typeMappings, "gozero")
	// Prepare template data
	data := struct {
		PackageName       string
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
		CopyPkgName       string
		ModelPkgName      string
		Fields            []Field
	}{
		PackageName:       packageName,
		RequestType:       requestType,
		ResponseType:      responseType,
		ModelStruct:       ToCamelCase(tableName),
		ModelInstanceName: modelInstanceName,
		ColumnSchema:      columnSchema,
		IdColumn:          idColumn,
		IdFieldName:       idFieldName,
		IdFieldType:       idFieldType,
		CopyPkgName:       filepath.Base(z.CopyDir),
		ModelPkgName:      z.ModelPkgName,
		Fields:            structInfo.Fields,
	}

	createAndGetImport := []string{}
	if z.ImportPathPrefix != "" {
		createAndGetImport = append(createAndGetImport, z.ImportPathPrefix+"/types")
		createAndGetImport = append(createAndGetImport, z.ImportPathPrefix+"/models")
		createAndGetImport = append(createAndGetImport, z.ImportPathPrefix+"/svc")
		if z.CopyDir != "" {
			createAndGetImport = append(createAndGetImport, z.ImportPathPrefix+"/"+filepath.Base(z.CopyDir))
		}
	}
	deleteImport := []string{}
	if z.ImportPathPrefix != "" {
		deleteImport = append(deleteImport, z.ImportPathPrefix+"/types")
		deleteImport = append(deleteImport, z.ImportPathPrefix+"/models")
		deleteImport = append(deleteImport, z.ImportPathPrefix+"/svc")
	}

	updateImport := []string{}
	if z.ImportPathPrefix != "" {
		updateImport = append(updateImport, z.ImportPathPrefix+"/types")
		updateImport = append(updateImport, z.ImportPathPrefix+"/models")
		updateImport = append(updateImport, z.ImportPathPrefix+"/svc")
		updateImport = append(updateImport, "github.com/fatih/structs")
	}

	data.Imports = createAndGetImport
	// Execute all templates
	var buf bytes.Buffer
	var all bytes.Buffer
	logicName := ToCamelCase(tableName)
	// Create
	data.LogicName = "Create" + logicName
	data.RequestType = "Create" + requestType
	data.ResponseType = "types.Create" + ToCamelCase(tableName) + "Resp"
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
	all.Write(formatted)
	buf.Reset()

	if z.CrudLogicDir != "" {
		WriteToFile(z.CrudLogicDir, "create_"+tableName+"_logic.go", formatted)
	}

	// Delete
	data.LogicName = "Delete" + logicName
	data.RequestType = "Delete" + requestType
	data.ResponseType = "types.BaseResp"
	data.Imports = deleteImport
	tmpl, err = template.New("deleteTemplate").Parse(string(deleteTmpl))
	if err != nil {
		return "", fmt.Errorf("failed to parse delete template: %w", err)
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute delete template: %w", err)
	}
	formatted, err = format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("failed to format delete template: %w", err)
	}
	all.Write(formatted)
	buf.Reset()
	if z.CrudLogicDir != "" {
		WriteToFile(z.CrudLogicDir, "delete_"+tableName+"_logic.go", formatted)
	}

	// GetById
	data.LogicName = "Get" + logicName + "ById"
	data.RequestType = "Get" + ToCamelCase(tableName) + "ByIdReq"
	data.ResponseType = "types.Get" + ToCamelCase(tableName) + "ByIdResp"
	data.Imports = createAndGetImport
	tmpl, err = template.New("getByIdTemplate").Parse(string(getByIdTmpl))
	if err != nil {
		return "", fmt.Errorf("failed to parse getbyid template: %w", err)
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute getbyid template: %w", err)
	}
	formatted, err = format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("failed to format getbyid template: %w", err)
	}

	all.Write(formatted)
	buf.Reset()
	if z.CrudLogicDir != "" {
		WriteToFile(z.CrudLogicDir, "get_"+tableName+"_by_id_logic.go", formatted)
	}
	// GetList
	data.LogicName = "Get" + logicName + "List"
	data.RequestType = "types.Get" + ToCamelCase(tableName) + "ListReq"
	data.ResponseType = "types.Get" + ToCamelCase(tableName) + "ListResp"
	data.Imports = createAndGetImport
	tmpl, err = template.New("getListTemplate").Funcs(funcMap).Parse(string(getListTmpl))
	if err != nil {
		return "", fmt.Errorf("failed to parse getlist template: %w", err)
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute getlist template: %w", err)
	}
	formatted, err = format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("failed to format getlist template: %w", err)
	}
	all.Write(formatted)
	buf.Reset()
	if z.CrudLogicDir != "" {
		WriteToFile(z.CrudLogicDir, "get_"+tableName+"_list_logic.go", formatted)
	}

	// Update
	data.LogicName = "Update" + logicName
	data.RequestType = "Update" + ToCamelCase(tableName) + "Req"
	data.ResponseType = "types.BaseResp"
	data.Imports = updateImport
	tmpl, err = template.New("updateTemplate").Funcs(funcMap).Parse(string(updateTmpl))
	if err != nil {
		return "", fmt.Errorf("failed to parse update template: %w", err)
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute update template: %w", err)
	}

	formatted, err = format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("failed to format update template: %w", err)
	}
	all.Write(formatted)
	buf.Reset()
	if z.CrudLogicDir != "" {
		WriteToFile(z.CrudLogicDir, "update_"+tableName+"_logic.go", formatted)
	}

	return all.String(), nil
}
