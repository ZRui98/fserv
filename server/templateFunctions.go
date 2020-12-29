package server

import (
	"html/template"
)

var templateFunctions = template.FuncMap{
	"sub": sub,
}

func sub(a, b int) int {
	return a - b
}
