package view

import (
	"fmt"
	"html/template"

	"ticket-sync-mvc/framework/support"
)

func helpers() template.FuncMap {
	return template.FuncMap{
		"initial": func(value any) string {
			name := []rune(support.Str(value))
			if len(name) == 0 {
				return "?"
			}
			return string(name[0])
		},
		"dict": func(args ...any) support.M {
			values := support.M{}
			for i := 0; i+1 < len(args); i += 2 {
				values[fmt.Sprint(args[i])] = args[i+1]
			}
			return values
		},
	}
}
