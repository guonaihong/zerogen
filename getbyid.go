package zerogen

import (
	"bytes"
	"fmt"
	"text/template"
)

func GenerateGetByIdLogic(
	homeDir string,
	packageName,
	serviceImport,
	modelImport,
	typeImport,
	logicName,
	logicDescription,
	logicFunctionName,
	requestType,
	responseType,
	modelStruct,
	modelInstanceName,
	errorMessage string) (string, error) {
	tmplContent, err := GetGetByIdTemplate(homeDir)
	if err != nil {
		return "", fmt.Errorf("failed to get getbyid template: %w", err)
	}

	tmpl, err := template.New("getByIdTemplate").Parse(string(tmplContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	data := struct {
		PackageName       string
		ServiceImport     string
		ModelImport       string
		TypeImport        string
		LogicName         string
		LogicDescription  string
		LogicFunctionName string
		RequestType       string
		ResponseType      string
		ModelStruct       string
		ModelInstanceName string
		ErrorMessage      string
	}{
		PackageName:       packageName,
		ServiceImport:     serviceImport,
		ModelImport:       modelImport,
		TypeImport:        typeImport,
		LogicName:         logicName,
		LogicDescription:  logicDescription,
		LogicFunctionName: logicFunctionName,
		RequestType:       requestType,
		ResponseType:      responseType,
		ModelStruct:       modelStruct,
		ModelInstanceName: modelInstanceName,
		ErrorMessage:      errorMessage,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
