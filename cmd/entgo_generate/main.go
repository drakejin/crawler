package main

import (
	"strings"
	"text/template"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/rs/zerolog/log"
)

func main() {
	err := entc.Generate("./internal/storage/db/ent/schema", &gen.Config{
		Features: []gen.Feature{
			gen.FeatureUpsert,
			gen.FeatureExecQuery,
		},
		Templates: []*gen.Template{
			gen.MustParse(gen.NewTemplate("static").
				Funcs(template.FuncMap{"title": strings.ToTitle}).
				ParseFiles("./internal/storage/db/ent/templates/debug.tmpl")),
		},
	})
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
