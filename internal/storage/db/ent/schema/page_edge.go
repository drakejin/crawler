package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"

	"github.com/drakejin/crawler/internal/storage/db/ent/validate"
)

// PageReferred holds the schema definition for the Keyword entity.
type PageReferred struct {
	ent.Schema
}

// Fields of the Word.
func (PageReferred) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique(),

		field.UUID("source_id", uuid.UUID{}),
		field.UUID("target_id", uuid.UUID{}),

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
func (PageReferred) Edges() []ent.Edge {
	return []ent.Edge{}
}

// Indexes of the Word
func (PageReferred) Indexes() []ent.Index {
	return []ent.Index{
		// unique index.
		index.
			Fields("source_id", "target_id").
			StorageKey("ux_source_id_and_target_id").
			Unique(),
		index.
			Fields("source_id"),
		index.
			Fields("target_id"),
	}
}

// Annotations of the Word.
func (PageReferred) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "page_referred",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_0900_ai_ci",
		},
	}
}
