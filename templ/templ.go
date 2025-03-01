package templ

import (
	"html/template"
	"sync"
)

var templates *template.Template
var initOne sync.Once

func NewTemplBlob(pattern string) *template.Template {
	var err error
	initOne.Do(func() {
		templates, err = template.ParseGlob(pattern)

		if err != nil {
			panic(err)
		}
	})

	return templates
}
