package base

import (
	"fmt"
	"html"
	"html/template"
)

var FuncMap = template.FuncMap{
	"raw": func(s string) (r string) {
		r = s
		// fmt.Println(fmt.Sprintf("s:\t%+v", s))
		r = html.EscapeString(s)
		fmt.Println(fmt.Sprintf("r:\t%+v", r))
		return
	},
}
