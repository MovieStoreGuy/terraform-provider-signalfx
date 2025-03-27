package main

import (
	"embed"
	"os"
	"text/template"

	"github.com/splunk-terraform/terraform-provider-signalfx/internal/feature"
)

//go:embed templates/*.tmpl
var markdowns embed.FS

func main() {
	tpl, err := template.New("content").ParseFS(markdowns, "templates/*.tmpl")
	if err != nil {
		panic(err)
	}
	for _, output := range os.Args[1:] {
		f, err := os.Create(output)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		err = tpl.ExecuteTemplate(f, "feature-preview.md.tmpl", map[string]any{
			"features": feature.GetGlobalRegistry().All(),
		})
		if err != nil {
			panic(err)
		}
	}
}
