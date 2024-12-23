package zerogen

import (
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"hasPrefix": strings.HasPrefix,
	"hasSuffix": strings.HasSuffix,
	"contains":  strings.Contains,
	"toLower":   strings.ToLower,
}
