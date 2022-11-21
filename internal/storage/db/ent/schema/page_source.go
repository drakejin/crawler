package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// PageSource holds the schema definition for the Keyword entity.
type PageSource struct {
	ent.Schema
}

// Fields of the PageSource.
func (PageSource) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique(),

		field.Text("source").
			Default("").
			SchemaType(map[string]string{
				dialect.MySQL: "mediumtext", // maximum 1mb
			}).
			Comment("html view source code"),
	}
}

// Edges of the PageSource.
func (PageSource) Edges() []ent.Edge {
	return []ent.Edge{}
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
