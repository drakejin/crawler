package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/drakejin/crawler/internal/storage/db/ent/validate"
	"time"
)

// PageLink holds the schema definition for the Keyword entity.
type PageLink struct {
	ent.Schema
}

// Fields of the Word.
func (PageLink) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique(),

		field.Int64("source_id"),
		field.Int64("target_id"),

		field.Time("created_at").
			Default(time.Now().UTC).
			Annotations(&entsql.Annotation{
				Default: "CURRENT_TIMESTAMP",
			}).
			Immutable().
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).
			Comment("first indexed time"),

		field.String("created_by").
			Annotations(entsql.Annotation{
				Size: 300,
			}).
			Validate(validate.MaxRuneCount(100)).
			Comment("first indexed time by which system"),

		field.Time("updated_at").
			Default(time.Now().UTC).
			UpdateDefault(time.Now().UTC).
			Annotations(&entsql.Annotation{
				Default: "CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP",
			}).
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).
			Comment("modified time"),

		field.String("updated_by").
			Annotations(entsql.Annotation{
				Size: 300,
			}).
			Validate(validate.MaxRuneCount(100)).
			Comment("modified by which system"),
	}
}

// Edges of the Word.
func (PageLink) Edges() []ent.Edge {
	return []ent.Edge{
		//edge.From("page_info", PageInfo.Type).
		//	Ref("page_link").
		//	Field("source_id").
		//	Required(),
		//edge.From("page_info", PageInfo.Type).
		//	Ref("page_link").
		//	Field("target_id").
		//	Required(),
	}
}

// Indexes of the Word
func (PageLink) Indexes() []ent.Index {
	return []ent.Index{
		// unique index.
		index.
			Fields("source_id", "target_id").
			Unique(),
		// unique index.
		index.
			Fields("source_id"),
		index.
			Fields("target_id"),
	}
}

// Annotations of the Word.
func (PageLink) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "page_links",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_0900_ai_ci",
		},
	}
}
