package zerogen

import (
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"hasPrefix": strings.HasPrefix,
}
