package schema

import (
	"context"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"

	myent "github.com/drakejin/crawler/internal/storage/db/ent"
	"github.com/drakejin/crawler/internal/storage/db/ent/hook"
)

// PageSource holds the schema definition for the Keyword entity.
type PageSource struct {
	ent.Schema
}

// Fields of the PageSource.
func (PageSource) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),

		field.String("page_id"),
		field.String("referred_page_id"),

		field.Text("url").
			Default("").
			SchemaType(map[string]string{
				dialect.MySQL: "text",
			}).
			Comment("this mean url"),

		field.Text("referred_url").
			Default("").
			SchemaType(map[string]string{
				dialect.MySQL: "text",
			}).
			Comment("this mean previous referred_url"),

		field.Text("source").
			Default("").
			SchemaType(map[string]string{
				dialect.MySQL: "mediumtext", // maximum 1mb
			}).
			Comment("html view source code"),
	}
}

// Hooks of the Page.
func (PageSource) Hooks() []ent.Hook {
	return []ent.Hook{
		// First hook.
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return hook.PageSourceFunc(func(ctx context.Context, m *myent.PageSourceMutation) (ent.Value, error) {
					m.SetID(uuid.NewString())
					return next.Mutate(ctx, m)
				})
			},
			// Limit the hook only for these operations.
			ent.OpCreate,
		),
	}
}

// Edges of the PageSource.
func (PageSource) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("page", Page.Type).
			Ref("page_source").
			Field("page_id").
			Required().
			Unique(),
	}
}

// Indexes of the PageSource
func (PageSource) Indexes() []ent.Index {
	return []ent.Index{}
}

// Annotations of the PageSource.
func (PageSource) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "page_source",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_0900_ai_ci",
		},
	}
}
