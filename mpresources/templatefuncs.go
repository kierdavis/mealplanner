package mpresources

import (
	"html/template"
	"strings"
)

// Additional functions available within templates.
var TemplateFuncMap = template.FuncMap{
	"join": strings.Join,
}

func init() {
	Templates.Funcs(TemplateFuncMap)
}

